// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"sort"

	"tchaik.com/index/attr"
)

// FilterItem is an interface which defines behaviour for creating arbitrary filters where each
// filtered item is a list of Paths.
type FilterItem interface {
	Name() string
	Fields() map[string]interface{}
	Paths() []Path
}

// FilterItemSlice is a convenience type which implements sort.Interface and is used to sort
// slices of FilterItem.
type FilterItemSlice []FilterItem

func (f FilterItemSlice) Len() int           { return len(f) }
func (f FilterItemSlice) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f FilterItemSlice) Less(i, j int) bool { return f[i].Name() < f[j].Name() }

type filterItem struct {
	name   string
	fields map[string]interface{}
	paths  []Path
}

// Name implements FilterItem.
func (f *filterItem) Name() string { return f.name }

// Fields implements FilterItem.
func (f *filterItem) Fields() map[string]interface{} { return f.fields }

// Paths implements FilterItem.
func (f *filterItem) Paths() []Path { return f.paths }

// Filter is an iterface which defines the Items method.
type Filter interface {
	// Items returns a list of FilterItems.
	Items() []FilterItem
}

type filter struct {
	items []FilterItem
}

func (f filter) Items() []FilterItem {
	return f.items
}

// FilterCollection creates Filter of the Collection using fields to partition
// Tracks in a collection.
func FilterCollection(c Collection, field attr.Interface) Filter {
	m := make(map[string][]Path)
	walkfn := func(t Track, p Path) error {
		f := field.Value(t)
		switch f := f.(type) {
		case string:
			m[f] = append(m[f], p)
		case []string:
			for _, x := range f {
				m[x] = append(m[x], p)
			}
		}
		return nil
	}
	Walk(c, Path([]Key{"Root"}), walkfn)

	items := make([]FilterItem, 0, len(m))
	for k, v := range m {
		items = append(items, &filterItem{
			name:   k,
			fields: make(map[string]interface{}),
			paths:  Union(v),
		})
	}
	sort.Sort(FilterItemSlice(items))
	return filter{items}
}
