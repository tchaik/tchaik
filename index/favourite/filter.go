package favourite

import "tchaik.com/index"

type Filter struct {
	Store
}

func (f *Filter) Filter(paths []index.Path) []index.Path {
	result := make([]index.Path, 0, len(paths))
	for _, p := range paths {
		if f.Get(p) {
			result = append(result, p)
		}
	}
	return result[:len(result):len(result)]
}
