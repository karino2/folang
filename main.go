package main

import (
	"bytes"
	"fmt"
)

type FType interface {
	ftype()
	ToGo() string // Goでの型を表す文字列、表せないものをどうするかはあとで考える
}

func (*FPrimitive) ftype() {}
func (*FFunc) ftype()      {}

type FPrimitive struct {
	Name string // "int", "string" etc.
}

func (p *FPrimitive) ToGo() string {
	// unitはなにも無しとしておく。関数のreturnなどはそれで動くので。
	if p.Name == "unit" {
		return ""
	}
	return p.Name
}

var (
	FInt    = &FPrimitive{"int"}
	FString = &FPrimitive{"string"}
	FUnit   = &FPrimitive{"unit"}
)

type FFunc struct {
	// カリー化されている前提で、引数も戻りも区別せず持つ
	Targets []FType
}

func (p *FFunc) Args() []FType {
	return p.Targets[0 : len(p.Targets)-1]
}

func (p *FFunc) ReturnType() FType {
	return p.Targets[len(p.Targets)-1]
}

func (p *FFunc) ToGo() string {
	var buf bytes.Buffer
	buf.WriteString("func (")
	for i, tp := range p.Args() {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(tp.ToGo())
	}

	last := p.ReturnType()
	if last != FUnit {
		buf.WriteString(" ")
		buf.WriteString(last.ToGo())
	}
	return buf.String()
}

type Expr interface {
	expr()
	FType() FType
	ToGo() string
}

func (*GoEval) expr()        {}
func (*StringLiteral) expr() {}
func (*FunCall) expr()       {}
func (*Var) expr()           {}

type StringLiteral struct {
	Value string
}

func (*StringLiteral) FType() FType { return FString }

// TODO: エスケープ
func (s *StringLiteral) ToGo() string { return fmt.Sprintf(`"%s"`, s.Value) }

// Goのコードを直接持つinline asm的な抜け穴
type GoEval struct {
	GoStmt string
}

func (*GoEval) FType() FType   { return FUnit }
func (e *GoEval) ToGo() string { return e.GoStmt }

// 変数。仮引数などの場合と変数自身の参照の場合の両方をこれで賄う。
type Var struct {
	Name string
	Type FType
}

func (v *Var) FType() FType { return v.Type }
func (v *Var) ToGo() string { return v.Name }

type FunCall struct {
	Func Var
	Args []Expr
}

func (fc *FunCall) FuncType() *FFunc {
	return fc.Func.Type.(*FFunc)
}

func (fc *FunCall) ArgTypes() []FType {
	return fc.FuncType().Args()
}

func (fc *FunCall) FType() FType {
	ftype := fc.FuncType()
	if len(fc.Args) == len(ftype.Args()) {
		return ftype.ReturnType()
	}
	panic("partial apply, NYI")
}
func (fc *FunCall) ToGo() string {
	var buf bytes.Buffer
	if len(fc.Args) != len(fc.ArgTypes()) {
		panic("partial apply, NYI")
	}
	buf.WriteString(fc.Func.Name)
	buf.WriteString("(")
	for i, arg := range fc.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.ToGo())
	}
	buf.WriteString(")")
	return buf.String()
}

type Stmt interface {
	stmt()
	ToGo() string
}

func (*FuncDef) stmt() {}
func (*Import) stmt()  {}
func (*Package) stmt() {}

type Import struct {
	PackageName string
}

func (im *Import) ToGo() string {
	return fmt.Sprintf("import \"%s\"", im.PackageName)
}

type Package struct {
	Name string
}

func (p *Package) ToGo() string {
	return fmt.Sprintf("package %s", p.Name)
}

type FuncDef struct {
	Name string
	// Unitはパース時点で0引数argsに変換済みの想定
	Args []Var
	Body Expr
}

func (fd *FuncDef) ToGoParams(buf *bytes.Buffer) {
	for i, arg := range fd.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.Name)
		buf.WriteString(" ")
		buf.WriteString(arg.Type.ToGo())
	}
}

func (fd *FuncDef) ToGo() string {
	var buf bytes.Buffer
	buf.WriteString("func ")
	buf.WriteString(fd.Name)
	buf.WriteString("(")
	fd.ToGoParams(&buf)
	buf.WriteString(") ")
	buf.WriteString(fd.Body.FType().ToGo())
	buf.WriteString("{\n")
	buf.WriteString(fd.Body.ToGo())
	buf.WriteString("\n}")
	return buf.String()
}

type Program struct {
	Stmts []Stmt
}

func (p *Program) AddStmt(stmt Stmt) {
	p.Stmts = append(p.Stmts, stmt)
}

func (p *Program) ToGo() string {
	var buf bytes.Buffer
	for _, stmt := range p.Stmts {
		buf.WriteString(stmt.ToGo())
		buf.WriteString("\n\n")
	}
	return buf.String()
}

func main() {
	var p Program
	p.AddStmt(&Package{"main"})
	p.AddStmt(&Import{"fmt"})
	p.AddStmt(&FuncDef{"hello", []Var{{"msg", FString}}, &GoEval{"fmt.Printf(\"Hello %s\\n\", msg)"}})
	p.AddStmt(&FuncDef{"main", nil,
		&FunCall{
			Var{"hello", &FFunc{[]FType{FString, FUnit}}},
			[]Expr{&StringLiteral{"Hoge"}},
		},
	})
	fmt.Println(p.ToGo())
}

/*
type ElemType int

const (
	TypeSymbol ElemType = iota
)

type Elem interface {
	Type() ElemType

	Parent() Elem
	FirstChild() Elem
	NextSibling() Elem
}

// Treeのノードの共通部分。
type Node struct {
	elemType ElemType

	parent      Elem
	firstChild  Elem
	nextSibling Elem
}

func (n *Node) Type() ElemType {
	return n.elemType
}

func (n *Node) Parent() Elem {
	return n.parent
}

func (n *Node) FirstChild() Elem {
	return n.firstChild
}

func (n *Node) NextSibling() Elem {
	return n.nextSibling
}

type SymbolElem struct {
	common Node
	Id     string
}
*/
