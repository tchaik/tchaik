// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import "sort"

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

// Filter creates a slice of FilterItems, each FilterItem is a
func Filter(c Collection, field string) []FilterItem {
	m := make(map[string][]Path)
	walkfn := func(t Track, p Path) error {
		f := t.GetString(field)
		m[f] = append(m[f], p)
		return nil
	}
	Walk(c, Path([]Key{"Root"}), walkfn)

	result := make([]FilterItem, 0, len(m))
	for k, v := range m {
		result = append(result, &filterItem{
			name:   k,
			fields: make(map[string]interface{}),
			paths:  Union(v),
		})
	}
	sort.Sort(FilterItemSlice(result))
	return result
}
