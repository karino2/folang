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
		{`"hoge(\"ika\")"`, &Token{ttype: STRING, begin: 0, len: 15, stringVal: `hoge(\"ika\")`}},
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

	// Type register test.
	got := parser.scope.LookupType("Hoge")
	_, ok = got.(*FRecord)
	if !ok {
		t.Errorf("cannot find FRecord by name")
	}

	got2 := parser.scope.LookupRecord([]string{"X", "Y"})
	if got2 == nil {
		t.Errorf("cannot find FRecord by field")
	}

	got3 := parser.scope.LookupRecord([]string{"X", "DEF"})
	if got3 != nil {
		t.Errorf("Wrongly matched with different name")
	}

	got4 := parser.scope.LookupRecord([]string{"X", "Y", "Z"})
	if got4 != nil {
		t.Errorf("Wrongly matched with extra field")
	}

	got5 := parser.scope.LookupRecord([]string{"X"})
	if got5 != nil {
		t.Errorf("Wrongly matched with few fields")
	}
}

func TestParserRecordExpression(t *testing.T) {
	src := `type Hoge = {X: string; Y: string}

let hoge () =
  {X="hoge"; Y="ika"}
`

	parser := NewParser()
	res := parser.Parse("", []byte(src))
	if len(res) != 2 {
		t.Errorf("expect record gen expression")
	}

	rg, ok := res[1].(*FuncDef).Body.FinalExpr.(*RecordGen)

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
		t.Errorf("want 2 rules, but %d", len(me.rules))
	}
}

func TestPkgInfoTypeDef(t *testing.T) {
	src := `package_info slice =
  type Buffer

`
	parser := NewParser()
	got := parser.Parse("", []byte(src))

	pi, ok := got[0].(*PackageInfo)
	if !ok {
		t.Errorf("Not pkg info: %v", got)
	}

	ctp, ok2 := pi.typeInfo["Buffer"]
	if !ok2 {
		t.Error("Buffer type does not exist")
	}

	if ctp.name != "slice.Buffer" {
		t.Errorf("want slice.Buffer, got %s", ctp.name)
	}
}

func TestPkfInfoTypeFun(t *testing.T) {
	src := `package_info slice =
  type Buffer
  let WriteString: Buffer->()
`
	parser := NewParser()
	got := parser.Parse("", []byte(src))

	pi, ok := got[0].(*PackageInfo)
	if !ok {
		t.Errorf("Not pkg info: %v", got)
	}

	if len(pi.typeInfo) != 1 {
		t.Errorf("want 1 type, but %v", pi.typeInfo)
	}

	ft, ok2 := pi.funcInfo["WriteString"]
	if !ok2 {
		t.Error("WriteString does not exist")
	}

	if len(ft.Targets) != 2 {
		t.Errorf("WriteString Targets is not 2. %v", ft.Targets)
	}

	if ft.Targets[1] != FUnit {
		t.Errorf("Return value must be FUnit, but not. %v", ft.Targets[1])
	}

	if ft.Targets[0] != pi.typeInfo["Buffer"] {
		t.Errorf("Want argument be Buffer type, but not. arg: %v, Buffer type %v", ft.Targets[0], pi.typeInfo["Buffer"])
	}

}

func TestPkgInfoUnderScoreDef(t *testing.T) {
	src := `package_info _ =
  type Buffer
  let WriteString: Buffer->()

`
	parser := NewParser()
	parser.Parse("", []byte(src))

	got := parser.scope.LookupType("Buffer")
	if _, ok := got.(*FExtType); !ok {
		t.Errorf("Buffer type is not FExtType: %v", got)
	}

	got2 := parser.scope.LookupVar("WriteString")
	if _, ok := got2.Type.(*FFunc); !ok {
		t.Errorf("WriteString type is not FFunc: %v", got)
	}
}

func TestBlockComment(t *testing.T) {
	src := `package main

/*
  This is comment, never found.
*/

let ika () =
  123
`

	got := transpile(src)
	if strings.Contains(got, "never found") {
		t.Errorf("comment is not skipped: '%s'", got)
	}
}

func TestLineComment(t *testing.T) {
	src := `package main

let ika () =
  123 // this is line comment, never found.
`

	got := transpile(src)
	if strings.Contains(got, "never found") {
		t.Errorf("comment is not skipped: '%s'", got)
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
		{
			`package_info buf =
  type Buffer
  let WriteString: Buffer->string->()

let main () =
  let b = GoEval<buf.Buffer> "buf.Buffer{}"
  buf.WriteString b "hogehoge"
`,
			"buf.WriteString(b, \"hogehoge\")",
		},
		{
			`package_info buf =
  type Buffer
  let WriteString: Buffer->string->()
  let New: ()->Buffer

let main () =
  let b = buf.New ()
  buf.WriteString b "hogehoge"
`,
			"buf.New()",
		},
		// resolve type parameter test.
		{
			`package_info slice =
  let Length<T>: []T -> int
  let Take<T> : int->[]T->[]T

let ika () =
  let s = GoEval<[]int> "int[]{1, 2}"
  slice.Take 2 s
`,
			"ika() []int{", // Take T must becomes int
		},
		{
			`package_info slice =
  let Take<T> : int->[]T->[]T

let ika () =
  let s = GoEval<[]int> "int[]{1, 2}"
  slice.Take 2 s
`,
			"return slice.Take(2, s)",
		},
		{
			`package_info slice =
  let Map<T, U> : (T->U)->[]T->[]U

let conv (i:int) =
  GoEval<string> "fmt.Sprintf(\"a %d\", i)"

let ika () =
  let s = GoEval<[]int> "int[]{1, 2}"
  slice.Map conv s
`,
			"ika() []string",
		},
		{
			`import slice
`,
			`import "github.com/karino2/folang/pkg/slice"`,
		},
		// left assoc
		{
			`package main
			
			let hoge () =
				5-7+1+2
			`,
			"(((5-7)+1)+2)",
		},
		// left assoc + ()
		{
			`package main

let hoge () =
  5-7+(1+2+3)
`,
			"((5-7)+((1+2)+3))",
		},
		{
			// once match parse wrongly move one token after end.
			// So this parse is failed.
			// check whether this is not failed. result string is not important.
			`package main

type AorB =
 | A
 | B

let ika (ab:AorB) =
  match ab with
  | A -> "a match"
  | B -> "b match"

/*
this is test
*/
`,
			"AorB", // whatever.
		},
		{
			`let ika (s1:string) (s2:string) =
  s1 <> s2

`,
			"frt.OpNotEqual(s1, s2)",
		},
		{
			`package main

let ika (s1:string) (s2:string) =
  if s1 = s2 then
    123
  else
    456
`,
			"frt.OpEqual",
		},
		// comment handling of block.
		{
			`package main

let ika (s1:string) (s2:string) =
  if s1 = s2 then
    123
  else
    // this line is comment
    456
`,
			"frt.OpEqual",
		},
		// type resolve for non-funcall function
		{
			`
package main

package_info slice =
  let Sort<T>: []T -> []T


let ika (fields: []string) =
  fields |> slice.Sort

`,
			"frt.Pipe",
		},
		// comment inside union case. just pass parse is enough.
		{
			`package main

type AorB =
 | A
 // comment here.
 // comment here2.
 | B

let ika (ab:AorB) =
  match ab with
  | A -> "a match"
  | B -> "b match"

`,
			"a match", // whatever.
		},
		// comment inside match case. just parse is enough for test.
		{
			`package main

type AorB =
 | A
 | B

let ika (ab:AorB) =
  match ab with
  | A -> "a match"
  // comment here.
  | B -> "b match"

`,
			"b match",
		},
		{
			`package main

package_info slice =
  let Zip<T, U>: []T->[]U->[]T*U

let ika () =
  let s1 = [1; 2; 3]
	let s2 = ["a"; "b"; "c"]
  slice.Zip s1 s2
`,
			"frt.Tuple2[int, string]",
		},
		{
			`package main

let ika () =
  "\n"
`,
			`return "\n"`,
		},
		{
			`package main

let ika (a:int) =
  if a = 0 && a = 2 then
    "abc"
	else
	  "def"
`,
			"(frt.OpEqual(a, 0)&&frt.OpEqual(a, 2))",
		},
		{
			`let ika (a:int) =
  if not (a = 0) then
    "abc"
	else
	  "def"

`,
			"frt.OpNot(",
		},
		// comment inside block.
		{
			`package main

let ika () =
   let a = 1
	 // this is comment
	 a
`,
			"ika", // whatever
		},
		// last field semicolon ending.
		{
			`
type ParseState = {
  tkz: Tokenizer;
  offsideCol: []int;
}
`,
			"ParseState", // whatever
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

/*
Check whether multiple string want contains.
*/
func TestParserContainMultiple(t *testing.T) {
	var tests = []struct {
		input string
		wants []string
	}{
		{
			`package_info slice =
  let Take<T> : int->[]T->[]T

let ika () =
  let s = GoEval<[]int> "[]int{1, 2, 3}"
  s |> slice.Take 2
`,
			[]string{"func (_r0 []int) []int", "frt.Pipe("},
		},
		{
			`package main

type Hoge = {X: string; Y: string}

let hoge () =
  let rec = {X="hoge"; Y="ika"}
  rec.Y
`,
			[]string{"hoge() string", "return rec.Y"},
		},
		{
			`package main

type RecordType = {name: string; ival: int}

type IorRec =
| Int of int
| Rec of RecordType

let hoge () =
  let rec = Rec {name="hoge"; ival=123}
  match rec with
	| Int i -> i
	| Rec r-> r.ival
`,
			[]string{"hoge() int", "r.ival"},
		},
		{
			`package main

type FType =
| FInt
| FSlice of SliceType
and SliceType = {elemType: FType}

let hoge () =
  let rec = FSlice {elemType=FInt}
  match rec with
	| FSlice s-> s.elemType
	| _ -> FInt
`,
			[]string{"Value SliceType", "elemType FType", "hoge() FType", "s.elemType"},
		},
		// pipe to unit func test.
		{
			`package main

package_info buf =
  type Buffer
  let New: ()->Buffer
  let Write: Buffer->string->()

let hoge () =
  let bw = buf.New ()
	"abc" |> buf.Write bw
`,
			[]string{"frt.PipeUnit", "{ buf.Write" /* no return */},
		},
		{
			`let ika (s1:string) (s2:string) =
  s1 = s2

`,
			[]string{"frt.OpEqual(s1, s2)", "s2 string) bool"},
		},
		{
			`package main

let ika (s1:string) (s2:string) =
  if s1 = s2 then 123 else 456
`,
			[]string{"return 123", "123}),"}, // no call.
		},
		{
			`package main

let ika () =
  [1; 2; 3]

`,
			[]string{"return ([]int{1,2,3})", "ika() []int"},
		},
		{
			`package main

let ika () =
  (1, 2)

`,
			[]string{") frt.Tuple2[int, int]", "frt.NewTuple2(1,2)"},
		},
		{
			`package main

package_info pair =
  let Fst<T, U> : T*U->T

let ika () =
  pair.Fst (1, "s")
`,
			[]string{"ika() int{", "pair.Fst"},
		},
		// unit arg parse.
		{
			`package main

let ika () =
  "123"

let hoge () =
   let s = ika ()
	 s
`,
			[]string{"ika() string", "= ika()"},
		},
		{
			`package main

type Inner = {Name: string}
type Nested = {Name: string; Elem: Inner}

let ika (a:Nested) =
   a.Elem.Name
`,
			[]string{") string", "a.Elem.Name"},
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

func TestParserAddhook(t *testing.T) {
	src := `package main

let ika () =
   let a = 1
	 // this is comment
	 a

`

	got := transpile(src)
	// t.Error(got)
	if got == "" {
		t.Error(got)
	}
}
