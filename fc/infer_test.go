package main

import (
	"testing"
)

func newTV(name string) TypeVar {
	return TypeVar{name}
}

func newTVF(name string) FType {
	return New_FType_FTypeVar(newTV(name))
}

func newFInt() FType {
	return New_FType_FInt
}

func TestBuildResolver(t *testing.T) {
	rels := []UniRel{
		{"_T3", newTVF("_T1")},
		{"_T2", newTVF("_T1")},
		{"_T2", newFInt()},
	}
	resv := buildResolver(rels)
	got := rsResolveType(resv, "_T3")
	if _, ok := got.(FType_FInt); !ok {
		t.Errorf("want int, got %T", got)
	}
}
