package favourite

import (
	"fmt"

	"tchaik.com/index"
)

// Filter is a wrapper for Store which defines the Filter method.
type Filter struct {
	Store
}

// Filter combs through the list of paths, returning paths which are favourites (as
// determined by the underlying Store).
func (f *Filter) Filter(paths []index.Path) []index.Path {
	result := make([]index.Path, 0, len(paths))
	for _, p := range paths {
		if f.Get(p) {
			result = append(result, p)
		}
	}
	return result[:len(result):len(result)]
}

// RootFilter is a wrapper for Store which defines the Filter method.
type RootFilter struct {
	Store
}

// Filter combs through the list of paths, return paths which are/contain favourites
// as determined by the underlying Store.
func (f *RootFilter) Filter(paths []index.Path) []index.Path {
	exp := make(map[string]bool)
	for _, p := range f.Store.List() {
		if len(p) > 1 {
			exp[fmt.Sprintf("%v", p[:2])] = true
		}
	}

	result := make([]index.Path, 0, len(paths))
	for _, p := range paths {
		if exp[fmt.Sprintf("%v", p)] {
			result = append(result, p)
		}
	}
	return result
}
