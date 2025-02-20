package main

import "github.com/karino2/folang/pkg/frt"

func ApplyL[T0 any, T1 any, T2 any](fn func(T0) T1, tup frt.Tuple2[T0, T2]) frt.Tuple2[T1, T2] {
	nl := frt.Pipe(frt.Fst(tup), fn)
	return frt.NewTuple2(nl, frt.Snd(tup))
}

func add(a int, b int) int {
	return (a + b)
}

func main() {
	frt.PipeUnit(frt.Pipe(frt.NewTuple2(123, "hoge"), (func(_r0 frt.Tuple2[int, string]) frt.Tuple2[int, string] {
		return ApplyL((func(_r0 int) int { return add(456, _r0) }), _r0)
	})), (func(_r0 frt.Tuple2[int, string]) { frt.Printf1("%v\n", _r0) }))
}
