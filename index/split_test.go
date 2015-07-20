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
			[]string{""},
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
	}

	for ii, tt := range tests {
		res := splitMultiple(tt.in, tt.seperators)
		if !reflect.DeepEqual(res, tt.out) {
			t.Errorf("[%d] splitMultiple(%#v, %#v) = %#v, expected: %#v", ii, tt.in, tt.seperators, res, tt.out)
		}
	}
}
