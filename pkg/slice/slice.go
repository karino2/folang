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

func Len[T any](s []T) int {
	return len(s)
}

func New[T any]() []T {
	return []T{}
}

func Item[T any](index int, s []T) T {
	return s[index]
}

func IsEmpty[T any](a []T) bool {
	return len(a) == 0
}

func IsNotEmpty[T any](a []T) bool {
	return len(a) != 0
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

func PopLast[T any](s []T) []T {
	return s[0:(len(s) - 1)]
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

func PushLast[T any](elem T, s []T) []T {
	return append(s, elem)
}

func PushHead[T any](elem T, s []T) []T {
	ret := []T{elem}
	return append(ret, s...)
}

func Collect[T any, U any](f func(T) []U, ss []T) []U {
	var res []U
	for _, e := range ss {
		one := f(e)
		res = append(res, one...)
	}
	return res
}

func Concat[T any](ss [][]T) []T {
	var res []T
	for _, s := range ss {
		res = append(res, s...)
	}
	return res
}

func Append[T any](s1 []T, s2 []T) []T {
	var res []T
	res = append(res, s1...)
	res = append(res, s2...)
	return res
}

func Distinct[T comparable](ss []T) []T {
	set := make(map[T]bool)
	res := []T{}
	for _, e := range ss {
		if _, ok := set[e]; !ok {
			res = append(res, e)
			set[e] = true
		}
	}
	return res
}

func TryFind[T any](pred func(T) bool, ss []T) frt.Tuple2[T, bool] {
	for _, e := range ss {
		if pred(e) {
			return frt.NewTuple2(e, true)
		}
	}
	var res T
	return frt.NewTuple2(res, false)
}

func Fold[T any, S any](folder func(S, T) S, iniS S, ss []T) S {
	stat := iniS
	for _, e := range ss {
		stat = folder(stat, e)
	}
	return stat
}
