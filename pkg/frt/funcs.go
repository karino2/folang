package frt

import (
	"fmt"

	gcmp "github.com/google/go-cmp/cmp"
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

// 1 arg sprintf. This is often used and enough for self host.
func Sprintf1[T any](fmtstr string, arg T) string {
	return fmt.Sprintf(fmtstr, arg)
}

func OpEqual[T any](e1 T, e2 T) bool {
	return gcmp.Equal(e1, e2)
}

func OpNotEqual[T any](e1 T, e2 T) bool {
	return !OpEqual(e1, e2)
}

// should handle short cut, NYI.
func OpAnd(e1 bool, e2 bool) bool {
	return e1 && e2
}

func OpNot(e1 bool) bool {
	return !e1
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

type Tuple2[T any, U any] struct {
	E0 T
	E1 U
}

func NewTuple2[T any, U any](e0 T, e1 U) Tuple2[T, U] {
	return Tuple2[T, U]{E0: e0, E1: e1}
}

// [Tuples - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/tuples)
func Fst[T any, U any](tup Tuple2[T, U]) T { return tup.E0 }
func Snd[T any, U any](tup Tuple2[T, U]) U { return tup.E1 }

// Destructuring.
func Destr[T any, U any](tup Tuple2[T, U]) (T, U) { return tup.E0, tup.E1 }

func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func Panic(msg string) {
	panic(msg)
}
