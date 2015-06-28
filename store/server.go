// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
)

// Request is a type which represents an incoming request.
type Request struct {
	Path, Label string
}

// Response is a type which represents a response to a Request.
type Response struct {
	Status  ResponseStatus
	Size    int64 // The size of the returned output
	ModTime time.Time
	Name    string
}

// ResponseStatus is an enumeration of possible response statuses.
type ResponseStatus string

// All defined ResponseStatus values.
const (
	StatusOK            ResponseStatus = "OK" // The response succeeded, and the file follows.
	StatusLabelNotFound                = "LF" // The label is invalid (does not exist).
	StatusPathError                    = "PE" // The path is errornous (could not be parsed).
	StatusInvalidPath                  = "IP" // The path is outside the root of the server.
	StatusNotFound                     = "NF" // The path is invalid (no file found).
	StatusFileError                    = "FE" // The path refers to a valid file, but there was a problem reading it.
	StatusDirectory                    = "ED" // The path refers to a directory, which cannot be transmitted.
)

// Implements Stringer.
func (r ResponseStatus) String() string {
	switch r {
	case StatusOK:
		return "OK"
	case StatusPathError:
		return "Path Error"
	case StatusInvalidPath:
		return "Invalid Path"
	case StatusNotFound:
		return "Not Found"
	case StatusFileError:
		return "File Error"
	case StatusDirectory:
		return "Directory"
	}
	return fmt.Sprintf("<INVALID ResponseStatus: %v>", string(r))
}

// Server represents a store server, which implements a simple protocol for transfering files
// to the local system whilst piping them to the requesting client.
type Server struct {
	addr string

	fileSystems map[string]FileSystem
}

// NewServer creates a new server listening on the given address.
func NewServer(addr string) *Server {
	return &Server{
		addr:        addr,
		fileSystems: make(map[string]FileSystem),
	}
}

// SetDefault sets the default (empty-name) file system
func (s *Server) SetDefault(fs FileSystem) {
	s.fileSystems[""] = fs
}

// SetFileSystem sets the underlying file system to use for a label.
func (s *Server) SetFileSystem(label string, fs FileSystem) {
	s.fileSystems[label] = fs
}

// Listen starts listening on s.Addr.  If there is an issue binding the
// listener, then an error is returned.  Any errors which occur due to
// individual connections are logged.
func (s *Server) Listen() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func(c net.Conn) {
			err := s.handle(c)
			if err != nil {
				log.Println(err)
			}
		}(conn)
	}
}

// handle the given connection
func (s *Server) handle(c net.Conn) (err error) {
	defer func() {
		if err1 := c.Close(); err == nil && err1 != nil {
			err = fmt.Errorf("error closing connection: %v", err1)
		}
	}()

	dec := json.NewDecoder(c)
	var r Request
	err = dec.Decode(&r)
	if err != nil {
		err = fmt.Errorf("error decoding request: %v", err)
		return
	}

	fs, ok := s.fileSystems[r.Label]
	if !ok {
		writeStatusResponse(c, StatusNotFound)
		err = fmt.Errorf("invalid label: %v", r.Label)
		return
	}

	// FIXME: Transfer the context from the request?
	f, err := fs.Open(context.Background(), r.Path)
	if err != nil {
		writeStatusResponse(c, StatusNotFound)
		err = fmt.Errorf("error opening file '%v': %v", r.Path, err)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		writeStatusResponse(c, StatusFileError)
		err = fmt.Errorf("error stating file: '%v': %v", r.Path, err)
		return
	}

	if stat.IsDir() {
		writeStatusResponse(c, StatusDirectory)
		err = fmt.Errorf("can't retrieve dir: '%v'", r.Path)
		return
	}

	resp := Response{
		Status:  StatusOK,
		ModTime: stat.ModTime(),
		Size:    stat.Size(),
		Name:    stat.Name(),
	}
	writeResponse(c, resp)
	n, err := io.Copy(c, f)
	if err != nil {
		err = fmt.Errorf("error copying data from file '%v': %v", r.Path, err)
		return
	}
	log.Printf("%#v: %v (%v, %d bytes)", r.Label, r.Path, stat.Name(), n)
	return
}

func writeStatusResponse(w io.Writer, status ResponseStatus) {
	writeResponse(w, Response{
		Status: status,
	})
}

func writeResponse(w io.Writer, resp Response) {
	enc := json.NewEncoder(w)
	err := enc.Encode(resp)
	if err != nil {
		log.Printf("error writing response '%#v': %v", resp, err)
	}
}
