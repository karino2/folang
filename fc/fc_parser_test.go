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

/*
Resolve mutual recursive in golang layer (NYI for and letfunc def).
*/
func parseLetFacade(ps ParseState) frt.Tuple2[ParseState, LetVarDef] {
	return parseLetVarDef(parseExprFacade, ps)
}

func parseBlockFacade(ps ParseState) frt.Tuple2[ParseState, Block] {
	return parseBlock(parseLetFacade, ps)
}

func parseExprFacade(ps ParseState) frt.Tuple2[ParseState, Expr] {
	return parseTerm(parseBlockFacade, ps)
}

func parseAll(ps ParseState) (ParseState, []Stmt) {
	res := parseStmts(parseExprFacade, ps)
	return frt.Destr(res)
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
	_, stmts := parseAll(ps)
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
	_, stmts := parseAll(ps)
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
		{
			`package main

type IntOrString =
  | I of int
  | S of string

let ika () =
    I 123

`,
			"New_IntOrString_I(123)",
		},
		// no content case
		{
			`package main

type AorB =
 | A
 | B

let ika (ab:AorB) =
  match ab with
  | A -> "a match"
  | B -> "b match"
`, "switch (ab).(type)",
		},
		{
			`package main

type AorB =
 | A
 | B

let main () =
  match A with
  | A -> GoEval "fmt.Println(\"A match\")"
  | B -> GoEval "fmt.Println(\"B match\")"
`, "switch (New_AorB_A).(type)",
		},
		{
			`let ika() =
  GoEval<int> "3+4"
`,
			"func ika() int{",
		},
		{
			`type IorS =
  | IT of int
  | ST of string

  let ika () =
  match IT 3 with
  | IT _ -> "i match"
  | _ -> "default"
`,
			"default:\nreturn \"default\"",
		},
		// bool test
		{
			`type IorS =
  | IT of int
  | ST of string

  let ika () =
  match IT 3 with
  | IT _ -> true
  | _ -> false
`,
			"func ika() bool{",
		},
		{
			`let hoge () =
  let a = "abc"
  a
`,
			`func hoge() string{
a := "abc"
return a
}
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

func TestTranspileContainsMulti(t *testing.T) {
	var tests = []struct {
		input string
		wants []string
	}{
		{
			`package main

type IntOrString =
  | I of int
  | S of string

let ika () =
  match I 123 with
  | I i -> i
  | S _ -> 456

`,
			[]string{"i := _v1.Value", "switch _v1 :="},
		},
	}

	for _, test := range tests {
		got := transpile(test.input)

		for _, want := range test.wants {
			if !strings.Contains(got, want) {
				t.Errorf("want to contain '%s', but got '%s'", want, got)
			}
		}
	}
}
func TestParseAddhook(t *testing.T) {
	src := `package main

let ika() =
  GoEval<int> "3+4"

`

	got := transpile(src)
	// t.Error(got)
	if got == "dummy" {
		t.Error(got)
	}
}
