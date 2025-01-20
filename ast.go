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
func (*IntImm) expr()        {}
func (*FunCall) expr()       {}
func (*Var) expr()           {}
func (*RecordGen) expr()     {}

type StringLiteral struct {
	Value string
}

func (*StringLiteral) FType() FType { return FString }

// TODO: エスケープ
func (s *StringLiteral) ToGo() string { return fmt.Sprintf(`"%s"`, s.Value) }

type IntImm struct {
	Value int
}

func (*IntImm) FType() FType { return FInt }

func (s *IntImm) ToGo() string { return fmt.Sprintf("%d", s.Value) }

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
	Func *Var
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

type RecordGen struct {
	fieldNames  []string
	fieldValues []Expr
	recordType  FType
}

func NewRecordGen(fieldNames []string, fieldValues []Expr) *RecordGen {
	return &RecordGen{fieldNames, fieldValues, &FUnresolved{}}
}

func (rg *RecordGen) FType() FType {
	return rg.recordType
}

func (rg *RecordGen) ToGo() string {
	rtype, ok := rg.recordType.(*FRecord)
	if !ok {
		panic("Unresolved record type.")
	}

	var buf bytes.Buffer
	buf.WriteString("&")

	buf.WriteString(rtype.StructName())
	buf.WriteString("{")
	for i, fname := range rg.fieldNames {
		if i != 0 {
			buf.WriteString(", ")
		}
		fval := rg.fieldValues[i]
		buf.WriteString(fname)
		buf.WriteString(": ")
		buf.WriteString(fval.ToGo())
	}
	buf.WriteString("}")
	return buf.String()
}

type Stmt interface {
	Node
	stmt()
	ToGo() string
}

func (*FuncDef) stmt()   {}
func (*Import) stmt()    {}
func (*Package) stmt()   {}
func (*RecordDef) stmt() {}
func (*UnionDef) stmt()  {}

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
	Params []*Var
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
	// TODO: multiple stmt support.
	if fd.Body.FType() != FUnit {
		buf.WriteString("return ")
	}
	buf.WriteString(fd.Body.ToGo())
	buf.WriteString("\n}")
	return buf.String()
}

type RecordDef struct {
	Name   string
	Fields []NameTypePair
}

// Use for type definition, type XXX struct {...}
func (rd *RecordDef) ToGo() string {
	var buf bytes.Buffer
	buf.WriteString("type ")
	buf.WriteString(rd.Name)
	buf.WriteString(" struct {\n")
	for _, field := range rd.Fields {
		buf.WriteString("  ")
		buf.WriteString(field.Name)
		buf.WriteString(" ")
		buf.WriteString(field.Type.ToGo())
		buf.WriteString("\n")
	}
	buf.WriteString("}")
	return buf.String()
}

func (rd *RecordDef) ToFType() *FRecord {
	return &FRecord{rd.Name, rd.Fields}
}

/*
Union implementation.
For following code:
type IntOrString =

	| I of int
	| B of bool

The result becomes following three types.

- IntOrString interface
- IntOrString_I struct (with Value int)
- IntOrString_B struct (with Value bool)

We call IntOrString_I "case struct of I".
*/
type UnionDef struct {
	Name  string
	Cases []NameTypePair
}

// Return IntOrString_I
func (ud *UnionDef) CaseStructName(index int) string {
	return fmt.Sprintf("%s_%s", ud.Name, ud.Cases[index].Name)
}

/*
	type IntOrString interface {
	  IntOrString_Union()
	}
*/
func (ud *UnionDef) buildUnionDef(buf *bytes.Buffer) {
	buf.WriteString("type ")
	buf.WriteString(ud.Name)
	buf.WriteString(" interface {\n")
	buf.WriteString("  ")
	buf.WriteString(ud.Name)
	buf.WriteString("_Union()\n")
	buf.WriteString("}\n")
}

/*
func (*IntOrString_I) IntOrString_Union(){}
func (*IntOrString_B) IntOrString_Union(){}
*/
func (ud *UnionDef) buildCaseStructConformMethod(buf *bytes.Buffer) {
	method := ud.Name + "_Union(){}\n"
	for i := range ud.Cases {
		buf.WriteString("func (*")
		buf.WriteString(ud.CaseStructName(i))
		buf.WriteString(") ")
		buf.WriteString(method)
	}
}

/*
	type IntOrString_I struct {
	   Value int
	}
*/
func (ud *UnionDef) buildCaseStructDef(buf *bytes.Buffer, index int) {
	buf.WriteString("type ")
	buf.WriteString(ud.CaseStructName(index))
	buf.WriteString(" struct {\n")
	buf.WriteString("  Value ")
	buf.WriteString(ud.Cases[index].Type.ToGo())
	buf.WriteString("\n}\n")
}

// New_IntOrString_I
func (ud *UnionDef) caseStructConstructorName(index int) string {
	return "New_" + ud.CaseStructName(index)
}

/*
func New_IntOrString_I(v int) IntOrString { return &IntOrString_I{v} }
*/
func (ud *UnionDef) buildCaseStructConstructor(buf *bytes.Buffer, index int) {
	buf.WriteString("func ")
	buf.WriteString(ud.caseStructConstructorName(index))
	buf.WriteString("(v ")
	buf.WriteString(ud.Cases[index].Type.ToGo())
	buf.WriteString(") ")
	buf.WriteString(ud.Name)
	buf.WriteString(" { return &")
	buf.WriteString(ud.CaseStructName(index))
	buf.WriteString("{v} }\n")
}

func (ud *UnionDef) ToGo() string {
	var buf bytes.Buffer
	ud.buildUnionDef(&buf)
	buf.WriteString("\n")
	ud.buildCaseStructConformMethod(&buf)
	buf.WriteString("\n")

	for i := range ud.Cases {
		ud.buildCaseStructDef(&buf, i)
		buf.WriteString("\n")
		ud.buildCaseStructConstructor(&buf, i)
		buf.WriteString("\n")
	}

	return buf.String()
}

func (ud *UnionDef) UnionFType() *FUnion {
	return &FUnion{ud.Name, ud.Cases}
}

func (ud *UnionDef) RegisterConstructorAlias(resolver *TypeResolver) {
	utype := ud.UnionFType()

	for i, cs := range ud.Cases {
		tps := []FType{cs.Type, utype}
		ftype := &FFunc{tps}
		resolver.AliasMap[cs.Name] = &Var{ud.caseStructConstructorName(i), ftype}
	}
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
			Walk(pm, f)
		}
		Walk(n.Body, f)
	case *Import, *Package, *RecordDef, *UnionDef:
		// no-op
	// ここからexpr
	case *GoEval, *StringLiteral, *Var, *IntImm:
		// no-op
	case *FunCall:
		Walk(n.Func, f)
		for _, arg := range n.Args {
			Walk(arg, f)
		}
	case *RecordGen:
		for _, fval := range n.fieldValues {
			Walk(fval, f)
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
