package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/buf"

import "github.com/karino2/folang/pkg/strings"

type GoEval struct {
	goStmt  string
	typeArg FType
}

type Var struct {
	name  string
	ftype FType
}

type FieldAccess struct {
	targetName string
	targetType RecordType
	fieldName  string
}

type MatchPattern struct {
	caseId  string
	varName string
}

type Expr interface {
	Expr_Union()
}

func (Expr_GoEval) Expr_Union()         {}
func (Expr_StringLiteral) Expr_Union()  {}
func (Expr_IntImm) Expr_Union()         {}
func (Expr_Unit) Expr_Union()           {}
func (Expr_BoolLiteral) Expr_Union()    {}
func (Expr_FunCall) Expr_Union()        {}
func (Expr_FieldAccess) Expr_Union()    {}
func (Expr_Var) Expr_Union()            {}
func (Expr_RecordGen) Expr_Union()      {}
func (Expr_ReturnableExpr) Expr_Union() {}

type Expr_GoEval struct {
	Value GoEval
}

func New_Expr_GoEval(v GoEval) Expr { return Expr_GoEval{v} }

type Expr_StringLiteral struct {
	Value string
}

func New_Expr_StringLiteral(v string) Expr { return Expr_StringLiteral{v} }

type Expr_IntImm struct {
	Value int
}

func New_Expr_IntImm(v int) Expr { return Expr_IntImm{v} }

type Expr_Unit struct {
}

var New_Expr_Unit Expr = Expr_Unit{}

type Expr_BoolLiteral struct {
	Value bool
}

func New_Expr_BoolLiteral(v bool) Expr { return Expr_BoolLiteral{v} }

type Expr_FunCall struct {
	Value FunCall
}

func New_Expr_FunCall(v FunCall) Expr { return Expr_FunCall{v} }

type Expr_FieldAccess struct {
	Value FieldAccess
}

func New_Expr_FieldAccess(v FieldAccess) Expr { return Expr_FieldAccess{v} }

type Expr_Var struct {
	Value Var
}

func New_Expr_Var(v Var) Expr { return Expr_Var{v} }

type Expr_RecordGen struct {
	Value RecordGen
}

func New_Expr_RecordGen(v RecordGen) Expr { return Expr_RecordGen{v} }

type Expr_ReturnableExpr struct {
	Value ReturnableExpr
}

func New_Expr_ReturnableExpr(v ReturnableExpr) Expr { return Expr_ReturnableExpr{v} }

type FunCall struct {
	targetFunc Var
	args       []Expr
}
type RecordGen struct {
	fieldNames  []string
	fieldValues []Expr
	recordType  RecordType
}
type Block struct {
	stms      []Stmt
	finalExpr Expr
	asFunc    bool
}
type MatchRule struct {
	pattern MatchPattern
	body    Block
}
type MatchExpr struct {
	target Expr
	rules  []MatchRule
}
type ReturnableExpr interface {
	ReturnableExpr_Union()
}

func (ReturnableExpr_Block) ReturnableExpr_Union()     {}
func (ReturnableExpr_MatchExpr) ReturnableExpr_Union() {}

type ReturnableExpr_Block struct {
	Value Block
}

func New_ReturnableExpr_Block(v Block) ReturnableExpr { return ReturnableExpr_Block{v} }

type ReturnableExpr_MatchExpr struct {
	Value MatchExpr
}

func New_ReturnableExpr_MatchExpr(v MatchExpr) ReturnableExpr { return ReturnableExpr_MatchExpr{v} }

type Stmt interface {
	Stmt_Union()
}

func (Stmt_Import) Stmt_Union()       {}
func (Stmt_Package) Stmt_Union()      {}
func (Stmt_LetFuncDef) Stmt_Union()   {}
func (Stmt_LetVarDef) Stmt_Union()    {}
func (Stmt_ExprStmt) Stmt_Union()     {}
func (Stmt_DefStmt) Stmt_Union()      {}
func (Stmt_MultipleDefs) Stmt_Union() {}

type Stmt_Import struct {
	Value string
}

func New_Stmt_Import(v string) Stmt { return Stmt_Import{v} }

type Stmt_Package struct {
	Value string
}

func New_Stmt_Package(v string) Stmt { return Stmt_Package{v} }

type Stmt_LetFuncDef struct {
	Value LetFuncDef
}

func New_Stmt_LetFuncDef(v LetFuncDef) Stmt { return Stmt_LetFuncDef{v} }

type Stmt_LetVarDef struct {
	Value LetVarDef
}

func New_Stmt_LetVarDef(v LetVarDef) Stmt { return Stmt_LetVarDef{v} }

type Stmt_ExprStmt struct {
	Value Expr
}

func New_Stmt_ExprStmt(v Expr) Stmt { return Stmt_ExprStmt{v} }

type Stmt_DefStmt struct {
	Value DefStmt
}

func New_Stmt_DefStmt(v DefStmt) Stmt { return Stmt_DefStmt{v} }

type Stmt_MultipleDefs struct {
	Value MultipeDefs
}

func New_Stmt_MultipleDefs(v MultipeDefs) Stmt { return Stmt_MultipleDefs{v} }

type LetFuncDef struct {
	name   string
	params []Var
	body   Block
}
type LetVarDef struct {
	name string
	rhs  Expr
}
type RecordDef struct {
	name   string
	fields []NameTypePair
}
type UnionDef struct {
	name  string
	cases []NameTypePair
}
type DefStmt interface {
	DefStmt_Union()
}

func (DefStmt_RecordDef) DefStmt_Union() {}
func (DefStmt_UnionDef) DefStmt_Union()  {}

type DefStmt_RecordDef struct {
	Value RecordDef
}

func New_DefStmt_RecordDef(v RecordDef) DefStmt { return DefStmt_RecordDef{v} }

type DefStmt_UnionDef struct {
	Value UnionDef
}

func New_DefStmt_UnionDef(v UnionDef) DefStmt { return DefStmt_UnionDef{v} }

type MultipeDefs struct {
	defs []DefStmt
}

func faToType(fa FieldAccess) FType {
	rt := fa.targetType
	field := frGetField(rt, fa.fieldName)
	return field.ftype
}

func blockReturnType(toT func(Expr) FType, block Block) FType {
	return toT(block.finalExpr)
}

func blockToType(toT func(Expr) FType, b Block) FType {
	rtype := blockReturnType(toT, b)
	return frt.IfElse(b.asFunc, (func() FType {
		return New_FType_FFunc(FuncType{targets: ([]FType{New_FType_FUnit, rtype})})
	}), (func() FType {
		return rtype
	}))
}

func blockToExpr(block Block) Expr {
	return New_Expr_ReturnableExpr(New_ReturnableExpr_Block(block))
}

func meToType(toT func(Expr) FType, me MatchExpr) FType {
	frule := frt.Pipe(me.rules, slice.Head)
	return frt.Pipe(frt.Pipe(frule.body, blockToExpr), toT)
}

func returnableToType(toT func(Expr) FType, rexpr ReturnableExpr) FType {
	switch _v46 := (rexpr).(type) {
	case ReturnableExpr_Block:
		b := _v46.Value
		return blockToType(toT, b)
	case ReturnableExpr_MatchExpr:
		me := _v46.Value
		return meToType(toT, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func fcToFuncType(fc FunCall) FuncType {
	tfv := fc.targetFunc
	ft := tfv.ftype
	switch _v47 := (ft).(type) {
	case FType_FFunc:
		ft := _v47.Value
		return ft
	default:
		return FuncType{}
	}
}

func fcArgTypes(fc FunCall) []FType {
	return frt.Pipe(fcToFuncType(fc), fargs)
}

func fcToType(fc FunCall) FType {
	ft := fcToFuncType(fc)
	tlen := frt.Pipe(fargs(ft), slice.Length)
	alen := slice.Length(fc.args)
	return frt.IfElse(frt.OpEqual(alen, tlen), (func() FType {
		return freturn(ft)
	}), (func() FType {
		if alen > tlen {
			panic("too many arugments")
		}
		newts := slice.Skip(alen, ft.targets)
		return New_FType_FFunc(FuncType{targets: newts})
	}))
}

func ExprToType(expr Expr) FType {
	switch _v48 := (expr).(type) {
	case Expr_GoEval:
		ge := _v48.Value
		return ge.typeArg
	case Expr_StringLiteral:
		return New_FType_FString
	case Expr_IntImm:
		return New_FType_FInt
	case Expr_Unit:
		return New_FType_FUnit
	case Expr_BoolLiteral:
		return New_FType_FBool
	case Expr_FieldAccess:
		fa := _v48.Value
		return faToType(fa)
	case Expr_Var:
		v := _v48.Value
		return v.ftype
	case Expr_RecordGen:
		rg := _v48.Value
		return New_FType_FRecord(rg.recordType)
	case Expr_ReturnableExpr:
		re := _v48.Value
		return returnableToType(ExprToType, re)
	case Expr_FunCall:
		fc := _v48.Value
		return fcToType(fc)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func rgFVToGo(toGo func(Expr) string, fvPair frt.Tuple2[string, Expr]) string {
	fn := frt.Fst(fvPair)
	fv := frt.Snd(fvPair)
	fvGo := toGo(fv)
	return frt.OpPlus(frt.OpPlus(fn, ": "), fvGo)
}

func rgToGo(toGo func(Expr) string, rg RecordGen) string {
	rtype := rg.recordType
	b := buf.New()
	frt.PipeUnit(frStructName(rtype), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "{")
	fvGo := frt.Pipe(frt.Pipe(slice.Zip(rg.fieldNames, rg.fieldValues), (func(_r0 []frt.Tuple2[string, Expr]) []string {
		return slice.Map((func(_r0 frt.Tuple2[string, Expr]) string { return rgFVToGo(toGo, _r0) }), _r0)
	})), (func(_r0 []string) string { return strings.Concat(", ", _r0) }))
	buf.Write(b, fvGo)
	buf.Write(b, "}")
	return buf.String(b)
}

func ExprToGo(expr Expr) string {
	switch _v49 := (expr).(type) {
	case Expr_BoolLiteral:
		b := _v49.Value
		return frt.Sprintf1("%t", b)
	case Expr_GoEval:
		ge := _v49.Value
		return ge.goStmt
	case Expr_StringLiteral:
		s := _v49.Value
		return frt.Sprintf1("\"%s\"", s)
	case Expr_IntImm:
		i := _v49.Value
		return frt.Sprintf1("%d", i)
	case Expr_Unit:
		return ""
	case Expr_FieldAccess:
		fa := _v49.Value
		return frt.OpPlus(frt.OpPlus(fa.targetName, "."), fa.fieldName)
	case Expr_Var:
		v := _v49.Value
		return v.name
	case Expr_RecordGen:
		rg := _v49.Value
		return rgToGo(ExprToGo, rg)
	case Expr_ReturnableExpr:
		return "NYI"
	case Expr_FunCall:
		return "NYI"
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
