package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/dict"

func transTypeVarFType(transTV func(TypeVar) FType, ftp FType) FType {
	recurse := (func(_r0 FType) FType { return transTypeVarFType(transTV, _r0) })
	switch _v1 := (ftp).(type) {
	case FType_FTypeVar:
		tv := _v1.Value
		return transTV(tv)
	case FType_FSlice:
		ts := _v1.Value
		et := recurse(ts.ElemType)
		return New_FType_FSlice(SliceType{ElemType: et})
	case FType_FTuple:
		ftup := _v1.Value
		nts := slice.Map(recurse, ftup.ElemTypes)
		return frt.Pipe(TupleType{ElemTypes: nts}, New_FType_FTuple)
	case FType_FFieldAccess:
		fa := _v1.Value
		nrec := recurse(fa.RecType)
		return frt.Pipe(FieldAccessType{RecType: nrec, FieldName: fa.FieldName}, faResolve)
	case FType_FFunc:
		fnt := _v1.Value
		nts := slice.Map(recurse, fnt.Targets)
		return frt.Pipe(FuncType{Targets: nts}, New_FType_FFunc)
	case FType_FParamd:
		pt := _v1.Value
		nts := slice.Map(recurse, pt.Targs)
		return frt.Pipe(ParamdType{Name: pt.Name, Targs: nts}, New_FType_FParamd)
	default:
		return ftp
	}
}

func transOneVar(transTV func(TypeVar) FType, v Var) Var {
	ntp := transTypeVarFType(transTV, v.Ftype)
	return Var{Name: v.Name, Ftype: ntp}
}

func collectTVarFType(ft FType) []string {
	recurse := collectTVarFType
	switch _v2 := (ft).(type) {
	case FType_FTypeVar:
		tv := _v2.Value
		return ([]string{tv.Name})
	case FType_FSlice:
		ts := _v2.Value
		return recurse(ts.ElemType)
	case FType_FTuple:
		ftup := _v2.Value
		return slice.Collect(recurse, ftup.ElemTypes)
	case FType_FFieldAccess:
		fa := _v2.Value
		return recurse(fa.RecType)
	case FType_FFunc:
		fnt := _v2.Value
		return slice.Collect(recurse, fnt.Targets)
	default:
		return []string{}
	}
}

func tpReplaceOne(tdic dict.Dict[string, FType], tv TypeVar) FType {
	tp, ok := frt.Destr(dict.TryFind(tdic, tv.Name))
	return frt.IfElse(ok, (func() FType {
		return tp
	}), (func() FType {
		return New_FType_FTypeVar(tv)
	}))
}

func tpreplace(tdic dict.Dict[string, FType], ft FType) FType {
	return transTypeVarFType((func(_r0 TypeVar) FType { return tpReplaceOne(tdic, _r0) }), ft)
}

func emptyFtps() []FType {
	return []FType{}
}

func tpname2tvtp(tvgen func() TypeVar, slist []FType, i int, tpname string) frt.Tuple2[string, FType] {
	return frt.IfElse((slice.Len(slist) > i), (func() frt.Tuple2[string, FType] {
		item := slice.Item(i, slist)
		return frt.NewTuple2(tpname, item)
	}), (func() frt.Tuple2[string, FType] {
		tv := frt.Pipe(tvgen(), New_FType_FTypeVar)
		return frt.NewTuple2(tpname, tv)
	}))
}

func GenFunc(ff FuncFactory, stlist []FType, tvgen func() TypeVar) FuncType {
	frt.IfOnly((slice.Len(stlist) > slice.Len(ff.Tparams)), (func() {
		frt.Panic("Too many type specified.")
	}))
	tdic := frt.Pipe(slice.Mapi((func(_r0 int, _r1 string) frt.Tuple2[string, FType] { return tpname2tvtp(tvgen, stlist, _r0, _r1) }), ff.Tparams), dict.ToDict)
	ntargets := slice.Map((func(_r0 FType) FType { return tpreplace(tdic, _r0) }), ff.Targets)
	return FuncType{Targets: ntargets}
}

func GenFuncVar(vname string, ff FuncFactory, stlist []FType, tvgen func() TypeVar) VarRef {
	funct := GenFunc(ff, stlist, tvgen)
	ft := New_FType_FFunc(funct)
	v := Var{Name: vname, Ftype: ft}
	return frt.IfElse(slice.IsEmpty(stlist), (func() VarRef {
		return New_VarRef_VRVar(v)
	}), (func() VarRef {
		return New_VarRef_VRSVar(SpecVar{Var: v, SpecList: stlist})
	}))
}

func genBuiltinFunCall(tvgen func() TypeVar, fname string, tpnames []string, targetTPs []FType, args []Expr) Expr {
	ff := FuncFactory{Tparams: tpnames, Targets: targetTPs}
	fvar := GenFuncVar(fname, ff, emptyFtps(), tvgen)
	return frt.Pipe(FunCall{TargetFunc: fvar, Args: args}, New_Expr_EFunCall)
}

func newTvf(name string) FType {
	return frt.Pipe(TypeVar{Name: name}, New_FType_FTypeVar)
}

func GenType(tfd TypeFactoryData, targs []FType) FType {
	frt.IfOnly(frt.OpNotEqual(slice.Len(tfd.Tparams), slice.Len(targs)), (func() {
		frt.Panic("wrong type param num for instantiate.")
	}))
	return frt.Pipe(ParamdType{Name: tfd.Name, Targs: targs}, New_FType_FParamd)
}

type TypeDefCtx struct {
	tva         TypeVarAllocator
	insideTD    bool
	defined     dict.Dict[string, FType]
	allocedDict dict.Dict[string, string]
}

type EquivInfo struct {
	eset    EquivSet
	resType FType
}

type Resolver struct {
	eid dict.Dict[string, EquivInfo]
}

func newResolver() Resolver {
	neid := dict.New[string, EquivInfo]()
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
	defdict := dict.New[string, FType]()
	adict := dict.New[string, string]()
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
	ntd := dict.New[string, FType]()
	nald := dict.New[string, string]()
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
	dict.Add(tdctx.allocedDict, tvar.Name, name)
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

func psNextIs(expectTT TokenType, ps ParseState) bool {
	return frt.OpEqual(psNextTT(ps), expectTT)
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

func psIsNeighborLT(ps ParseState) bool {
	return tkzIsNeighborLT(ps.tkz)
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

func piRegEType(pi PackageInfo, tname string, tparams []string) TypeFactoryData {
	fullName := piFullName(pi, tname)
	tfd := TypeFactoryData{Name: fullName, Tparams: tparams}
	dict.Add(pi.TypeInfo, tname, tfd)
	return tfd
}

func scRegFunFac(sc Scope, fname string, ff FuncFactory) {
	scRegisterVarFac(sc, fname, (func(_r0 []FType, _r1 func() TypeVar) VarRef { return GenFuncVar(fname, ff, _r0, _r1) }))
}

func scRegTFData(sc Scope, tname string, tfd TypeFactoryData) {
	scRegisterTypeFac(sc, tname, (func(_r0 []FType) FType { return GenType(tfd, _r0) }))
}

func piRegFF(pi PackageInfo, fname string, ff FuncFactory, ps ParseState) ParseState {
	dict.Add(pi.FuncInfo, fname, ff)
	scRegFunFac(ps.scope, fname, ff)
	return ps
}

func regFF(pi PackageInfo, sc Scope, sff frt.Tuple2[string, FuncFactory]) {
	ffname, ff := frt.Destr(sff)
	fullName := piFullName(pi, ffname)
	scRegFunFac(sc, fullName, ff)
}

func regTF(pi PackageInfo, sc Scope, etp frt.Tuple2[string, TypeFactoryData]) {
	tfd := frt.Snd(etp)
	scRegTFData(sc, tfd.Name, tfd)
}

func piRegAll(pi PackageInfo, sc Scope) {
	frt.PipeUnit(dict.KVs(pi.FuncInfo), (func(_r0 []frt.Tuple2[string, FuncFactory]) {
		slice.Iter((func(_r0 frt.Tuple2[string, FuncFactory]) { regFF(pi, sc, _r0) }), _r0)
	}))
	frt.PipeUnit(dict.KVs(pi.TypeInfo), (func(_r0 []frt.Tuple2[string, TypeFactoryData]) {
		slice.Iter((func(_r0 frt.Tuple2[string, TypeFactoryData]) { regTF(pi, sc, _r0) }), _r0)
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
	switch _v3 := (rht).(type) {
	case FType_FFunc:
		ft := _v3.Value
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
	tps := ([]FType{New_FType_FBool, ft, New_FType_FUnit})
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
	dict.Add(ps.tdctx.defined, recT.Name, New_FType_FRecord(recT))
}

func psRegUdToTDCtx(ud UnionDef, ps ParseState) {
	sc := ps.scope
	udRegisterCsCtors(sc, ud)
	fut := udToFUt(ud)
	scRegisterType(sc, ud.Name, fut)
	dict.Add(ps.tdctx.defined, ud.Name, fut)
}

func transVNTPair(transV func(TypeVar) FType, ntp NameTypePair) NameTypePair {
	nt := transTypeVarFType(transV, ntp.Ftype)
	return NameTypePair{Name: ntp.Name, Ftype: nt}
}

func transDefStmt(transV func(TypeVar) FType, df DefStmt) DefStmt {
	switch _v4 := (df).(type) {
	case DefStmt_DRecordDef:
		rd := _v4.Value
		nfields := slice.Map((func(_r0 NameTypePair) NameTypePair { return transVNTPair(transV, _r0) }), rd.Fields)
		return frt.Pipe(RecordDef{Name: rd.Name, Fields: nfields}, New_DefStmt_DRecordDef)
	case DefStmt_DUnionDef:
		ud := _v4.Value
		ncases := slice.Map((func(_r0 NameTypePair) NameTypePair { return transVNTPair(transV, _r0) }), ud.Cases)
		return frt.Pipe(UnionDef{Name: ud.Name, Cases: ncases}, New_DefStmt_DUnionDef)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func transVByTDCtx(tdctx TypeDefCtx, tv TypeVar) FType {
	rname, ok := frt.Destr(dict.TryFind(tdctx.allocedDict, tv.Name))
	return frt.IfElse(ok, (func() FType {
		nt, ok2 := frt.Destr(dict.TryFind(tdctx.defined, rname))
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
	switch _v5 := (df).(type) {
	case DefStmt_DRecordDef:
		rd := _v5.Value
		frt.PipeUnit(rdToRecType(rd), (func(_r0 RecordType) { scRegisterRecType(sc, _r0) }))
	case DefStmt_DUnionDef:
		ud := _v5.Value
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
