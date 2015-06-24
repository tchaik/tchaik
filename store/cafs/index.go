package cafs

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/tchaik/tchaik/store"
)

// Index is an interface which contains methods for implementing
// a content addressable file system.
type Index interface {
	// Get returns the real path for the given filename, with true if and only
	// if the path exists in the index.
	Get(path string) (string, bool)

	// Add adds the path to the index, and returns the path to the file
	// and whether the path/content already exists.
	Add(path string, sum string) bool

	// Exists returns true of the sum is in the index.
	Exists(sum string) bool
}

type index struct {
	sync.RWMutex

	files map[string]string // path -> sha1
	index map[string]bool   // {sha1}
}

// NewIndex creates a new file system index.
func NewIndex() *index {
	return &index{
		files: make(map[string]string),
		index: make(map[string]bool),
	}
}

func (i *index) MarshalJSON() ([]byte, error) {
	i.RLock()
	defer i.RUnlock()

	exp := struct {
		Files map[string]string `json:"files"`
	}{
		Files: i.files,
	}
	return json.Marshal(exp)
}

func (i *index) UnmarshalJSON(b []byte) error {
	i.Lock()
	defer i.Unlock()

	var exp struct {
		Files map[string]string
	}
	err := json.Unmarshal(b, &exp)
	if err != nil {
		return err
	}

	i.files = make(map[string]string)
	i.index = make(map[string]bool)
	for k, v := range exp.Files {
		i.index[v] = true
		i.files[k] = v
	}
	return nil
}

// Get implements Index.
func (i *index) Get(path string) (string, bool) {
	i.RLock()
	defer i.RUnlock()

	x, ok := i.files[path]
	return x, ok
}

// Add implements Index.
func (i *index) Add(path, sum string) bool {
	i.Lock()
	defer i.Unlock()

	i.files[path] = sum
	old := i.index[sum]
	i.index[sum] = true
	return old
}

// Exists implements Index.
func (i *index) Exists(sum string) bool {
	i.RLock()
	defer i.RUnlock()

	return i.index[sum]
}

// Add the path+data to the index.  Returns the content sum and true if the content
// already existed in the index, false otherwise.
func AddContent(idx Index, path string, content []byte) (string, bool) {
	if x, ok := idx.Get(path); ok {
		return x, true
	}
	s := fmt.Sprintf("%x", sha1.Sum(content))
	return s, idx.Add(path, s)
}

// FileSystem is a type which defines a content addressable filesystem.
type FileSystem struct {
	idx Index

	fs store.RWFileSystem
}

// Open the file with the given path.  Uses the internal index to identify
// which file to open.  NB: Stat on the returned file will refer to the
// internal file.
func (s *FileSystem) Open(path string) (http.File, error) {
	realPath, ok := s.idx.Get(path)
	if !ok {
		return nil, fmt.Errorf("no such file: %v", path)
	}
	return s.open(realPath)
}

// Wait implements RWFileSystem.
func (s *FileSystem) Wait() error { return nil }

func (s *FileSystem) open(path string) (http.File, error) {
	return s.fs.Open(path)
}

func (s *FileSystem) create(path string) (io.WriteCloser, error) {
	return s.fs.Create(path)
}

type file struct {
	bytes.Buffer

	fs   *FileSystem
	path string
}

// Close acts as a signal that all the information has been written to
// the underlying buffer, and the file can be written to the RWFileSystem.
func (a *file) Close() error {
	_, ok := a.fs.idx.Get(a.path)
	if ok {
		return fmt.Errorf("file already exists: %v", a.path)
	}

	path, ok := AddContent(a.fs.idx, a.path, a.Bytes())
	if !ok {
		fmt.Println("creating", path)
		f, err := a.fs.create(path)
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer f.Close()
		_, err = io.Copy(f, a)
		if err != nil {
			return fmt.Errorf("error copying data into file '%v': %v", path, err)
		}
	}
	err := a.fs.persist()
	if err != nil {
		return err
	}
	return nil
}

// Create a new file with path. We buffer the contents written to the io.WriteCloser
// so that the content can be hashed and then written to the underlying RWFileSystem.
func (s *FileSystem) Create(path string) (io.WriteCloser, error) {
	_, ok := s.idx.Get(path)
	if ok {
		return nil, fmt.Errorf("file already exists for '%v'", path)
	}

	return &file{
		Buffer: bytes.Buffer{},
		path:   path,
		fs:     s,
	}, nil
}

// initIndex initialises the cafs index.
func (s *FileSystem) initIndex() error {
	f, err := s.open(".idx")
	if err != nil {
		// FIXME: Improve this
		// Can't guarantee that we will get an IsNotExist(err) here
		return nil
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error reading index: %v", err)
	}
	err = json.Unmarshal(b, s.idx)
	if err != nil {
		return fmt.Errorf("error decoding index: %v", err)
	}
	fmt.Printf("Index initialised: %d files (%d paths)", len(s.idx.(*index).index), len(s.idx.(*index).files))
	return nil
}

func (s *FileSystem) persist() error {
	f, err := s.fs.Create(".idx")
	if err != nil {
		return fmt.Errorf("error creating index: %v", err)
	}
	defer f.Close()

	b, err := json.Marshal(s.idx)
	if err != nil {
		return fmt.Errorf("error encoding index: %v", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("error writing index: %v", err)
	}
	return nil
}

// New creates a new content addressable RWFileSystem.
func New(fs store.RWFileSystem) (*FileSystem, error) {
	s := &FileSystem{
		idx: NewIndex(),
		fs:  fs,
	}
	err := s.initIndex()
	if err != nil {
		return nil, err
	}
	return s, nil
}
