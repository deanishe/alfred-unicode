// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-06

package main

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNormaliseHexColour verifies normalisation of CSS colours.
func TestNormaliseHexColour(t *testing.T) {
	tests := []struct {
		in    string
		x     string
		valid bool
	}{
		{"#aabbcc", "AABBCC", true},
		{"#abc", "AABBCC", true},
		{"#AABBCC", "AABBCC", true},
		{"#ABC", "AABBCC", true},
		{"aabbcc", "AABBCC", true},
		{"abc", "AABBCC", true},
		{"AABBCC", "AABBCC", true},
		{"ABC", "AABBCC", true},
		{"#aabbccdd", "AABBCCDD", true},
		{"#abcd", "AABBCCDD", true},
		{"#AABBCCDD", "AABBCCDD", true},
		{"#ABCD", "AABBCCDD", true},
		{"aabbccdd", "AABBCCDD", true},
		{"abcd", "AABBCCDD", true},
		{"AABBCCDD", "AABBCCDD", true},
		{"ABCD", "AABBCCDD", true},

		{" aabbccdd", "AABBCCDD", true},
		{"abcd  ", "AABBCCDD", true},
		{" AABBCCDD\n", "AABBCCDD", true},

		{"", "", false},
		{"dean", "", false},
		{"a", "", false},
		{"ab", "", false},
		{"abcde", "", false},
		{"#ab", "", false},
	}
	for _, td := range tests {
		td := td
		t.Run(td.in, func(t *testing.T) {
			s, err := NormaliseHexColour(td.in)
			if !td.valid {
				assert.NotNil(t, err, "invalid colour unexpectedly parsed")
			} else {
				assert.Nil(t, err, "unexpected error")
				assert.Equal(t, td.x, s, "unexpected colour")
			}
		})
	}
}

// TestParseHexColour verifies CSS colour parsing.
func TestParseHexColour(t *testing.T) {
	tests := []struct {
		in    string
		valid bool
		x     color.RGBA
	}{
		{"000000", true, color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}},
		{"ffffff", true, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}},
		{"#000000", true, color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}},
		{"#ffffff", true, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}},
		{"#FFFFFF", true, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}},
		{"#abcdef", true, color.RGBA{R: 0xab, G: 0xcd, B: 0xef, A: 0xff}},
		{"#101112", true, color.RGBA{R: 0x10, G: 0x11, B: 0x12, A: 0xff}},
		{"#d3ad3d", true, color.RGBA{R: 0xd3, G: 0xad, B: 0x3d, A: 0xff}},

		// 3 characters
		{"#FFF", true, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}},
		{"#fff", true, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}},
		{"abc", true, color.RGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 0xff}},
		{"123", true, color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xff}},

		// with alpha
		{"000000cc", true, color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xcc}},
		{"ffffffcc", true, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xcc}},
		{"#FFFd", true, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xdd}},
		{"#fffa", true, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xaa}},
		{"abcd", true, color.RGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 0xdd}},
		{"1234", true, color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 0x44}},

		// invalid
		{"dean", false, color.RGBA{}},
		{"fffff", false, color.RGBA{}},
		{"ff", false, color.RGBA{}},
		{"ffffnn", false, color.RGBA{}},
		{"ffn", false, color.RGBA{}},
	}
	for _, td := range tests {
		td := td
		t.Run(td.in, func(t *testing.T) {
			c, err := ParseHexColour(td.in)
			if !td.valid {
				assert.NotNil(t, err, "invalid colour unexpectedly parsed")
			} else {
				assert.Nil(t, err, "unexpected error")
				assert.Equal(t, td.x, c, "unexpected colour")
			}
		})
	}
}

// TestSlugify verifies generation of FS- and URL-safe slugs.
func TestSlugify(t *testing.T) {
	tests := []struct {
		in, x string
	}{
		// unchanged
		{"", ""},
		{"-", "-"},
		// trimming
		{" ", ""},
		{"  ", ""},
		{" a ", "a"},
		// collapsing dashes
		{"--", "-"},
		{"-a-", "a"},
		{"-A-", "a"},
		{"-A--B-", "a-b"},
		// case
		{"ABC", "abc"},
		{"abc", "abc"},
		{"aBc", "abc"},
		// bad characters
		{"a/b/c", "a-b-c"},
		{":a:b:c/", "a-b-c"},
		{"??***???", "-"},
		// non-ASCII
		{"Jürgen", "jurgen"},
		{"dóómed", "doomed"},
		{"💀Jürgen is dóómed💀", "jurgen-is-doomed"},
		// {"", ""},
		// {"", ""},
		// {"", ""},
	}
	for _, td := range tests {
		td := td
		t.Run(td.in, func(t *testing.T) {
			v := slugify(td.in)
			assert.Equal(t, td.x, v, "unexpected slug")
		})
	}
}
