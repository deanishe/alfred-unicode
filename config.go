// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-13

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
	"github.com/friendsofgo/errors"
)

const (
	defaultTextColour = "FF2FA9"
	defaultFontName   = "Arial Unicode MS"
	defaultIconSize   = 128
	defaultDelimiter  = "\u3009"
)

// user configuration & input
type options struct {
	FontName         string
	TextSize         float64
	ResizeFactor     float64
	IconSize         int
	TextColour       string
	BackgroundColour string
	Delimiter        string

	ActiveApp      string `env:"active_app"`
	SpinnerCounter int    `env:"spinner_counter"`
	CopyAction     string `env:"-"`

	Character    string `env:"-"`
	ConfigChoice string `env:"-"`
	Query        string `env:"-"`
}

func (opts *options) prepare() (err error) {
	if err = fs.Parse(wf.Args()); err != nil {
		return
	}

	if err = wf.Config.To(opts); err != nil {
		return
	}

	if opts.TextColour == "" {
		opts.TextColour = defaultTextColour
	}
	if opts.FontName == "" {
		opts.FontName = defaultFontName
	}
	if opts.IconSize == 0 {
		opts.IconSize = defaultIconSize
	}
	if opts.Delimiter == "" {
		opts.Delimiter = defaultDelimiter
	}

	if opts.TextColour, err = NormaliseHexColour(opts.TextColour); err != nil {
		return
	}

	if opts.BackgroundColour != "" {
		if opts.BackgroundColour, err = NormaliseHexColour(opts.BackgroundColour); err != nil {
			return
		}
	}

	// Tell Copy to Clipboard element to automatically paste.
	// If active_app is set, workflow was called from its Snippet Trigger.
	if opts.ActiveApp != "" {
		opts.CopyAction = "Paste"
		wf.Var("autopaste", "true")
	} else {
		opts.CopyAction = "Copy"
	}

	opts.ConfigChoice = *flagChoose
	opts.Query = fs.Arg(0)
	return
}

// workflow configuration Script Filter.
func runConfig() {
	if opts.Query == "" {
		// force default order if there's no query
		wf.Configure(aw.SuppressUIDs(true))
	}

	switch opts.ConfigChoice {
	case "font":
		runFonts()
		return
	case "text-colour", "background-colour":
		runColour()
		return
	}

	// update status
	var (
		title = "Workflow Is Up To Date"
		sub   = "↩ to check for update now"
		icon  = iconUpdateOK
	)
	if wf.UpdateAvailable() {
		title = "Update Available!"
		sub = "↩ to download and install"
		icon = iconUpdateAvailable
	}

	wf.NewItem(title).
		Subtitle(sub).
		Autocomplete("workflow:update").
		UID("update").
		Icon(icon).
		Valid(false)

	// font
	title = "Font: " + opts.FontName
	if opts.FontIsDefault() {
		title += " (default)"
	}
	it := wf.NewItem(title).
		Subtitle("Change the font used for search result icons").
		Arg("font").
		UID("font").
		Valid(true).
		Icon(iconFont).
		Var("action", "").
		Var("hide_alfred", "").
		Var("trigger", "choose").
		Var("trigger_query", "font")

	if !opts.FontIsDefault() {
		it.NewModifier(aw.ModCmd).
			Subtitle("Reset").
			Icon(iconTrash).
			Arg("").
			Var("action", "setconfig").
			Var("varname", "FONT_NAME")
	}

	// text colour
	title = "Text Colour: #" + opts.TextColour
	if opts.TextColourIsDefault() {
		title += " (default)"
	}
	it = wf.NewItem(title).
		Subtitle("Change the text colour used for search result icons").
		Arg("text-colour").
		UID("text-colour").
		Valid(true).
		Icon(iconColour).
		Var("action", "").
		Var("hide_alfred", "").
		Var("trigger", "choose").
		Var("trigger_query", "text-colour")

	if !opts.TextColourIsDefault() {
		it.NewModifier(aw.ModCmd).
			Subtitle("Reset").
			Icon(iconTrash).
			Arg("").
			Var("action", "setconfig").
			Var("varname", "TEXT_COLOUR")
	}

	// background colour
	title = "Background Colour: #" + opts.BackgroundColour
	if opts.BackgroundColourIsDefault() {
		title = "Background Colour: transparent (default)"
	}
	it = wf.NewItem(title).
		Subtitle("Change the background colour used for search result icons").
		Arg("background-colour").
		UID("background-colour").
		Valid(true).
		Icon(iconColour).
		Var("action", "").
		Var("hide_alfred", "").
		Var("trigger", "choose").
		Var("trigger_query", "background-colour")

	if !opts.BackgroundColourIsDefault() {
		it.NewModifier(aw.ModCmd).
			Subtitle("Reset").
			Icon(iconTrash).
			Arg("").
			Var("action", "setconfig").
			Var("varname", "BACKGROUND_COLOUR")
	}

	// link to docs
	wf.NewItem("Documentation").
		Subtitle("Open workflow website in your browser").
		Arg(website).
		UID("documentation").
		Valid(true).
		Icon(iconWebsite).
		Var("action", "open").
		Var("hide_alfred", "true")

	// link to issue tracker
	wf.NewItem("Report Bug").
		Subtitle("Open workflow issue tracker in your browser").
		Arg(helpURL).
		UID("issue").
		Valid(true).
		Icon(iconIssue).
		Var("action", "open").
		Var("hide_alfred", "true")

	// clear cache
	wf.NewItem("Delete Cached Data").
		Subtitle("Force update of Unicode database and preview images").
		UID("cache").
		Autocomplete("workflow:delcache").
		Valid(false).
		Icon(iconTrash)

	if opts.Query != "" {
		wf.Filter(opts.Query)
	}

	addNavigationItems("config")

	wf.WarnEmpty("No Matching Items", "Try a different query?")
	wf.SendFeedback()
}

func runColour() {
	var (
		name    = "Foreground"
		varname = "TEXT_COLOUR"
	)
	if opts.ConfigChoice == "background-colour" {
		name = "Background"
		varname = "BACKGROUND_COLOUR"
	}
	if opts.Query == "" {
		wf.NewItem("Enter New Colour").
			Subtitle("Enter an HTML/CSS hex colour (3/6 and 4/8 chars supported)").
			Icon(iconColour)
		wf.SendFeedback()
		return
	}

	colour, err := NormaliseHexColour(opts.Query)
	if err != nil {
		wf.NewWarningItem("Invalid Colour: "+opts.Query, "Enter 3/6 or 4/8 hexadecimal characters")
		wf.SendFeedback()
		return
	}

	wf.NewItem(fmt.Sprintf("Set %s Colour To #%s", name, colour)).
		Subtitle("↩ to set").
		Arg(colour).
		Valid(true).
		Icon(colouredIcon(iconColour, colour)).
		Var("action", "setconfig").
		Var("varname", varname).
		Var("hide_alfred", "")

	// wf.WarnEmpty("No Matching Items", "Try a different query?")
	wf.SendFeedback()
}

// Font is a font face supported by macOS.
type Font struct {
	Family string `json:"family"`
	Face   string `json:"face"`
}

// String returns full font name.
func (f Font) String() string {
	s := f.Family
	if f.Face != "" && f.Face != "Regular" {
		s += " " + f.Face
	}
	return s
}

func runFonts() {
	fonts, err := loadFontsCached()
	checkErr(err)
	log.Printf("%d font(s) available", len(fonts))

	for _, f := range fonts {
		var (
			name   = f.String()
			sub    = "↩ to select this font"
			icon   = iconFont
			action = "setconfig"
			valid  = true
		)
		if name == opts.FontName {
			valid = false
			sub = "this font is selected"
			icon = iconFontSelected
			action = ""
		}
		wf.NewItem(name).
			Subtitle(sub).
			Arg(name).
			UID(name).
			Icon(icon).
			Valid(valid).
			Var("action", action).
			Var("varname", "FONT_NAME").
			Var("trigger", "").
			Var("hide_alfred", "")
	}

	if opts.Query != "" {
		wf.Filter(opts.Query)
	}

	addNavigationItems()

	wf.WarnEmpty("No Matching Fonts", "Try a different query?")
	wf.SendFeedback()
}

func loadFonts() ([]Font, error) {
	var fonts []Font
	out, err := util.Run("icongen", "--fonts")
	if err != nil {
		return nil, errors.Wrap(err, "load fonts")
	}

	if err := json.Unmarshal(out, &fonts); err != nil {
		return nil, errors.Wrap(err, "load fonts")
	}

	return fonts, nil
}

func loadFontsCached() ([]Font, error) {
	var (
		reload = func() (interface{}, error) { return loadFonts() }
		fonts  []Font
	)
	if err := wf.Session.LoadOrStoreJSON("fonts.json", reload, &fonts); err != nil {
		return nil, errors.Wrap(err, "cached fonts")
	}
	return fonts, nil
}

// FontIsDefault returns true if font is the default.
func (opts *options) FontIsDefault() bool { return os.Getenv("FONT_NAME") == "" }

// TextColourIsDefault returns true if font is the default.
func (opts *options) TextColourIsDefault() bool { return os.Getenv("TEXT_COLOUR") == "" }

// BackgroundColourIsDefault returns true if font is the default.
func (opts *options) BackgroundColourIsDefault() bool { return os.Getenv("BACKGROUND_COLOUR") == "" }
