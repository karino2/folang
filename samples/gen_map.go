package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func conv(i int) string {
	return frt.Sprintf1("a %d", i)
}

func main() {
	s := ([]int{5, 6, 7, 8})
	s2 := slice.Map(conv, s)
	frt.Printf1("%v\n", s2)
}
