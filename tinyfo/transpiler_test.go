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
				&FuncDef{"main", nil, &Block{nil, NewGoEval("fmt.Println(\"Hello World\")")}},
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
				&FuncDef{"hello", []*Var{{"msg", FString}}, &Block{nil, NewGoEval("fmt.Printf(\"Hello %s\\n\", msg)")}},
				&FuncDef{"main", nil,
					&Block{nil,
						&FunCall{
							&Var{"hello", &FFunc{[]FType{FString, FUnit}}},
							[]Expr{&StringLiteral{"Hoge"}},
						},
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
				&Package{"main"},
				&Import{"fmt"},
				&FuncDef{"hello", []*Var{{"msg", FString}}, &Block{nil, NewGoEval("fmt.Printf(\"Hello %s\\n\", msg)")}},
				&FuncDef{"main", nil,
					&Block{nil,
						&FunCall{
							// 型解決が動くか？
							&Var{"hello", &FUnresolved{}},
							[]Expr{&StringLiteral{"Hoge"}},
						},
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
	} {
		tp := NewTranspiler()
		f := NewFile(test.stmts)
		res := tp.Transpile(f)
		if test.want != res {
			t.Errorf("got %s, want %s", res, test.want)
		}
	}
}
func TestResolver(t *testing.T) {
	funCall := &FunCall{
		&Var{"hello", &FUnresolved{}},
		[]Expr{&StringLiteral{"Hoge"}},
	}

	tp := NewTranspiler()
	f := NewFile([]Stmt{
		&Package{"main"},
		&Import{"fmt"},
		&FuncDef{"hello", []*Var{{"msg", FString}}, &Block{nil, NewGoEval("fmt.Printf(\"Hello %s\\n\", msg)")}},
		&FuncDef{"main", nil,
			&Block{nil, funCall},
		},
	})

	tp.resolveAndRegisterType(f.Stmts)

	got := funCall.Func.Type.ToGo()
	want := "func (string)"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestRecordDefLookup(t *testing.T) {
	f := NewFile([]Stmt{
		&RecordDef{"hoge", []NameTypePair{{"X", FString}, {"Y", FString}}},
	})
	tp := NewTranspiler()
	tp.resolveAndRegisterType(f.Stmts)
	got := tp.Resolver.LookupByTypeName("hoge")
	_, ok := got.(*FRecord)
	if !ok {
		t.Errorf("cannot find FRecord by name")
	}

	got2 := tp.Resolver.LookupRecord([]string{"X", "Y"})
	if got2 == nil {
		t.Errorf("cannot find FRecord by field")
	}

	got3 := tp.Resolver.LookupRecord([]string{"X", "DEF"})
	if got3 != nil {
		t.Errorf("Wrongly matched with different name")
	}

	got4 := tp.Resolver.LookupRecord([]string{"X", "Y", "Z"})
	if got4 != nil {
		t.Errorf("Wrongly matched with extra field")
	}

	got5 := tp.Resolver.LookupRecord([]string{"X"})
	if got5 != nil {
		t.Errorf("Wrongly matched with few fields")
	}
}

func TestRecordGen(t *testing.T) {
	recGen := NewRecordGen(
		[]string{"X", "Y"},
		[]Expr{&StringLiteral{"abc"}, &IntImm{123}},
	)

	f := NewFile([]Stmt{
		&RecordDef{"hoge", []NameTypePair{{"X", FString}, {"Y", FInt}}},
		&FuncDef{"ika", nil,
			&Block{
				nil,
				recGen,
			},
		},
	})
	tp := NewTranspiler()
	got := tp.Transpile(f)
	want := `type hoge struct {
  X string
  Y int
}

func ika() hoge{
return hoge{X: "abc", Y: 123}
}

`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
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

func TestUnionDefConstructorHandling(t *testing.T) {
	funCall := &FunCall{
		&Var{"I", &FUnresolved{}},
		[]Expr{&IntImm{123}},
	}

	f := NewFile([]Stmt{
		&UnionDef{"IntOrString", []NameTypePair{{"I", FInt}, {"S", FString}}},
		&FuncDef{"test_func", nil, &Block{nil, funCall}},
	})

	tp := NewTranspiler()
	tp.resolveAndRegisterType(f.Stmts)

	fn := funCall.Func
	if fn.Name != "New_IntOrString_I" {
		t.Errorf("want New_IntOrString_I, got %v", fn)
	}
	ft, ok := fn.Type.(*FFunc)
	if !ok {
		t.Errorf("want FFunc type, got %v", fn)
	}
	got := ft.String()
	if got != "int -> IntOrString" {
		t.Errorf("want 'int -> IntOrString', got %s", got)
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
