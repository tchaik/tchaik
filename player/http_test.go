package player

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreatePlayer(t *testing.T) {
	ps := NewPlayers()
	h := NewHTTPHandler(ps)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "", nil)
	if err != nil {
		t.Errorf("unexpected error creating request: %v", err)
	}
	r.Body = ioutil.NopCloser(strings.NewReader(""))

	h.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("w.Code = %d, expected %d", w.Code, http.StatusBadRequest)
	}

	if !bytes.HasPrefix(w.Body.Bytes(), []byte("error parsing JSON")) {
		t.Errorf("w.Body = %#v, expected %#v", string(w.Body.Bytes()), "error parsing JSON...")
	}
}
