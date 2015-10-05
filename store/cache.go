// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/context"
)

// RWFileSystem is an interface which includes http.FileSystem and a Create method
// for creating files.
type RWFileSystem interface {
	FileSystem

	// Create a file with associated path, returns an io.WriteCloser.  Only when Close()
	// returns can it be assumed that the file has been written.
	Create(ctx context.Context, path string) (io.WriteCloser, error)

	// Wait blocks until any pending write calls have been completed.
	Wait() error
}

// dir implements RWFileSystem by extending the behaviour of http.Dir to include a Create
// method which creates files under the root.
type dir struct {
	FileSystem

	root string
}

// Dir creates a new RWFileSystem with the specified root (similar to http.Dir)
func Dir(root string) RWFileSystem {
	return &dir{
		NewFileSystem(http.Dir(root), fmt.Sprintf("local (%v)", root)),
		root,
	}
}

func (d *dir) absPath(path string) (string, error) {
	cleanPath := filepath.Clean(d.root + "/" + path)
	path, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("error finding absolute path for: '%v' ('%v'): %v", path, cleanPath, err)
	}

	absRoot, err := filepath.Abs(d.root)
	if err != nil {
		return "", fmt.Errorf("error finding absolute path for root '%v': %v", d.root, err)
	}

	if !strings.HasPrefix(filepath.Dir(path), absRoot) {
		return "", fmt.Errorf("invalid path ('%v' is outside '%v'): %v", filepath.Dir(path), absRoot, path)
	}
	return path, nil
}

// Create a file rooted in the Dir file system.
func (d *dir) Create(ctx context.Context, path string) (io.WriteCloser, error) {
	absPath, err := d.absPath(path)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(filepath.Dir(absPath), os.ModePerm)
	if err != nil {
		return nil, err
	}
	return os.Create(absPath)
}

// Wait implements RWFileSystem.
func (d *dir) Wait() error { return nil }

// CachedError is an error returned by CachedErrorFileSystems when Open errors are cached
// rather than live.
type CachedError struct {
	Err error
}

// Error implements error.
func (c *CachedError) Error() string {
	return fmt.Sprintf("cached error: %v", c.Err)
}

// CachedErrorFileSystem provides an error cache to prevent erroring FileSystem requests
// from being repeated. See open for more details.
type CachedErrorFileSystem struct {
	FileSystem

	sync.RWMutex
	m map[string]error
}

func (c *CachedErrorFileSystem) setError(path string, err error) {
	c.Lock()
	defer c.Unlock()

	c.m[path] = err
}

func (c *CachedErrorFileSystem) getError(path string) (error, bool) {
	c.RLock()
	defer c.RUnlock()

	err, ok := c.m[path]
	return err, ok
}

// Open implements FileSystem, and caches errors from the underlying FileSystem.  The first time
// an error is encountered it is returned unchanged. Subsequent calls with an erroring path
// return a CachedError-wrapped version of the original error.
func (c CachedErrorFileSystem) Open(ctx context.Context, path string) (http.File, error) {
	err, ok := c.getError(path)
	if ok {
		return nil, &CachedError{
			Err: err,
		}
	}

	f, err := c.FileSystem.Open(ctx, path)
	if err != nil {
		c.setError(path, err)
		return nil, err
	}
	return f, nil
}

// CachedFileSystem is an implemetation of http.FileServer which caches the results of
// calls to src in a RWFileSystem.
type CachedFileSystem struct {
	src   FileSystem
	cache RWFileSystem

	errCh chan<- error
	wg    sync.WaitGroup
}

// Open implements FileSystem.  If the required file isn't in the cache
// then the file is opened from the src, and then concurrently copied into the
// cache (with errors passed back on the filesystem error channel).
func (c *CachedFileSystem) Open(ctx context.Context, path string) (http.File, error) {
	f, err := c.cache.Open(ctx, path)
	if err == nil {
		return f, nil
	}

	f, err = c.src.Open(ctx, path)
	if err != nil {
		return nil, err
	}

	go func() { // TODO: improve this so that we don't have to fetch the file again!
		c.wg.Add(1)
		defer c.wg.Done()

		src, err := c.src.Open(ctx, path)
		if err != nil {
			c.errCh <- fmt.Errorf("error opening file for second time: %v", err)
			return
		}
		defer func() {
			err := src.Close()
			if err != nil {
				c.errCh <- err
			}
		}()

		cache, err := c.cache.Create(ctx, path)
		if err != nil {
			c.errCh <- fmt.Errorf("error creating file in cache: %v", err)
			return
		}
		defer func() {
			err := cache.Close()
			if err != nil {
				c.errCh <- err
			}
		}()

		_, err = io.Copy(cache, src)
		if err != nil {
			c.errCh <- fmt.Errorf("error copying src file data into cache: %v", err)
		}
	}()

	return f, nil
}

// Wait implements RWFileSystem.
func (c *CachedFileSystem) Wait() error {
	c.wg.Wait()
	return nil
}

// NewCachedFileSystem implements http.FileSystem and caches every request made to
// src in cache.  The returned error channel passes back any errors which occur when
// files are being concurrently copied into the cache.
func NewCachedFileSystem(src FileSystem, cache RWFileSystem) (*CachedFileSystem, <-chan error) {
	errCh := make(chan error)
	return &CachedFileSystem{
		src:   src,
		cache: cache,
		errCh: errCh,
	}, errCh
}
