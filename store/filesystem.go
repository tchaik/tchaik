// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
)

// FileSystem is an interface which defines an open method similar to http.FileSystem,
// but which also includes a context parameter.
type FileSystem interface {
	Open(ctx context.Context, path string) (http.File, error)
}

// NewFileSystem creates a new FileSystem using an http.FileSystem as the underlying
// storage.
func NewFileSystem(fs http.FileSystem) FileSystem {
	return &fileSystem{fs}
}

type fileSystem struct {
	http.FileSystem
}

// Open implements FileSystem.
func (cfs *fileSystem) Open(ctx context.Context, path string) (http.File, error) {
	return cfs.FileSystem.Open(path)
}

// RemoteFileSystem is an extension of the http.FileSystem interface
// which includes the RemoteOpen method.
type RemoteFileSystem interface {
	http.FileSystem

	// RemoteOpen returns a File which
	RemoteOpen(string) (*File, error)
}

// remoteFileSystem implements RemoteFileSystem
type remoteFileSystem struct {
	client Client
}

// NewRemoteFileSystem creates a new file system using the given Client to handle
// file requests.
func NewRemoteFileSystem(c Client) RemoteFileSystem {
	return &remoteFileSystem{
		client: c,
	}
}

// file is a basic representation of a remote file such that all operations (i.e.
// seeking) will work correctly.
type file struct {
	io.ReadSeeker
	stat *fileInfo
}

// RemoteOpen returns a *File which represents the remote file, and implements
// io.ReadCloser which reads the file contents from the remote system.
func (fs *remoteFileSystem) RemoteOpen(path string) (*File, error) {
	rf, err := fs.client.Get(path)
	if err != nil {
		return nil, err
	}
	return rf, nil
}

// Open the given file and return an http.File implementation representing it.  This method
// will block until the file has been completely fetched (http.File implements io.Seeker
// which means that for a trivial implementation we need all the underlying data).
func (fs *remoteFileSystem) Open(path string) (http.File, error) {
	rf, err := fs.RemoteOpen(path)
	if err != nil {
		return nil, err
	}
	defer rf.Close()

	buf, err := ioutil.ReadAll(rf)
	if err != nil {
		return nil, err
	}

	return &file{
		ReadSeeker: bytes.NewReader(buf),
		stat: &fileInfo{
			name:    rf.Name,
			size:    rf.Size,
			modTime: rf.ModTime,
		},
	}, nil
}

// Close is a nop as we have already closed the original file.
func (f *file) Close() error {
	return nil
}

// Implements http.File.
func (f *file) Readdir(int) ([]os.FileInfo, error) {
	return nil, nil
}

// FileInfo is a simple implementation of os.FileInfo.
type fileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (f *fileInfo) Name() string       { return f.name }
func (f *fileInfo) Size() int64        { return f.size }
func (f *fileInfo) Mode() os.FileMode  { return os.FileMode(0777) }
func (f *fileInfo) ModTime() time.Time { return f.modTime }
func (f *fileInfo) IsDir() bool        { return false }
func (f *fileInfo) Sys() interface{}   { return nil }

func (f *file) Stat() (os.FileInfo, error) {
	return f.stat, nil
}
