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
				&RecordDef{"hoge", []RecordField{{"X", FString}, {"Y", FString}}},
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
		&RecordDef{"hoge", []RecordField{{"X", FString}, {"Y", FString}}},
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
		&RecordDef{"hoge", []RecordField{{"X", FString}, {"Y", FInt}}},
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
