package slice

import (
	"cmp"
	"slices"

	"github.com/karino2/folang/pkg/frt"
)

/*
Try to make similar to F# list.
[List (FSharp.Core) - FSharp.Core](https://fsharp.github.io/fsharp-core-docs/reference/fsharp-collections-listmodule.html)
*/

func Length[T any](s []T) int {
	return len(s)
}

func Last[T any](s []T) T {
	return s[len(s)-1]
}

func Head[T any](s []T) T {
	if len(s) == 0 {
		panic("call Head to empty list")
	}
	return s[0]
}

func Tail[T any](s []T) []T {
	if len(s) == 0 {
		panic("call Tail to empty list")
	}
	return s[1:]
}

func Take[T any](num int, s []T) []T {
	var res []T

	for i := 0; i < num; i++ {
		res = append(res, s[i])
	}

	return res
}

func Skip[T any](count int, s []T) []T {
	var res []T

	for i := count; i < len(s); i++ {
		res = append(res, s[i])
	}

	return res
}

func Map[T any, U any](f func(T) U, s []T) []U {
	var res []U
	for _, e := range s {
		res = append(res, f(e))
	}
	return res
}

func Mapi[T any, U any](f func(int, T) U, s []T) []U {
	var res []U
	for i, e := range s {
		res = append(res, f(i, e))
	}
	return res
}

func Iter[T any](action func(T), s []T) {
	for _, e := range s {
		action(e)
	}
}

func Filter[T any](pred func(T) bool, s []T) []T {
	var res []T
	for _, e := range s {
		if pred(e) {
			res = append(res, e)
		}
	}
	return res
}

func Sort[T cmp.Ordered](s []T) []T {
	// non destructive sort. copy arg first.
	res := append(s[:0:0], s...)
	slices.SortFunc(res, cmp.Compare)
	return res
}

func SortBy[T any, U cmp.Ordered](proj func(T) U, s []T) []T {
	// non destructive sort. copy arg first.
	res := append(s[:0:0], s...)
	slices.SortFunc(res, func(s1, s2 T) int { return cmp.Compare(proj(s1), proj(s2)) })
	return res
}

func Zip[T, U any](s1 []T, s2 []U) []frt.Tuple2[T, U] {
	var ret []frt.Tuple2[T, U]
	if len(s1) != len(s2) {
		panic("zip with different length slices.")
	}
	for i, e1 := range s1 {
		ret = append(ret, frt.NewTuple2(e1, s2[i]))
	}
	return ret
}

func Forall[T any](pred func(T) bool, s []T) bool {
	for _, e := range s {
		if !pred(e) {
			return false
		}
	}
	return true
}

func Forany[T any](pred func(T) bool, s []T) bool {
	for _, e := range s {
		if pred(e) {
			return true
		}
	}
	return false
}

func Append[T any](elem T, s []T) []T {
	return append(s, elem)
}

func Prepend[T any](elem T, s []T) []T {
	ret := []T{elem}
	return append(ret, s...)
}
