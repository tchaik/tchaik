package player

import (
	"encoding/json"
	"fmt"
	"testing"
)

type testPlayer string

func (p testPlayer) Key() string           { return string(p) }
func (testPlayer) Do(a Action) error       { return nil }
func (testPlayer) SetMute(bool) error      { return nil }
func (testPlayer) SetRepeat(bool) error    { return nil }
func (testPlayer) SetVolume(float64) error { return nil }
func (testPlayer) SetTime(float64) error   { return nil }

func TestPlayers(t *testing.T) {
	oneKey := "one"
	onePl := testPlayer(oneKey)

	ps := NewPlayers()

	list := ps.List()
	if len(list) != 0 {
		t.Errorf("len(ps.List()) = %d, expected: %d", len(list), 0)
	}

	ps.Add(onePl)

	p := ps.Get(oneKey)
	if p != onePl {
		t.Errorf("Get(%#v) = %#v, expected: %#v", oneKey, p, onePl)
	}

	b, err := json.Marshal(p)
	if err != nil {
		t.Errorf("json.Marshal(p) returned an unexpected error: %v", err)
	}

	list = ps.List()
	if len(list) != 1 {
		t.Errorf("len(ps.List()) = %d, expected: %d", len(list), 1)
	}

	ps.Remove(oneKey)

	p = ps.Get(oneKey)
	if p != nil {
		t.Errorf("Get(%#v) = %#v, expected: %#v", oneKey, p, nil)
	}

	list = ps.List()
	if len(list) != 0 {
		t.Errorf("len(ps.List()) = %d, expected: %d", len(list), 0)
	}
}
