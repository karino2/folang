package main

import (
	"fmt"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func MyExprToGo(expr Expr) string {
	return ExprToGo(func(s Stmt) string { return fmt.Sprintf("stmt: %v", s) }, expr)
}

func TestExprToGo(t *testing.T) {
	var tests = []struct {
		input Expr
		want  string
	}{
		{
			New_Expr_EStringLiteral("abc"),
			`"abc"`,
		},
		{
			New_Expr_EFieldAccess(FieldAccess{"rec", RecordType{"MyRec", []NameTypePair{{"field1", New_FType_FInt}}}, "field1"}),
			"rec.field1",
		},
		{
			New_Expr_ERecordGen(
				RecordGen{[]NEPair{{"hoge", New_Expr_EStringLiteral("sval")},
					{"ika", New_Expr_EIntImm(123)}},
					RecordType{"MyRec", []NameTypePair{{"hoge", New_FType_FString}, {"ika", New_FType_FInt}}}}),
			"MyRec{hoge: \"sval\", ika: 123}",
		},
	}

	for _, test := range tests {
		got := MyExprToGo(test.input)
		if got != test.want {
			t.Errorf("want %s, got %s.", test.want, got)
		}
	}
}

func newBlock(expr Expr) Block {
	return Block{[]Stmt{}, expr}
}

func TestMatchExprToGo(t *testing.T) {
	resetUniqueTmpCounter()
	defer resetUniqueTmpCounter()
	unionType := UnionType{"IntOrString", []NameTypePair{{"I", New_FType_FInt}, {"S", New_FType_FString}}}
	target := New_Expr_EVar(Var{"udata", New_FType_FUnion(unionType)})
	matchExpr := MatchExpr{
		target,
		[]MatchRule{
			{
				MatchPattern{"I", "ival"},
				newBlock(New_Expr_EStringLiteral("I match")),
			},
			{
				MatchPattern{"S", "sval"},
				newBlock(New_Expr_EStringLiteral("S match")),
			},
		},
	}

	want :=
		`(func () string {
switch _v1 := (udata).(type){
case IntOrString_I:
ival := _v1.Value
return "I match"
case IntOrString_S:
sval := _v1.Value
return "S match"
default:
panic("Union pattern fail. Never reached here.")
}})()`

	got := MyExprToGo(meToExpr(matchExpr))
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestRecordDefToGo(t *testing.T) {
	rd := RecordDef{"hoge", []NameTypePair{{"X", New_FType_FString}, {"Y", New_FType_FString}}}
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
	ud := UnionDef{"IntOrString", []NameTypePair{{"I", New_FType_FInt}, {"S", New_FType_FString}}}
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
