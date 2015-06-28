// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cmdflag unifies the configuration of stores using command line flags across
// several tools.
//
// FIXME: Need a more coherent way of doing this: it's now a huge mess.
package cmdflag

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mitchellh/goamz/aws"

	"github.com/tchaik/tchaik/store"
	"github.com/tchaik/tchaik/store/cafs"
)

var localStore, remoteStore string
var mediaFileSystemCache, artworkFileSystemCache string
var trimPathPrefix, addPathPrefix string

func init() {
	flag.StringVar(&localStore, "local-store", "/", "local media store, full local path /path/to/root")
	flag.StringVar(&remoteStore, "remote-store", "", "remote media store, tchstore server address <hostname>:<port>, or s3://<bucket>/path/to/root for S3")

	flag.StringVar(&artworkFileSystemCache, "artwork-cache", "", "path to local artwork cache (content addressable)")
	flag.StringVar(&mediaFileSystemCache, "media-cache", "", "path to local media cache")

	flag.StringVar(&trimPathPrefix, "trim-path-prefix", "", "remove prefix from every path")
	flag.StringVar(&addPathPrefix, "add-path-prefix", "", "add prefix to every path")
}

type stores struct {
	media, artwork store.FileSystem
}

func buildRemoteStore(s *stores) (err error) {
	if remoteStore == "" {
		return nil
	}
	var c store.Client
	if strings.HasPrefix(remoteStore, "s3://") {
		path := strings.TrimPrefix(remoteStore, "s3://")
		bucketPathSplit := strings.Split(path, "/")

		if len(bucketPathSplit) == 0 {
			return fmt.Errorf("invalid S3 path: %#v\n", remoteStore)
		}
		bucket := bucketPathSplit[0]
		var auth aws.Auth
		auth, err = aws.GetAuth("", "") // Extract credentials from the current instance.
		if err != nil {
			return fmt.Errorf("error getting AWS credentials: %v", err)
		}
		c = store.NewS3Client(bucket, auth, aws.APSoutheast2)
	} else {
		c = store.NewClient(remoteStore, "")
		s.artwork = store.NewFileSystem(store.NewRemoteFileSystem(store.NewClient(remoteStore, "artwork")), "artwork")
	}

	s.media = store.NewFileSystem(store.NewRemoteChunkedFileSystem(c, 32*1024), fmt.Sprintf("remoteStore(%v)", remoteStore))
	if s.artwork == nil {
		s.artwork = store.Trace(store.ArtworkFileSystem(s.media), "artwork")
	}
	return nil
}

func buildLocalStore(s *stores) {
	if localStore != "" {
		fs := store.NewFileSystem(http.Dir(localStore), fmt.Sprintf("localstore (%v)", localStore))
		if s.media != nil {
			s.media = store.MultiFileSystem(fs, s.media)
		} else {
			s.media = fs
		}

		afs := store.Trace(store.ArtworkFileSystem(fs), "local artworkstore")
		if s.artwork != nil {
			s.artwork = store.MultiFileSystem(afs, s.artwork)
		} else {
			s.artwork = afs
		}
	}
}

func buildMediaCache(s *stores) {
	if mediaFileSystemCache != "" {
		var errCh <-chan error
		localCache := store.Dir(mediaFileSystemCache)
		s.media, errCh = store.NewCachedFileSystem(s.media, localCache)
		go func() {
			for err := range errCh {
				// TODO: pull this out!
				log.Printf("mediaFileSystem cache: %v", err)
			}
		}()
	}
}

func buildArtworkCache(s *stores) error {
	if artworkFileSystemCache != "" {
		cfs, err := cafs.New(store.Dir(artworkFileSystemCache))
		if err != nil {
			return fmt.Errorf("error creating artwork cafs: %v", err)
		}

		var errCh <-chan error
		s.artwork, errCh = store.NewCachedFileSystem(
			s.artwork,
			cfs,
		)
		go func() {
			for err := range errCh {
				// TODO: pull this out!
				log.Printf("artwork cache: %v", err)
			}
		}()
	}
	return nil
}

// Stores returns a media and artwork filesystem as defined by the command line flags.
func Stores() (media, artwork store.FileSystem, err error) {
	s := &stores{}
	err = buildRemoteStore(s)
	if err != nil {
		return nil, nil, err
	}

	buildLocalStore(s)
	buildMediaCache(s)

	err = buildArtworkCache(s)
	if err != nil {
		return nil, nil, err
	}

	if trimPathPrefix != "" || addPathPrefix != "" {
		s.media = store.NewPathRewrite(s.media, trimPathPrefix, addPathPrefix)
		s.artwork = store.NewPathRewrite(s.artwork, trimPathPrefix, addPathPrefix)
	}
	return s.media, s.artwork, nil
}
