package main

import (
	"testing"
)

func TestGetSteps(t *testing.T) {
	r := GetSteps("")
	if len(r) != 0 {
		t.Error("GetSteps from empty string should be empty")
	}
}
