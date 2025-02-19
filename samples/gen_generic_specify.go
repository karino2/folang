package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func main() {
	s := slice.New[string]()
	ss := slice.PushHead("hoge", s)
	e := slice.Head(ss)
	frt.Println(e)
}
