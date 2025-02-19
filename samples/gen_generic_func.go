package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func hoge[T0 any](a []T0) T0 {
	return slice.Head(a)
}

func main() {
	b := ([]int{1, 2, 3})
	c := hoge(b)
	frt.Printf1("%d\n", c)
}
