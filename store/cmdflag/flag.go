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

	"tchaik.com/store"
	"tchaik.com/store/cafs"
)

var localStore, remoteStore string
var mediaFileSystemCache, artworkFileSystemCache string
var trimPathPrefix, addPathPrefix string

func init() {
	flag.StringVar(&localStore, "local-store", "/", "local media store, full local path /path/to/root")
	flag.StringVar(&remoteStore, "remote-store", "", "remote media store, tchstore server address <hostname>:<port>, s3://<bucket>/path/to/root for S3, or gs://<project-id>:<bucket>/path/to/root for Google Cloud Storage")

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
	switch {
	case strings.HasPrefix(remoteStore, "s3://"):
		path := strings.TrimPrefix(remoteStore, "s3://")
		bucketPathSplit := strings.Split(path, "/")

		if len(bucketPathSplit) == 0 {
			return fmt.Errorf("invalid S3 path: %#v", remoteStore)
		}
		bucket := bucketPathSplit[0]
		var auth aws.Auth
		auth, err = aws.GetAuth("", "") // Extract credentials from the current instance.
		if err != nil {
			return fmt.Errorf("error getting AWS credentials: %v", err)
		}
		c = store.TraceClient(store.NewS3Client(bucket, auth, aws.APSoutheast2), "S3")

	case strings.HasPrefix(remoteStore, "gs://"):
		path := strings.TrimPrefix(remoteStore, "gs://")
		bucketPathSplit := strings.Split(path, "/")

		if len(bucketPathSplit) == 0 {
			return fmt.Errorf("invalid Google Cloud Storage path: %#v", remoteStore)
		}
		projIDBucket := bucketPathSplit[0]
		if len(projIDBucket) == 0 {
			return fmt.Errorf("invalid Google Cloud Storage path: %#v", remoteStore)
		}

		projIDBucketSplit := strings.Split(projIDBucket, ":")
		if len(projIDBucketSplit) != 2 {
			return fmt.Errorf("invalid Google Cloud Storage path prefix (<project-id>:<bucket>): %#v", projIDBucket)
		}
		if len(projIDBucketSplit[0]) == 0 {
			return fmt.Errorf("invalid Google Cloud Storage path prefix (<project-id>:<bucket>): empty project-id: %#v", projIDBucket)
		}
		c = store.TraceClient(store.NewCloudStorageClient(projIDBucketSplit[0], projIDBucketSplit[1]), "CloudStorage")

	default:
		c = store.TraceClient(store.NewClient(remoteStore, ""), "tchstore")
		s.artwork = store.NewRemoteFileSystem(store.NewClient(remoteStore, "artwork"))
	}

	s.media = store.NewRemoteChunkedFileSystem(c, 32*1024)
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
		s.media = store.PathRewrite(s.media, trimPathPrefix, addPathPrefix)
		s.artwork = store.PathRewrite(s.artwork, trimPathPrefix, addPathPrefix)
	}
	return s.media, s.artwork, nil
}
