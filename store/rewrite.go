package store

import (
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

// PathRewrite is a convenience type which wraps an http.FileSystem and rewrites paths
// which are passed to the Open method.
type PathRewrite struct {
	FileSystem

	trimPrefix, addPrefix string
}

// NewPathRewrite creates a new PathRewrite which will trim prefixes and then add prefixes
// to paths.
func NewPathRewrite(fs FileSystem, trimPrefix, addPrefix string) FileSystem {
	return &PathRewrite{
		FileSystem: fs,
		trimPrefix: trimPrefix,
		addPrefix:  addPrefix,
	}
}

// Open implements FileSystem.
func (p *PathRewrite) Open(ctx context.Context, path string) (http.File, error) {
	return p.FileSystem.Open(ctx, p.addPrefix+strings.TrimPrefix(path, p.trimPrefix))
}
