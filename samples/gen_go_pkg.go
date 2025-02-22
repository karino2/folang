package main

import "github.com/karino2/folang/pkg/frt"

import "path/filepath"

func main() {
	target := "/home/karino2/src/folang/samples/README.md"
	frt.PipeUnit(filepath.Dir(target), frt.Println)
	frt.PipeUnit(filepath.Base(target), frt.Println)
}
