package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func main() {
	s := ([]int{5, 6, 7, 8})
	s2 := slice.Take(2, s)
	frt.Printf1("%v\n", s2)
}
