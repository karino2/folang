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
func (*UnitVal) expr()       {}
func (*BoolLiteral) expr()   {}
func (*FunCall) expr()       {}
func (*FieldAccess) expr()   {}
func (*Var) expr()           {}
func (*RecordGen) expr()     {}
func (*Block) expr()         {}
func (*MatchExpr) expr()     {}
func (*SliceExpr) expr()     {}
func (*TupleExpr) expr()     {}

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

type BoolLiteral struct {
	Value bool
}

func (*BoolLiteral) FType() FType { return FBool }

func (s *BoolLiteral) ToGo() string { return fmt.Sprintf("%t", s.Value) }

type IntImm struct {
	Value int
}

func (*IntImm) FType() FType { return FInt }

func (s *IntImm) ToGo() string { return fmt.Sprintf("%d", s.Value) }

type UnitVal struct{}

func (*UnitVal) FType() FType { return FUnit }
func (*UnitVal) ToGo() string { return "" }

var gUnitVal = &UnitVal{}

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

type FieldAccess struct {
	targetName string
	targetType *FRecord
	fieldName  string
}

func (fa *FieldAccess) FType() FType { return fa.targetType.GetField(fa.fieldName).Type }
func (fa *FieldAccess) ToGo() string { return fmt.Sprintf("%s.%s", fa.targetName, fa.fieldName) }

type ResolvedTypeParam struct {
	Name         string
	ResolvedType FType // nil if not resoled.
}

type FunCall struct {
	Func       *Var
	Args       []Expr
	TypeParams []ResolvedTypeParam
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

	if len(fc.Args) > len(ftype.Args()) {
		panic("too many arguments")
	}

	// partial apply
	newtypes := ftype.Targets[len(fc.Args):]
	// For partial apply, Type parameters might be resolved.
	// But I don't know how to handle this here, so I ignore until I need it.
	return NewFFunc(newtypes...)
}

func (fc *FunCall) toGoPartialApply() string {
	var buf bytes.Buffer
	restArgTypes := fc.ArgTypes()[len(fc.Args):]
	buf.WriteString("(func (")
	for i, rt := range restArgTypes {
		if i != 0 {
			buf.WriteString(", ")
		}
		// arg name is _r0, _r1, _r2, ...
		buf.WriteString(fmt.Sprintf("_r%d ", i))
		buf.WriteString(rt.ToGo())

	}
	buf.WriteString(") ")
	rttype := fc.FuncType().ReturnType()
	if rttype == FUnit {
		buf.WriteString("{ ")
	} else {
		buf.WriteString(rttype.ToGo())
		buf.WriteString(" { return ")
	}
	buf.WriteString(fc.Func.Name)
	buf.WriteString("(")

	for i, arg := range fc.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.ToGo())
	}
	for i := range restArgTypes {
		buf.WriteString(fmt.Sprintf(", _r%d ", i))
	}
	buf.WriteString(") })")
	return buf.String()
}

func (fc *FunCall) writeTypeParam(buf *bytes.Buffer) {
	if len(fc.TypeParams) > 0 {
		firstTypeParam := true
		for _, tp := range fc.TypeParams {
			if tp.ResolvedType != nil {
				if firstTypeParam {
					buf.WriteString("[")
					firstTypeParam = false
				} else {
					buf.WriteString(", ")
				}
				buf.WriteString(tp.ResolvedType.ToGo())
			}
		}
		if !firstTypeParam {
			buf.WriteString("]")
		}
	}
}

func (fc *FunCall) ToGo() string {
	if len(fc.Args) > len(fc.ArgTypes()) {
		panic("too many argument")
	}

	if len(fc.Args) < len(fc.ArgTypes()) {
		return fc.toGoPartialApply()
	}
	var buf bytes.Buffer
	buf.WriteString(fc.Func.Name)
	fc.writeTypeParam(&buf)
	buf.WriteString("(")

	oneUnitArg := len(fc.Args) == 1 && fc.Args[0] == gUnitVal
	if !oneUnitArg {
		for i, arg := range fc.Args {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(arg.ToGo())
		}
	}
	buf.WriteString(")")
	return buf.String()
}

/*
If target is generic type, register that type in paramInfo as exprType.
*/
func resolveOneType(target FType, exprType FType, paramInfo map[string]FType) {
	switch gt := target.(type) {
	case *FParametrized:
		paramInfo[gt.name] = exprType
	case *FSlice:
		// only check one level of []T
		if pet, ok := gt.elemType.(*FParametrized); ok {
			// exprType must be slice to.
			est := exprType.(*FSlice)
			paramInfo[pet.name] = est.elemType
		}
	case *FTuple:
		etup := exprType.(*FTuple)
		// only check one level of T*U
		for i, tt := range gt.Elems {
			if ptt, ok := tt.(*FParametrized); ok {
				paramInfo[ptt.name] = etup.Elems[i]
			}
		}
	case *FFunc:
		// only check one level of T->U
		eft, ok := exprType.(*FFunc)
		if !ok {
			panic("can't infer from argtype to expr, NYI.")
		}
		for i, tt := range gt.Targets {
			if ptt, ok := tt.(*FParametrized); ok {
				paramInfo[ptt.name] = eft.Targets[i]
			}
		}
	default:
		return
	}
}

func (fc *FunCall) buildResolvedInfo() map[string]FType {
	tinfo := make(map[string]FType)
	for _, rt := range fc.TypeParams {
		if rt.ResolvedType != nil {
			tinfo[rt.Name] = rt.ResolvedType
		}
	}
	return tinfo
}

func updateType(old FType, tinfo map[string]FType) FType {
	switch grt := old.(type) {
	case *FParametrized:
		if nt, ok := tinfo[grt.name]; ok {
			return nt
		} else {
			return old
		}
	case *FSlice:
		// only check one level of []T
		if pet, ok := grt.elemType.(*FParametrized); ok {
			if nt, ok := tinfo[pet.name]; ok {
				return &FSlice{nt}
			} else {
				return old
			}
		} else {
			return old
		}
	case *FTuple:
		// only check one level of T*U
		var newEts []FType
		for _, et := range grt.Elems {
			if pet, ok := et.(*FParametrized); ok {
				if nt, ok := tinfo[pet.name]; ok {
					newEts = append(newEts, nt)
				} else {
					newEts = append(newEts, et)
				}
			} else {
				newEts = append(newEts, et)
			}
		}
		return &FTuple{newEts}

		// TODO: *FFunc support.
	default:
		return old
	}

}

func (fc *FunCall) buildNewTypeParams(tinfo map[string]FType) (resolved []ResolvedTypeParam, notResolved []string) {
	ft := fc.FuncType()
	for _, pname := range ft.TypeParams {
		if rt, ok := tinfo[pname]; ok {
			resolved = append(resolved, ResolvedTypeParam{pname, rt})
		} else {
			notResolved = append(notResolved, pname)
			resolved = append(resolved, ResolvedTypeParam{pname, nil})
		}
	}
	return
}

/*
Resolve type param from arguments.
*/
func (fc *FunCall) ResolveTypeParamByArgs() {
	tinfo := fc.buildResolvedInfo()
	ft := fc.FuncType()
	fargs := ft.Args()
	var newTypes []FType
	for i, at := range fargs {
		if i >= len(fc.Args) {
			newTypes = append(newTypes, at)
		} else {
			realT := fc.Args[i].FType()
			resolveOneType(at, realT, tinfo)
			newTypes = append(newTypes, realT)
		}
	}
	rt := ft.ReturnType()
	nt := updateType(rt, tinfo)
	newTypes = append(newTypes, nt)

	resolvedTypeParams, notResoledTypeParams := fc.buildNewTypeParams(tinfo)

	fc.Func = &Var{fc.Func.Name, &FFunc{newTypes, notResoledTypeParams}}
	fc.TypeParams = resolvedTypeParams
}

/*
T, int -> paramName=T, matchType=int, ok=true
[]T, []int -> paramName=T, matchType=int, ok=true
int*T, int*string -> paramName=T matchType=string ok=true, currently, there is no way to inform two match, so only either of one param is matched. NYI for both.
*/
func matchTypeParam(target FType, realType FType) (paramName string, matchType FType, matched bool) {
	switch tt := target.(type) {
	case *FParametrized:
		paramName = tt.name
		matchType = realType
		matched = true
		return
	case *FSlice:
		if elemPt, ok := tt.elemType.(*FParametrized); ok {
			realSlice := realType.(*FSlice)
			paramName = elemPt.name
			matchType = realSlice.elemType
			matched = true
			return
		}
	case *FTuple:
		// NYI: support both match case.
		for i, targetElemType := range tt.Elems {
			if elemPt, ok := targetElemType.(*FParametrized); ok {
				realTupleType := realType.(*FTuple)
				paramName = elemPt.name
				matchType = realTupleType.Elems[i]
				matched = true
				return
			}
		}
	}
	matched = false
	return
}

/*
Try to resolve like:
ss |> slice.Take 2

In this case, argType is the type of 'ss'.
f is 'slice.Take 2', which is FunCall.

For f, there is also another possibility like:

ss |> slice.Length

For this case, f is "&Var".
In this case, the result go code becomes

frt.Pipe(ss, slice.Length)

and resolved by go lang. So only the type of &Var is important.
*/
func resolveFuncTypeByArgType(f Expr, argType FType) Expr {
	switch fe := f.(type) {
	case *FunCall:
		/*
			FunCall, in this case, This FunCall must be partial apply.
			And something like
			slice.Take 2
			The only important type is whole expr type, not arg type or slice.Take type.
			And this must be: argTyp->X
		*/
		patType, ok := fe.FType().(*FFunc)
		if !ok {
			panic("func application but not func")
		}
		fargType := patType.Targets[0]

		if matchPname, matchType, ok := matchTypeParam(fargType, argType); ok {
			// can resolve.
			tinfo := fe.buildResolvedInfo()
			tinfo[matchPname] = matchType

			oldFtype := fe.FuncType() // head Func var type.
			var newTypes []FType
			for _, old := range oldFtype.Targets {
				newTypes = append(newTypes, updateType(old, tinfo))
			}
			resolvedTypeParams, notResoledTypeParams := fe.buildNewTypeParams(tinfo)

			fe.Func = &Var{fe.Func.Name, &FFunc{newTypes, notResoledTypeParams}}
			fe.TypeParams = resolvedTypeParams
			return fe
		}
		return f
	case *Var:
		// for this case, only need to update FType.
		// golang will resolve type parameter for this function automatically.
		patType, ok := fe.FType().(*FFunc)
		if !ok {
			panic("func application but not func(2)")
		}
		fargType := patType.Targets[0]

		if matchPname, matchType, ok := matchTypeParam(fargType, argType); ok {
			// can resolve.
			tinfo := make(map[string]FType)
			tinfo[matchPname] = matchType

			var newTypes []FType
			for _, old := range patType.Targets {
				newTypes = append(newTypes, updateType(old, tinfo))
			}
			var notResolved []string
			for _, tp := range patType.TypeParams {
				if _, ok := tinfo[tp]; !ok {
					notResolved = append(notResolved, tp)
				}
			}

			return &Var{fe.Name, &FFunc{newTypes, notResolved}}
		}
		return f
	default:
		return f
	}
}

func NewIfOnlyCall(cond Expr, tbody *Block) *FunCall {
	tbody.asFunc = true
	if tbody.ReturnType() != FUnit {
		panic("if only but not unit. compile error")
	}
	pvar := &Var{"frt.IfOnly", NewFFunc(FBool, NewFFunc(FUnit, FUnit), FUnit)}
	return &FunCall{pvar, []Expr{cond, tbody}, []ResolvedTypeParam{}}
}

func NewIfElseCall(cond Expr, tbody *Block, fbody *Block) *FunCall {
	var funcName string
	tbody.asFunc = true
	fbody.asFunc = true
	retType := tbody.ReturnType()
	if retType == FUnit {
		funcName = "frt.IfElseUnit"
	} else {
		funcName = "frt.IfElse"
	}
	pvar := &Var{funcName, NewFFunc(FBool, NewFFunc(FUnit, retType), NewFFunc(FUnit, retType), retType)}
	return &FunCall{pvar, []Expr{cond, tbody, fbody}, []ResolvedTypeParam{}}
}

func NewBinOpCall(btype TokenType, binfo binOpInfo, lhs Expr, rhs Expr) *FunCall {
	if btype == PIPE {
		// PIPE needs different inference pattern.
		return NewPipeCall(lhs, rhs)
	}

	// normal arithmetic.
	// all func like OpPlus has one type argument.
	// And type is T->T->(T|bool)
	// currently I just assume both side type is already resolved.
	// And return type is also the same. So use it as type parameter result.

	t1name := UniqueTmpTypeParamName()
	resolvedType := lhs.FType()
	var typeparams []ResolvedTypeParam
	typeparams = append(typeparams, ResolvedTypeParam{t1name, resolvedType})
	retType := resolvedType

	// It's better those login in binOpInfo, but just handle here for a while.
	switch btype {
	case EQ:
		retType = FBool
	case BRACKET:
		retType = FBool
	}

	pvar := &Var{binfo.goFuncName, &FFunc{[]FType{resolvedType, resolvedType, retType}, []string{t1name}}}
	return &FunCall{pvar, []Expr{lhs, rhs}, typeparams}

}

// frt.Pipe<T1, T2> : T1->(T1->T2)->T2
//
// lhs |> rhs
// rhs must be T1->T2 and we might resolve T1 or T2 by lhs T1 type.
func NewPipeCall(lhs Expr, rhs Expr) *FunCall {
	t1name := UniqueTmpTypeParamName()
	t2name := UniqueTmpTypeParamName()

	// resolve type by lhs ftype.
	rhs2 := resolveFuncTypeByArgType(rhs, lhs.FType())

	rt1 := lhs.FType()
	rt2 := rhs2.FType().(*FFunc).ReturnType()

	var typeparams []ResolvedTypeParam
	if rt2 == FUnit {
		typeparams = append(typeparams, ResolvedTypeParam{t1name, rt1})

		pvar := &Var{"frt.PipeUnit", &FFunc{[]FType{rt1, &FFunc{[]FType{rt1, FUnit}, []string{t1name}}, FUnit}, []string{t1name}}}
		return &FunCall{pvar, []Expr{lhs, rhs2}, typeparams}
	} else {
		typeparams = append(typeparams, ResolvedTypeParam{t1name, rt1})
		if _, ok := rt2.(*FParametrized); ok {
			typeparams = append(typeparams, ResolvedTypeParam{t2name, nil})
		} else {
			typeparams = append(typeparams, ResolvedTypeParam{t2name, rt2})
		}

		pvar := &Var{"frt.Pipe", &FFunc{[]FType{rt1, &FFunc{[]FType{rt1, rt2}, []string{t1name, t2name}}, rt2}, []string{t1name, t2name}}}
		return &FunCall{pvar, []Expr{lhs, rhs2}, typeparams}

	}
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
	scope     *Scope // block has own scope.
	asFunc    bool   // ToGo contains last () or not.
}

func NewBlock(retExpr Expr, stmts ...Stmt) *Block {
	return &Block{stmts, retExpr, nil, false}
}

func (b *Block) ReturnType() FType {
	return b.FinalExpr.FType()
}

func (b *Block) FType() FType {
	if b.asFunc {
		return NewFFunc(FUnit, b.ReturnType())
	} else {
		return b.ReturnType()
	}
}

func wrapFunc(returnType FType, goReturnBody string) string {
	var buf bytes.Buffer
	buf.WriteString("(func () ")
	buf.WriteString(returnType.ToGo())
	buf.WriteString(" {\n")
	buf.WriteString(goReturnBody)
	buf.WriteString("})")
	return buf.String()
}

func wrapFunCall(returnType FType, goReturnBody string) string {
	var buf bytes.Buffer
	buf.WriteString(wrapFunc(returnType, goReturnBody))
	buf.WriteString("()")
	return buf.String()
}

func (b *Block) ToGo() string {
	if b.asFunc {
		return wrapFunc(b.ReturnType(), b.ToGoReturn())
	} else {
		return wrapFunCall(b.FType(), b.ToGoReturn())
	}
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

func UniqueTmpTypeParamName() string {
	uniqueId++
	return fmt.Sprintf("_T%d", uniqueId)
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

type SliceExpr struct {
	exprs []Expr
}

func (se *SliceExpr) FType() FType {
	if len(se.exprs) == 0 {
		panic("empty slice, can't resolve type, illegal.")
	}
	return &FSlice{se.exprs[0].FType()}
}

func (se *SliceExpr) ToGo() string {
	var buf bytes.Buffer

	buf.WriteString("(")
	buf.WriteString(se.FType().ToGo())
	buf.WriteString("{")
	for i, e := range se.exprs {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(e.ToGo())
	}
	buf.WriteString("}")
	buf.WriteString(")")
	return buf.String()
}

type TupleExpr struct {
	Elems []Expr
}

// Support only two element for a while.
func (te *TupleExpr) FType() FType {
	if len(te.Elems) != 2 {
		panic("Only two element tuple support for a while.")
	}
	return NewFTuple(te.Elems[0].FType(), te.Elems[1].FType())
}

func (te *TupleExpr) ToGo() string {
	var buf bytes.Buffer

	buf.WriteString("frt.NewTuple2(")
	buf.WriteString(te.Elems[0].ToGo())
	buf.WriteString(",")
	buf.WriteString(te.Elems[1].ToGo())
	buf.WriteString(")")
	return buf.String()
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

func (*FuncDef) stmt()      {}
func (*LetVarDef) stmt()    {}
func (*Import) stmt()       {}
func (*Package) stmt()      {}
func (*ExprStmt) stmt()     {}
func (*RecordDef) stmt()    {}
func (*UnionDef) stmt()     {}
func (*MultipleDefs) stmt() {}
func (*PackageInfo) stmt()  {}

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

func varsToFTypes(vars []*Var) []FType {
	if len(vars) == 0 {
		return []FType{FUnit}
	}

	var fts []FType
	for _, arg := range vars {
		fts = append(fts, arg.Type)
	}
	return fts
}

func (fd *FuncDef) FuncFType() FType {
	retType := fd.Body.FType()
	if IsUnresolved(retType) {
		return retType
	}

	fts := varsToFTypes(fd.Params)
	fts = append(fts, retType)
	// type parameter NYI.
	return NewFFunc(fts...)
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
For 'and' type def, there are multiple type definition in one statement.
defs is either RecordDef or UnionDef.
*/
type MultipleDefs struct {
	defs []Stmt
}

func (md *MultipleDefs) ToGo() string {
	var buf bytes.Buffer
	for i, def := range md.defs {
		if i != 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(def.ToGo())
	}
	return buf.String()
}

/*
let a = expr

Name: "a"
Rhs: expr
*/
type LetVarDef struct {
	Name string
	Rhs  Expr
}

func (lvd *LetVarDef) ToGo() string {
	var buf bytes.Buffer
	buf.WriteString(lvd.Name)
	buf.WriteString(" := ")
	buf.WriteString(lvd.Rhs.ToGo())
	return buf.String()
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

func (ud *UnionDef) registerConstructor(scope *Scope) {
	utype := ud.UnionFType()

	for i, cs := range ud.Cases {
		if cs.Type == FUnit {
			scope.DefineVar(cs.Name, &Var{ud.caseStructConstructorName(i), utype})
		} else {
			tps := []FType{cs.Type, utype}
			ftype := NewFFunc(tps...)
			scope.DefineVar(cs.Name, &Var{ud.caseStructConstructorName(i), ftype})
		}
	}
}

func (ud *UnionDef) registerUnionTypeInfo(scope *Scope) {
	scope.typeMap[ud.Name] = ud.UnionFType()
}

func (ud *UnionDef) registerToScope(scope *Scope) {
	ud.registerConstructor(scope)
	ud.registerUnionTypeInfo(scope)
}

/*
External package info.
This emit no go code, but treat as dummy Stmt.
*/
type PackageInfo struct {
	name     string
	funcInfo map[string]*FFunc
	typeInfo map[string]*FExtType
}

func (*PackageInfo) ToGo() string { return "" }
func NewPackageInfo(name string) *PackageInfo {
	pi := &PackageInfo{name: name}
	pi.funcInfo = make(map[string]*FFunc)
	pi.typeInfo = make(map[string]*FExtType)
	return pi
}

func (pi *PackageInfo) registerExtType(name string) *FExtType {
	ret := &FExtType{pi.name + "." + name}
	pi.typeInfo[name] = ret
	return ret
}

func (pi *PackageInfo) registerToScope(scope *Scope) {
	for _, tp := range pi.typeInfo {
		fullName := tp.name
		scope.typeMap[fullName] = tp
	}

	for fname, ftp := range pi.funcInfo {
		fullName := pi.name + "." + fname
		scope.DefineVar(fullName, &Var{fullName, ftp})
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
	case *Block:
		for _, stmt := range n.Stmts {
			Walk(stmt, f)
		}
		Walk(n.FinalExpr, f)
	case *LetVarDef:
		Walk(n.Rhs, f)
	case *MultipleDefs:
		for _, stmt := range n.defs {
			Walk(stmt, f)
		}
	case *Import, *Package, *RecordDef, *UnionDef, *PackageInfo:
		// no-op
	// ここからexpr
	case *GoEval, *StringLiteral, *Var, *IntImm, *BoolLiteral, *FieldAccess:
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
