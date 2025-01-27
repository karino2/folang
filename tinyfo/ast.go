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
func (*Block) expr()         {}
func (*MatchExpr) expr()     {}

/*
Some expr is OK for just return in some situation (like body of function).
For example, golang switch is not a expression, to use it as expression, we need to wrap (func()...{...})().
But if we use it in function body, just return is enough.
For those situation, Use ToGoReturn() for non-wrapped go code.
*/
type ReturnableExpr interface {
	Expr
	ToGoReturn() string
}

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
	GoStmt  string
	TypeArg FType
}

func (e *GoEval) FType() FType     { return e.TypeArg }
func (e *GoEval) ToGo() string     { return e.GoStmt }
func NewGoEval(src string) *GoEval { return &GoEval{src, FUnit} }

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
	if fc.Func.IsUnresolved() {
		return fc.Func.Type // return FUnresolved.
	}
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

/*
Block contains stmts and final return expr.
*/
type Block struct {
	Stmts     []Stmt
	FinalExpr Expr
}

func (b *Block) FType() FType { return b.FinalExpr.FType() }

func wrapFunCall(returnType FType, goReturnBody string) string {
	var buf bytes.Buffer
	buf.WriteString("(func () ")
	buf.WriteString(returnType.ToGo())
	buf.WriteString(" {\n")
	buf.WriteString(goReturnBody)
	buf.WriteString("})()")
	return buf.String()
}

func (b *Block) ToGo() string {
	return wrapFunCall(b.FType(), b.ToGoReturn())
}

func (b *Block) ToGoReturn() string {
	var buf bytes.Buffer
	for _, s := range b.Stmts {
		buf.WriteString(s.ToGo())
		buf.WriteString("\n")
	}
	last := b.FinalExpr

	if lastre, ok := last.(ReturnableExpr); ok {
		// return is generated inside lastre.ToGoReturn().
		buf.WriteString(lastre.ToGoReturn())
	} else {
		if last.FType() != FUnit {
			buf.WriteString("return ")
		}
		buf.WriteString(last.ToGo())
	}

	return buf.String()
}

var uniqueId = 0

func UniqueTmpVarName() string {
	uniqueId++
	return fmt.Sprintf("_v%d", uniqueId)
}

func ResetUniqueTmpCounter() {
	uniqueId = 0
}

/*
Union case matching.
Only support variable match for a while:
| I i -> ...
| Record r -> ...
*/
type MatchPattern struct {
	caseId  string
	varName string
}

type MatchRule struct {
	pattern *MatchPattern
	body    *Block
}

/*
	produce

case IntOrBool_I:

	[varName] := tmpV
	body...

or

default:

	...body...
*/
func (mr *MatchRule) ToGo(uname string, tmpVarName string) string {
	var buf bytes.Buffer
	pat := mr.pattern
	if pat.caseId == "_" {
		buf.WriteString("default:\n")
	} else {
		buf.WriteString("case ")
		buf.WriteString(UnionCaseStructName(uname, pat.caseId))
		buf.WriteString(":\n")
		if pat.varName != "_" && pat.varName != "" {
			buf.WriteString(pat.varName)
			buf.WriteString(" := ")
			buf.WriteString(tmpVarName)
			buf.WriteString(".Value")
			buf.WriteString("\n")
		}
	}
	buf.WriteString(mr.body.ToGoReturn())
	buf.WriteString("\n")
	return buf.String()
}

type MatchExpr struct {
	target Expr
	rules  []*MatchRule
}

func (me *MatchExpr) FType() FType {
	// return type must be the same for all rules, so I use first one.
	return me.rules[0].body.FType()
}

/*
At least one rule has "of" and not "_"
*/
func (me *MatchExpr) ruleHasContent() bool {
	for _, rule := range me.rules {
		if rule.pattern.caseId != "_" && rule.pattern.varName != "" && rule.pattern.varName != "_" {
			return true
		}
	}
	return false
}

func (me *MatchExpr) ToGoReturn() string {
	ut, ok := me.target.FType().(*FUnion)
	if !ok {
		panic("NYI, non union match expr")
	}

	var buf bytes.Buffer
	tmpV := UniqueTmpVarName()

	buf.WriteString("switch ")
	if me.ruleHasContent() {
		buf.WriteString(tmpV)
		buf.WriteString(" := ")
	}
	buf.WriteString("(")
	buf.WriteString(me.target.ToGo())

	buf.WriteString(").(type){\n")
	hasDefault := false
	for _, rule := range me.rules {
		if rule.pattern.caseId == "_" {
			hasDefault = true
		}
		buf.WriteString(rule.ToGo(ut.name, tmpV))
	}
	if !hasDefault {
		buf.WriteString("default:\npanic(\"Union pattern fail. Never reached here.\")\n")
	}

	buf.WriteString("}")
	return buf.String()
}

func (me *MatchExpr) ToGo() string {
	return wrapFunCall(me.FType(), me.ToGoReturn())
}

/*
	End of Expression.

	Begin Stmt related code.
*/

type Stmt interface {
	Node
	stmt()
	ToGo() string
}

func (*FuncDef) stmt()   {}
func (*Import) stmt()    {}
func (*Package) stmt()   {}
func (*ExprStmt) stmt()  {}
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

// Just contain 1 expression as Stmt.
type ExprStmt struct {
	Expr Expr
}

func (es *ExprStmt) ToGo() string {
	return es.Expr.ToGo()
}

type FuncDef struct {
	Name string
	// Unitはパース時点で0引数paramsに変換済みの想定
	Params []*Var
	Body   *Block
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
	retType := fd.Body.FType()
	if IsUnresolved(retType) {
		return retType
	}

	var fts []FType
	for _, arg := range fd.Params {
		fts = append(fts, arg.Type)
	}
	fts = append(fts, retType)
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
	buf.WriteString(fd.Body.ToGoReturn())
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
	return UnionCaseStructName(ud.Name, ud.Cases[index].Name)
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
func (IntOrString_I) IntOrString_Union(){}
func (IntOrString_B) IntOrString_Union(){}
*/
func (ud *UnionDef) buildCaseStructConformMethod(buf *bytes.Buffer) {
	method := ud.Name + "_Union(){}\n"
	for i := range ud.Cases {
		buf.WriteString("func (")
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
	contentTp := ud.Cases[index].Type
	if contentTp != FUnit {
		buf.WriteString("  Value ")
		buf.WriteString(contentTp.ToGo())
		buf.WriteString("\n")
	}
	buf.WriteString("}\n")
}

// New_IntOrString_I
func (ud *UnionDef) caseStructConstructorName(index int) string {
	return "New_" + ud.CaseStructName(index)
}

/*
func New_IntOrString_I(v int) IntOrString { return &IntOrString_I{v} }
*/
func (ud *UnionDef) buildCaseStructConstructorContent(buf *bytes.Buffer, index int) {
	buf.WriteString("func ")
	buf.WriteString(ud.caseStructConstructorName(index))
	buf.WriteString("(v ")
	buf.WriteString(ud.Cases[index].Type.ToGo())
	buf.WriteString(") ")
	buf.WriteString(ud.Name)
	buf.WriteString(" { return ")
	buf.WriteString(ud.CaseStructName(index))
	buf.WriteString("{v} }\n")
}

/*
No arg case constructor case.
In this case, folang regard no arg func as variable.
So the result must be following:

New_IntOrString_I IntOrString = &IntOrString_I{}
*/
func (ud *UnionDef) buildCaseStructConstructorAsVar(buf *bytes.Buffer, index int) {
	buf.WriteString("var ")
	buf.WriteString(ud.caseStructConstructorName(index))
	buf.WriteString(" ")
	buf.WriteString(ud.Name)
	buf.WriteString(" = ")
	buf.WriteString("")
	buf.WriteString(ud.CaseStructName(index))
	buf.WriteString("{}\n")
}

func (ud *UnionDef) buildCaseStructConstructor(buf *bytes.Buffer, index int) {
	contentTp := ud.Cases[index].Type
	if contentTp == FUnit {
		ud.buildCaseStructConstructorAsVar(buf, index)
	} else {
		ud.buildCaseStructConstructorContent(buf, index)
	}
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

func (ud *UnionDef) registerConstructor(resolver *TypeResolver) {
	utype := ud.UnionFType()

	for i, cs := range ud.Cases {
		if cs.Type == FUnit {
			resolver.AliasMap[cs.Name] = &Var{ud.caseStructConstructorName(i), utype}
		} else {
			tps := []FType{cs.Type, utype}
			ftype := &FFunc{tps}
			resolver.AliasMap[cs.Name] = &Var{ud.caseStructConstructorName(i), ftype}
		}
	}
}

func (ud *UnionDef) ResiterUnionTypeInfo(resolver *TypeResolver) {
	ud.registerConstructor(resolver)
	resolver.RegisterType(ud.Name, ud.UnionFType())
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
	case *Block:
		for _, stmt := range n.Stmts {
			Walk(stmt, f)
		}
		Walk(n.FinalExpr, f)
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
	case *ExprStmt:
		Walk(n.Expr, f)
	case *MatchExpr:
		Walk(n.target, f)
		for _, rule := range n.rules {
			// pattern is currently only identifier, so just walk body only.
			Walk(rule.body, f)
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
