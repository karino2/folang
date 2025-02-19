package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func tpname2tvtp(tvgen func() TypeVar, tpname string) frt.Tuple2[string, TypeVar] {
	tv := tvgen()
	return frt.NewTuple2(tpname, tv)
}

func transTypeVarFType(transTV func(TypeVar) FType, ftp FType) FType {
	recurse := (func(_r0 FType) FType { return transTypeVarFType(transTV, _r0) })
	switch _v254 := (ftp).(type) {
	case FType_FTypeVar:
		tv := _v254.Value
		return transTV(tv)
	case FType_FSlice:
		ts := _v254.Value
		et := recurse(ts.ElemType)
		return New_FType_FSlice(SliceType{ElemType: et})
	case FType_FTuple:
		ftup := _v254.Value
		nts := slice.Map(recurse, ftup.ElemTypes)
		return frt.Pipe(TupleType{ElemTypes: nts}, New_FType_FTuple)
	case FType_FFieldAccess:
		fa := _v254.Value
		nrec := recurse(fa.RecType)
		return frt.Pipe(FieldAccessType{RecType: nrec, FieldName: fa.FieldName}, faResolve)
	case FType_FFunc:
		fnt := _v254.Value
		nts := slice.Map(recurse, fnt.Targets)
		return frt.Pipe(FuncType{Targets: nts}, New_FType_FFunc)
	default:
		return ftp
	}
}

func tpReplaceOne(tvd TypeVarDict, tv TypeVar) FType {
	return frt.Pipe(tvdLookupNF(tvd, tv.Name), New_FType_FTypeVar)
}

func tpreplace(tvd TypeVarDict, ft FType) FType {
	return transTypeVarFType((func(_r0 TypeVar) FType { return tpReplaceOne(tvd, _r0) }), ft)
}

func GenFunc(ff FuncFactory, tvgen func() TypeVar) FuncType {
	tvd := frt.Pipe(slice.Map((func(_r0 string) frt.Tuple2[string, TypeVar] { return tpname2tvtp(tvgen, _r0) }), ff.Tparams), toTVDict)
	ntargets := slice.Map((func(_r0 FType) FType { return tpreplace(tvd, _r0) }), ff.Targets)
	return FuncType{Targets: ntargets}
}

func GenFuncVar(vname string, ff FuncFactory, tvgen func() TypeVar) Var {
	funct := GenFunc(ff, tvgen)
	ft := New_FType_FFunc(funct)
	return Var{Name: vname, Ftype: ft}
}

func genBuiltinFunCall(tvgen func() TypeVar, fname string, tpnames []string, targetTPs []FType, args []Expr) Expr {
	ff := FuncFactory{Tparams: tpnames, Targets: targetTPs}
	fvar := GenFuncVar(fname, ff, tvgen)
	return frt.Pipe(FunCall{TargetFunc: fvar, Args: args}, New_Expr_EFunCall)
}

func newTvf(name string) FType {
	return frt.Pipe(TypeVar{Name: name}, New_FType_FTypeVar)
}

type TypeDefCtx struct {
	tva         TypeVarAllocator
	insideTD    bool
	defined     TypeDict
	allocedDict SDict
}

type Resolver struct {
	eid EquivInfoDict
}

func newResolver() Resolver {
	neid := NewEquivInfoDict()
	return Resolver{eid: neid}
}

type TypeVarCtx struct {
	tva      TypeVarAllocator
	resolver Resolver
}

func newTypeVarCtx() TypeVarCtx {
	tva := NewTypeVarAllocator("_T")
	res := newResolver()
	return TypeVarCtx{tva: tva, resolver: res}
}

func tvcToTypeVarGen(tvc TypeVarCtx) func() TypeVar {
	return tvaToTypeVarGen(tvc.tva)
}

type ParseState struct {
	tkz        Tokenizer
	scope      Scope
	offsideCol []int
	tvc        TypeVarCtx
	tdctx      TypeDefCtx
}

func newParse(tkz Tokenizer, scope Scope, offCols []int, tvc TypeVarCtx, tdctx TypeDefCtx) ParseState {
	return ParseState{tkz: tkz, scope: scope, offsideCol: offCols, tvc: tvc, tdctx: tdctx}
}

func psWithTkz(org ParseState, tkz Tokenizer) ParseState {
	return newParse(tkz, org.scope, org.offsideCol, org.tvc, org.tdctx)
}

func psWithScope(org ParseState, nsc Scope) ParseState {
	return newParse(org.tkz, nsc, org.offsideCol, org.tvc, org.tdctx)
}

func psWithOffside(org ParseState, offs []int) ParseState {
	return newParse(org.tkz, org.scope, offs, org.tvc, org.tdctx)
}

func psWithTDCtx(org ParseState, ntdctx TypeDefCtx) ParseState {
	return newParse(org.tkz, org.scope, org.offsideCol, org.tvc, ntdctx)
}

func psWithTVCtx(org ParseState, ntvctx TypeVarCtx) ParseState {
	return newParse(org.tkz, org.scope, org.offsideCol, ntvctx, org.tdctx)
}

func initParse(src string) ParseState {
	tkz := newTkz(src)
	scope := NewScope()
	offside := ([]int{0})
	tva2 := NewTypeVarAllocator("_P")
	defdict := newTD()
	adict := newSD()
	tvctx := newTypeVarCtx()
	tdctx := TypeDefCtx{tva: tva2, insideTD: false, defined: defdict, allocedDict: adict}
	return newParse(tkz, scope, offside, tvctx, tdctx)
}

func psSetNewSrc(src string, ps ParseState) ParseState {
	tkz := newTkz(src)
	return psWithTkz(ps, tkz)
}

func psTypeVarGen(ps ParseState) func() TypeVar {
	return tvcToTypeVarGen(ps.tvc)
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

func psEnterTypeDef(ps ParseState) ParseState {
	old := ps.tdctx
	ntd := newTD()
	nald := newSD()
	ntdctx := TypeDefCtx{tva: old.tva, insideTD: true, defined: ntd, allocedDict: nald}
	tvaReset(ntdctx.tva)
	return psWithTDCtx(ps, ntdctx)
}

func psLeaveTypeDef(ps ParseState) ParseState {
	old := ps.tdctx
	ntdctx := TypeDefCtx{tva: old.tva, insideTD: false, defined: old.defined, allocedDict: old.allocedDict}
	return psWithTDCtx(ps, ntdctx)
}

func psInsideTypeDef(ps ParseState) bool {
	return ps.tdctx.insideTD
}

func tdctxTVFAlloc(tdctx TypeDefCtx, name string) FType {
	gen := tvaToTypeVarGen(tdctx.tva)
	tvar := gen()
	sdPut(tdctx.allocedDict, tvar.Name, name)
	return frt.Pipe(tvar, New_FType_FTypeVar)
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
	return frt.Pipe(newTypeVarCtx(), (func(_r0 TypeVarCtx) ParseState { return psWithTVCtx(ps, _r0) }))
}

func udToUt(ud UnionDef) UnionType {
	return UnionType(ud)
}

func udToFUt(ud UnionDef) FType {
	return frt.Pipe(udToUt(ud), New_FType_FUnion)
}

func csRegisterCtor(sc Scope, ud UnionDef, cas NameTypePair) {
	ctorName := csConstructorName(ud.Name, cas)
	ut := udToFUt(ud)
	v := (func() Var {
		switch (cas.Ftype).(type) {
		case FType_FUnit:
			return Var{Name: ctorName, Ftype: ut}
		default:
			tps := ([]FType{cas.Ftype, ut})
			funcTp := New_FType_FFunc(FuncType{Targets: tps})
			return Var{Name: ctorName, Ftype: funcTp}
		}
	})()
	scDefVar(sc, cas.Name, v)
}

func udRegisterCsCtors(sc Scope, ud UnionDef) {
	frt.PipeUnit(ud.Cases, (func(_r0 []NameTypePair) { slice.Iter((func(_r0 NameTypePair) { csRegisterCtor(sc, ud, _r0) }), _r0) }))
}

func piFullName(pi PackageInfo, name string) string {
	return frt.IfElse(frt.OpEqual(pi.Name, "_"), (func() string {
		return name
	}), (func() string {
		return ((pi.Name + ".") + name)
	}))
}

func piRegEType(pi PackageInfo, tname string) FType {
	fullName := piFullName(pi, tname)
	etype := New_FType_FExtType(fullName)
	etdPut(pi.TypeInfo, tname, fullName)
	return etype
}

func piRegFF(pi PackageInfo, fname string, ff FuncFactory, ps ParseState) ParseState {
	ffdPut(pi.FuncInfo, fname, ff)
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
	frt.PipeUnit(ffdKVs(pi.FuncInfo), (func(_r0 []frt.Tuple2[string, FuncFactory]) {
		slice.Iter((func(_r0 frt.Tuple2[string, FuncFactory]) { regFF(pi, sc, _r0) }), _r0)
	}))
	frt.PipeUnit(etdKVs(pi.TypeInfo), (func(_r0 []frt.Tuple2[string, string]) {
		slice.Iter((func(_r0 frt.Tuple2[string, string]) { regET(sc, _r0) }), _r0)
	}))
}

type BinOpInfo struct {
	Precedence int
	GoFuncName string
	IsBoolOp   bool
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
	secFncT := New_FType_FFunc(FuncType{Targets: ([]FType{t1type, t2type})})
	names := ([]string{t1name, t2name})
	tps := ([]FType{t1type, secFncT, t2type})
	args := ([]Expr{lhs, rhs})
	return genBuiltinFunCall(tvgen, "frt.Pipe", names, tps, args)
}

func newPipeCallUnit(tvgen func() TypeVar, lhs Expr, rhs Expr) Expr {
	t1name := "T1"
	t1type := newTvf(t1name)
	secFncT := New_FType_FFunc(FuncType{Targets: ([]FType{t1type, New_FType_FUnit})})
	names := ([]string{t1name})
	tps := ([]FType{t1type, secFncT, New_FType_FUnit})
	args := ([]Expr{lhs, rhs})
	return genBuiltinFunCall(tvgen, "frt.PipeUnit", names, tps, args)
}

func newPipeCall(tvgen func() TypeVar, lhs Expr, rhs Expr) Expr {
	rht := ExprToType(rhs)
	switch _v256 := (rht).(type) {
	case FType_FFunc:
		ft := _v256.Value
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
	rtype := frt.IfElse(binfo.IsBoolOp, (func() FType {
		return New_FType_FBool
	}), (func() FType {
		return ExprToType(rhs)
	}))
	return frt.Pipe(BinOpCall{Op: binfo.GoFuncName, Rtype: rtype, Lhs: lhs, Rhs: rhs}, New_Expr_EBinOpCall)
}

func newBinOpCall(tvgen func() TypeVar, tk TokenType, binfo BinOpInfo, lhs Expr, rhs Expr) Expr {
	switch (tk).(type) {
	case TokenType_PIPE:
		return newPipeCall(tvgen, lhs, rhs)
	case TokenType_EQ:
		return newEqNeq(tvgen, binfo.GoFuncName, lhs, rhs)
	case TokenType_BRACKET:
		return newEqNeq(tvgen, binfo.GoFuncName, lhs, rhs)
	default:
		return newBinOpNormal(binfo, lhs, rhs)
	}
}

func newFnTp(argType FType, retType FType) FType {
	tgs := ([]FType{argType, retType})
	return frt.Pipe(FuncType{Targets: tgs}, New_FType_FFunc)
}

func emptySS() []string {
	return []string{}
}

func newIfElseCall(tvgen func() TypeVar, cond Expr, tbody Block, fbody Block) Expr {
	ltbody := frt.Pipe(LazyBlock{Block: tbody}, New_Expr_ELazyBlock)
	lfbody := frt.Pipe(LazyBlock{Block: fbody}, New_Expr_ELazyBlock)
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
	ltbody := frt.Pipe(LazyBlock{Block: tbody}, New_Expr_ELazyBlock)
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

func rdToRecType(rd RecordDef) RecordType {
	return RecordType(rd)
}

func psRegRecDefToTDCtx(rd RecordDef, ps ParseState) {
	recT := rdToRecType(rd)
	scRegisterRecType(ps.scope, recT)
	tdPut(ps.tdctx.defined, recT.Name, New_FType_FRecord(recT))
}

func psRegUdToTDCtx(ud UnionDef, ps ParseState) {
	sc := ps.scope
	udRegisterCsCtors(sc, ud)
	fut := udToFUt(ud)
	scRegisterType(sc, ud.Name, fut)
	tdPut(ps.tdctx.defined, ud.Name, fut)
}

func transVNTPair(transV func(TypeVar) FType, ntp NameTypePair) NameTypePair {
	nt := transTypeVarFType(transV, ntp.Ftype)
	return NameTypePair{Name: ntp.Name, Ftype: nt}
}

func transDefStmt(transV func(TypeVar) FType, df DefStmt) DefStmt {
	switch _v260 := (df).(type) {
	case DefStmt_DRecordDef:
		rd := _v260.Value
		nfields := slice.Map((func(_r0 NameTypePair) NameTypePair { return transVNTPair(transV, _r0) }), rd.Fields)
		return frt.Pipe(RecordDef{Name: rd.Name, Fields: nfields}, New_DefStmt_DRecordDef)
	case DefStmt_DUnionDef:
		ud := _v260.Value
		ncases := slice.Map((func(_r0 NameTypePair) NameTypePair { return transVNTPair(transV, _r0) }), ud.Cases)
		return frt.Pipe(UnionDef{Name: ud.Name, Cases: ncases}, New_DefStmt_DUnionDef)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func transVByTDCtx(tdctx TypeDefCtx, tv TypeVar) FType {
	rname, ok := frt.Destr(sdLookup(tdctx.allocedDict, tv.Name))
	return frt.IfElse(ok, (func() FType {
		nt, ok2 := frt.Destr(tdLookup(tdctx.defined, rname))
		return frt.IfElse(ok2, (func() FType {
			return nt
		}), (func() FType {
			frt.Panic("Unresolved foward decl type")
			return nt
		}))
	}), (func() FType {
		return New_FType_FTypeVar(tv)
	}))
}

func resolveFwrdDecl(md MultipleDefs, ps ParseState) MultipleDefs {
	transV := (func(_r0 TypeVar) FType { return transVByTDCtx(ps.tdctx, _r0) })
	transD := (func(_r0 DefStmt) DefStmt { return transDefStmt(transV, _r0) })
	ndefs := slice.Map(transD, md.Defs)
	return MultipleDefs{Defs: ndefs}
}

func scRegDefStmtType(sc Scope, df DefStmt) {
	switch _v261 := (df).(type) {
	case DefStmt_DRecordDef:
		rd := _v261.Value
		frt.PipeUnit(rdToRecType(rd), (func(_r0 RecordType) { scRegisterRecType(sc, _r0) }))
	case DefStmt_DUnionDef:
		ud := _v261.Value
		udRegisterCsCtors(sc, ud)
		fut := udToFUt(ud)
		scRegisterType(sc, ud.Name, fut)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func psRegMdTypes(md MultipleDefs, ps ParseState) {
	slice.Iter((func(_r0 DefStmt) { scRegDefStmtType(ps.scope, _r0) }), md.Defs)
}
