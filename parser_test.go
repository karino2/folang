package main

import (
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
	parser := &Parser{}
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
			`type hoge = {X: string; Y: string}

let ika () =
    {X="abc"; Y="def"}
`,
			`type hoge struct {
  X string
  Y string
}

func ika() *hoge{
return &hoge{X: "abc", Y: "def"}
}

`,
		},
	}

	for _, test := range tests {
		parser := &Parser{}
		res := parser.Parse("", []byte(test.input))
		tp := NewTranspiler()
		f := NewFile(res)
		got := tp.Transpile(f)

		if got != test.want {
			t.Errorf("got [%s], want [%s]", got, test.want)
		}
	}
}

func TestParserRecordDef(t *testing.T) {
	src := `type Hoge = {X: string; Y: string}
`

	parser := &Parser{}
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

	parser := &Parser{}
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
