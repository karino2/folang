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
		t.Errorf("expect 2 stmt, only %d, %v", len(res), res)
	}

	wants := []string{`package main`, `import "fmt"`}
	for i, want := range wants {
		got := res[i].ToGo()
		if got != want {
			t.Errorf("want '%s', got '%s'", want, got)
		}
	}
}

func TestParserHello(t *testing.T) {
	src :=
		`package main
import "fmt"

let main () =
    GoEval "fmt.Println(\"Hello World\")"
`
	parser := &Parser{}
	res := parser.Parse("", []byte(src))
	if len(res) != 3 {
		t.Errorf("expect 2 stmt, only %d, %v", len(res), res)
	}
	p := NewProgram(res)

	want := `package main

import "fmt"

func main() {
fmt.Println("Hello World")
}

`
	got := p.Compile()
	if got != want {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}
