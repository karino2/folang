package slice

import (
	"testing"
)

/*
func Sort[T cmp.Ordered](s []T) []T {
	// non destructive sort. copy arg first.
	res := append(s[:0:0], s...)
	sort.Slice(res, cmp.Less)
	return res
}
*/

func TestSort(t *testing.T) {
	org := []string{"ika", "hoge", "def", "abc", "zzz"}
	got := Sort(org)
	if org[0] != "ika" {
		t.Errorf("original slice is modified. %v", org)
	}
	if got[0] != "abc" {
		t.Errorf("result not sorted: %v", got)
	}
}

func TestZip(t *testing.T) {
	s1 := ([]int{1, 2, 3})
	s2 := ([]string{"a", "b", "c"})
	got := Zip[int, string](s1, s2)

	if len(got) != 3 {
		t.Errorf("got %v", got)
	}
}
