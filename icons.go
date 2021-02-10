// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-05

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
	"github.com/friendsofgo/errors"
	"golang.org/x/text/unicode/norm"
)

// icons
var (
	iconUpdateAvailable = &aw.Icon{Value: "icons/update-available.png"}
	iconConfig          = &aw.Icon{Value: "icons/config.png"}
	iconError           = &aw.Icon{Value: "icons/error.png"}
	iconWarning         = &aw.Icon{Value: "icons/warning.png"}
	spinnerIcons        = []*aw.Icon{
		{Value: "icons/spinner-0.png"},
		{Value: "icons/spinner-1.png"},
		{Value: "icons/spinner-2.png"},
	}
)

// IconPath returns the relative path for the icon.
func (cp Codepoint) IconPath() string {
	// return cp.Hex[0:2] + "/" + cp.Hex[2:4] + "/" + cp.Name + ".png"
	return fmt.Sprintf("%s/U+%04X - %s.png", cp.Hex[0:2], cp.Dec, cp.Name)
}

// spinnerIcon returns a spinner icon. It rotates by 15 deg on every
// subsequent call. Use with wf.Reload(0.1) to implement an animated
// spinner.
func spinnerIcon() *aw.Icon {
	n := opts.SpinnerCounter
	wf.Var("spinner_counter", fmt.Sprintf("%d", n+1))
	return spinnerIcons[n%3]
}

// NormaliseHexColour converts a hex colour (#abc, #AABBCC, a1b2c3 etc.) to
// a 6- or 8-character form without preceding hash sign.
func NormaliseHexColour(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", errors.New("empty string")
	}
	if s[0] == '#' {
		s = s[1:]
	}

	if len(s) == 3 || len(s) == 4 {
		var s2 string
		for _, c := range s {
			s2 += string(c) + string(c)
		}
		s = s2
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("invalid colour: %q", s)
	}

	if len(b) == 3 || len(b) == 4 {
		return strings.ToUpper(s), nil
	}

	return "", fmt.Errorf("colour hex must be 3, 4, 6 or 8 characters: %s", s)
}

// ParseHexColour parses a CSS hex (e.g. #ffffff) to RGBA.
//
// Input must be with preceding # and have 3/6 (RGB) or 4/8 (RGBA) characters.
func ParseHexColour(s string) (color.RGBA, error) {
	var (
		// default to 100% opaque
		c   = color.RGBA{A: 0xff}
		err error
	)

	if s[0] == '#' {
		s = s[1:]
	}

	if len(s) == 3 || len(s) == 4 {
		var s2 string
		for _, c := range s {
			s2 += string(c) + string(c)
		}
		s = s2
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return c, fmt.Errorf("invalid colour: %q", s)
	}

	switch len(b) {
	case 3:
		c.R = b[0]
		c.G = b[1]
		c.B = b[2]
	case 4:
		c.R = b[0]
		c.G = b[1]
		c.B = b[2]
		c.A = b[3]
	default:
		err = fmt.Errorf("invalid colour: %q", s)
	}

	return c, err
}

func runGenerate() error {
	gen := newIconGenerator(iconDir, aw.IconWorkflow)
	return gen.Generate()
}

type iconGenerator struct {
	Dir         string
	DefaultIcon *aw.Icon
	queue       []Codepoint
}

func newIconGenerator(dir string, defaultIcon *aw.Icon) *iconGenerator {
	return &iconGenerator{
		Dir:         dir,
		DefaultIcon: defaultIcon,
	}
}

// Get returns icon for a Codepoint.
func (gen *iconGenerator) Get(cp Codepoint) *aw.Icon {
	path := gen.iconPath(cp)
	if util.PathExists(path) {
		return &aw.Icon{Value: path}
	}
	gen.queue = append(gen.queue, cp)
	return gen.DefaultIcon
}

// HasQueue returns true if generator's queue isn't empty.
func (gen *iconGenerator) HasQueue() bool {
	return len(gen.queue) > 0
}

// Generate processes the queue.
func (gen *iconGenerator) Generate() error {
	path := gen.queuePath()
	if !util.PathExists(path) {
		return nil
	}

	out, err := util.Run("icongen", path)
	if err != nil {
		return errors.Wrap(err, "generate icons")
	}

	log.Print("icongen output:\n" + string(out))
	if err := os.Remove(path); err != nil {
		return errors.Wrap(err, "remove queue")
	}
	return nil
}

type queueIcon struct {
	IconSize         int     `json:"icon_size"`
	Text             string  `json:"text"`
	Path             string  `json:"filepath"`
	FontName         string  `json:"font_name,omitempty"`
	TextSize         float64 `json:"text_size,omitempty"`
	ResizeFactor     float64 `json:"resize_factor,omitempty"`
	TextColour       string  `json:"text_colour,omitempty"`
	BackgroundColour string  `json:"background_colour,omitempty"`
}

// Saves saves queue if not empty.
func (gen *iconGenerator) Save() error {
	if len(gen.queue) == 0 {
		return nil
	}
	payload := make([]queueIcon, len(gen.queue))

	for i, cp := range gen.queue {
		payload[i] = queueIcon{
			IconSize:         opts.IconSize,
			Text:             string(cp.Dec),
			Path:             gen.iconPath(cp),
			FontName:         opts.FontName,
			TextSize:         opts.TextSize,
			ResizeFactor:     opts.ResizeFactor,
			TextColour:       opts.TextColour,
			BackgroundColour: opts.BackgroundColour,
		}
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshal queue to JSON")
	}

	if err := os.MkdirAll(gen.Dir, 0700); err != nil {
		return errors.Wrap(err, "create icons directory")
	}

	if err := util.WriteFile(gen.queuePath(), data, 0600); err != nil {
		return errors.Wrap(err, "save queue")
	}

	log.Printf("[icons] queued %d icon(s) for generation", len(gen.queue))
	gen.queue = []Codepoint{}
	return nil
}

// iconPath returns icon path for Codepoint.
func (gen *iconGenerator) iconPath(cp Codepoint) string {
	return filepath.Join(gen.Dir, opts.ThemeDir(), cp.IconPath())
}

// queuePath returns path to queue file.
func (gen *iconGenerator) queuePath() string {
	return filepath.Join(gen.Dir, "queue.json")
}

// return "transparent" if alpha is 0, 6-char hex if alpha is 1,
// and 8-char hex if 0 < alpha < 1.
func colourName(hex string) string {
	if hex == "" {
		return "transparent"
	}
	c, _ := ParseHexColour(hex)
	if c.A == 0 {
		return "transparent"
	}
	if len(hex) == 6 {
		return hex
	}
	if c.A == 255 {
		return hex[0:6]
	}
	return hex
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return s
	}

	s = norm.NFD.String(s)

	clean := []rune{}
	var last rune
	for _, r := range s {
		if r > 0x7F { // skip anything above ASCII
			continue
		}
		if (r > 0x2F && r < 0x3A) ||
			(r > 0x60 && r < 0x7A) {
			clean = append(clean, r)
			last = r
		} else if last != '-' {
			clean = append(clean, '-')
			last = '-'
		}
	}
	s = string(clean)
	s = strings.Trim(s, "-")
	if s == "" {
		return "-"
	}
	return s
}
