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
	ud := NewUnionDef("IntOrString", []string{}, []NameTypePair{{"I", New_FType_FInt}, {"S", New_FType_FString}})
	got := udfToGo(ud)
	want := `type IntOrString interface {
  IntOrString_Union()
}

func (IntOrString_I) IntOrString_Union(){}
func (IntOrString_S) IntOrString_Union(){}

func (v IntOrString_I) String() string { return frt.Sprintf1("(I: %v)", v.Value) }
func (v IntOrString_S) String() string { return frt.Sprintf1("(S: %v)", v.Value) }

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

func TestUnionDefGenToGo(t *testing.T) {
	ud := NewUnionDef("Result", []string{"T0"}, []NameTypePair{{"Success", New_FType_FTypeVar(TypeVar{"T0"})}, {"Failure", New_FType_FUnit}})
	got := udfToGo(ud)
	want := `type Result[T0 any] interface {
  Result_Union()
}

func (Result_Success[T0]) Result_Union(){}
func (Result_Failure[T0]) Result_Union(){}

func (v Result_Success[T0]) String() string { return frt.Sprintf1("(Success: %v)", v.Value) }
func (v Result_Failure[T0]) String() string { return "(Failure)" }

type Result_Success[T0 any] struct {
  Value T0
}

func New_Result_Success[T0 any](v T0) Result[T0] { return Result_Success[T0]{v} }

type Result_Failure[T0 any] struct {
}

func New_Result_Failure[T0 any]() Result[T0] { return Result_Failure[T0]{} }

`
	// t.Error(got)
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(got, want, false)

	if len(diffs) != 1 {
		t.Errorf("diff found: %s", dmp.DiffPrettyText(diffs))
	}
}
