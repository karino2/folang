package main

import "testing"

func TestIsUnresolve(t *testing.T) {
	got := IsUnresolved(New_FType_FUnresolved)
	if got != true {
		t.Errorf("want true, got %v", got)
	}
}
