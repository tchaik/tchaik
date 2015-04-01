// Copyright (c) 2012 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//   * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//   * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//   * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package store

import (
	"io"
	"sort"
)

// This code is taken from https://talks.golang.org/2013/oscon-dl/sizereaderat.go

// A SizeReaderAt is a ReaderAt with a Size method.
//
// An io.SectionReader implements SizeReaderAt.
type SizeReaderAt interface {
	Size() int64
	io.ReaderAt
}

// NewMultiReaderAt is like io.MultiReader but produces a ReaderAt
// (and Size), instead of just a reader.
func NewMultiReaderAt(parts ...SizeReaderAt) SizeReaderAt {
	m := &multi{
		parts: make([]offsetAndSource, 0, len(parts)),
	}
	var off int64
	for _, p := range parts {
		m.parts = append(m.parts, offsetAndSource{off, p})
		off += p.Size()
	}
	m.size = off
	return m
}

type offsetAndSource struct {
	off int64
	SizeReaderAt
}

type multi struct {
	parts []offsetAndSource
	size  int64
}

func (m *multi) Size() int64 { return m.size }

func (m *multi) ReadAt(p []byte, off int64) (n int, err error) {
	wantN := len(p)

	// Skip past the requested offset.
	skipParts := sort.Search(len(m.parts), func(i int) bool {
		// This function returns whether parts[i] will
		// contribute any bytes to our output.
		part := m.parts[i]
		return part.off+part.Size() > off
	})
	parts := m.parts[skipParts:]

	// How far to skip in the first part.
	needSkip := off
	if len(parts) > 0 {
		needSkip -= parts[0].off
	}

	for len(parts) > 0 && len(p) > 0 {
		readP := p
		partSize := parts[0].Size()
		if int64(len(readP)) > partSize-needSkip {
			readP = readP[:partSize-needSkip]
		}
		pn, err0 := parts[0].ReadAt(readP, needSkip)
		if err0 != nil {
			return n, err0
		}
		n += pn
		p = p[pn:]
		if int64(pn)+needSkip == partSize {
			parts = parts[1:]
		}
		needSkip = 0
	}

	if n != wantN {
		err = io.ErrUnexpectedEOF
	}
	return
}
