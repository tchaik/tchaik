// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"reflect"
	"testing"
)

func TestSplitMultiple(t *testing.T) {
	tests := []struct {
		in         string
		seperators []string
		out        []string
	}{
		{
			"",
			nil,
			nil,
		},
		{
			"One",
			nil,
			[]string{"One"},
		},
		{
			"One",
			[]string{","},
			[]string{"One"},
		},
		{
			"One, Two",
			[]string{","},
			[]string{"One", "Two"},
		},
		{
			"One, Two & Three",
			[]string{",", "&"},
			[]string{"One", "Two", "Three"},
		},
		{
			"Vernon Handley",
			ListSeparators,
			[]string{"Vernon Handley"},
		},
	}

	for ii, tt := range tests {
		res := splitMultiple(tt.in, tt.seperators)
		if !reflect.DeepEqual(res, tt.out) {
			t.Errorf("[%d] splitMultiple(%#v, %#v) = %#v, expected: %#v", ii, tt.in, tt.seperators, res, tt.out)
		}
	}
}
