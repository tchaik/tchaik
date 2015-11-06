package player

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreatePlayerEmptyRequest(t *testing.T) {
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

func TestCreatePlayer(t *testing.T) {
	ps := NewPlayers()
	ps.Add(testPlayer("1"))

	h := NewHTTPHandler(ps)

	in := struct {
		Key        string   `json:"key"`
		PlayerKeys []string `json:"playerKeys"`
	}{
		Key:        "2",
		PlayerKeys: []string{"1"},
	}

	b, err := json.Marshal(in)
	if err != nil {
		t.Errorf("unexpected error in json.Marshal(): %v", err)
	}

	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "", bytes.NewReader(b))
	if err != nil {
		t.Errorf("unexpected error creating request: %v", err)
	}

	h.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Errorf("w.Code = %d, expected %d", w.Code, http.StatusCreated)
	}
}

func TestRemovePlayer(t *testing.T) {
	ps := NewPlayers()
	ps.Add(testPlayer("1"))

	h := NewHTTPHandler(ps)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("DELETE", "1", nil)
	if err != nil {
		t.Errorf("unexpected error creating request: %v", err)
	}

	h.ServeHTTP(w, r)
	if w.Code != http.StatusNoContent {
		t.Errorf("w.Code = %d, expected %d", w.Code, http.StatusNoContent)
	}

	n := len(ps.List())
	if n != 0 {
		t.Errorf("len(ps.List()) = %d, expected %d", n, 0)
	}
}
