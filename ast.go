package main

import (
	"bytes"
	"fmt"
)

// ExprとStmtの共通インターフェース。今の所空。
type Node interface{}

type Expr interface {
	Node
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

func (v *Var) FType() FType       { return v.Type }
func (v *Var) ToGo() string       { return v.Name }
func (v *Var) IsUnresolved() bool { return IsUnresolved(v.Type) }

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
	Node
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
	// Unitはパース時点で0引数paramsに変換済みの想定
	Params []Var
	Body   Expr
}

func (fd *FuncDef) ToGoParams(buf *bytes.Buffer) {
	for i, arg := range fd.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.Name)
		buf.WriteString(" ")
		buf.WriteString(arg.Type.ToGo())
	}
}

func (fd *FuncDef) FuncFType() FType {
	var fts []FType
	for _, arg := range fd.Params {
		fts = append(fts, arg.Type)
	}
	fts = append(fts, fd.Body.FType())
	return &FFunc{fts}
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

/*
StmtとExprをトラバースしていく。
f(node)がtrueを返すと子どもを辿っていく。
最後はf(nil)を呼ぶ。
*/
func Walk(n Node, f func(Node) bool) {
	if n == nil {
		panic("nil")
	}
	if !f(n) {
		return
	}

	switch n := n.(type) {
	case *FuncDef:
		for _, pm := range n.Params {
			Walk(&pm, f)
		}
		Walk(n.Body, f)
	case *Import, *Package:
		// no-op
	// ここからexpr
	case *GoEval, *StringLiteral, *Var:
		// no-op
	case *FunCall:
		Walk(&n.Func, f)
		for _, arg := range n.Args {
			Walk(arg, f)
		}
	default:
		panic(n)
	}
	f(nil)
}

func WalkStmts(stmts []Stmt, f func(Node) bool) {
	for _, stmt := range stmts {
		Walk(stmt, f)
	}
}
