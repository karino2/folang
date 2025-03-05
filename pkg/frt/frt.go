package frt

import (
	"fmt"
	"reflect"

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

func Sprintf2[T any, U any](fmtstr string, arg0 T, arg1 U) string {
	return fmt.Sprintf(fmtstr, arg0, arg1)
}

func Printf1[T any](fmtstr string, arg T) {
	fmt.Printf(fmtstr, arg)
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
func Destr2[T any, U any](tup Tuple2[T, U]) (T, U) { return tup.E0, tup.E1 }

// backward compat, obsolete.
func Destr[T any, U any](tup Tuple2[T, U]) (T, U) { return tup.E0, tup.E1 }

type Tuple3[T any, U any, V any] struct {
	E0 T
	E1 U
	E2 V
}

func NewTuple3[T any, U any, V any](e0 T, e1 U, e2 V) Tuple3[T, U, V] {
	return Tuple3[T, U, V]{E0: e0, E1: e1, E2: e2}
}

func Destr3[T any, U any, V any](tup Tuple3[T, U, V]) (T, U, V) { return tup.E0, tup.E1, tup.E2 }

func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func Panic(msg string) {
	panic(msg)
}

func Panicf1[T any](fmt string, arg T) {
	msg := Sprintf1(fmt, arg)
	panic(msg)
}

func Panicf2[T any, U any](fmt string, arg0 T, arg1 U) {
	msg := Sprintf2(fmt, arg0, arg1)
	panic(msg)
}

func Empty[T any]() T {
	var res T
	return res
}

func toS(arg any) string {
	rval := reflect.ValueOf(arg)
	switch rval.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", rval.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", rval.Float())
	case reflect.String:
		return rval.String()
	default:
		return fmt.Sprintf("%v", arg)
	}
}

func SInterP(fmt1 string, args ...any) string {
	var sargs []any
	for _, arg := range args {
		sargs = append(sargs, toS(arg))
	}
	return fmt.Sprintf(fmt1, sargs...)
}
