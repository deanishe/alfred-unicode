// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/deanishe/awgo/util"
	"github.com/deanishe/awgo/util/build"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	blackfriday "github.com/russross/blackfriday/v2"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

var (
	info     *build.Info
	env      map[string]string
	ldflags  string
	workDir  string
	buildDir = "./build"
	distDir  = "./dist"
	iconsDir = "./icons"
)

func init() {
	var err error
	if info, err = build.NewInfo(); err != nil {
		panic(err)
	}
	if workDir, err = os.Getwd(); err != nil {
		panic(err)
	}
	env = info.Env()
	env["VERSION"] = info.Version
	env["PKG"] = "main"
	env["TAGS"] = "fts5"
	ldflags = `-X "$PKG.version=$VERSION"`
}

func mod(args ...string) error {
	argv := append([]string{"mod"}, args...)
	return sh.RunWith(env, "go", argv...)
}

// Aliases are mage command aliases.
var Aliases = map[string]interface{}{
	"b": Build,
	"c": Clean,
	"d": Dist,
	"l": Link,
}

// Helper build swift icon generator
func Helper() error {
	return sh.Run("/usr/bin/swift", "build")
}

// Executable build alfred-unicode in ./build
func Executable() error {
	fmt.Println("building workflow executable ...")
	return sh.RunWith(env,
		"go", "build",
		"-tags", "$TAGS",
		"-ldflags", ldflags,
		"-o", "./build/alfred-unicode",
		".",
	)
}

// Docs generate README.html
func Docs() error {
	fmt.Println("generating README.html from README.md ...")
	md, err := ioutil.ReadFile("README.md")
	if err != nil {
		return err
	}
	// r := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
	// 	Flags: blackfriday.CommonHTMLFlags | blackfriday.CompletePage,
	// })
	html := blackfriday.Run(md,
		blackfriday.WithExtensions(
			blackfriday.CommonExtensions,
		),
		// blackfriday.WithRenderer(r),
	)
	header := fmt.Sprintf(htmlHeader, css)
	html = append([]byte(header), html...)
	html = append(html, []byte(htmlFooter)...)
	// []byte(htmlFooter)
	file, err := os.Create("README.html")
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(html); err != nil {
		return err
	}
	return nil
}

// Assets link assets to ./build
func Assets() error {
	fmt.Println("building assets ...")
	globs := build.Globs(
		"*.png",
		"info.plist",
		"*.html",
		"README.md",
		"LICENCE.txt",
		"icons/*.png",
	)

	if err := build.SymlinkGlobs(buildDir, globs...); err != nil {
		return err
	}

	other := []struct {
		src, dst string
	}{
		{"icons/config.png", "12890D35-22B2-47AB-A54A-3BF0C1F8841D.png"},
		{"icons/config.png", "B564A9DA-3F70-471A-8DF3-8A537CBF2975.png"},
		{"icons/font.png", "F3FBD6A3-BEA9-4D64-8F7C-AC67C1E8B5EC.png"},
		{"../.build/debug/icongen", "icongen"},
	}

	for _, f := range other {
		src, dst := filepath.Join(buildDir, f.src), filepath.Join(buildDir, f.dst)
		if err := build.Symlink(dst, src, true); err != nil {
			return err
		}
	}
	return nil
}

// Build build workflow in ./build
func Build() error {
	mg.Deps(cleanBuild)
	mg.Deps(Executable, Assets)
	return nil
}

// Run run workflow
func Run() error {
	mg.Deps(Build)
	fmt.Println("running ...")
	return sh.RunWith(env, buildDir+"/alfred-unicode", "-h")
}

// Dist build an .alfredworkflow file in ./dist
func Dist() error {
	mg.SerialDeps(Clean)
	mg.Deps(Executable, Helper)
	mg.Deps(Assets)
	p, err := build.Export(buildDir, distDir)
	if err != nil {
		return err
	}

	fmt.Printf("built workflow file %q\n", p)
	return nil
}

//Format run gofmt -s on source files.
func Format() error {
	return sh.Run("gofmt", "-s", "-w", ".")
}

// Config display configuration
func Config() {
	fmt.Println("     Workflow name:", info.Name)
	fmt.Println("         Bundle ID:", info.BundleID)
	fmt.Println("  Workflow version:", info.Version)
	fmt.Println("  Preferences file:", info.AlfredPrefsBundle)
	fmt.Println("       Sync folder:", info.AlfredSyncDir)
	fmt.Println("Workflow directory:", info.AlfredWorkflowDir)
	fmt.Println("    Data directory:", info.DataDir)
	fmt.Println("   Cache directory:", info.CacheDir)
}

// Link symlinks ./build directory to Alfred's workflow directory.
func Link() error {
	mg.Deps(Build)

	fmt.Println("linking ./build to workflow directory ...")
	target := filepath.Join(info.AlfredWorkflowDir, info.BundleID)
	// fmt.Printf("target: %s\n", target)

	if util.PathExists(target) {
		fmt.Println("removing existing workflow ...")
	}
	// try to remove it anyway, as dangling symlinks register as existing
	if err := os.RemoveAll(target); err != nil && !os.IsNotExist(err) {
		return err
	}

	src, err := filepath.Abs(buildDir)
	if err != nil {
		return err
	}
	return build.Symlink(target, src, true)
}

// Deps ensure dependencies
func Deps() error {
	mg.Deps(cleanDeps)
	fmt.Println("downloading deps ...")
	return mod("download")
}

// Vendor copy dependencies to ./vendor
func Vendor() error {
	mg.Deps(Deps)
	fmt.Println("vendoring deps ...")
	return mod("vendor")
}

// Clean remove build files
func Clean() {
	fmt.Println("cleaning ...")
	mg.Deps(cleanBuild, cleanHelper, cleanMage, cleanDeps)
}

func cleanDeps() error {
	return mod("tidy", "-v")
}

// remove & recreate directory
func cleanDir(name string) error {
	if err := sh.Rm(name); err != nil {
		return err
	}
	return os.MkdirAll(name, 0755)
}

func cleanBuild() error {
	return cleanDir(buildDir)
}

func cleanHelper() error {
	return cleanDir("icongen/.build")
}

func cleanMage() error {
	return sh.Run("mage", "-clean")
}

// CleanIcons delete all generated icons from ./icons
func CleanIcons() error {
	return cleanDir(iconsDir)
}

const (
	htmlHeader = `<!DOCTYPE html>
<html>
<head>
	<title>Alfred Unicode Search</title>
	<style>
	%s
	</style>
</head>
<body>
	<div id="wrapper">
`
	htmlFooter = `
	</div>
</body>
</html>
`
	css = `
/*** Start Vanilla CSS 1.0.2 */
body{font:1em/1.5em sans-serif;}pre,code,tt{font:1em/1.5em ‘Andale Mono’,’Lucida Console’,monospace;}h1,h2,h3,h4,h5,h6,b,strong{font-weight:bold;}em,i,dfn{font-style:italic;}dfn{font-weight:bold;}p,code,pre,kbd{margin:0 0 1.5em 0;}blockquote{margin:0 1.5em 1.5em 1.5em;}cite{font-style:italic;}li ul,li ol{margin:0 1.5em;}ul,ol{margin:0 1.5em 1.5em 1.5em;}ul{list-style-type:disc;}ol{list-style-type:decimal;}ol ol{list-style:upper-alpha;}ol ol ol{list-style:lower-roman;}ol ol ol ol{list-style:lower-alpha;}dl{margin:0 0 1.5em 0;}dl dt{font-weight:bold;}dd{margin-left:1.5em;}table{margin-bottom:1.4em;width:100%;}th{font-weight:bold;}th,td,caption{padding:4px 10px 4px 5px;}tfoot{font-style:italic;}sup,sub{line-height:0;}abbr,acronym{border-bottom:1px dotted;}address{margin:0 0 1.5em;font-style:italic;}del{text-decoration:line-through;}pre{margin:1.5em 0;white-space:pre;}img.centered,.aligncenter,div.aligncenter{display:block;margin-left:auto;margin-right:auto;}img.alignright{display:inline;}img.alignleft{display:inline;}.alignright{float:right;margin-left:10px;}.alignleft{float:left;margin-right:10px;}img { max-width: 100%; }* html .clearfix{height:1%;}*+html .clearfix{display:inline-block;}.clearfix:after{content:”.”;display:block;height:0;clear:both;visibility:hidden;}* html .group{height:1%;}*+html .group{display:inline-block;}.group:after{content:”.”;display:block;height:0;clear:both;visibility:hidden;}

#wrapper {
	margin: 2em auto;
	max-width: 900px;
}
code {
	background-color: #f0f0f8;
	border-radius: 0.05em;
	padding: 0.15em 0.3em;
}
ul {
	padding-left: 0;
}
`
)
