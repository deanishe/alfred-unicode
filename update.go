// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-12

package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deanishe/awgo/util"
)

func runUpdate() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	log.Printf("checking for update ...")
	// 	if err := wf.CheckForUpdate(); err != nil {
	// 		log.Printf("[error] update check failed: %v", err)
	// 	}
	// }()

	fileDepth := func(base, path string) int {
		p, err := filepath.Rel(base, path)
		checkErr(err)
		if p == "." || p == "/" {
			return 0
		}
		n := len(strings.Split(p, "/"))
		return n
	}

	go func() {
		defer wg.Done()
		themeDir := filepath.Join(iconDir, opts.ThemeDir())
		themeDepth := fileDepth(iconDir, themeDir)
		log.Printf("[cache] deleting unused icons ...")
		err := filepath.Walk(iconDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() || path == themeDir {
				return filepath.SkipDir
			}

			if path == iconDir {
				return nil
			}

			// ignore everything below theme level (i.e. XX/0xXXXX - ICON.png)
			if fileDepth(iconDir, path) > themeDepth {
				return filepath.SkipDir
			}

			if isParent(path, themeDir) {
				return nil
			}

			log.Printf("[cache] deleting %q ...", util.PrettyPath(path))
			return filepath.SkipDir
		})
		if err != nil {
			log.Printf("[error] cleaning cache failed: %v", err)
		}
	}()

	wg.Wait()
}

func isParent(parent, child string) bool {
	p := child
	for {
		if p == parent {
			return true
		}
		if p == "/" || p == "." {
			return false
		}
		p = filepath.Dir(p)
	}
}
