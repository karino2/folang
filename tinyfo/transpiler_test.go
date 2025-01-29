package main

import (
	"testing"
)

func TestCompile(t *testing.T) {
	for _, test := range []struct {
		stmts []Stmt
		want  string
	}{
		{
			[]Stmt{
				&Import{"fmt"},
				&FuncDef{"main", nil, &Block{nil, NewGoEval("fmt.Println(\"Hello World\")"), nil}},
			},
			`import "fmt"

func main() {
fmt.Println("Hello World")
}

`,
		},
		{
			[]Stmt{
				&Package{"main"},
				&Import{"fmt"},
				&FuncDef{"hello", []*Var{{"msg", FString}}, &Block{nil, NewGoEval("fmt.Printf(\"Hello %s\\n\", msg)"), nil}},
				&FuncDef{"main", nil,
					&Block{nil,
						&FunCall{
							&Var{"hello", &FFunc{[]FType{FString, FUnit}}},
							[]Expr{&StringLiteral{"Hoge"}},
						},
						nil,
					},
				},
			},
			`package main

import "fmt"

func hello(msg string) {
fmt.Printf("Hello %s\n", msg)
}

func main() {
hello("Hoge")
}

`,
		},
		{
			[]Stmt{
				&RecordDef{"hoge", []NameTypePair{{"X", FString}, {"Y", FString}}},
			},
			`type hoge struct {
  X string
  Y string
}

`,
		},
		{
			[]Stmt{
				&LetVarDef{"hoge", &StringLiteral{"ABC"}},
			},
			`hoge := "ABC"

`,
		},
	} {
		tp := NewTranspiler()
		f := NewFile(test.stmts)
		res := tp.Transpile(f)
		if test.want != res {
			t.Errorf("got %s, want %s", res, test.want)
		}
	}
}

func TestUnionDef(t *testing.T) {
	unionDef := &UnionDef{"IntOrString", []NameTypePair{{"I", FInt}, {"S", FString}}}
	got := unionDef.ToGo()

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
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestPatternMatchUnion(t *testing.T) {
	ResetUniqueTmpCounter()
	defer ResetUniqueTmpCounter()
	unionFT := &FUnion{"IntOrString", []NameTypePair{{"I", FInt}, {"S", FString}}}
	target := &Var{"udata", unionFT}
	matchExpr := &MatchExpr{
		target,
		[]*MatchRule{
			{
				&MatchPattern{
					"I",
					"ival",
				},
				&Block{
					nil,
					&StringLiteral{"I match."},
					nil,
				},
			},
			{
				&MatchPattern{
					"S",
					"sval",
				},
				&Block{
					nil,
					&StringLiteral{"s match."},
					nil,
				},
			},
		},
	}
	got1 := matchExpr.ToGoReturn()
	want1 :=
		`switch _v1 := (udata).(type){
case IntOrString_I:
ival := _v1.Value
return "I match."
case IntOrString_S:
sval := _v1.Value
return "s match."
default:
panic("Union pattern fail. Never reached here.")
}`
	if got1 != want1 {
		t.Errorf("want: %s, got: %s", want1, got1)
	}
	want2 :=
		`(func () string {
switch _v2 := (udata).(type){
case IntOrString_I:
ival := _v2.Value
return "I match."
case IntOrString_S:
sval := _v2.Value
return "s match."
default:
panic("Union pattern fail. Never reached here.")
}})()`
	got2 := matchExpr.ToGo()
	if got2 != want2 {
		t.Errorf("want: %s, got: %s", want2, got2)
	}
}

func TestPatternMatchUnionUnusedVar(t *testing.T) {
	ResetUniqueTmpCounter()
	defer ResetUniqueTmpCounter()
	unionFT := &FUnion{"IntOrString", []NameTypePair{{"I", FInt}, {"S", FString}}}
	target := &Var{"udata", unionFT}
	matchExpr := &MatchExpr{
		target,
		[]*MatchRule{
			{
				&MatchPattern{
					"I",
					"_",
				},
				&Block{
					nil,
					&StringLiteral{"I match."},
					nil,
				},
			},
			{
				&MatchPattern{
					"S",
					"_",
				},
				&Block{
					nil,
					&StringLiteral{"s match."},
					nil,
				},
			},
		},
	}
	got1 := matchExpr.ToGoReturn()
	want1 :=
		`switch (udata).(type){
case IntOrString_I:
return "I match."
case IntOrString_S:
return "s match."
default:
panic("Union pattern fail. Never reached here.")
}`
	if got1 != want1 {
		t.Errorf("want: %s, got: %s", want1, got1)
	}
}
