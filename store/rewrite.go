package store

import (
	"net/http"
	"strings"
)

// PathRewrite is a convenience type which wraps an http.FileSystem and rewrites paths
// which are passed to the Open method.
type PathRewrite struct {
	http.FileSystem

	trimPrefix, addPrefix string
}

// NewPathRewrite creates a new PathRewrite which will trim prefixes and then add prefixes
// to paths.
func NewPathRewrite(fs http.FileSystem, trimPrefix, addPrefix string) http.FileSystem {
	return &PathRewrite{
		FileSystem: fs,
		trimPrefix: trimPrefix,
		addPrefix:  addPrefix,
	}
}

// Open implements http.FileSystem.
func (p *PathRewrite) Open(path string) (http.File, error) {
	return p.FileSystem.Open(p.addPrefix + strings.TrimPrefix(path, p.trimPrefix))
}
