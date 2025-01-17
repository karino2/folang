package main

import (
	"fmt"
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
				&FuncDef{"hello", []Var{{"msg", FString}}, &GoEval{"fmt.Printf(\"Hello %s\\n\", msg)"}},
				&FuncDef{"main", nil,
					&FunCall{
						Var{"hello", &FFunc{[]FType{FString, FUnit}}},
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
				&FuncDef{"hello", []Var{{"msg", FString}}, &GoEval{"fmt.Printf(\"Hello %s\\n\", msg)"}},
				&FuncDef{"main", nil,
					&FunCall{
						// 型解決が動くか？
						Var{"hello", &FUnresolved{}},
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
	} {
		p := NewProgram(test.stmts)
		res := p.Compile()
		if test.want != res {
			t.Errorf("got %s, want %s", res, test.want)
		}
	}
}
func TestResolver(t *testing.T) {
	resolver := NewResolver()
	resolver.Register("hello", &FFunc{[]FType{FString, FUnit}})

	funCall := &FunCall{
		Var{"hello", &FUnresolved{}},
		[]Expr{&StringLiteral{"Hoge"}},
	}

	p := NewProgram([]Stmt{
		&Package{"main"},
		&Import{"fmt"},
		&FuncDef{"hello", []Var{{"msg", FString}}, &GoEval{"fmt.Printf(\"Hello %s\\n\", msg)"}},
		&FuncDef{"main", nil,
			funCall,
		},
	})

	p.ResolveAndRegisterType()

	got := funCall.Func.Type.ToGo()
	want := "func (string)"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAddhookCompile(t *testing.T) {
	p := NewProgram([]Stmt{
		&Import{"fmt"},
		&FuncDef{"hello", []Var{{"msg", FString}}, &GoEval{"fmt.Printf(\"Hello %s\\n\", msg)"}},
		&FuncDef{"main", nil, &GoEval{"fmt.Println(\"Hello World\")"}},
	})
	fmt.Println(p.Compile())
}
