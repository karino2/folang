package main

import (
	"strings"
	"testing"

	"github.com/karino2/folang/pkg/frt"
)

func TestParsePackage(t *testing.T) {
	src := `package main
`
	ps := initParse(src)
	gotPair := parsePackage(ps)
	got := gotPair.E1

	switch tgot := got.(type) {
	case Stmt_Package:
		if tgot.Value != "main" {
			t.Errorf("expect main, got %s", tgot.Value)
		}
	default:
		t.Errorf("Unexpected stmt. Expect Package, got %T", got)
	}
}

func TestParseParams(t *testing.T) {
	src := `(a:int) (b:int) =
`
	ps := initParse(src)
	gotPair := parseParams(ps)
	got := gotPair.E1
	if len(got) != 2 {
		t.Errorf("want 2 param, got %d: %v", len(got), got)
	}
	// t.Errorf("%d, %T, %T", len(got), got[0], got[1])

	ps2 := gotPair.E0
	tt := psCurrentTT(ps2)
	if tt != New_TokenType_EQ {
		t.Errorf("want EQ, got %T", tt)
	}
}

func TestParseLetFuncDef(t *testing.T) {
	src := `let hoge () =
  123
`
	ps := initParse(src)
	gotPair := parseLetFuncDef(ps)
	_, stmt := frt.Destr(gotPair)

	if lfdS, ok := stmt.(Stmt_LetFuncDef); ok {
		lfd := lfdS.Value
		if lfd.name != "hoge" {
			t.Errorf("name is not hoge, %v", lfd)
		}
	} else {
		t.Errorf("not stmt not lfd, %T", stmt)
	}

	got := StmtToGo(stmt)
	want := `func hoge() int{
return 123
}`

	if got != want {
		t.Errorf("want %s, got %s", want, got)
	}

}

func TestParseTwoFunc(t *testing.T) {
	src := `package main
import "fmt"

let hello (msg:string) = 
    GoEval "fmt.Printf(\"Hello %s\\n\", msg)"

let main () =
   hello "World"

`
	ps := initParse(src)
	gotPair := parseStmts(ps)
	_, stmts := frt.Destr(gotPair)
	got := StmtsToGo(stmts)

	want :=
		`package main

import "fmt"

func hello(msg string) {
fmt.Printf("Hello %s\n", msg)
}

func main() {
hello("World")
}
`
	if got != want {
		t.Errorf("want %s, got %s.", want, got)
	}
}

func transpile(src string) string {
	ps := initParse(src)
	gotPair := parseStmts(ps)
	_, stmts := frt.Destr(gotPair)
	return StmtsToGo(stmts)
}

func TestTranspileContain(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{
			`type hoge = {X: string; Y: int}

let ika () =
    {X="abc"; Y=123}

`,
			`type hoge struct {
  X string
  Y int
}

func ika() hoge{
return hoge{X: "abc", Y: 123}
}
`,
		},
		{
			`type IntOrString =
  | I of int
  | S of string
`,
			`type IntOrString interface {
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
`,
		},
		{
			`type AorB =
  | A
  | B
`,
			`type AorB interface {
  AorB_Union()
}

func (AorB_A) AorB_Union(){}
func (AorB_B) AorB_Union(){}

type AorB_A struct {
}

var New_AorB_A AorB = AorB_A{}

type AorB_B struct {
}

var New_AorB_B AorB = AorB_B{}
`,
		},
	}
	for _, test := range tests {
		got := transpile(test.input)

		if !strings.Contains(got, test.want) {
			t.Errorf("want to contain '%s', but got '%s'", test.want, got)
		}
	}

}
func TestParseAddhook(t *testing.T) {
	src := `package main

type hoge = {X: string; Y: string}

let ika () =
    {X="abc"; Y="def"}

`

	got := transpile(src)
	// t.Error(got)
	if got == "dummy" {
		t.Error(got)
	}
}
