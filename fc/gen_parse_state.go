package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func tpname2tvtp(tvgen func() TypeVar, tpname string) frt.Tuple2[string, TypeVar] {
	tv := tvgen()
	return frt.NewTuple2(tpname, tv)
}

func transTypeVarFType(transTV func(TypeVar) FType, ftp FType) FType {
	recurse := (func(_r0 FType) FType { return transTypeVarFType(transTV, _r0) })
	switch _v231 := (ftp).(type) {
	case FType_FTypeVar:
		tv := _v231.Value
		return transTV(tv)
	case FType_FSlice:
		ts := _v231.Value
		et := recurse(ts.elemType)
		return New_FType_FSlice(SliceType{elemType: et})
	case FType_FTuple:
		ftup := _v231.Value
		nts := slice.Map(recurse, ftup.elemTypes)
		return frt.Pipe(TupleType{elemTypes: nts}, New_FType_FTuple)
	case FType_FFunc:
		fnt := _v231.Value
		nts := slice.Map(recurse, fnt.targets)
		return frt.Pipe(FuncType{targets: nts}, New_FType_FFunc)
	default:
		return ftp
	}
}

func tpReplaceOne(tvd TypeVarDict, tv TypeVar) FType {
	return frt.Pipe(tvdLookupNF(tvd, tv.name), New_FType_FTypeVar)
}

func tpreplace(tvd TypeVarDict, ft FType) FType {
	return transTypeVarFType((func(_r0 TypeVar) FType { return tpReplaceOne(tvd, _r0) }), ft)
}

func GenFunc(ff FuncFactory, tvgen func() TypeVar) FuncType {
	tvd := frt.Pipe(slice.Map((func(_r0 string) frt.Tuple2[string, TypeVar] { return tpname2tvtp(tvgen, _r0) }), ff.tparams), toTVDict)
	ntargets := slice.Map((func(_r0 FType) FType { return tpreplace(tvd, _r0) }), ff.targets)
	return FuncType{targets: ntargets}
}

func GenFuncVar(vname string, ff FuncFactory, tvgen func() TypeVar) Var {
	funct := GenFunc(ff, tvgen)
	ft := New_FType_FFunc(funct)
	return Var{name: vname, ftype: ft}
}

func genBuiltinFunCall(tvgen func() TypeVar, fname string, tpnames []string, targetTPs []FType, args []Expr) Expr {
	ff := FuncFactory{tparams: tpnames, targets: targetTPs}
	fvar := GenFuncVar(fname, ff, tvgen)
	return frt.Pipe(FunCall{targetFunc: fvar, args: args}, New_Expr_EFunCall)
}

func newTvf(name string) FType {
	return frt.Pipe(TypeVar{name: name}, New_FType_FTypeVar)
}

type ParseState struct {
	tkz        Tokenizer
	scope      Scope
	offsideCol []int
	tva        TypeVarAllocator
}

func newParse(tkz Tokenizer, scope Scope, offCols []int, tva TypeVarAllocator) ParseState {
	return ParseState{tkz: tkz, scope: scope, offsideCol: offCols, tva: tva}
}

func psWithTkz(org ParseState, tkz Tokenizer) ParseState {
	return newParse(tkz, org.scope, org.offsideCol, org.tva)
}

func psWithScope(org ParseState, nsc Scope) ParseState {
	return newParse(org.tkz, nsc, org.offsideCol, org.tva)
}

func psWithOffside(org ParseState, offs []int) ParseState {
	return newParse(org.tkz, org.scope, offs, org.tva)
}

func psTypeVarGen(ps ParseState) func() TypeVar {
	return tvaToTypeVarGen(ps.tva)
}

func psPushScope(org ParseState) ParseState {
	return frt.Pipe(newScope(org.scope), (func(_r0 Scope) ParseState { return psWithScope(org, _r0) }))
}

func psPopScope(org ParseState) ParseState {
	return frt.Pipe(popScope(org.scope), (func(_r0 Scope) ParseState { return psWithScope(org, _r0) }))
}

func psCurOffside(ps ParseState) int {
	return slice.Last(ps.offsideCol)
}

func psCurCol(ps ParseState) int {
	return ps.tkz.col
}

func psPushOffside(ps ParseState) ParseState {
	curCol := psCurCol(ps)
	frt.IfOnly((psCurOffside(ps) >= curCol), (func() {
		frt.Panic("Overrun offside rule")
	}))
	return frt.Pipe(slice.PushLast(curCol, ps.offsideCol), (func(_r0 []int) ParseState { return psWithOffside(ps, _r0) }))
}

func psPopOffside(ps ParseState) ParseState {
	return frt.Pipe(slice.PopLast(ps.offsideCol), (func(_r0 []int) ParseState { return psWithOffside(ps, _r0) }))
}

func initParse(src string) ParseState {
	tkz := newTkz(src)
	scope := NewScope()
	offside := ([]int{0})
	tva := NewTypeVarAllocator()
	return newParse(tkz, scope, offside, tva)
}

func psCurrent(ps ParseState) Token {
	return ps.tkz.current
}

func psCurrentTT(ps ParseState) TokenType {
	tk := psCurrent(ps)
	return tk.ttype
}

func psCurIs(expectTT TokenType, ps ParseState) bool {
	return frt.OpEqual(psCurrentTT(ps), expectTT)
}

func psNext(ps ParseState) ParseState {
	ntk := tkzNext(ps.tkz)
	return psWithTkz(ps, ntk)
}

func psNextTT(ps ParseState) TokenType {
	return frt.Pipe(psNext(ps), psCurrentTT)
}

func psNextNOL(ps ParseState) ParseState {
	ntk := tkzNextNOL(ps.tkz)
	return psWithTkz(ps, ntk)
}

func psSkipEOL(ps ParseState) ParseState {
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps), New_TokenType_EOL), (func() ParseState {
		return psNextNOL(ps)
	}), (func() ParseState {
		return ps
	}))
}

func psExpect(ttype TokenType, ps ParseState) {
	cur := psCurrent(ps)
	frt.IfOnly(frt.OpNotEqual(cur.ttype, ttype), (func() {
		frt.Panic("non expected token")
	}))

}

func psConsume(ttype TokenType, ps ParseState) ParseState {
	psExpect(ttype, ps)
	return psNext(ps)
}

func psIdentName(ps ParseState) string {
	psExpect(New_TokenType_IDENTIFIER, ps)
	cur := psCurrent(ps)
	return cur.stringVal
}

func psStringVal(ps ParseState) string {
	psExpect(New_TokenType_STRING, ps)
	cur := psCurrent(ps)
	return cur.stringVal
}

func psStrNx(f func(ParseState) string, ps ParseState) frt.Tuple2[ParseState, string] {
	s := f(ps)
	ps2 := psNext(ps)
	return frt.NewTuple2(ps2, s)
}

func psIdentNameNx(ps ParseState) frt.Tuple2[ParseState, string] {
	return psStrNx(psIdentName, ps)
}

func psStringValNx(ps ParseState) frt.Tuple2[ParseState, string] {
	return psStrNx(psStringVal, ps)
}

func psCurrentNx(ps ParseState) frt.Tuple2[ParseState, Token] {
	tk := psCurrent(ps)
	ps2 := psNext(ps)
	return frt.NewTuple2(ps2, tk)
}

func psCurrentTTNx(ps ParseState) frt.Tuple2[ParseState, TokenType] {
	tt := psCurrentTT(ps)
	ps2 := psNext(ps)
	return frt.NewTuple2(ps2, tt)
}

func psIdentNameNxL(ps ParseState) frt.Tuple2[ParseState, string] {
	return frt.Pipe(psIdentNameNx(ps), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return CnvL(psSkipEOL, _r0) }))
}

func psStringValNxL(ps ParseState) frt.Tuple2[ParseState, string] {
	return frt.Pipe(psStringValNx(ps), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return CnvL(psSkipEOL, _r0) }))
}

func psCurrentNxL(ps ParseState) frt.Tuple2[ParseState, Token] {
	return frt.Pipe(psCurrentNx(ps), (func(_r0 frt.Tuple2[ParseState, Token]) frt.Tuple2[ParseState, Token] { return CnvL(psSkipEOL, _r0) }))
}

func psCurrentTTNxL(ps ParseState) frt.Tuple2[ParseState, TokenType] {
	return frt.Pipe(psCurrentTTNx(ps), (func(_r0 frt.Tuple2[ParseState, TokenType]) frt.Tuple2[ParseState, TokenType] {
		return CnvL(psSkipEOL, _r0)
	}))
}

func psResetTmpCtx(ps ParseState) ParseState {
	resetUniqueTmpCounter()
	tvaReset(ps.tva)
	return ps
}

func udToUt(ud UnionDef) UnionType {
	return UnionType(ud)
}

func udToFUt(ud UnionDef) FType {
	return frt.Pipe(udToUt(ud), New_FType_FUnion)
}

func csRegisterCtor(sc Scope, ud UnionDef, cas NameTypePair) {
	ctorName := csConstructorName(ud.name, cas)
	ut := udToFUt(ud)
	v := (func() Var {
		switch (cas.ftype).(type) {
		case FType_FUnit:
			return Var{name: ctorName, ftype: ut}
		default:
			tps := ([]FType{cas.ftype, ut})
			funcTp := New_FType_FFunc(FuncType{targets: tps})
			return Var{name: ctorName, ftype: funcTp}
		}
	})()
	scDefVar(sc, cas.name, v)
}

func udRegisterCsCtors(sc Scope, ud UnionDef) {
	frt.PipeUnit(ud.cases, (func(_r0 []NameTypePair) { slice.Iter((func(_r0 NameTypePair) { csRegisterCtor(sc, ud, _r0) }), _r0) }))
}

func udRegisterToScope(sc Scope, ud UnionDef) {
	udRegisterCsCtors(sc, ud)
	frt.PipeUnit(udToFUt(ud), (func(_r0 FType) { scRegisterType(sc, ud.name, _r0) }))
}

func piFullName(pi PackageInfo, name string) string {
	return frt.IfElse(frt.OpEqual(pi.name, "_"), (func() string {
		return name
	}), (func() string {
		return ((pi.name + ".") + name)
	}))
}

func piRegEType(pi PackageInfo, tname string) FType {
	fullName := piFullName(pi, tname)
	etype := New_FType_FExtType(fullName)
	etdPut(pi.typeInfo, tname, fullName)
	return etype
}

func piRegFF(pi PackageInfo, fname string, ff FuncFactory, ps ParseState) ParseState {
	ffdPut(pi.funcInfo, fname, ff)
	scRegisterVarFac(ps.scope, fname, (func(_r0 func() TypeVar) Var { return GenFuncVar(fname, ff, _r0) }))
	return ps
}

func regFF(pi PackageInfo, sc Scope, sff frt.Tuple2[string, FuncFactory]) {
	ffname, ff := frt.Destr(sff)
	fullName := piFullName(pi, ffname)
	scRegisterVarFac(sc, fullName, (func(_r0 func() TypeVar) Var { return GenFuncVar(fullName, ff, _r0) }))
}

func regET(sc Scope, etp frt.Tuple2[string, string]) {
	fullName := frt.Snd(etp)
	frt.PipeUnit(New_FType_FExtType(fullName), (func(_r0 FType) { scRegisterType(sc, fullName, _r0) }))
}

func piRegAll(pi PackageInfo, sc Scope) {
	frt.PipeUnit(ffdKVs(pi.funcInfo), (func(_r0 []frt.Tuple2[string, FuncFactory]) {
		slice.Iter((func(_r0 frt.Tuple2[string, FuncFactory]) { regFF(pi, sc, _r0) }), _r0)
	}))
	frt.PipeUnit(etdKVs(pi.typeInfo), (func(_r0 []frt.Tuple2[string, string]) {
		slice.Iter((func(_r0 frt.Tuple2[string, string]) { regET(sc, _r0) }), _r0)
	}))
}

type BinOpInfo struct {
	precedence int
	goFuncName string
	isBoolOp   bool
}

func newEqNeq(tvgen func() TypeVar, goFname string, lhs Expr, rhs Expr) Expr {
	t1name := "T1"
	t1tp := newTvf(t1name)
	names := ([]string{t1name})
	tps := ([]FType{t1tp, t1tp, New_FType_FBool})
	args := ([]Expr{lhs, rhs})
	return genBuiltinFunCall(tvgen, goFname, names, tps, args)
}

func newPipeCallNormal(tvgen func() TypeVar, lhs Expr, rhs Expr) Expr {
	t1name := "T1"
	t1type := newTvf(t1name)
	t2name := "T2"
	t2type := newTvf(t2name)
	secFncT := New_FType_FFunc(FuncType{targets: ([]FType{t1type, t2type})})
	names := ([]string{t1name, t2name})
	tps := ([]FType{t1type, secFncT, t2type})
	args := ([]Expr{lhs, rhs})
	return genBuiltinFunCall(tvgen, "frt.Pipe", names, tps, args)
}

func newPipeCallUnit(tvgen func() TypeVar, lhs Expr, rhs Expr) Expr {
	t1name := "T1"
	t1type := newTvf(t1name)
	secFncT := New_FType_FFunc(FuncType{targets: ([]FType{t1type, New_FType_FUnit})})
	names := ([]string{t1name})
	tps := ([]FType{t1type, secFncT, New_FType_FUnit})
	args := ([]Expr{lhs, rhs})
	return genBuiltinFunCall(tvgen, "frt.PipeUnit", names, tps, args)
}

func newPipeCall(tvgen func() TypeVar, lhs Expr, rhs Expr) Expr {
	rht := ExprToType(rhs)
	switch _v233 := (rht).(type) {
	case FType_FFunc:
		ft := _v233.Value
		switch (freturn(ft)).(type) {
		case FType_FUnit:
			return newPipeCallUnit(tvgen, lhs, rhs)
		default:
			return newPipeCallNormal(tvgen, lhs, rhs)
		}
	default:
		return newPipeCallNormal(tvgen, lhs, rhs)
	}
}

func newBinOpNormal(binfo BinOpInfo, lhs Expr, rhs Expr) Expr {
	rtype := frt.IfElse(binfo.isBoolOp, (func() FType {
		return New_FType_FBool
	}), (func() FType {
		return ExprToType(rhs)
	}))
	return frt.Pipe(BinOpCall{op: binfo.goFuncName, rtype: rtype, lhs: lhs, rhs: rhs}, New_Expr_EBinOpCall)
}

func newBinOpCall(tvgen func() TypeVar, tk TokenType, binfo BinOpInfo, lhs Expr, rhs Expr) Expr {
	switch (tk).(type) {
	case TokenType_PIPE:
		return newPipeCall(tvgen, lhs, rhs)
	case TokenType_EQ:
		return newEqNeq(tvgen, binfo.goFuncName, lhs, rhs)
	case TokenType_BRACKET:
		return newEqNeq(tvgen, binfo.goFuncName, lhs, rhs)
	default:
		return newBinOpNormal(binfo, lhs, rhs)
	}
}

func newFnTp(argType FType, retType FType) FType {
	tgs := ([]FType{argType, retType})
	return frt.Pipe(FuncType{targets: tgs}, New_FType_FFunc)
}

func emptySS() []string {
	return []string{}
}

func newIfElseCall(tvgen func() TypeVar, cond Expr, tbody Block, fbody Block) Expr {
	ltbody := frt.Pipe(LazyBlock{block: tbody}, New_Expr_ELazyBlock)
	lfbody := frt.Pipe(LazyBlock{block: fbody}, New_Expr_ELazyBlock)
	retType := blockReturnType(ExprToType, tbody)
	fname := (func() string {
		switch (retType).(type) {
		case FType_FUnit:
			return "frt.IfElseUnit"
		default:
			return "frt.IfElse"
		}
	})()
	emptyS := emptySS()
	args := ([]Expr{cond, ltbody, lfbody})
	ft := newFnTp(New_FType_FUnit, retType)
	tps := ([]FType{New_FType_FBool, ft, ft, retType})
	return genBuiltinFunCall(tvgen, fname, emptyS, tps, args)
}

func newIfOnlyCall(tvgen func() TypeVar, cond Expr, tbody Block) Expr {
	ltbody := frt.Pipe(LazyBlock{block: tbody}, New_Expr_ELazyBlock)
	emptyS := emptySS()
	args := ([]Expr{cond, ltbody})
	ft := newFnTp(New_FType_FUnit, New_FType_FUnit)
	tps := ([]FType{New_FType_FBool, ft, ft, New_FType_FUnit})
	return genBuiltinFunCall(tvgen, "frt.IfOnly", emptyS, tps, args)
}

func newUnaryNotCall(tvgen func() TypeVar, cond Expr) Expr {
	emptyS := emptySS()
	args := ([]Expr{cond})
	tps := ([]FType{New_FType_FBool, New_FType_FBool})
	return genBuiltinFunCall(tvgen, "frt.OpNot", emptyS, tps, args)
}
