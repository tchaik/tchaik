// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
tchverify is a tool that verifies the "Location" field of tracks in an index by checking that
the associated media file exists.

All configuration is done through command line parameters, see --help flag for details.

NB: this tool cannot verify the content of remote stores, it only supports local stores.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"tchaik.com/index"
	"tchaik.com/index/itl"
)

var itlXML, tchLib string
var trimPathPrefix, addPathPrefix string

func init() {
	flag.StringVar(&itlXML, "itlXML", "", "iTunes Library XML `file`")
	flag.StringVar(&tchLib, "lib", "", "Tchaik library `file`")

	flag.StringVar(&trimPathPrefix, "trim-path-prefix", "", "remove `prefix` from every path")
	flag.StringVar(&addPathPrefix, "add-path-prefix", "", "add `prefix` to every path")
}

func main() {
	flag.Parse()

	l, err := readLibrary()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tracks := l.Tracks()
	fmt.Printf("Checking %d tracks...\n", len(tracks))

	var errCount int
	for _, t := range l.Tracks() {
		loc := t.GetString("Location")
		loc = rewritePath(loc)

		if _, err := os.Stat(loc); err != nil {
			errCount++
			if os.IsNotExist(err) {
				fmt.Printf("file not found: '%v'\n", loc)
				continue
			}
			fmt.Printf("could not stat file '%v': %v\n", loc, err)
		}
	}
	fmt.Printf("Completed: %d error(s).\n", errCount)
}

func rewritePath(path string) string {
	if trimPathPrefix != "" {
		path = strings.TrimPrefix(path, trimPathPrefix)
	}
	if addPathPrefix != "" {
		path = addPathPrefix + path
	}
	return path
}

func readLibrary() (index.Library, error) {
	if itlXML == "" && tchLib == "" {
		return nil, fmt.Errorf("must specify one library file (-itlXML or -lib)")
	}

	if itlXML != "" && tchLib != "" {
		return nil, fmt.Errorf("must only specify one library file (-itlXML or -lib)")
	}

	var l index.Library
	if itlXML != "" {
		f, err := os.Open(itlXML)
		if err != nil {
			return nil, fmt.Errorf("could open iTunes library file: %v", err)
		}
		defer f.Close()
		il, err := itl.ReadFrom(f)
		if err != nil {
			return nil, fmt.Errorf("error parsing iTunes library file: %v", err)
		}

		l = index.Convert(il, "ID")
		return l, nil
	}

	f, err := os.Open(tchLib)
	if err != nil {
		return nil, fmt.Errorf("could not open Tchaik library file: %v", err)
	}
	defer f.Close()

	l, err = index.ReadFrom(f)
	if err != nil {
		return nil, fmt.Errorf("error parsing Tchaik library file: %v\n", err)
	}
	return l, nil
}
