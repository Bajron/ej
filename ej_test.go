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
		t.Error("Single key should return one step")
	}

	r = GetSteps("key1.key2")
	if len(r) != 2 {
		t.Error("Single field reference should return 2 steps")
	}
	r = GetSteps("key1.key2.key3")
	if len(r) != 3 {
		t.Error("Double field reference should return 3 steps")
	}
}
