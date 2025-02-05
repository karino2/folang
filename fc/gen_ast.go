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
func (Expr_LazyBlock) Expr_Union()      {}
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

type Expr_LazyBlock struct {
	Value LazyBlock
}

func New_Expr_LazyBlock(v LazyBlock) Expr { return Expr_LazyBlock{v} }

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
	stmts     []Stmt
	finalExpr Expr
}
type LazyBlock struct {
	stmts     []Stmt
	finalExpr Expr
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
func (Stmt_PackageInfo) Stmt_Union()  {}
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

type Stmt_PackageInfo struct {
	Value PackageInfo
}

func New_Stmt_PackageInfo(v PackageInfo) Stmt { return Stmt_PackageInfo{v} }

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
type PackageInfo struct {
	name     string
	funcInfo funcTypeDict
	typeInfo extTypeDict
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

func lblockReturnType(toT func(Expr) FType, lb LazyBlock) FType {
	return toT(lb.finalExpr)
}

func lblockToType(toT func(Expr) FType, lb LazyBlock) FType {
	rtype := lblockReturnType(toT, lb)
	return New_FType_FFunc(FuncType{targets: ([]FType{New_FType_FUnit, rtype})})
}

func blockReturnType(toT func(Expr) FType, block Block) FType {
	return toT(block.finalExpr)
}

func blockToType(toT func(Expr) FType, b Block) FType {
	return blockReturnType(toT, b)
}

func blockToExpr(block Block) Expr {
	return New_Expr_ReturnableExpr(New_ReturnableExpr_Block(block))
}

func meToType(toT func(Expr) FType, me MatchExpr) FType {
	frule := frt.Pipe(me.rules, slice.Head)
	return frt.Pipe(frt.Pipe(frule.body, blockToExpr), toT)
}

func returnableToType(toT func(Expr) FType, rexpr ReturnableExpr) FType {
	switch _v77 := (rexpr).(type) {
	case ReturnableExpr_Block:
		b := _v77.Value
		return blockToType(toT, b)
	case ReturnableExpr_MatchExpr:
		me := _v77.Value
		return meToType(toT, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func fcToFuncType(fc FunCall) FuncType {
	tfv := fc.targetFunc
	ft := tfv.ftype
	switch _v78 := (ft).(type) {
	case FType_FFunc:
		ft := _v78.Value
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
	switch _v79 := (expr).(type) {
	case Expr_GoEval:
		ge := _v79.Value
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
		fa := _v79.Value
		return faToType(fa)
	case Expr_Var:
		v := _v79.Value
		return v.ftype
	case Expr_RecordGen:
		rg := _v79.Value
		return New_FType_FRecord(rg.recordType)
	case Expr_LazyBlock:
		lb := _v79.Value
		return lblockToType(ExprToType, lb)
	case Expr_ReturnableExpr:
		re := _v79.Value
		return returnableToType(ExprToType, re)
	case Expr_FunCall:
		fc := _v79.Value
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

func buildReturn(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, stmts []Stmt, lastExpr Expr) string {
	stmtGos := frt.Pipe(slice.Map(sToGo, stmts), (func(_r0 []string) string { return strings.Concat("\n", _r0) }))
	lastGo := (func() string {
		switch _v80 := (lastExpr).(type) {
		case Expr_ReturnableExpr:
			re := _v80.Value
			return reToGoRet(re)
		default:
			mayReturn := frt.IfElse(frt.OpEqual(ExprToType(lastExpr), New_FType_FUnit), (func() string {
				return ""
			}), (func() string {
				return "return "
			}))
			lg := eToGo(lastExpr)
			return frt.OpPlus(mayReturn, lg)
		}
	})()
	return frt.OpPlus(frt.OpPlus(stmtGos, "\n"), lastGo)
}

func wrapFunc(toGo func(FType) string, rtype FType, goReturnBody string) string {
	b := buf.New()
	buf.Write(b, "(func () ")
	frt.PipeUnit(toGo(rtype), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, " {\n")
	buf.Write(b, goReturnBody)
	buf.Write(b, "})")
	return buf.String(b)
}

func wrapFunCall(toGo func(FType) string, rtype FType, goReturnBody string) string {
	wf := wrapFunc(toGo, rtype, goReturnBody)
	return frt.OpPlus(wf, "()")
}

func lbToGo(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, lb LazyBlock) string {
	returnBody := buildReturn(sToGo, eToGo, reToGoRet, lb.stmts, lb.finalExpr)
	rtype := ExprToType(lb.finalExpr)
	return wrapFunc(FTypeToGo, rtype, returnBody)
}

func blockToGoReturn(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, block Block) string {
	return buildReturn(sToGo, eToGo, reToGoRet, block.stmts, block.finalExpr)
}

func blockToGo(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, block Block) string {
	goRet := blockToGoReturn(sToGo, eToGo, reToGoRet, block)
	rtype := ExprToType(block.finalExpr)
	return wrapFunCall(FTypeToGo, rtype, goRet)
}

func mpToCaseHeader(uname string, mp MatchPattern, tmpVarName string) string {
	b := buf.New()
	buf.Write(b, "case ")
	frt.PipeUnit(unionCaseStructName(uname, mp.caseId), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ":\n")
	frt.IfOnly(frt.OpAnd(frt.OpNotEqual(mp.varName, "_"), frt.OpNotEqual(mp.varName, "")), (func() {
		buf.Write(b, mp.varName)
		buf.Write(b, " := ")
		buf.Write(b, tmpVarName)
		buf.Write(b, ".Value")
		buf.Write(b, "\n")
	}))
	return buf.String(b)
}

func mrToCase(btogRet func(Block) string, uname string, tmpVarName string, mr MatchRule) string {
	b := buf.New()
	mp := mr.pattern
	cheader := frt.IfElse(frt.OpEqual(mp.caseId, "_"), (func() string {
		return "default:\n"
	}), (func() string {
		return mpToCaseHeader(uname, mp, tmpVarName)
	}))
	buf.Write(b, cheader)
	frt.PipeUnit(btogRet(mr.body), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n")
	return buf.String(b)
}

func mrHasCaseVar(mr MatchRule) bool {
	pat := mr.pattern
	return frt.OpAnd(frt.OpAnd(frt.OpNotEqual(pat.caseId, "_"), frt.OpNotEqual(pat.varName, "")), frt.OpNotEqual(pat.varName, "_"))
}

func meHasCaseVar(me MatchExpr) bool {
	return slice.Forall(mrHasCaseVar, me.rules)
}

func mrIsDefault(mr MatchRule) bool {
	pat := mr.pattern
	return frt.OpEqual(pat.caseId, "_")
}

func meToGoReturn(toGo func(Expr) string, btogRet func(Block) string, me MatchExpr) string {
	ttype := ExprToType(me.target)
	uttype := ttype.(FType_FUnion).Value
	hasCaseVar := meHasCaseVar(me)
	hasDefault := slice.Forany(mrIsDefault, me.rules)
	tmpVarName := frt.IfElse(hasCaseVar, (func() string {
		return uniqueTmpVarName()
	}), (func() string {
		return ""
	}))
	mrtocase := (func(_r0 MatchRule) string { return mrToCase(btogRet, uttype.name, tmpVarName, _r0) })
	b := buf.New()
	buf.Write(b, "switch ")
	frt.IfOnly(hasCaseVar, (func() {
		buf.Write(b, tmpVarName)
		buf.Write(b, " := ")
	}))
	buf.Write(b, "(")
	frt.PipeUnit(toGo(me.target), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ").(type){\n")
	frt.PipeUnit(frt.Pipe(slice.Map(mrtocase, me.rules), (func(_r0 []string) string { return strings.Concat("", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	frt.IfOnly(frt.OpNot(hasDefault), (func() {
		buf.Write(b, "default:\npanic(\"Union pattern fail. Never reached here.\")\n")
	}))
	buf.Write(b, "}")
	return buf.String(b)
}

func meToExpr(me MatchExpr) Expr {
	return frt.Pipe(New_ReturnableExpr_MatchExpr(me), New_Expr_ReturnableExpr)
}

func meToGo(toGo func(Expr) string, btogRet func(Block) string, me MatchExpr) string {
	goret := meToGoReturn(toGo, btogRet, me)
	rtype := ExprToType(meToExpr(me))
	return wrapFunCall(FTypeToGo, rtype, goret)
}

func reToGoReturn(sToGo func(Stmt) string, eToGo func(Expr) string, rexpr ReturnableExpr) string {
	rtgr := (func(_r0 ReturnableExpr) string { return reToGoReturn(sToGo, eToGo, _r0) })
	btogoRet := (func(_r0 Block) string { return blockToGoReturn(sToGo, eToGo, rtgr, _r0) })
	switch _v81 := (rexpr).(type) {
	case ReturnableExpr_Block:
		b := _v81.Value
		return blockToGoReturn(sToGo, eToGo, rtgr, b)
	case ReturnableExpr_MatchExpr:
		me := _v81.Value
		return meToGoReturn(eToGo, btogoRet, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func reToGo(sToGo func(Stmt) string, eToGo func(Expr) string, rexpr ReturnableExpr) string {
	rtgr := (func(_r0 ReturnableExpr) string { return reToGoReturn(sToGo, eToGo, _r0) })
	btogRet := (func(_r0 Block) string { return blockToGoReturn(sToGo, eToGo, rtgr, _r0) })
	switch _v82 := (rexpr).(type) {
	case ReturnableExpr_Block:
		b := _v82.Value
		return blockToGo(sToGo, eToGo, rtgr, b)
	case ReturnableExpr_MatchExpr:
		me := _v82.Value
		return meToGo(eToGo, btogRet, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func ExprToGo(sToGo func(Stmt) string, expr Expr) string {
	eToGo := (func(_r0 Expr) string { return ExprToGo(sToGo, _r0) })
	switch _v83 := (expr).(type) {
	case Expr_BoolLiteral:
		b := _v83.Value
		return frt.Sprintf1("%t", b)
	case Expr_GoEval:
		ge := _v83.Value
		return ge.goStmt
	case Expr_StringLiteral:
		s := _v83.Value
		return frt.Sprintf1("\\\"%s\\\"", s)
	case Expr_IntImm:
		i := _v83.Value
		return frt.Sprintf1("%d", i)
	case Expr_Unit:
		return ""
	case Expr_FieldAccess:
		fa := _v83.Value
		return frt.OpPlus(frt.OpPlus(fa.targetName, "."), fa.fieldName)
	case Expr_Var:
		v := _v83.Value
		return v.name
	case Expr_RecordGen:
		rg := _v83.Value
		return rgToGo(eToGo, rg)
	case Expr_ReturnableExpr:
		re := _v83.Value
		return reToGo(sToGo, eToGo, re)
	case Expr_FunCall:
		return "NYI"
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
