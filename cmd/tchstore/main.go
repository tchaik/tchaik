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
	"strings"

	"github.com/mitchellh/goamz/aws"

	"github.com/dhowden/tchaik/store"
	"github.com/dhowden/tchaik/store/cafs"
)

var listen string
var root string

var mediaFileSystemCache, artworkFileSystemCache string

var debug bool

func init() {
	flag.StringVar(&listen, "listen", "localhost:1844", "<addr>:<port> to listen on")
	flag.StringVar(&root, "root", "", "the media root, full local path /path/to/root, or s3://<bucket>/path/to/root for S3")

	flag.StringVar(&artworkFileSystemCache, "artwork-cache", "", "path to artwork cache (content addressable)")
	flag.StringVar(&mediaFileSystemCache, "media-cache", "", "path to media cache")

	flag.BoolVar(&debug, "debug", false, "output extra debugging information")
}

func main() {
	flag.Parse()
	if listen == "" || root == "" {
		flag.Usage()
		return
	}

	// src is the main source of all data in the store
	var src http.FileSystem
	if strings.HasPrefix(root, "s3://") {
		path := strings.TrimPrefix(root, "s3://")
		bucketPathSplit := strings.Split(path, "/")

		if len(bucketPathSplit) == 0 {
			fmt.Println("invalid S3 path: %#v\n", root)
			os.Exit(1)
		}

		bucket := bucketPathSplit[0]
		auth, err := aws.GetAuth("", "") // Extract credentials from the current instance.
		if err != nil {
			log.Printf("error getting credentials: %v", err)
			os.Exit(1)
		}

		c := store.NewS3Client(bucket, auth, aws.APSoutheast2)
		src = store.NewRemoteChunkedFileSystem(c, 32*1024)
	} else {
		src = http.Dir(root)
	}

	mediaFileSystem := src
	if mediaFileSystemCache != "" {
		var errCh <-chan error
		localCache := store.NewDir(mediaFileSystemCache)
		mediaFileSystem, errCh = store.NewCachedFileSystem(
			mediaFileSystem,
			localCache,
		)
		go func() {
			for err := range errCh {
				log.Printf("mediaFileSystem cache: %v", err)
			}
		}()
	}

	var artworkFileSystem http.FileSystem
	artworkFileSystem = store.LogFileSystem{"Extracting artwork", store.ArtworkFileSystem{mediaFileSystem}}
	if artworkFileSystemCache != "" {
		cfs, err := cafs.New(store.NewDir(artworkFileSystemCache))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var errCh <-chan error
		artworkFileSystem, errCh = store.NewCachedFileSystem(
			artworkFileSystem,
			cfs,
		)
		go func() {
			for err := range errCh {
				log.Printf("artwork cache: %v", err)
			}
		}()
	}

	if debug {
		mediaFileSystem = store.LogFileSystem{"Media", mediaFileSystem}
		artworkFileSystem = store.LogFileSystem{"Artwork", artworkFileSystem}
	}

	s := store.NewServer(listen)
	s.SetDefault(mediaFileSystem)
	s.SetFileSystem("artwork", artworkFileSystem)
	log.Fatal(s.Listen())
}
