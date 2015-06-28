package store

import (
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

// PathRewrite creates a FileSystem wrapper which will trim a prefix and add a prefix
// to all paths passed to Open on the underlying FileSystem.
func PathRewrite(fs FileSystem, trimPrefix, addPrefix string) FileSystem {
	return &pathRewrite{
		FileSystem: fs,
		trimPrefix: trimPrefix,
		addPrefix:  addPrefix,
	}
}

type pathRewrite struct {
	FileSystem

	trimPrefix, addPrefix string
}

// Open implements FileSystem.
func (p *pathRewrite) Open(ctx context.Context, path string) (http.File, error) {
	return p.FileSystem.Open(ctx, p.addPrefix+strings.TrimPrefix(path, p.trimPrefix))
}
