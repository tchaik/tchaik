package store

import (
	"net/http"
	"testing"

	"golang.org/x/net/context"
)

type pathRecordFS struct {
	path string
}

func (r *pathRecordFS) Open(ctx context.Context, path string) (http.File, error) {
	r.path = path
	return nil, nil
}

func TestPathRewrite(t *testing.T) {
	tests := []struct {
		trimPrefix, addPrefix string
		in, out                  string
	}{
		{
			"", "", "test", "test",
		},
		{
			"t", "", "test", "est",
		},
		{
			"", "t", "test", "ttest",
		},
		{
			"t", "t", "test", "test",
		},
	}

	for ii, tt := range tests {
		rfs := &pathRecordFS{}
		fs := PathRewrite(rfs, tt.trimPrefix, tt.addPrefix)
		fs.Open(context.TODO(), tt.in) // ignore the response

		if tt.out != rfs.path {
			t.Errorf("[%d] recorded path: %#v, expected %#v", ii, rfs.path, tt.out)
		}
	}
}
