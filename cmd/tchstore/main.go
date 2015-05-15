// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
tchstore is a tool to create a remote store for tchaik files, including media and artwork files (TODO: index).

It is assumed that tchstore is run relatively local to the data it is serving (i.e. in EC2 for S3, or on a
fileserver).

All configuration is done through command line parameters, use --help flag for full details.

To fetch files from the `tchaik-store` S3 bucket (and store media fetches in local cache under `media-cache`
directory, and artwork in a content addressable cache under the directory `artwork-cache`):

  tchstore -listen 0.0.0.0:1844 -root s3://tchaik-store/ -artwork-cache artwork-cache -media-cache media-cache

Implementation note: by default the artwork store will try to:
  1) retrieve files from the local artwork cache (if the cache exists)
  2) retrieve data from media files in a local store (media-cache or root if local)
  3) fetch the whole media file from S3 (if configured), which will in turn add the
     media file to (media-cache).

Set the environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to pass credentials to the S3 client.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tchaik/tchaik/store"
	"github.com/tchaik/tchaik/store/cmdflag"
)

var listen string
var debug bool

func init() {
	flag.StringVar(&listen, "listen", "localhost:1844", "<addr>:<port> to listen on")
	flag.BoolVar(&debug, "debug", false, "output extra debugging information")
}

func main() {
	flag.Parse()
	if listen == "" {
		flag.Usage()
		return
	}

	mediaFileSystem, artworkFileSystem, err := cmdflag.Stores()
	if err != nil {
		fmt.Println("error setting up stores:", err)
		os.Exit(1)
	}

	if debug {
		mediaFileSystem = store.LogFileSystem("Media", mediaFileSystem)
		artworkFileSystem = store.LogFileSystem("Artwork", artworkFileSystem)
	}

	s := store.NewServer(listen)
	s.SetDefault(mediaFileSystem)
	s.SetFileSystem("artwork", artworkFileSystem)
	log.Fatal(s.Listen())
}
