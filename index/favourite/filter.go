package favourite

import "tchaik.com/index"

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
