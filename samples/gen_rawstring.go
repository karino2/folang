package main

import "github.com/karino2/folang/pkg/frt"

func main() {
	rawstr := "This\nis raw\nstring. \"Yes!\" \\ backslash is not esape.\n"
	a := 123
	rawinterp := frt.SInterP("Raw String with\nstring interpolation. a=\"%s\".\nThere is no way to escape brace in this expression.", a)
	frt.PipeUnit((rawstr + rawinterp), frt.Println)
}
