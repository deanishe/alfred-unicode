// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-05

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/mattn/go-sqlite3"
)

const (
	// workflow's GitHub repo (for updates)
	repo    = "deanishe/alfred-unicode"
	website = "https://github.com/deanishe/alfred-unicode/"
	helpURL = "https://github.com/deanishe/alfred-unicode/issues"
)

// CLI
var (
	wf   *aw.Workflow
	opts = &options{}

	fs       = flag.NewFlagSet("alfred-unicode", flag.ContinueOnError)
	flagHelp = fs.Bool("h", false, "show this message and exit")

	// Script Filters
	flagConfig = fs.Bool("config", false, "show workflow configuration")
	flagChoose = fs.String("choose", "", "show configuration option")

	// Other actions
	flagGenerate = fs.Bool("generate", false, "generate queued icons")
	flagUpdate   = fs.Bool("update", false, "check for newer version of the workflow")
	flagSetup    = fs.Bool("setup", false, "download and process Unicode data")
	flagForce    = fs.Bool("f", false, "force update data")
	flagNoop     = fs.Bool("noop", false, "do nothing")
)

func usage() {
	fmt.Fprint(fs.Output(), `alfred-unicode [options...]

Explore Unicode characters in Alfred

`)
	fs.PrintDefaults()
}

func init() {
	fs.SetOutput(os.Stderr)

	aw.IconError = iconError
	aw.IconWarning = iconWarning

	wf = aw.New(
		update.GitHub(repo),
		aw.HelpURL(helpURL),
	)

	zipCachePath = filepath.Join(wf.Cache.Dir, zipCachePath)
	xmlCachePath = filepath.Join(wf.Cache.Dir, xmlCachePath)
	dbPath = filepath.Join(wf.Cache.Dir, dbPath)
	queuePath = filepath.Join(wf.Cache.Dir, queuePath)
	iconDir = filepath.Join(wf.Cache.Dir, iconDir)
	helpFile = filepath.Join(wf.Dir(), helpFile)
}

// main entry point
func run() {
	checkErr(opts.prepare())

	if *flagHelp {
		usage()
		return
	}

	// do nothing operation for script
	if *flagNoop {
		return
	}

	if *flagUpdate {
		wf.Configure(aw.TextErrors(true))
		runUpdate()
		return
	}

	if *flagSetup {
		wf.Configure(aw.TextErrors(true))
		runSetup(*flagForce)
		return
	}

	if *flagGenerate {
		wf.Configure((aw.TextErrors(true)))
		checkErr(runGenerate())
		return
	}

	log.Printf("opts=%#v", opts)

	// common functions for interactive commands
	if wf.UpdateCheckDue() {
		if !wf.IsRunning("update") {
			checkErr(wf.RunInBackground("update", exec.Command(os.Args[0], "-update")))
		}
	}

	if *flagConfig {
		runConfig()
		return
	}

	if opts.Query == "" && wf.UpdateAvailable() {
		wf.Configure(aw.SuppressUIDs(true))
		wf.NewItem("Update Available").
			Subtitle("⇥ or ↩ to update workflow").
			Autocomplete("workflow:update").
			Valid(false).
			Icon(iconUpdateAvailable)
	}

	var settingUp bool
	if setupRequired() {
		if !wf.IsRunning("setup") {
			checkErr(wf.RunInBackground("setup", exec.Command(os.Args[0], "-setup")))
			settingUp = true
		}
	}

	if settingUp || wf.IsRunning("setup") {
		wf.Rerun(0.1)
		wf.NewItem("Building Unicode Database…").
			Subtitle("Please wait a few moments").
			Valid(false).
			Icon(spinnerIcon())
		wf.SendFeedback()
		return
	}

	// parse query
	char, query, backup := parseQuery(opts.Query)
	log.Printf("parsed query: char=%q, query=%q, backup=%v", char, query, backup)
	if backup {
		if b, _ := wf.Cache.Load("previousQuery"); b != nil {
			query = string(b)
		}
		log.Printf("backing up to previous query %q ...", query)
		checkErr(wf.Alfred.RunTrigger("search", query))
		return
	}
	if char != "" {
		opts.Query = query
		opts.Character = char
		runDetails()
		return
	}

	runSearch()
}

// Unicode search Script Filter
func runSearch() {
	var (
		gen   = newIconGenerator(iconDir, spinnerIcon())
		rerun bool
		repo  *Repo
		cps   []Codepoint
		err   error
	)

	// user-configured actions
	checkErr(loadModifiers())

	repo, err = Open(dbPath)
	checkErr(err)
	defer checkClose(repo)

	cps, err = repo.Search(opts.Query)
	checkErr(err)

	for _, cp := range cps {
		icon := gen.Get(cp)
		arg := fmt.Sprintf("0x%s / ", cp.Hex)
		it := wf.NewItem(cp.Name).
			Subtitle(fmt.Sprintf("U+%04X  |  %s", cp.Dec, cp.CategoryName())).
			Autocomplete(arg).
			Arg(arg).
			UID(cp.Hex).
			Valid(false).
			Icon(icon).
			Quicklook(icon.Value).
			Copytext(string(cp.Dec))

		addModifiers(cp, it)
	}

	addNavigationItems("search")

	running := wf.IsRunning("generate")
	if !running && gen.HasQueue() {
		checkErr(gen.Save())
		checkErr(wf.RunInBackground("generate", exec.Command(os.Args[0], "-generate")))
		rerun = true
	}

	if rerun || running {
		wf.Rerun(0.1)
	}

	checkErr(wf.Cache.Store("previousQuery", []byte(opts.Query)))
	wf.WarnEmpty("No Matching Unicode Codepoints", "Try a different query?")
	wf.SendFeedback()
}

// internal navigation items to move between Script Filters.
var navigationItems = []struct {
	keyword  string
	trigger  string
	title    string
	subtitle string
	uid      string
	icon     *aw.Icon
}{
	// {"config", "config", "Configuration", "View/Edit workflow configuration", "configuration", iconConfig},
	{"prefs", "config", "Preferences", "View/Edit workflow configuration", "configuration", iconConfig},
	{"search", "search", "Search Unicode", "Search for Unicode codepoints", "search", aw.IconWorkflow},
}

// Add above actions to results if query is at least 3 characters long and
// matches the actions keyword. Pass keywords of items to be ignored (e.g. the
// keyword of the script filter you're calling the function from).
func addNavigationItems(ignore ...string) {
	if len(opts.Query) < 3 {
		return
	}
	query := strings.ToLower(opts.Query)
	ignoreMap := make(map[string]struct{}, len(ignore))
	for _, s := range ignore {
		ignoreMap[s] = struct{}{}
	}
	for _, ni := range navigationItems {
		if _, ok := ignoreMap[ni.keyword]; ok {
			continue
		}
		if strings.HasPrefix(ni.keyword, query) {
			wf.NewItem(ni.title).
				Subtitle(ni.subtitle).
				UID(ni.uid).
				Arg(ni.trigger).
				Valid(true).
				Icon(ni.icon).
				Var("trigger", ni.trigger)
		}
	}
	if strings.HasPrefix("help", query) {
		wf.NewItem("Workflow Documentation").
			Subtitle("Open workflow documentation in your browser").
			UID("help").
			Arg(helpFile).
			Valid(true).
			Icon(iconDocs).
			Var("action", "open").
			Var("hide_alfred", "true").
			Var("trigger", "")
	}
}

func main() {
	wf.Run(run)
}

func checkErr(err error) {
	if err != nil {
		// If it's a syntax error from SQLite, try to  show a warning
		// instead of an error. This should be safe, as only Script Filters
		// use the database.
		var syntaxError sqlite3.Error
		if errors.As(err, &syntaxError) {
			s := strings.TrimPrefix(syntaxError.Error(), "fts5: ")
			wf.NewWarningItem("Invalid Search Query", s)
			wf.SendFeedback()
		}
		panic(err)
	}
}

func checkClose(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}
