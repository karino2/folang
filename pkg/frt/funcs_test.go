package frt

import (
	"testing"
)

func add1(a int) int { return a + 1 }

func TestPipe(t *testing.T) {
	got := Pipe(5, add1)
	if got != 6 {
		t.Errorf("Want 6, got %d", got)
	}
}
