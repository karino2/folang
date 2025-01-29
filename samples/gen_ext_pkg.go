package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/buf"

func main() {
	bb := buf.New()
	buf.Write(bb, "hello")
	buf.Write(bb, "world")
	res := buf.String(bb)
	frt.Println(res)
}
