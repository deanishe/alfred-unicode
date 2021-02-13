// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-05

package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/deanishe/awgo/util"
)

// list of all Unicode codepoints
const xmlURL = "https://www.unicode.org/Public/UCD/latest/ucdxml/ucd.nounihan.flat.zip"

// Codepoint is a Unicode codepoint.
type Codepoint struct {
	Dec            int32   `db:"dec"`
	Hex            string  `db:"hex"`
	Category       string  `db:"category"`
	Name           string  `db:"name"`
	Aliases        string  `db:"aliases"`
	HTMLEntityName string  `db:"html_entity"`
	Rank           float64 `db:"rank"`
}

// String implements Stringer.
func (cp Codepoint) String() string {
	return fmt.Sprintf("%08X %q", cp.Dec, cp.Name)
}

// IsZero returns true if Codepoint is zero.
func (cp Codepoint) IsZero() bool { return cp.Hex == "" }

func setupRequired() bool { return !util.PathExists(dbPath) }

// ensure Unicode data and database exist.
func runSetup(force bool) {
	var (
		db  *Repo
		err error
	)
	if !util.PathExists(zipCachePath) || force {
		log.Printf("[setup] downloading Unicode data ...")
		checkErr(download(xmlURL, zipCachePath))
	}

	if newerThan(zipCachePath, dbPath) {
		log.Printf("[setup] extracting Unicode data ...")
		if err := unzip(xmlCacheName, xmlCachePath, zipCachePath); err != nil {
			log.Printf("unzip failed: %v", err)
			if strings.Contains(err.Error(), "not a valid zip file") {
				_ = os.Remove(zipCachePath)
				_ = os.Remove(xmlCachePath)
			}
			panic(err)
		}

		log.Printf("[setup] building database ...")
		db, err = Open(dbPath)
		checkErr(err)
		defer db.Close()

		var cps []Codepoint
		cps, err = parseXML(xmlCachePath)
		checkErr(err)
		checkErr(db.Import(cps))
		checkErr(os.Remove(xmlCachePath))
	}
}

type xmlPayload struct {
	Description string `xml:"description"`
	Repetoire   struct {
		Chars []struct {
			Name      string `xml:"na,attr"`
			Name1     string `xml:"na1,attr"`
			Category  string `xml:"gc,attr"`
			Codepoint string `xml:"cp,attr"`
			Aliases   []struct {
				Alias string `xml:"alias,attr"`
				Type  string `xml:"type,attr"`
			} `xml:"name-alias"`
		} `xml:"char"`
		XMLName xml.Name `xml:"repertoire"`
	}
	XMLName xml.Name `xml:"ucd"`
}

func parseXML(path string) ([]Codepoint, error) {
	log.Printf("[setup] loading %q ...", util.PrettyPath(path))
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open XML file: %w", err)
	}
	defer f.Close()

	var (
		dec = xml.NewDecoder(f)
		pl  = xmlPayload{}
		cps = []Codepoint{}
	)
	if err := dec.Decode(&pl); err != nil {
		return nil, fmt.Errorf("parse XML: %w", err)
	}

	for i, v := range pl.Repetoire.Chars {
		if v.Codepoint == "" || v.Name == "CJK UNIFIED IDEOGRAPH-#" {
			continue
		}
		aliases := make([]string, len(v.Aliases))
		for i, a := range v.Aliases {
			aliases[i] = a.Alias
		}
		name := v.Name
		if name == "" {
			name = v.Name1
		}

		if name == "" {
			continue
		}

		n, err := parseHex(v.Codepoint)
		if err != nil {
			return nil, err
		}
		cp := Codepoint{
			Hex:            fmt.Sprintf("%04X", n),
			Dec:            n,
			Category:       v.Category,
			Name:           name,
			Aliases:        strings.Join(aliases, ", "),
			HTMLEntityName: NamedHTMLEntities[fmt.Sprintf("%05X", n)],
		}
		if i == 0 {
			log.Printf("%#v", cp)
		}
		cps = append(cps, cp)
	}
	log.Printf("[setup] %d codepoints for %q from: %s",
		len(cps), pl.Description, util.PrettyPath(path))
	return cps, nil
}

func unzip(name, dest, zipfile string) error {
	r, err := zip.OpenReader(zipfile)
	if err != nil {
		return fmt.Errorf("unzip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		log.Printf("[setup] %s: %s", zipfile, f.Name)
		if f.Name == name {
			var (
				out *os.File
				rc  io.ReadCloser
			)
			if out, err = os.Create(dest); err != nil {
				return fmt.Errorf("unzip: extract file: %w", err)
			}
			defer out.Close()
			if rc, err = f.Open(); err != nil {
				return fmt.Errorf("unzip: extract file: %w", err)
			}
			defer rc.Close()
			if _, err = io.Copy(out, rc); err != nil {
				return fmt.Errorf("unzip: write data: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("unzip: file %q not found in %q", name, zipfile)
}

func newerThan(name1, name2 string) bool {
	var (
		fi1, fi2 os.FileInfo
		err      error
	)
	if fi1, err = os.Stat(name1); err != nil {
		return false
	}
	if fi2, err = os.Stat(name2); err != nil {
		return true
	}
	return fi1.ModTime().After(fi2.ModTime())
}

// download saves file at URL to path.
func download(url, path string) error {
	r, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("retrieve URL %q: %w", url, err)
	}
	log.Printf("[http] [%d] %s", r.StatusCode, url)
	defer r.Body.Close()
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create XML cache file: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, r.Body); err != nil {
		return fmt.Errorf("write XML cache file: %w", err)
	}
	log.Printf("[http] saved %q", util.PrettyPath(path))
	return nil
}

func parseHex(hex string) (int32, error) {
	n, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}

// CategoryName returns the full name of Codepoint's category.
func (cp Codepoint) CategoryName() string {
	switch cp.Category {
	case "Cc":
		return "Control"
	case "Cf":
		return "Format"
	case "Co":
		return "Private Use"
	case "Cs":
		return "Surrrogate"
	case "Ll":
		return "Lowercase Letter"
	case "Lm":
		return "Modifier Letter"
	case "Lo":
		return "Other Letter"
	case "Lt":
		return "Titlecase Letter"
	case "Lu":
		return "Uppercase Letter"
	case "Mc":
		return "Spacing Mark"
	case "Me":
		return "Enclosing Mark"
	case "Mn":
		return "Nonspacing Mark"
	case "Nd":
		return "Decimal Number"
	case "Nl":
		return "Letter Number"
	case "No":
		return "Other Number"
	case "Pc":
		return "Connector Punctuation"
	case "Pd":
		return "Dash Punctuation"
	case "Pe":
		return "Close Punctuation"
	case "Pf":
		return "Final Punctuation"
	case "Pi":
		return "Initial Punctuation"
	case "Po":
		return "Other Punctuation"
	case "Ps":
		return "Open Punctuation"
	case "Sc":
		return "Currency Symbol"
	case "Sk":
		return "Modifier Symbol"
	case "Sm":
		return "Math Symbol"
	case "So":
		return "Other Symbol"
	case "Zl":
		return "Line Separator"
	case "Zp":
		return "Paragraph Separator"
	case "Zs":
		return "Space Separator"
	default:
		return ""
	}
}
