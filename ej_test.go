package main

import (
	"testing"
)

func TestGetSteps(t *testing.T) {
	r := GetSteps("")
	if len(r) != 0 {
		t.Error("GetSteps from empty string should be empty")
	}

	r = GetSteps("key")
	if len(r) != 1 {
		t.Errorf("Single key should return one step. Got %d", len(r))
	}
	if r[0].kind != KEY {
		t.Error("Key should have proper kind")
	}
	if r[0].name != "key" {
		t.Errorf("Key should have name as provided. Got %s", r[0].name)
	}

	r = GetSteps("key1.key2")
	if len(r) != 2 {
		t.Error("Single field reference should return 2 steps")
	}
	r = GetSteps("key1.key2.key3")
	if len(r) != 3 {
		t.Error("Double field reference should return 3 steps")
	}
	if r[1].name != "key2" {
		t.Errorf("expected key2, got %s", r[1].name)
	}

	r = GetSteps("[0]")
	if len(r) != 1 {
		t.Errorf("Single index should return one step. Got %d", len(r))
	}
	if r[0].kind != INDEX {
		t.Error("Index step should be marked as such")
	}

	r = GetSteps("[0].key")
	if l := len(r); l != 2 {
		t.Errorf("Key from index should return 2 steps. Got %d", l)
	}

	r = GetSteps("key[2]")
	if l := len(r); l != 2 {
		t.Errorf("Index from key should return 2 steps. Got %d", l)
	}
	if r[1].index != 2 {
		t.Errorf("Index should have value provided. Got %d", r[1].index)
	}
}
