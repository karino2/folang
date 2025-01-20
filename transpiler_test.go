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
				&FuncDef{"main", nil, &GoEval{"fmt.Println(\"Hello World\")"}},
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
				&FuncDef{"hello", []*Var{{"msg", FString}}, &GoEval{"fmt.Printf(\"Hello %s\\n\", msg)"}},
				&FuncDef{"main", nil,
					&FunCall{
						&Var{"hello", &FFunc{[]FType{FString, FUnit}}},
						[]Expr{&StringLiteral{"Hoge"}},
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
				&FuncDef{"hello", []*Var{{"msg", FString}}, &GoEval{"fmt.Printf(\"Hello %s\\n\", msg)"}},
				&FuncDef{"main", nil,
					&FunCall{
						// 型解決が動くか？
						&Var{"hello", &FUnresolved{}},
						[]Expr{&StringLiteral{"Hoge"}},
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
		&FuncDef{"hello", []*Var{{"msg", FString}}, &GoEval{"fmt.Printf(\"Hello %s\\n\", msg)"}},
		&FuncDef{"main", nil,
			funCall,
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
			recGen,
		},
	})
	tp := NewTranspiler()
	got := tp.Transpile(f)
	want := `type hoge struct {
  X string
  Y int
}

func ika() *hoge{
return &hoge{X: "abc", Y: 123}
}

`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUnionDef(t *testing.T) {
	unionDef := &UnionDef{"IntOrBool", []NameTypePair{{"I", FInt}, {"S", FString}}}
	got := unionDef.ToGo()

	want := `type IntOrBool interface {
  IntOrBool_Union()
}

func (*IntOrBool_I) IntOrBool_Union(){}
func (*IntOrBool_S) IntOrBool_Union(){}

type IntOrBool_I struct {
  Value int
}

func New_IntOrBool_I(v int) IntOrBool { return &IntOrBool_I{v} }

type IntOrBool_S struct {
  Value string
}

func New_IntOrBool_S(v string) IntOrBool { return &IntOrBool_S{v} }

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
		&UnionDef{"IntOrBool", []NameTypePair{{"I", FInt}, {"S", FString}}},
		&FuncDef{"test_func", nil, funCall},
	})

	tp := NewTranspiler()
	tp.resolveAndRegisterType(f.Stmts)

	fn := funCall.Func
	if fn.Name != "New_IntOrBool_I" {
		t.Errorf("want New_IntOrBool_I, got %v", fn)
	}
	ft, ok := fn.Type.(*FFunc)
	if !ok {
		t.Errorf("want FFunc type, got %v", fn)
	}
	got := ft.String()
	if got != "int -> IntOrBool" {
		t.Errorf("want 'int -> IntOrBool', got %s", got)
	}

}
