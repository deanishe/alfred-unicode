// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deanishe/awgo/util"
	"github.com/deanishe/awgo/util/build"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
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

// BuildHelper build swift icon generator
func BuildHelper() error {
	return sh.Run("/usr/bin/swift", "build")
}

// BuildExecutable build alfred-unicode in ./build
func BuildExecutable() error {
	fmt.Println("building workflow executable ...")
	return sh.RunWith(env,
		"go", "build",
		"-tags", "$TAGS",
		"-ldflags", ldflags,
		"-o", "./build/alfred-unicode",
		".",
	)
}

// BuildAssets link assets to ./build
func BuildAssets() error {
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
		// {"../IconGenerator/bin/IconGenerator", "IconGenerator"},
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
	mg.Deps(BuildExecutable, BuildAssets)
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
	mg.Deps(BuildExecutable, BuildHelper)
	mg.Deps(BuildAssets)
	p, err := build.Export(buildDir, distDir)
	if err != nil {
		return err
	}

	fmt.Printf("built workflow file %q\n", p)
	return nil
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
