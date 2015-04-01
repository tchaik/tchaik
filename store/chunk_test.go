// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNewChunkedReaderAt(t *testing.T) {
	input := `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

	chunkSizes := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 13, len(input),
	}

	for _, chunkSize := range chunkSizes {
		r := NewChunkedReaderAt(ioutil.NopCloser(strings.NewReader(input)), int64(len(input)), int64(chunkSize))

		output, err := ioutil.ReadAll(io.NewSectionReader(r, 0, r.Size()))
		if err != nil {
			t.Errorf("unexpected error on chunkSize: %d: %v", chunkSize, err)
		}
		if string(output) != input {
			t.Errorf("output = %s\nexpected %s", output, input)
		}
	}
}

// return err after reading n bytes
type errorReader struct {
	io.Reader
	err error
	n   int

	i int
}

func (er *errorReader) Read(b []byte) (int, error) {
	if er.i <= er.n && er.n < er.i+len(b) {
		n, _ := er.Reader.Read(b[:er.n-er.i])
		return n, er.err
	}
	er.i += len(b)
	return er.Reader.Read(b)
}

func TestErrorChunkedReaderAt(t *testing.T) {
	input := `Lorem ipsum dolor sit amet, consectetur adipisicing elit.`
	expectedErr := errors.New("chunking error")

	tests := []struct {
		n         int
		chunkSize int
		output    string
	}{

		// All chunks die
		{
			0,
			1,
			"",
		},

		// First chunk survives
		{
			1,
			1,
			"L",
		},

		// First two chunks survive
		{
			2,
			1,
			"Lo",
		},

		// All chunks die
		{
			0,
			2,
			"",
		},

		// All chunks die
		{
			1,
			2,
			"",
		},

		// First chunk survives, rest die.
		{
			2,
			2,
			"Lo",
		},
	}

	for ii, tt := range tests {
		r := NewChunkedReaderAt(ioutil.NopCloser(&errorReader{
			Reader: strings.NewReader(input),
			n:      tt.n,
			err:    expectedErr,
		}), int64(len(input)), int64(tt.chunkSize))
		output, err := ioutil.ReadAll(io.NewSectionReader(r, 0, r.Size()))
		if err != expectedErr {
			t.Errorf("[%d] err = %v, expected: %v", ii, err, expectedErr)
		}
		if string(output) != tt.output {
			t.Errorf("[%d] output = %#v, expected %#v", ii, string(output), tt.output)
		}
	}
}
