package main

import "github.com/karino2/folang/pkg/frt"

func ika() frt.Tuple2[int, string] {
	return frt.NewTuple2(123, "abc")
}

func main() {
	a, _ := frt.Destr(ika())
	frt.PipeUnit(frt.Sprintf1("a=%d", a), frt.Println)
}
