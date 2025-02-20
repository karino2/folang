package dict

import "github.com/karino2/folang/pkg/frt"

/*
  https://fsharp.github.io/fsharp-core-docs/reference/fsharp-collections-fsharpmap-2.html
  FSharp is Map. but constructor is dict. We use Dict for type name.
*/

type Dict[K comparable, V any] struct {
	Fdict map[K]V
}

func New[K comparable, V any]() Dict[K, V] {
	res := Dict[K, V]{}
	res.Fdict = make(map[K]V)
	return res
}

func Add[K comparable, V any](d Dict[K, V], key K, v V) {
	d.Fdict[key] = v
}

func ContainsKey[K comparable, V any](d Dict[K, V], key K) bool {
	_, ok := d.Fdict[key]
	return ok
}

func TryFind[K comparable, V any](d Dict[K, V], key K) frt.Tuple2[V, bool] {
	e, ok := d.Fdict[key]
	return frt.NewTuple2(e, ok)
}

func Item[K comparable, V any](d Dict[K, V], key K) V {
	e := d.Fdict[key]
	return e
}

func KVs[K comparable, V any](d Dict[K, V]) []frt.Tuple2[K, V] {
	var res []frt.Tuple2[K, V]
	for k, v := range d.Fdict {
		res = append(res, frt.NewTuple2(k, v))
	}
	return res
}

func ToDict[K comparable, V any](ss []frt.Tuple2[K, V]) Dict[K, V] {
	dic := New[K, V]()
	for _, tp := range ss {
		k, v := frt.Destr(tp)
		Add(dic, k, v)
	}
	return dic
}
