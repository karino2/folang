package main

import (
	"fmt"

	"github.com/karino2/folang/pkg/frt"
)

func dictPut[T any](dict map[string]T, key string, v T) {
	dict[key] = v
}

func dictKeyValues[K comparable, V any](dict map[K]V) []frt.Tuple2[K, V] {
	var res []frt.Tuple2[K, V]
	for k, v := range dict {
		res = append(res, frt.NewTuple2(k, v))
	}
	return res
}

// currently, map is NYI  and generic exxt type is also NYI.
// wrap to standard type for each.
type funcTypeDict = map[string]FuncType
type extTypeDict = map[string]string

func newFTD() funcTypeDict {
	return make(map[string]FuncType)
}

func newETD() extTypeDict {
	return make(map[string]string)
}

func ftdPut(dic funcTypeDict, key string, v FuncType) {
	dictPut(dic, key, v)
}

func ftdKVs(dic funcTypeDict) []frt.Tuple2[string, FuncType] {
	return dictKeyValues(dic)
}

func etdPut(dic extTypeDict, key string, v string) {
	dictPut(dic, key, v)
}

func etdKVs(dic extTypeDict) []frt.Tuple2[string, string] {
	return dictKeyValues(dic)
}

var uniqueId = 0

func uniqueTmpVarName() string {
	uniqueId++
	return fmt.Sprintf("_v%d", uniqueId)
}

/*
func uniqueTmpTypeParamName() string {
	uniqueId++
	return fmt.Sprintf("_T%d", uniqueId)
}
*/

func resetUniqueTmpCounter() {
	uniqueId = 0
}
