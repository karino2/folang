package main

import (
	"fmt"
	"testing"
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
			New_Expr_StringLiteral("abc"),
			`"abc"`,
		},
		{
			New_Expr_FieldAccess(FieldAccess{"rec", RecordType{"MyRec", []NameTypePair{{"field1", New_FType_FInt}}}, "field1"}),
			"rec.field1",
		},
		{
			New_Expr_RecordGen(RecordGen{[]string{"hoge", "ika"}, []Expr{New_Expr_StringLiteral("sval"), New_Expr_IntImm(123)}, RecordType{"MyRec", []NameTypePair{{"hoge", New_FType_FString}, {"ika", New_FType_FInt}}}}),
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
	target := New_Expr_Var(Var{"udata", New_FType_FUnion(unionType)})
	matchExpr := MatchExpr{
		target,
		[]MatchRule{
			{
				MatchPattern{"I", "ival"},
				newBlock(New_Expr_StringLiteral("I match")),
			},
			{
				MatchPattern{"S", "sval"},
				newBlock(New_Expr_StringLiteral("S match")),
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
