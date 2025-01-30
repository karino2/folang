package main

import "fmt"

import "github.com/karino2/folang/pkg/slice"

func main() {
	s := []int{5, 6, 7, 8}
	s2 := slice.Take[int](2, s)
	fmt.Printf("%v", s2)
}
