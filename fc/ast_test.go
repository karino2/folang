package main

import (
	"fmt"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func MyExprToGo(expr Expr) string {
	return ExprToGo(func(s Stmt) string { return fmt.Sprintf("stmt: %v", s) }, expr)
}

func TestRecordDefToGo(t *testing.T) {
	rd := RecordDef{"hoge", []string{}, []NameTypePair{{"X", New_FType_FString}, {"Y", New_FType_FString}}}
	got := rdfToGo(rd)
	want := `type hoge struct {
  X string
  Y string
}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUnionDefToGo(t *testing.T) {
	ud := NewUnionDef("IntOrString", []NameTypePair{{"I", New_FType_FInt}, {"S", New_FType_FString}})
	got := udfToGo(ud)
	want := `type IntOrString interface {
  IntOrString_Union()
}

func (IntOrString_I) IntOrString_Union(){}
func (IntOrString_S) IntOrString_Union(){}

type IntOrString_I struct {
  Value int
}

func New_IntOrString_I(v int) IntOrString { return IntOrString_I{v} }

type IntOrString_S struct {
  Value string
}

func New_IntOrString_S(v string) IntOrString { return IntOrString_S{v} }

`
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(got, want, false)

	if len(diffs) != 1 {
		t.Errorf("diff found: %s", dmp.DiffPrettyText(diffs))
	}
}
