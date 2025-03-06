package main

import (
	"testing"
)

func TestFArgs(t *testing.T) {
	target := FuncType{[]FType{New_FType_FInt, New_FType_FInt, New_FType_FString, New_FType_FString}}
	got := fargs(target)
	if len(got) != 3 {
		t.Errorf("want [i, i, s], but got %v", got)
	}
	if New_FType_FInt != got[0] {
		t.Errorf("want FInt in 0, but got %v", got)
	}
	if New_FType_FString != got[2] {
		t.Errorf("want FString in 2, but got %v", got)
	}
}

func TestFFuncToGo(t *testing.T) {
	target := FuncType{[]FType{New_FType_FInt, New_FType_FInt, New_FType_FString, New_FType_FString}}
	got := funcTypeToGo(target, FTypeToGo)

	if got != "func (int,int,string) string" {
		t.Errorf("got %s", got)
	}
}

func TestToGoFunc(t *testing.T) {
	target := FuncType{[]FType{New_FType_FInt, New_FType_FInt, New_FType_FString, New_FType_FString}}
	got := FTypeToGo(New_FType_FFunc(target))

	if got != "func (int,int,string) string" {
		t.Errorf("got %s", got)
	}
}

func TestToGoSlice(t *testing.T) {
	target := SliceType{New_FType_FString}
	got := FTypeToGo(New_FType_FSlice(target))
	if got != "[]string" {
		t.Errorf("got %s", got)
	}
}

func TestUnionCaseStructName(t *testing.T) {
	got := unionCSName("IntOrString", "I")
	if got != "IntOrString_I" {
		t.Errorf("got %s", got)
	}
}
