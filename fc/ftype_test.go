package main

import (
	"fmt"
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

func toGoSimple(ftype FType) string { return fmt.Sprintf("%v", ftype) }

func TestFFuncToGo(t *testing.T) {
	target := FuncType{[]FType{New_FType_FInt, New_FType_FInt, New_FType_FString, New_FType_FString}}
	got := FFuncToGo(target, toGoSimple)

	if got != "func ({},{},{}) {}" {
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

func TestRecordGetField(t *testing.T) {
	rec := RecordType{"MyRec", []NameTypePair{{"hoge", New_FType_FString}, {"ika", New_FType_FInt}}}
	hpair := frGetField(rec, "hoge")
	ipair := frGetField(rec, "ika")
	if hpair.name != "hoge" || hpair.ftype != New_FType_FString {
		t.Errorf("wrong hoge field: %v", hpair)
	}

	if ipair.name != "ika" || ipair.ftype != New_FType_FInt {
		t.Errorf("wrong ika field: %v", hpair)
	}
}

func TestRecordMatch(t *testing.T) {
	rec := RecordType{"MyRec", []NameTypePair{{"hoge", New_FType_FString}, {"ika", New_FType_FInt}}}
	var tests = []struct {
		input []string
		want  bool
	}{
		{
			[]string{"hoge"},
			false,
		},
		{
			[]string{"hoge", "ika"},
			true,
		},
		// different order.
		{
			[]string{"ika", "hoge"},
			true,
		},
	}

	for _, test := range tests {
		got := frMatch(rec, test.input)
		if got != test.want {
			t.Errorf("got %t, want %t with inputs %v", got, test.want, test.input)
		}
	}

}

func TestUnionCaseStructName(t *testing.T) {
	got := unionCaseStructName("IntOrString", "I")
	if got != "IntOrString_I" {
		t.Errorf("got %s", got)
	}
}
