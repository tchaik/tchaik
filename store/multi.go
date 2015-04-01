// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import "net/http"

// multiFileSystem implements FileSystem wrapping a list of FileSystems.
type multiFileSystem struct {
	fileSystems []http.FileSystem
}

// Open implements http.FileSystem
func (mfs *multiFileSystem) Open(name string) (http.File, error) {
	var err error
	var f http.File
	for _, fs := range mfs.fileSystems {
		f, err = fs.Open(name)
		if err == nil {
			return f, err
		}
	}
	return nil, err
}

// MultiFileSystem implements FileSystem and wraps an ordered list of FileSystem
// implementations. With each call to Open, the file systems are tried in turn until
// one returns without error. If all return errors, then we pass the result back to
// the caller.
func MultiFileSystem(fs ...http.FileSystem) http.FileSystem {
	return &multiFileSystem{
		fileSystems: fs,
	}
}
