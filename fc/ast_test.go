package main

import (
	"testing"
)

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
		got := ExprToGo(test.input)
		if got != test.want {
			t.Errorf("want %s, got %s.", test.want, got)
		}
	}
}
