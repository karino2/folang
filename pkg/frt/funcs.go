package frt

import (
	"cmp"
	"fmt"

	"golang.org/x/exp/constraints"
)

func Pipe[T any, U any](elem T, f func(T) U) U {
	return f(elem)
}

func PipeUnit[T any](elem T, f func(T)) {
	f(elem)
}

func Println(str string) {
	fmt.Println(str)
}

func OpPlus[T cmp.Ordered](e1 T, e2 T) T {
	return e1 + e2
}

type Numeric interface {
	constraints.Integer | constraints.Float
}

func OpMinus[T Numeric](e1 T, e2 T) T {
	return e1 - e2
}

func OpEqual[T comparable](e1 T, e2 T) bool {
	return e1 == e2
}

func OpNotEqual[T comparable](e1 T, e2 T) bool {
	return e1 != e2
}
