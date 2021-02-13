// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-13

package main

import (
	"encoding/json"
	"log"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
	"github.com/friendsofgo/errors"
)

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

	if opts.Query == "" {
		// force default order if there's no query
		wf.Configure(aw.SuppressUIDs(true))
	}

	for _, f := range fonts {
		var (
			name   = f.String()
			sub    = "↩ to select this font"
			icon   = iconFont
			action = "setfont"
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
			Var("trigger", "").
			Var("hide_alfred", "false")
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
