package main

import "fmt"

type IorS interface {
	IorS_Union()
}

func (*IorS_IT) IorS_Union() {}
func (*IorS_ST) IorS_Union() {}

type IorS_IT struct {
	Value int
}

func New_IorS_IT(v int) IorS { return &IorS_IT{v} }

type IorS_ST struct {
	Value string
}

func New_IorS_ST(v string) IorS { return &IorS_ST{v} }

func main() {
	switch _v1 := (New_IorS_IT(3)).(type) {
	case *IorS_IT:
		ival := _v1.Value
		fmt.Printf("ival=%d\n", ival)
	case *IorS_ST:
		sval := _v1.Value
		fmt.Printf("ival=%s\n", sval)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
