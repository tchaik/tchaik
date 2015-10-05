// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"reflect"
	"testing"
)

func stringSet(s []string) map[string]bool {
	ss := make(map[string]bool)
	for _, x := range s {
		ss[x] = true
	}
	return ss
}

func TestPrefixMultiExpander(t *testing.T) {
	tests := []struct {
		words []string
		n     int
		in    string
		out   []string
	}{
		// Single word, prefix 3
		{
			[]string{"hello"},
			3,
			"hel",
			[]string{"hello"},
		},

		// Single word, prefix up to 4
		{
			[]string{"hello"},
			4,
			"hel",
			[]string{"hello"},
		},

		{
			[]string{"hello"},
			3,
			"he",
			[]string{"he"},
		},

		{
			[]string{"hello"},
			3,
			"x",
			[]string{"x"},
		},

		{
			[]string{"hello", "helloworld"},
			3,
			"hel",
			[]string{"hello", "helloworld"},
		},

		// // Two words, the same, different
		// {
		// 	[]string{"hello", "hellothere"},
		// 	10,
		// 	[]string{"", "h", "he", "hel", "hell", "hello", "hellot", "helloth", "hellothere"},
		// 	[]string{"", "", "", "", "", "", "hellothere", "hellothere", "hellothere"},
		// 	[]bool{false, false, false, false, false, false, true, true, true},
		// },
		//

		// Multiple unique words
		{
			[]string{"prokofiev", "shostakovich", "tchaikovsky", "rachmaninov", "rachmaninoff", "xenakis"},
			3,
			"pro",
			[]string{"prokofiev"},
		},

		// Multiple unique words
		{
			[]string{"prokofiev", "shostakovich", "tchaikovsky", "rachmaninov", "rachmaninoff", "xenakis"},
			3,
			"tch",
			[]string{"tchaikovsky"},
		},

		// Multiple unique words
		{
			[]string{"prokofiev", "shostakovich", "tchaikovsky", "rachmaninov", "rachmaninoff", "xenakis"},
			3,
			"rachmanin",
			[]string{"rachmaninoff", "rachmaninov"},
		}, // Multiple unique words
	}

	for ii, tt := range tests {
		pm := BuildPrefixMultiExpander(tt.words, tt.n)
		got := pm.Expand(tt.in)
		if !reflect.DeepEqual(stringSet(tt.out), stringSet(got)) {
			t.Errorf("[%d] Expand(%#v) = %#v expected: %#v (compared unordered)", ii, tt.in, got, tt.out)
		}
	}
}

func TestWordIndex(t *testing.T) {
	tests := []struct {
		in  map[string][]Path
		out map[string][]Path
	}{
		{
			map[string][]Path{
				"gustav mahler": []Path{
					Path{"Root", "Gustav Mahler", "Symphony No 1", "0"},
					Path{"Root", "Gustav Mahler", "Symphony No 1", "1"},
					Path{"Root", "Gustav Mahler", "Symphony No 1", "2"},
				},
			},
			map[string][]Path{
				"gustav": []Path{
					Path{"Root", "Gustav Mahler", "Symphony No 1", "0"},
					Path{"Root", "Gustav Mahler", "Symphony No 1", "1"},
					Path{"Root", "Gustav Mahler", "Symphony No 1", "2"},
				},
				"mahler": []Path{
					Path{"Root", "Gustav Mahler", "Symphony No 1", "0"},
					Path{"Root", "Gustav Mahler", "Symphony No 1", "1"},
					Path{"Root", "Gustav Mahler", "Symphony No 1", "2"},
				},
			},
		},
	}

	for ii, tt := range tests {
		w := &wordIndex{
			fields: []string{},
			words:  make(map[string][]Path),
		}

		for s, ps := range tt.in {
			for _, p := range ps {
				w.AddString(s, p)
			}
		}

		for k, v := range tt.out {
			ps := w.Search(k)
			if !reflect.DeepEqual(ps, v) {
				t.Errorf("[%d] does't match: %#v, %#v", ii, ps, v)
			}
		}
	}
}
