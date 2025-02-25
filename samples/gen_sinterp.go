package main

import "github.com/karino2/folang/pkg/frt"

func main() {
	a := 123
	b := "str val"
	frt.PipeUnit(frt.SInterP("a is :%s, b is %s", a, b), frt.Println)
}
