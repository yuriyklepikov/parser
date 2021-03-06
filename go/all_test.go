// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func dbg(s string, va ...interface{}) {
	_, fn, fl, _ := runtime.Caller(1)
	fmt.Printf("%s:%d: ", path.Base(fn), fl)
	fmt.Printf(s, va...)
	fmt.Println()
}

var (
	std       = filepath.Join(runtime.GOROOT(), "src")
	tests     = filepath.Join(runtime.GOROOT(), "test")
	whitelist = []string{}
)

func test(t *testing.T, root string) {
	var (
		count int
		size  int64
	)

	if err := filepath.Walk(
		root,
		func(pth string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			ok, err := filepath.Match("*.go", filepath.Base(pth))
			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				return nil
			}

			q := pth[len(root)+1:]
			for _, v := range whitelist {
				if q == v {
					return nil
				}
			}

			dbg("%s", pth)
			_, eerr := parser.ParseFile(token.NewFileSet(), pth, nil, 0)
			_, gerr := ParseFile(pth, nil)

			if g, e := gerr == nil, eerr == nil; g != e {
				t.Fatalf("%q\ng: %v\ne: %v", pth, gerr, eerr)
			}

			count++
			size += info.Size()
			return nil
		}); err != nil {
		t.Fatal(err)
	}

	t.Logf("%d .go files, %d bytes\n", count, size)
}

func TestTestData(t *testing.T) {
	test(t, "testdata")
}

func TestStdlib(t *testing.T) {
	test(t, std)
}

func TestTests(t *testing.T) {
	test(t, tests)
}
