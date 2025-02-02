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

func sfirst[T any](s []T) T {
	return s[0]
}

func TestPipeGenerics(t *testing.T) {
	goti := Pipe([]int{5, 4, 3}, sfirst)
	if goti != 5 {
		t.Errorf("goti %v", goti)
	}

	gots := Pipe([]string{"d", "c", "b"}, sfirst)
	if gots != "d" {
		t.Errorf("gots %s", gots)
	}
}

func Take[T any](num int, s []T) []T {
	var res []T

	for i := 0; i < num; i++ {
		res = append(res, s[i])
	}

	return res
}

func TestPipeTake(t *testing.T) {
	s := []int{1, 2, 3}
	got := Pipe[[]int, []int](s, (func(_r0 []int) []int { return Take(2, _r0) }))
	if len(got) != 2 || got[1] != 2 {
		t.Errorf("got %v", got)
	}
}

func TestOpPlus(t *testing.T) {
	got := OpPlus[string]("abc", "def")
	if got != "abcdef" {
		t.Errorf("got %v", got)
	}
}

func TestOpEaualSlice(t *testing.T) {
	s1 := []string{"1", "2"}
	s2 := []string{"1", "2"}
	s3 := []string{"2", "1"}
	got := OpEqual(s1, s2)
	if !got {
		t.Error("not equal slice.")
	}
	got = OpEqual(s1, s3)
	if got {
		t.Error("wrongly equal the same slice.")
	}
}

func TestSprintf1(t *testing.T) {
	got := Sprintf1("a=%d", 123)
	if got != "a=123" {
		t.Errorf("got %s", got)
	}
	got = Sprintf1("b=%s", "hoge")
	if got != "b=hoge" {
		t.Errorf("got %s", got)
	}
}
