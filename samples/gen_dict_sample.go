package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/dict"

func main() {
	d := dict.New[string, int]()
	dict.Add(d, "hoge", 123)
	dict.Add(d, "ika", 456)
	i1 := dict.Item(d, "hoge")
	i2 := dict.Item(d, "ika")
	frt.Printf1("sum=%d\n", (i1 + i2))
}
