package main

import (
	"strings"
	"testing"
)

func newIdentifierToken(begin int, len int, sval string) *Token {
	res := NewToken(IDENTIFIER, begin, len)
	res.stringVal = sval
	return res
}

func TestTokenizerAnalyzeCur(t *testing.T) {
	var tests = []struct {
		input string
		want  *Token
	}{
		{"", NewToken(EOF, 0, 0)},
		{" ", NewToken(SPACE, 0, 1)},
		{" a", NewToken(SPACE, 0, 1)},
		{"  ", NewToken(SPACE, 0, 2)},
		{"a", newIdentifierToken(0, 1, "a")},
		{"=", &Token{ttype: EQ, begin: 0, len: 1, stringVal: "="}},
		{"let", &Token{ttype: LET, begin: 0, len: 3, stringVal: "let"}},
		{"\n", &Token{ttype: EOL, begin: 0, len: 1, stringVal: "\n"}},
		{"package", &Token{ttype: PACKAGE, begin: 0, len: 7, stringVal: "package"}},
		{"import", &Token{ttype: IMPORT, begin: 0, len: 6, stringVal: "import"}},
		{"(", &Token{ttype: LPAREN, begin: 0, len: 1, stringVal: "("}},
		{")", &Token{ttype: RPAREN, begin: 0, len: 1, stringVal: ")"}},
		{`"hoge"`, &Token{ttype: STRING, begin: 0, len: 6, stringVal: "hoge"}},
		{`"hoge(\"ika\")"`, &Token{ttype: STRING, begin: 0, len: 15, stringVal: `hoge("ika")`}},
	}

	for _, test := range tests {
		tkzr := NewTokenizer("", []byte(test.input))
		tkzr.analyzeCur()

		got := tkzr.Current()
		if *got != *test.want {
			t.Errorf("got %v, want %v", got, test.want)
		}
	}
}

func TestTokenizerAnalyzeBaseScenario(t *testing.T) {
	src :=
		`package main
import "fmt"
`
	tkzr := NewTokenizer("", []byte(src))
	tkzr.Setup()

	cur := tkzr.Current()
	if cur.ttype != PACKAGE {
		t.Errorf("expect package, got %v", cur)
		return
	}

	tkzr.GotoNext()
	cur = tkzr.Current()
	if cur.ttype != SPACE {
		t.Errorf("expect space, got %v", cur)
		return
	}
	tkzr.GotoNext()
	cur = tkzr.Current()

	if cur.ttype != IDENTIFIER || cur.stringVal != "main" {
		t.Errorf("expect main identifier, got %v", cur)
		return
	}

	tkzr.GotoNext()
	cur = tkzr.Current()
	if cur.ttype != EOL {
		t.Errorf("expect eol, got %v", cur)
		return
	}

	tkzr.GotoNext()
	cur = tkzr.Current()
	if cur.ttype != IMPORT {
		t.Errorf("expect import, got %v", cur)
		return
	}

	tkzr.GotoNext()
	cur = tkzr.Current()
	if cur.ttype != SPACE {
		t.Errorf("expect space, got %v", cur)
		return
	}

	tkzr.GotoNext()
	cur = tkzr.Current()
	if cur.ttype != STRING || cur.stringVal != "fmt" {
		t.Errorf("expect string literal 'fmt', got %v", cur)
		return
	}

	tkzr.GotoNext()
	cur = tkzr.Current()
	if cur.ttype != EOL {
		t.Errorf("expect eol, got %v", cur)
		return
	}

	tkzr.GotoNext()
	cur = tkzr.Current()
	if cur.ttype != EOF {
		t.Errorf("expect eof, got %v", cur)
		return
	}
}

func TestParserEasiest(t *testing.T) {
	src :=
		`package main
import "fmt"
`
	parser := NewParser()
	res := parser.Parse("", []byte(src))
	if len(res) != 2 {
		t.Errorf("expect 2 stmt, got %d, %v", len(res), res)
	}

	wants := []string{`package main`, `import "fmt"`}
	for i, want := range wants {
		got := res[i].ToGo()
		if got != want {
			t.Errorf("want '%s', got '%s'", want, got)
		}
	}
}

func transpile(src string) string {
	ResetUniqueTmpCounter()
	defer ResetUniqueTmpCounter()

	parser := NewParser()
	res := parser.Parse("", []byte(src))
	tp := NewTranspiler()
	f := NewFile(res)
	return tp.Transpile(f)

}

func TestParserAndTranspile(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{
			`package main
import "fmt"

let main () =
    GoEval "fmt.Println(\"Hello World\")"
`,
			`package main

import "fmt"

func main() {
fmt.Println("Hello World")
}

`,
		},
		{
			`package main
import "fmt"

let hello (msg:string) = 
    GoEval "fmt.Printf(\"Hello %s\\n\", msg)"

let main () =
    hello "World"
`, `package main

import "fmt"

func hello(msg string) {
fmt.Printf("Hello %s\n", msg)
}

func main() {
hello("World")
}

`,
		},
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
			`let ika() =
  GoEval<int> "3+4"
`,
			`func ika() int{
return 3+4
}

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
			`let a = "abc"
`,
			`a := "abc"

`,
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

		if got != test.want {
			t.Errorf("got [%s], want [%s]", got, test.want)
		}
	}
}

func TestParserRecordDef(t *testing.T) {
	src := `type Hoge = {X: string; Y: string}
`

	parser := NewParser()
	res := parser.Parse("", []byte(src))
	if len(res) != 1 {
		t.Errorf("expect 1 stmt, got %d, %v", len(res), res)
	}
	rd, ok := res[0].(*RecordDef)
	if !ok {
		t.Error("expect RecordDef but not")
	}
	if rd.Name != "Hoge" ||
		len(rd.Fields) != 2 ||
		rd.Fields[0].Name != "X" ||
		rd.Fields[0].Type != FString ||
		rd.Fields[1].Name != "Y" ||
		rd.Fields[1].Type != FString {
		t.Errorf("unexpected rd: %v", rd)
	}
}

func TestParserRecordExpression(t *testing.T) {
	src := `{X="hoge"; Y="ika"}
`

	parser := NewParser()
	parser.Setup("", []byte(src))
	res := parser.parseExpr()
	if res == nil {
		t.Errorf("expect record gen expression")
	}
	rg, ok := res.(*RecordGen)
	if !ok {
		t.Error("expect RecordGen but not")
	}
	if len(rg.fieldNames) != 2 ||
		rg.fieldNames[0] != "X" ||
		rg.fieldNames[1] != "Y" {
		t.Errorf("unexpected rg: %v", rg)
	}
}

func TestParserMatchExpression(t *testing.T) {
	src :=
		`   match x with
    | I ival -> "i match"
    | S sval -> "s match"
`

	unionDef := &UnionDef{"IntOrString", []NameTypePair{{"I", FInt}, {"S", FString}}}
	varX := &Var{"x", unionDef.UnionFType()}

	parser := NewParser()
	parser.Setup("", []byte(src))
	parser.scope.DefineVar("x", varX)
	parser.skipSpace()
	parser.pushOffside()
	res := parser.parseExpr()

	me, ok := res.(*MatchExpr)
	if !ok {
		t.Error("parse result is not MatchExpr.")
	}
	v, ok := me.target.(*Var)
	if !ok {
		t.Error("target is not variable.")
	}
	if v.Name != "x" {
		t.Errorf("wrong target: %v", v)
	}
	if len(me.rules) != 2 {
		t.Errorf("want 2 rules, but %v", me.rules)
	}
}

/*
The results of this category are too complex to assert whole.
Check only part of it.
*/
func TestParserContainTest(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{
			`package main

type IorS =
 | IT of int
 | ST of string

 let ika () =
  match IT 3 with
  | IT ival -> "i match"
  | ST sval -> "s match"
`,
			"case IorS_IT:",
		},
		// no var case.
		{`package main

type IorS =
  | IT of int
  | ST of string

  let ika () =
  match IT 3 with
  | IT _ -> "i match"
  | ST sval -> "s match"
`,
			`case IorS_IT:
return "i match"`,
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
			`
type NameTypePair = {Name: string; Type: string}

type RecordType = {name: string; fiedls: []NameTypePair}
`,
			"fiedls []NameTypePair",
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
	}

	for _, test := range tests {
		got := transpile(test.input)

		// t.Error(got)
		if !strings.Contains(got, test.want) {
			t.Errorf("want to contain '%s', but got '%s'", test.want, got)
		}
	}
}

func TestParserAddhook(t *testing.T) {
	src := `type IorS =
  | IT of int
  | ST of string

  let ika () =
  match IT 3 with
  | IT _ -> true
  | _ -> false
`

	got := transpile(src)
	// t.Error(got)

	want := "ika() bool{"
	if !strings.Contains(got, want) {
		t.Errorf("want to contains(%s), but got %s", want, got)
	}

}
