// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"sync"

	"tchaik.com/index"
	"tchaik.com/index/attr"
)

// newBootstrapSearch creates a new index.Searcher which builds the search index
// on the first call to Search.
func newBootstrapSearcher(root index.Collection) index.Searcher {
	return &bootstrapSearcher{
		root: root,
	}
}

type bootstrapSearcher struct {
	once sync.Once
	root index.Collection

	index.Searcher
}

func (b *bootstrapSearcher) bootstrap() {
	wi := index.BuildCollectionWordIndex(b.root, []string{"Composer", "Artist", "Album", "Name"})
	b.Searcher = index.FlatSearcher{
		Searcher: index.WordsIntersectSearcher(index.BuildPrefixExpandSearcher(wi, wi, 10)),
	}
}

// Search implements index.Searcher.
func (b *bootstrapSearcher) Search(input string) []index.Path {
	b.once.Do(b.bootstrap)
	return b.Searcher.Search(input)
}

// newBootstrapFilter creates a new index.Filter which initialises the filter on the
// first call to Filter.
func newBootstrapFilter(root index.Collection, field attr.Interface) index.Filter {
	return &bootstrapFilter{
		root:  root,
		field: field,
	}
}

type bootstrapFilter struct {
	once  sync.Once
	root  index.Collection
	field attr.Interface

	index.Filter
}

func (b *bootstrapFilter) bootstrap() {
	b.Filter = index.FilterCollection(b.root, b.field)
}

// Items implements index.Filter.
func (b *bootstrapFilter) Items() []index.FilterItem {
	b.once.Do(b.bootstrap)
	return b.Filter.Items()
}

type bootstrapRecent struct {
	once sync.Once
	root index.Collection
	n    int

	list []index.Path
}

func (b *bootstrapRecent) bootstrap() {
	b.list = index.Recent(b.root, b.n)
}

// List implements index.Lister.
func (b *bootstrapRecent) List() []index.Path {
	b.once.Do(b.bootstrap)
	return b.list
}
