package slice

import (
	"cmp"
	"slices"
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

func Take[T any](num int, s []T) []T {
	var res []T

	for i := 0; i < num; i++ {
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
