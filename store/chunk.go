// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"bytes"
	"io"
	"sync"
)

type chunk struct {
	io.ReaderAt
	sync.RWMutex

	err  error
	size int64
}

func newChunk(size int64) *chunk {
	c := &chunk{size: size}
	c.Lock() // block all calls to ReadAt until we have the underlying data
	return c
}

// ReadFrom will read c.size bytes from the reader into the buffer and then unlock
// the mutex to allow for reads via ReadAt.
func (c *chunk) ReadFrom(r io.Reader) (int64, error) {
	defer c.Unlock()

	buf := make([]byte, c.size)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		c.err = err
		return int64(n), err
	}
	c.ReaderAt = bytes.NewReader(buf)
	return int64(n), err
}

// setError sets the error which will be returned by all subsequent calls to
// ReadAt.  Typically this is used to record an error in chunks which follow
// a chunk whose ReadFrom caused an error.
func (c *chunk) setError(err error) {
	c.err = err
	c.Unlock()
}

// ReadAt implements ReaderAt, and will block until the underlying buffer
// has been filled (see ReadFrom).
func (c *chunk) ReadAt(b []byte, offset int64) (int, error) {
	c.RLock()
	defer c.RUnlock()

	if c.err != nil {
		return 0, c.err
	}
	return c.ReaderAt.ReadAt(b, offset)
}

// Size returns the size of the chunk.
func (c *chunk) Size() int64 {
	return c.size
}

// NewChunkedReaderAt reads 'size' bytes of data from the reader, buffering
// into chunks of size 'chunkSize'. The returned SizeReaderAt allows access
// to chunks which have been downloaded in full, and will block on any
// partially downloaded chunks until they have completed.
func NewChunkedReaderAt(r io.ReadCloser, size, chunkSize int64) SizeReaderAt {
	n := size / chunkSize    // number of whole sized chunks
	last := size % chunkSize // size of trailing (non-whole) chunk

	chunks := make([]*chunk, 0, n)
	for i := 0; i < int(n); i++ {
		chunks = append(chunks, newChunk(chunkSize))
	}
	if last != 0 {
		chunks = append(chunks, newChunk(last))
	}

	go func() {
		defer r.Close()

		i := 0
		var err error
		for _, c := range chunks {
			_, err = c.ReadFrom(r)
			if err != nil {
				break
			}
			i++
		}
		if err != nil && i < len(chunks)-1 {
			for _, c := range chunks[i+1:] {
				c.setError(err)
			}
		}
	}()

	l := make([]SizeReaderAt, len(chunks))
	for i := 0; i < len(chunks); i++ {
		l[i] = chunks[i]
	}
	return NewMultiReaderAt(l...)
}
