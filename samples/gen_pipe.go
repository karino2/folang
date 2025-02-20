package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func main() {
	s := []int{5, 6, 7, 8}
	s2 := frt.Pipe(s, (func(_r0 []int) []int { return slice.Take(2, _r0) }))
	frt.Printf1("%v\n", s2)
}
