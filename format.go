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
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
)

// formats Codepoint
type formatFunc func(Codepoint) string

type formatter struct {
	uid         string
	description string
	format      formatFunc
	enabled     bool
}

var formatters = []formatter{
	{"character", "UTF-8 character", func(cp Codepoint) string { return string(cp.Dec) }, true},
	{"name", "Codepoint name", func(cp Codepoint) string { return cp.Name }, true},
	{"unicode", "Unicode ID", func(cp Codepoint) string { return fmt.Sprintf("U+%04X", cp.Dec) }, true},
	{"hexadecimal", "Hexadecimal literal", func(cp Codepoint) string { return fmt.Sprintf("0x%X", cp.Dec) }, true},
	{"decimal", "Decimal literal", func(cp Codepoint) string { return fmt.Sprintf("%d", cp.Dec) }, true},
	{"python", "Python or Go escaped string", formatPython, true},
	{"urlencoded", "URL-encoded", func(cp Codepoint) string { return url.PathEscape(string(cp.Dec)) }, true},
	{"javascript", "JavaScript escaped string", formatJS, true},
	{"css", "CSS entity", func(cp Codepoint) string { return fmt.Sprintf(`\%06X`, cp.Dec) }, true},
	{"htmlhex", "HTML entity, hexadecimal", func(cp Codepoint) string { return fmt.Sprintf("&#x%X;", cp.Dec) }, true},
	{"htmldec", "HTML entity, decimal", func(cp Codepoint) string { return fmt.Sprintf("&#%d;", cp.Dec) }, true},
	{"htmlnamed", "HTML entity, named", formatHTMLEntityNamed, true},
}

var formattersMap map[string]int

func init() {
	disabled := map[string]bool{}
	for _, s := range strings.Split(os.Getenv("DISABLED"), ",") {
		disabled[s] = true
	}
	formattersMap = make(map[string]int, len(formatters))
	for i, f := range formatters {
		formattersMap[f.uid] = i
		formatters[i].enabled = !disabled[f.uid]
	}
}

/*
func loadFormatters() {
	disabled := map[string]bool{}
	for _, s := range strings.Split(opts.Disabled, ",") {
		disabled[s] = true
	}
	for i, fm := range formatters {
		formatters[i].enabled = disabled[fm.uid]
	}
}
*/

var modifiers = map[string]string{
	"cmd":   "character",
	"opt":   "python",
	"ctrl":  "hexadecimal",
	"shift": "unicode",
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
	if cp.Dec <= 0xFF {
		return fmt.Sprintf(`\x%02X`, cp.Dec)
	}
	if cp.Dec <= 0xFFFF {
		return fmt.Sprintf(`\u%04X`, cp.Dec)
	}
	if opts.ES6Mode {
		return fmt.Sprintf(`\u{%X}`, cp.Dec)
	}
	return jsSurrogatePair(cp.Dec)
}

// encode rune to JS surrogate pair
// https://unicodebook.readthedocs.io/unicode_encodings.html#utf-16-surrogate-pairs
func jsSurrogatePair(r rune) string {
	if r <= 0xFFFF {
		return fmt.Sprintf(`\u%04X`, r)
	}
	code := r - 0x10000
	one := 0xD800 | (code >> 10)
	two := 0xDC00 | (code & 0x3FF)
	return fmt.Sprintf(`\u%04X\u%04X`, one, two)
}

// match various escape formats
var (
	// contain hex/decimal values
	rxUnicode    = regexp.MustCompile(`^(?i)U\+([A-F0-9]{4,8})$`)
	rxPython4    = regexp.MustCompile(`^\\u([a-fA-F0-9]{4,4})$`)
	rxPython8    = regexp.MustCompile(`^\\U([a-fA-F0-9]{8,8})$`)
	rxJSHex      = regexp.MustCompile(`^\\x([a-fA-F0-9]{2,2})$`)
	rxES6        = regexp.MustCompile(`^\\u\{([a-fA-F0-9]{1,8})\}$`)
	rxHTMLHex    = regexp.MustCompile(`^(?i)&#x([A-F0-9]+);?$`)
	rxHTMLDec    = regexp.MustCompile(`^(?i)&#([0-9]+);?$`)
	rxCSS        = regexp.MustCompile(`^(?i)\\([A-F0-9]{4,6})$`)
	rxHexLiteral = regexp.MustCompile(`^(?i)0x([A-F0-9]{2,8})$`)
	// other
	rxHTMLNamed = regexp.MustCompile(`^(?i)&([A-Z]+);?$`)
	// rxFieldSearch = regexp.MustCompile(`^(?i)dec:[A-Z0-9]+$`)
)

// Parse specially-formatted queries. If stop is true, the returned query
// should not be further processed.
func specialQuery(query string) (newQuery string, stop bool) {
	// contain hex values
	rxs := []*regexp.Regexp{
		rxUnicode,
		rxES6,
		rxJSHex,
		rxPython4,
		rxPython8,
		rxHTMLHex,
		rxHTMLDec,
		rxCSS,
		rxHexLiteral,
	}
	s := strings.TrimSpace(query)
	for _, rx := range rxs {
		if m := rx.FindStringSubmatch(s); m != nil {
			if n, err := parseHex(m[1]); err == nil {
				return fmt.Sprintf("hex:%04X", n), true
			}
		}
	}
	// named HTML entity
	if m := rxHTMLNamed.FindStringSubmatch(s); m != nil {
		return fmt.Sprintf("entity:%s", m[1]), true
	}
	var prefix string
	// if s is valid hex, add it as a normalised OR query
	if n, err := parseHex(s); err == nil {
		prefix += fmt.Sprintf("hex:%04X OR ", n)
	}
	// if s is valid decimal, add it as a normalised OR query
	if n, err := strconv.ParseUint(s, 10, 32); err == nil {
		prefix += fmt.Sprintf("dec:%d OR ", n)
	}
	s = prefix + query

	/*
		// append space to prevent wildcard being added to digits
		if rxFieldSearch.MatchString(s) {
			s += " "
		}
	*/
	return s, false
}

// Script Filter to show character details
func runDetails() {
	if opts.Query == "" {
		wf.Configure(aw.SuppressUIDs(true))
	}

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

	for _, f := range formatters {
		if !f.enabled {
			continue
		}

		v := f.format(cp)
		if v == "" {
			continue
		}

		wf.NewItem(v).
			Subtitle(f.description).
			Match(f.description).
			UID(f.uid).
			Arg(v).
			Valid(true).
			Copytext(f.uid)
	}

	if opts.Query != "" {
		wf.Filter(opts.Query)
	}

	wf.WarnEmpty("No Matches", "Try a different query?")
	wf.SendFeedback()
}

func loadModifiers() error {
	for _, s := range os.Environ() {
		i := strings.Index(s, "=")
		k, v := strings.ToLower(s[:i]), s[i+1:]
		if !strings.HasPrefix(k, "mod_") || v == "" {
			continue
		}
		k = strings.Replace(k[4:], "_", "+", -1)
		log.Printf("user-defined action: %s -> %s", k, v)
		modifiers[k] = v
	}
	return nil
}

func addModifiers(cp Codepoint, it *aw.Item) {
	for key, uid := range modifiers {
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
