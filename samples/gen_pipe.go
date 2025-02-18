package main

import "fmt"

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func main() {
	s := []int{5, 6, 7, 8}
	s2 := frt.Pipe(s, (func(_r0 []int) []int { return slice.Take(2, _r0) }))
	fmt.Printf("%v", s2)
}
