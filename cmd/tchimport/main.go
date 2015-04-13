// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
tchimport is a tool which converts iTunes Libraries (plist XML) into a Tchaik libraries.

This is particularly useful if you have a large iTunes Library file and don't want to have
it all loaded into memory.  The Tchaik library has a much smaller set of data attributes
stored for each track, and the serialised form is gzipped JSON rather than plist.

  tchimport -xml <itunes-library> -out <output-file>
*/
package main

import (
	"flag"
	"log"
	"os"

	"github.com/dhowden/tchaik/index"
	"github.com/dhowden/itl"
)

var itlXML string
var out string

func init() {
	flag.StringVar(&itlXML, "xml", "", "iTunes Music Library XML file")
	flag.StringVar(&out, "out", "data.tch", "output file (Tchaik library binary format)")
}

func main() {
	flag.Parse()

	if itlXML == "" || out == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(itlXML)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	l, err := itl.ReadFromXML(f)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	itl := index.NewITunesLibrary(&l)
	nl := index.Convert(itl, "TrackID")

	nf, err := os.Create(out)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer nf.Close()

	err = index.WriteTo(nl, nf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
