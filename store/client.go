// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

// Client is an interface which defines the Get method used to fetch files
// from remote hosts.
type Client interface {
	// Get reaches out to a remote server with a request for the given path.
	Get(ctx context.Context, path string) (*File, error)

	// Put reaches out to a remote server with a put (write) request.
	// Put(path string)
}

// NewClient initialises the default Client implementation with the given remote
// addr and filesystem label.
func NewClient(addr, label string) *client {
	return &client{
		addr:  addr,
		label: label,
	}
}

type client struct {
	addr  string
	label string
}

// File contains meta data for a remote file, and implements io.ReadCloser.
type File struct {
	io.ReadCloser
	Name    string
	ModTime time.Time
	Size    int64
}

type readCloser struct {
	io.Reader
	io.Closer
}

// Get implements Client.
func (c *client) Get(ctx context.Context, path string) (f *File, err error) {
	if tr, ok := trace.FromContext(ctx); ok {
		tr.LazyPrintf("(%v, %#v) get '%v'", c.addr, c.label, path)
		defer func() {
			if err != nil {
				tr.LazyPrintf("(%v, %#v) error: %v", c.addr, c.label, err)
				return
			}
			tr.LazyPrintf("(%v, %#v) returned: %v", c.addr, c.label, f.Name)
		}()
	}

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return nil, err
	}

	enc := json.NewEncoder(conn)
	err = enc.Encode(Request{
		Path:  path,
		Label: c.label,
	})
	if err != nil {
		return nil, err
	}

	// Decode the Response
	dec := json.NewDecoder(conn)
	var resp Response
	err = dec.Decode(&resp)
	if err != nil {
		return nil, err
	}

	if resp.Status != StatusOK {
		return nil, fmt.Errorf("error from '%v' (%v): %v", c.addr, c.label, resp.Status)
	}

	r := readCloser{io.MultiReader(dec.Buffered(), conn), conn}

	// Remove the extra \n character added when the JSON is encoded
	b := make([]byte, 1)
	_, err = r.Read(b)
	if err != nil {
		return nil, fmt.Errorf("error reading '\n' after Response: %v", err)
	}
	if b[0] != '\n' {
		return nil, fmt.Errorf("expected to read '\n' after Response")
	}

	return &File{
		ReadCloser: r,
		Name:       resp.Name,
		ModTime:    resp.ModTime,
		Size:       resp.Size,
	}, nil
}

// TraceClient creates a convenience method adding a tracing wrapper around a Client.
func TraceClient(c Client, name string) Client {
	return &traceClient{
		Client: c,
		name:   name,
	}
}

// traceClient is a type which wraps a Client implementation.  If a trace is associated
// to the context.Context passed to Get, then the trace will be updated to reflect the call.
type traceClient struct {
	Client
	name string
}

// Open implements Client.
func (tc traceClient) Get(ctx context.Context, path string) (f *File, err error) {
	if tr, ok := trace.FromContext(ctx); ok {
		tr.LazyPrintf("%v: Get: %v", tc.name, path)
		defer func() {
			if err != nil {
				tr.LazyPrintf("%v: error opening '%v': %v", tc.name, path, err)
				return
			}
			tr.LazyPrintf("%v: got file: %v", tc.name, f.Name)
		}()
	}
	return tc.Client.Get(ctx, path)
}
