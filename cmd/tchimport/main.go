// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
tchimport is a tool which builds Tchaik Libraries (metadata indexes) from iTunes Library XML files or
alternatively by reading metadata from audio files within a directory tree.

Importing large iTunes XML Library files is recommended: the Tchaik library has a much smaller set of data attributes
for each track (so a much smaller memory footprint).

  tchimport -itlXML <itunes-library> -out lib.tch

Alternatively you can specify a path which will be transversed. All supported audio files within this path
(.mp3, .m4a, .flac - ID3.v1,2.{2,3,4}, MP4 and FLAC) will be scanned for metadata. Only tracks which have readable
metadata will be added to the library.  Any errors are logged to stdout. As no other unique identifying data is know,
the SHA1 sum of the file path is used as the TrackID.

  tchimport -path <directory-path> -out lib.tch
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tchaik/tchaik/index"
	"github.com/tchaik/tchaik/index/itl"
)

var itlXML, path string
var out string

func init() {
	flag.StringVar(&itlXML, "itlXML", "", "iTunes Music Library XML file")
	flag.StringVar(&path, "path", "", "directory path containing audio files")
	flag.StringVar(&out, "out", "data.tch", "output file (Tchaik library binary format)")
}

func main() {
	flag.Parse()

	if itlXML != "" && path != "" || itlXML == "" && path == "" {
		fmt.Println("must specify either 'itlXML' or 'path'")
		os.Exit(1)
	}

	if out == "" {
		fmt.Println("must specify 'out'")
		os.Exit(1)
	}

	var l index.Library
	var err error
	switch {
	case itlXML != "":
		l, err = importXML(itlXML)
	case path != "":
		l = importPath(path)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = writeLibrary(index.Convert(l, "TrackID"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func writeLibrary(l index.Library) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	return index.WriteTo(l, f)
}

func importXML(itlXML string) (index.Library, error) {
	f, err := os.Open(itlXML)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	l, err := itl.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func importPath(path string) index.Library {
	tracks := make(map[string]*Track)
	files := walk(path)
	for p := range files {
		if validExtension(p) {
			track, err := processPath(p)
			if err != nil {
				log.Printf("error processing '%v': %v\n", p, err)
				continue
			}
			tracks[p] = track
		}
	}

	return &Library{
		tracks: tracks,
	}
}
