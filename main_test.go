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
	} {
		var p Program
		for _, s := range test.stmts {
			p.AddStmt(s)
		}
		res := p.ToGo()
		if test.want != res {
			t.Errorf("got %s, want %s", res, test.want)
		}
	}
}

func TestFuncDefCompile(t *testing.T) {
	var p Program
	p.AddStmt(&Import{"fmt"})
	p.AddStmt(&FuncDef{"hello", []Var{{"msg", FString}}, &GoEval{"fmt.Println(\"Hello %s\", msg)"}})
	p.AddStmt(&FuncDef{"main", nil, &GoEval{"fmt.Println(\"Hello World\")"}})
	fmt.Println(p.ToGo())
}
