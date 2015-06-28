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
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"

	"github.com/tchaik/tchaik/store"
	"github.com/tchaik/tchaik/store/cmdflag"
)

var listen string
var debug bool

var traceListenAddr string

func init() {
	flag.StringVar(&listen, "listen", "localhost:1844", "<addr>:<port> to listen on")
	flag.BoolVar(&debug, "debug", false, "output extra debugging information")

	flag.StringVar(&traceListenAddr, "trace-listen", "", "bind address for trace HTTP server")
}

type rootTraceFS struct {
	store.FileSystem
	family string
}

// Open implements store.FileSystem
func (r *rootTraceFS) Open(ctx context.Context, path string) (http.File, error) {
	tr := trace.New("request", path)
	defer tr.Finish()

	return r.FileSystem.Open(trace.NewContext(ctx, tr), path)
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

	mediaFileSystem = &rootTraceFS{mediaFileSystem, "media"}
	artworkFileSystem = &rootTraceFS{artworkFileSystem, "artwork"}

	if debug {
		mediaFileSystem = store.LogFileSystem("Media", mediaFileSystem)
		artworkFileSystem = store.LogFileSystem("Artwork", artworkFileSystem)
	}

	if traceListenAddr != "" {
		fmt.Printf("Starting trace server on http://%v\n", traceListenAddr)
		go func() {
			log.Fatal(http.ListenAndServe(traceListenAddr, nil))
		}()
	}

	s := store.NewServer(listen)
	s.SetDefault(mediaFileSystem)
	s.SetFileSystem("artwork", artworkFileSystem)
	log.Fatal(s.Listen())
}
