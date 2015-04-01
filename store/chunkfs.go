// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
)

// wrapper around a fileSource so that multiple callers can use a single
// file (and get their own Seeking etc).
type chunkedFile struct {
	io.ReadSeeker

	src *source
}

// Implements http.File.
func (cf *chunkedFile) Close() error {
	cf.src.wg.Done()
	return nil
}

// Implements http.File.
func (cf *chunkedFile) Readdir(int) ([]os.FileInfo, error) {
	return nil, nil
}

// Implements http.File.
func (cf *chunkedFile) Stat() (os.FileInfo, error) {
	return cf.src.stat, nil
}

// source contains the underlying data source for files currently being fetched.
type source struct {
	wg   sync.WaitGroup
	sra  SizeReaderAt
	stat os.FileInfo
}

// remoteChunkedFileSystem implements http.FileSystem.
type remoteChunkedFileSystem struct {
	client Client

	sync.RWMutex // protects files
	files        map[string]*source
	chunkSize    int64
}

// chunkedFile returns a *chunkedFile if a file with the given path is known, otherwise
// nil, false.
func (rcfs *remoteChunkedFileSystem) chunkedFile(path string) (*chunkedFile, bool) {
	rcfs.RLock()
	defer rcfs.RUnlock()

	src, ok := rcfs.files[path]
	if !ok {
		return nil, false
	}

	src.wg.Add(1)
	return &chunkedFile{
		io.NewSectionReader(src.sra, 0, src.sra.Size()), // create a ReadSeeker
		src,
	}, true
}

func (rcfs *remoteChunkedFileSystem) setSource(path string, sra SizeReaderAt, stat os.FileInfo) {
	rcfs.Lock()
	defer rcfs.Unlock()

	rcfs.files[path] = &source{
		sra:  sra,
		stat: stat,
	}
}

func (rcfs *remoteChunkedFileSystem) removeSource(path string) {
	rcfs.Lock()
	defer rcfs.Unlock()

	delete(rcfs.files, path)
}

// Open the file identified by path from the remote file system and read it into
// a chunked local copy, so that it can be read immediately (i.e. before the fetch
// completes any completed chunks can be read).  NB: when a chunk has not been fetched
// any operations will block until it is available.  Multiple calls to Open with the same
// path will receive independant http.File implementations using the same underlying
// data source (the file will only be fetched once).
func (rcfs *remoteChunkedFileSystem) Open(path string) (http.File, error) {
	cf, ok := rcfs.chunkedFile(path)
	if ok {
		return cf, nil
	}

	f, err := rcfs.client.Get(path)
	if err != nil {
		return nil, err
	}

	stat := &fileInfo{
		name:    f.Name,
		size:    f.Size,
		modTime: f.ModTime,
	}

	sra := NewChunkedReaderAt(f, f.Size, rcfs.chunkSize)
	rcfs.setSource(path, sra, stat)

	cf, ok = rcfs.chunkedFile(path)
	if !ok {
		return nil, errors.New("could not create chunked file from source")
	}
	go func() {
		cf.src.wg.Wait()
		rcfs.removeSource(path)
	}()
	return cf, nil
}

// NewRemoteChunkedFileSystem creates an implementation of http.FileSystem which fetches
// files from the given Client, and allows access to chunks of the file contents as they
// are retrieved.  See Open for more details.
func NewRemoteChunkedFileSystem(client Client, chunkSize int64) *remoteChunkedFileSystem {
	return &remoteChunkedFileSystem{
		client:    client,
		files:     make(map[string]*source),
		chunkSize: chunkSize,
	}
}
