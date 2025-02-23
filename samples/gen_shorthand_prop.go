package main

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/strings"

type Rec struct {
	SField string
	IField int
}

func main() {
	r1 := Rec{SField: "abc", IField: 123}
	r2 := Rec{SField: "def", IField: 456}
	frt.PipeUnit(frt.Pipe(frt.Pipe(([]Rec{r1, r2}), (func(_r0 []Rec) []string {
		return slice.Map(func(_v1 Rec) string {
			return _v1.SField
		}, _r0)
	})), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), frt.Println)
}
