// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-08

package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	aw "github.com/deanishe/awgo"
)

// formats Codepoint
type formatFunc func(Codepoint) string

var formatters = []struct {
	uid         string
	description string
	format      formatFunc
}{
	{"character", "UTF-8 character", func(cp Codepoint) string { return string(cp.Dec) }},
	{"name", "Codepoint name", func(cp Codepoint) string { return cp.Name }},
	{"unicode", "Unicode ID", func(cp Codepoint) string { return fmt.Sprintf("U+%04X", cp.Dec) }},
	{"css", "CSS entity", func(cp Codepoint) string { return fmt.Sprintf(`\%06X`, cp.Dec) }},
	{"html-hex", "HTML entity, hexadecimal", func(cp Codepoint) string { return fmt.Sprintf("&#x%X;", cp.Dec) }},
	{"html-dec", "HTML entity, decimal", func(cp Codepoint) string { return fmt.Sprintf("&#%d;", cp.Dec) }},
	{"html-named", "HTML entity, named", formatHTMLEntityNamed},
	{"url-encoded", "URL-encoded", func(cp Codepoint) string { return url.PathEscape(string(cp.Dec)) }},
	{"python", "Python or Go escaped string", formatPython},
	{"javascript", "JavaScript escaped string", formatJS},
}

var formattersMap map[string]int

func init() {
	formattersMap = make(map[string]int, len(formatters))
	for i, f := range formatters {
		formattersMap[f.uid] = i
	}
}

var actions = map[string]string{
	"cmd":   "character",
	"alt":   "python",
	"ctrl":  "unicode",
	"shift": "html-hex",
}

func formatHTMLEntityNamed(cp Codepoint) string {
	if cp.HTMLEntityName != "" {
		return fmt.Sprintf("&%s;", cp.HTMLEntityName)
	}
	return ""
}

func formatPython(cp Codepoint) string {
	if cp.Dec > 0xFFFF {
		return fmt.Sprintf(`\U%08X`, cp.Dec)
	}
	return fmt.Sprintf(`\u%04X`, cp.Dec)
}

func formatJS(cp Codepoint) string {
	if cp.Dec > 0xFFFF {
		return fmt.Sprintf(`\u{%04X}`, cp.Dec) // ES6 only
	}
	return fmt.Sprintf(`\u%04X`, cp.Dec)
}

// match input copied from source code
var (
	rxUnicode = regexp.MustCompile(`^(?i)U\+([A-F0-9]{4,8})$`)
	rxPython4 = regexp.MustCompile(`^\\u([a-fA-F0-9]{4,4})$`)
	rxPython8 = regexp.MustCompile(`^\\U([a-fA-F0-9]{8,8})$`)
	rxHTMLHex = regexp.MustCompile(`^(?i)&#x([A-F0-9]+);?$`)
	rxHTMLDec = regexp.MustCompile(`^(?i)&#([0-9]+);?$`)
	rxCSS     = regexp.MustCompile(`^(?i)\\([A-F0-9]{4,6})$`)
)

func specialQuery(s string) (string, bool) {
	rxs := []*regexp.Regexp{
		rxUnicode,
		rxPython4,
		rxPython8,
		rxHTMLHex,
		rxHTMLDec,
		rxCSS,
	}
	s = strings.TrimSpace(s)
	for _, rx := range rxs {
		if m := rx.FindStringSubmatch(s); m != nil {
			if n, err := parseHex(m[1]); err == nil {
				return fmt.Sprintf("%04X", n), true
			}
		}
	}
	return s, false
}

// Script Filter to show character details
func runDetails() {
	wf.Var("hide_alfred", "true").
		Var("trigger", "").
		Var("action", "copy")

	repo, err := Open(dbPath)
	checkErr(err)
	defer checkClose(repo)

	n, err := parseHex(strings.TrimPrefix(opts.Character, "U+"))
	checkErr(err)
	// if err != nil {
	// 	panic(fmt.Sprintf("unknown character: %s", opts.Character))
	// }
	cps, err := repo.Search(string(n))
	checkErr(err)
	if len(cps) != 1 {
		panic(fmt.Sprintf("unknown character: %s", opts.Character))
	}
	cp := cps[0]

	for _, t := range formatters {
		v := t.format(cp)
		if v == "" {
			continue
		}

		wf.NewItem(v).
			Subtitle(t.description).
			Match(t.description).
			UID(t.uid).
			Arg(v).
			Valid(true).
			Copytext(t.uid)
	}

	if opts.Query != "" {
		wf.Filter(opts.Query)
	}

	wf.WarnEmpty("No Matches", "Try a different query?")
	wf.SendFeedback()
}

func loadActions() error {
	for _, s := range os.Environ() {
		i := strings.Index(s, "=")
		k, v := strings.ToLower(s[:i]), s[i+1:]
		if !strings.HasPrefix(k, "action_") || v == "" {
			continue
		}
		k = strings.Replace(k[7:], "_", "+", -1)
		log.Printf("user-defined action: %s -> %s", k, v)
		actions[k] = v
	}
	return nil
}

func addActions(cp Codepoint, it *aw.Item) {
	for key, uid := range actions {
		i, ok := formattersMap[uid]
		if !ok {
			log.Printf("error: unknown action: %s", uid)
			continue
		}
		f := formatters[i]

		var keys []aw.ModKey
		for _, s := range strings.Split(key, "+") {
			keys = append(keys, aw.ModKey(s))
		}

		v := f.format(cp)
		sub := fmt.Sprintf("%s “%s” (%s)", opts.CopyAction, v, f.description)
		it.NewModifier(keys...).
			Subtitle(sub).
			Arg(v).
			Valid(true).
			Var("action", "copy").
			Var("hide_alfred", "true")
	}
}
