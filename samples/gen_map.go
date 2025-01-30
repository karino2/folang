package main

import "fmt"

import "github.com/karino2/folang/pkg/slice"

func conv(i int) string {
	return fmt.Sprintf("a %d", i)
}

func main() {
	s := []int{5, 6, 7, 8}
	s2 := slice.Map[int, string](conv, s)
	fmt.Printf("%v", s2)
}
