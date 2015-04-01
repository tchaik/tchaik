package cafs

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

// InvalidPathError is returned by cachine filesystem when a previous attempt has been
// made to fetch a file and failed.  Prevents further access to the source.
type InvalidPathError string

// Error implements error.
func (e *InvalidPathError) Error() string {
	return string(*e)
}

// CacheFileSystem is an implemetation of http.FileServer which assumes that both the
// cache and src is content addressable (i.e. file.Stat().Name() is a content hash).
// More efficient than the standard caching file system since any requested paths which
// aren't in the index are passed to the src and then their associated content is only
// downloaded if not already present.
type CachedFileSystem struct {
	src   http.FileSystem
	cache *FileSystem

	errCh chan<- error
	wg    sync.WaitGroup
}

// Open implements http.FileSystem.  If the required file isn't in the cache
// then the file is opened from the src, and then concurrently copied into the
// cache (with errors passed back on the filesystem error channel).
func (c *CachedFileSystem) Open(path string) (http.File, error) {
	f, err := c.cache.Open(path)
	if err == nil {
		return f, nil
	}

	if _, ok := err.(*InvalidPathError); ok {
		return nil, nil
	}

	f, err = c.src.Open(path)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// if the content is already on the local machine, then don't fetch it, just
	// update the index.
	if c.cache.idx.Exists(stat.Name()) {
		c.cache.idx.Add(path, stat.Name())
		f.Close()

		f, err = c.cache.Open(path)
		if err != nil {
			return nil, fmt.Errorf("error opening file after adding to cache: %v", err)
		}
		return f, nil
	}

	go func() {
		c.wg.Add(1)
		defer c.wg.Done()

		src, err := c.src.Open(path)
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

		cache, err := c.cache.Create(path)
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

// Wait blocks until any concurrent caching operations have been completed.
func (c *CachedFileSystem) Wait() error {
	c.wg.Wait()
	return nil
}

// NewCacheFileSystem implements http.FileSystem and caches every request made to
// src in cache.  The returned error channel passes back any errors which occur when
// files are being concurrently copied into the cache.  Both src and cache must be
// content addressable using the same hashing scheme.
func NewCachedFileSystem(src http.FileSystem, cache *FileSystem) (http.FileSystem, <-chan error) {
	errCh := make(chan error)
	return &CachedFileSystem{
		src:   src,
		cache: cache,
		errCh: errCh,
	}, errCh
}
