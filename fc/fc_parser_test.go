package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestParsePackage(t *testing.T) {
	src := `package main
`
	ps := initParse(src)
	gotPair := parsePackage(ps)
	got := gotPair.E1

	switch tgot := got.(type) {
	case RootStmt_RSPackage:
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
	got := RootStmtsToGo(stmts)

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
	return RootStmtsToGo(stmts)
}

func TestRawString(t *testing.T) {
	rawstring := "`abc\ndef \"dq\" backslash \\`"

	src := fmt.Sprintf(`package main

let hoge () =
  let a = %s
  a+"ghi"

`, rawstring)
	// t.Error(src)

	got := transpile(src)
	// t.Error(got)

	want := `"abc\ndef \"dq\" backslash \\"`
	if !strings.Contains(got, want) {
		t.Errorf("want to contain '%s', but got '%s'", want, got)
	}
}

func TestRawStringInterP(t *testing.T) {
	rawstring := "$`abc\ndef a=({a}) backslash \\`"

	src := fmt.Sprintf(`package main

let hoge () =
  let a = 123
	%s

`, rawstring)
	// t.Error(src)

	got := transpile(src)
	// t.Error(got)

	want := `frt.SInterP("abc\ndef a=(%s) backslash \\", a)`
	if !strings.Contains(got, want) {
		t.Errorf("want to contain '%s', but got '%s'", want, got)
	}
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

func (v IntOrString_I) String() string { return frt.Sprintf1("(I: %v)", v.Value) }
func (v IntOrString_S) String() string { return frt.Sprintf1("(S: %v)", v.Value) }

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

func (v AorB_A) String() string { return "(A)" }
func (v AorB_B) String() string { return "(B)" }

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
			`package_info buf =
  type Buffer
  let WriteString: Buffer->string->()

let main () =
  let b = GoEval<buf.Buffer> "buf.Buffer{}"
  buf.WriteString b "hogehoge"
`,
			"buf.WriteString(b, \"hogehoge\")",
		},
		// comment inside package_info
		{
			`package_info buf =
  type Buffer
  // comment test
  let WriteString: Buffer->string->()

let main () =
  let b = GoEval<buf.Buffer> "buf.Buffer{}"
  buf.WriteString b "hogehoge"
`,
			"buf.WriteString(b, \"hogehoge\")",
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
		{
			`
type NameTypePair = {Name: string; Type: string}

type RecordType = {name: string; fiedls: []NameTypePair}
`,
			"fiedls []NameTypePair",
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
			"ika() []int{",
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
			`package main

package_info slice =
  let Zip<T, U>: []T->[]U->[](T*U)

let ika () =
  let s1 = [1; 2; 3]
  let s2 = ["a"; "b"; "c"]
  slice.Zip s1 s2
`,
			"frt.Tuple2[int, string]",
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
  let Zip<T, U>: []T->[]U->[](T*U)

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
  tkz: string;
  offsideCol: []int;
}
`,
			"ParseState", // whatever
		},
		{
			`package main

type hoge = {X: string; Y: int}
type ika = {X: string; Y: int}

let fuga () =
   let h = {X="ab"; Y="de"}
   let i = {ika.X="gh"; Y="jk"}
   (h, i)
`,
			"ika{X: ",
		},
		{
			`package main

let ika () =
  (123, "abc")

let fuga () =
  let (a, b) = ika ()
  a+1

`,
			"a, b := frt.Destr2(ika())",
		},
		{
			`package main

let ika () =
  (123, "abc")

let fuga () =
  let (_, b) = ika ()
  b+"def"
`,
			"_, b :=",
		},
		{
			`package main

let ika () =
  (123, "abc")

let fuga () =
  let (a, _) = ika ()
  a+4

`,
			"a, _ :=",
		},
		{
			`package main

package_info _=
  let lookupFunc: string->(()->string)*bool

let ika () =
  lookupFunc "hoge"

`,
			"ika() frt.Tuple2[func () string, bool]{",
		},
		{
			`package main

let tpname2tvtp (tvgen: ()->string) (tpname:string) =
  let tv = tvgen ()
  (tpname, tv)

`,
			"tvgen func () string",
		},
		{
			`package main

package_info _ =
  let lookupVarFac: string->((()->string)->string)*bool


let hoge () =
  lookupVarFac "abc"

`,
			"[func (func () string) string, bool]",
		},
		{
			`package main

type IorS =
  | IT of int
  | ST of string

let nestMatch (lhs:IorS) (rhs:IorS) =
  match lhs with
  | IT ival ->
    match rhs with
    | IT i2 ->
      ival+i2
    | _ ->
      ival+456
  | _ ->
    123
`,
			"return (ival+456)\n}",
		},
		{
			`package main

package_info _ =
  let Concat<T>: [][]T -> []T

let hoge () =
  let s1 = GoEval<[]int> "[]int{1, 2}"
  let s2 = GoEval<[]int> "[]int{3, 4}"
  let s3 = [s1; s2]
  Concat s3

`,
			"hoge() []int{",
		},
		{
			`package main

package_info _ =
  let Concat<T>: [][]T -> []T

let hoge () =
  let s1 = GoEval<[]int> "[]int{1, 2}"
  let s2 = GoEval<[]int> "[]int{3, 4}"
  [s1; s2]
  |> Concat

`,
			"hoge() []int{",
		},
		// tuple two typevar resolution.
		{
			`package main

package_info _ =
  let Map<T, U> : (T->U)->[]T->[]U
  let Snd<T, U>: T*U->U

let hoge () =
  let s1 = GoEval<[]int> "[]int{1, 2}"
  let s2 = GoEval<[]int> "[]int{3, 4}"
  [(123, s1); (456, s2)]
  |> Map Snd

`,
			"[][]int{",
		},
		{
			`package main

package_info _ =
  let Map<T, U> : (T->U)->[]T->[]U
  let Snd<T, U> : T*U->U

let hoge () =
  let s1 = GoEval<[]int> "[]int{1, 2}"
  let s2 = GoEval<[]int> "[]int{3, 4}"
  let tups = [(123, s1); (456, s2)]
  Map Snd tups

`,
			"[][]int{",
		},
		{
			`package main

package_info _ =
  let Head<T> : []T->T

type Item = {f3:string; f4:int}

let hello (is:[]Item) = 
  let fr = is |> Head
  fr.f3
`,
			"string{",
		},
		// destructuring let inference.
		{
			`package main

let hoge (a:int) =
  (123, "abc")

let ika () = 
  let (a, b) = 123 |> hoge
  a+123
`,
			"ika() int",
		},
		{
			`package main

type Expr =
| EInt of int
| ESlice of []Expr

let hoge (e:Expr) =
  match e with
  | EInt i -> i
  | ESlice _ -> 123

`,
			"Value []Expr",
		},
		// IfOnly became wrong func type test.
		{
			`package main

package_info _ =
  let Println: string->()


let hoge (i:int) =
  if i <> 3 then
	  Println "hit"

`,
			"hoge(i int) {",
		},
		// BinOp resolve check.
		{
			`package main

package_info _ =
  let Print: string->()


let hoge (uname: string) =
  if uname <> "_"  then
	  Print "hoge"
  123
`,
			"hoge(uname",
		},
		{
			`package main

package_info slice =
  let New<T>: ()->[]T
  let IsEmpty<T>: []T->bool

let ika () =
  let s = slice.New<string> ()
  slice.IsEmpty s

`,
			"New[string]",
		},
		{
			`package main

package_info slice =
  let New<T>: ()->[]T
  let Head<T>: []T->T

let ika () =
  let s = slice.New<string> ()
  slice.Head s

`,
			"ika() string",
		},
		{
			`package main

let add (a:int) b = 
  a+b
`,
			"int{",
		},
		{
			`package main

package_info frt =
  let Sprintf1<T>: string->T->string

package_info slice =
  let Map<T, U> : (T->U)->[]T->[]U

let toS a =
  frt.Sprintf1 "%d" a

let hoge () =
  slice.Map toS [1; 2; 3]
`,
			"toS, ([]int{1,2,3})",
		},
		{
			`package main

package_info _ =
  let PushLast<T>: T->[]T->[]T

let hoge iToS =
	["abc"; "def"]
	|> PushLast (iToS 123)

`,
			"(iToS func (int) string) []string",
		},
		{
			`package main

type Rec={Slice: []string}

let hoge t =
  {Slice=t}

`,
			"t []string",
		},
		{
			`package main

let hoge () =
  3*2/3
`,
			"((3*2)/3)",
		},
		{
			`package main

package_info _ =
  let Head<T>: []T->T

let hoge () =
  [1; 2; 3]
	|> fun x-> Head x
`,
			"(x []int)int",
		},
		{
			`package main

type A =
| UB of B
and B = {Fc: C}
and C = {Fd: D}
and D =
| I of int
| S of string

let hoge c =
  {Fc=c}

`,
			"hoge(c C)",
		},
		{
			`package main

let hoge () =
  let a = 123
  $"This is a: {a}, \nnext line"

`,
			", \\nnext",
		},
		{
			`package main

type Result<T> =
| Success of T
| Failure

let hoge () =
  Failure<string> ()

`,
			"hoge() Result[string]",
		},
		{
			`package main

package_info _ =
  let GetF: ()->float

let hoge () =
  GetF ()

`,
			"hoge() float64",
		},
		// 3 tuple package_info test.
		// if inference pass, it must be parsed collectly.
		{
			`package main

type Expr =
| BinOp of (int*string*int)

let ika a b c =
  (a, b, c) |> BinOp

`,
			"frt.NewTuple3", // whatever.
		},
		{
			`package main

type AorB =
 | A
 | B

let ika () =
  let a =
    match A with
    | A -> "a match"
    | B -> "b match"
  a
`,
			"a := (func",
		},
		{
			`package main

let ika () =
  let a =
    match "abc" with
    | "abc" -> "hit"
    | _ -> "not hit"
  a
`,
			`switch ("abc")`,
		},
		{
			`package main

let ika () =
  let a =
    match "abc" with
    | "abc" -> "hit"
    | s -> s + " not hit"
  a
`,
			`switch s :=("abc"); s{`,
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

package_info buf =
  type Buffer
  let New: ()->Buffer
  let Write: Buffer->string->()

let hoge () =
  buf.New ()

`,
			[]string{"buf.Buffer", "return buf.New()"},
		},
		{
			`package main

type FType =
| FInt
| FSlice of SliceType
and SliceType = {ElemType: FType}

let hoge () =
  let rec = FSlice {ElemType=FInt}
  match rec with
  | FSlice s-> s.ElemType
  | _ -> FInt
`,
			[]string{"Value SliceType", "ElemType FType", "hoge() FType", "s.ElemType"},
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
			[]string{") frt.Tuple2[int, int]", "frt.NewTuple2(1, 2)"},
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
		{
			`package main

let ika (a:int) =
  if a = 1 then
	  "abc"
	elif a <> 5 then
	  "def"
	elif a = 3 then
	  "xxx"
	else
	  "ghi"

`,
			[]string{"frt.IfElse(frt.OpEqual(a, 1)", "return frt.IfElse(frt.OpNotEqual(a, 5)"},
		},
		{
			`package main

package_info _ =
  let IsEmpty<T> : []T->bool
  let Head<T>: []T->T
  let Tail<T>: []T->[]T

let sum (args: []int) :int =
  if IsEmpty args then
    0
  else
    let h = Head args
    let tail = Tail args
    h + (sum tail)

`,
			[]string{"sum(args []int) int{", "sum(tail)"},
		},
		{
			`package main

package_info _ =
  let Head<T>: []T->T


let hoge i =
  Head i

`,
			[]string{"[T0 any]", "T0{"},
		},
		{
			`package_info dict =
  type Dict<K, V>
  let New<K, V>: ()->Dict<K, V>
  let Put<K, V>: Dict<K, V>->K->V->()
  let Lookup<K, V>: Dict<K, V>->K->V*bool


let hoge () =
  let dic = dict.New<string, int> ()
  dict.Put dic "hoge" 123
  let (v, ok) = dict.Lookup dic "hoge"
  if ok then
    v
  else
    999

`,
			[]string{"hoge() int", "New[string, int]"},
		},
		{
			`package_info dict =
  type Dict<K, V>
  let Add<K, V>: Dict<K, V>->K->V->()

let setAddKeys d k =
  dict.Add d k true

`,
			[]string{"setAddKeys[T0 any]", "Dict[T0, bool]", "k T0"},
		},
		{
			`package main

package_info frt =
  let Fst<T, U> : T*U->T
  let Snd<T, U> : T*U->U


let CnvL fn tup =
  let nl = frt.Fst tup |> fn
  (nl, frt.Snd tup)
`,
			[]string{"[T0 any, T1 any, T2 any]", "fn func (T0) T1", "tup frt.Tuple2[T0, T2]"},
		},
		{
			`package main

package_info dict =
  type Dict<K, V>
  let New<K, V>: () -> Dict<K, V>


let hoge () =
  dict.New<any, string> ()
`,
			[]string{"Dict[any, string]", "New[any, string]"},
		},
		{
			`package main
			
			package_info _ =
				let Map<T, U> : (T->U)->[]T->[]U
			
			
			let hoge () =
				Map (fun x->x+2) [1; 2; 3]
			`,
			[]string{"func (x int)int", "return (x+2)"},
		},
		{
			`package main

type Rec = {S: string; I:int}

package_info _ =
	let Map<T, U> : (T->U)->[]T->[]U

let hoge (rs:[]Rec) =
  rs |> Map _.S

`,
			[]string{"func (_v1 Rec)string", "_v1.S"},
		},
		{
			`package main

type A<T> = {FA: T; FSA: []T}

let hoge () =
  {FA=123; FSA=[456]}

`,
			[]string{"type A[T any]", "hoge() A[int]", "return A[int]"},
		},
		{
			`package main

let hoge () =
  let a = 123
  let b = "abc"
  $"This is a: {a}, b: {b}"

`,
			[]string{`return frt.SInterP("This is a: %s, b: %s", a, b)`, "hoge() string"},
		},
		{
			`package main

let g_var = "abc"

let hoge () =
  g_var + "def"

`,
			[]string{`var g_var = "abc"`, "return (g_var"},
		},
		{
			`package main

type Result<T> =
| Success of T
| Failure

let hoge () =
  Success 123

`,
			[]string{"New_Result_Failure[T any]()", "hoge() Result[int]"},
		},
		{
			`package main

let hoge () =
  (1, "2", 3)

`,
			[]string{"frt.Tuple3[int, string, int]", `frt.NewTuple3(1, "2", 3)`},
		},
		{
			`package main

let ika () =
  let (a, b, c) = (1, "2", 3)
  a+c

`,
			[]string{"frt.Destr3", "frt.NewTuple3"},
		},
		{
			`package main

let ika (str:string) =
  let sfun b =
	  b + " append"
  sfun str
`,
			[]string{"sfun := func (b string)string{", "return sfun(str)"},
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

func TestNonExaustive(t *testing.T) {
	src := `package main

type Union =
| A
| B of string

let hoge (u:Union) =
  match u with
  | A -> ""
`

	var panicFound = false
	{
		defer func() {
			if r := recover(); r != nil {
				panicFound = true
			}
		}()
		transpile(src)
		// t.Error(got)
	}

	if !panicFound {
		t.Error("Not exausitve, but not check error.")
	}
}

func TestParseAddhook(t *testing.T) {
	src := `package main

let ika (str:string) =
  let sfun b =
	  b + " append"
  sfun str
`

	got := transpile(src)
	// t.Error(got)
	if got == "dummy" {
		t.Error(got)
	}
}
