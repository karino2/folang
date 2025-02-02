package frt

import (
	"cmp"
	"fmt"

	gcmp "github.com/google/go-cmp/cmp"
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

func OpEqual[T any](e1 T, e2 T) bool {
	return gcmp.Equal(e1, e2)
}

func OpNotEqual[T any](e1 T, e2 T) bool {
	return !OpEqual(e1, e2)
}

func IfElse[T any](cond bool, tbody func() T, fbody func() T) T {
	if cond {
		return tbody()
	} else {
		return fbody()
	}
}

func IfElseUnit(cond bool, tbody func(), fbody func()) {
	if cond {
		tbody()
	} else {
		fbody()
	}
}

// For no else condition, return type must be unit.
func IfOnly(cond bool, tbody func()) {
	if cond {
		tbody()
	}
}
