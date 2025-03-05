package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/dict"

type ScopeDict struct {
	VarFacMap  dict.Dict[string, func([]FType, func() TypeVar) VarRef]
	RecFacMap  dict.Dict[string, RecordFactory]
	TypeFacMap dict.Dict[string, func([]FType) FType]
}

func NewScopeDict() ScopeDict {
	fvm := dict.New[string, func([]FType, func() TypeVar) VarRef]()
	rfm := dict.New[string, RecordFactory]()
	tfm := dict.New[string, func([]FType) FType]()
	return ScopeDict{VarFacMap: fvm, RecFacMap: rfm, TypeFacMap: tfm}
}

func NewScope0() Scope {
	sd := NewScopeDict()
	return NewScopeImpl0(sd)
}

func NewScope(parent Scope) Scope {
	sd := NewScopeDict()
	return NewScopeImpl(sd, parent)
}

func vToVarFac[T0 any, T1 any](v Var, tlist T0, tvgen T1) VarRef {
	return New_VarRef_VRVar(v)
}

func scDefVar(s Scope, name string, v Var) {
	sdic := SCSDict(s)
	dict.Add(sdic.VarFacMap, name, (func(_r0 []FType, _r1 func() TypeVar) VarRef { return vToVarFac(v, _r0, _r1) }))
}

func scRegisterVarFac(s Scope, name string, fac func([]FType, func() TypeVar) VarRef) {
	sdic := SCSDict(s)
	dict.Add(sdic.VarFacMap, name, fac)
}

func emptyVarFac(tlist []FType, tgen func() TypeVar) VarRef {
	PanicNow("should never called")
	return frt.Empty[VarRef]()
}

func scLookupVarFac(s Scope, name string) frt.Tuple2[func([]FType, func() TypeVar) VarRef, bool] {
	sd := SCSDict(s)
	vfac, ok := frt.Destr2(dict.TryFind(sd.VarFacMap, name))
	return frt.IfElse(ok, (func() frt.Tuple2[func([]FType, func() TypeVar) VarRef, bool] {
		return frt.NewTuple2(vfac, ok)
	}), (func() frt.Tuple2[func([]FType, func() TypeVar) VarRef, bool] {
		return frt.IfElse(SCHasParent(s), (func() frt.Tuple2[func([]FType, func() TypeVar) VarRef, bool] {
			return scLookupVarFac(SCParent(s), name)
		}), (func() frt.Tuple2[func([]FType, func() TypeVar) VarRef, bool] {
			return frt.NewTuple2(emptyVarFac, false)
		}))
	}))
}

func scRegisterTypeFac(s Scope, name string, fac func([]FType) FType) {
	sdic := SCSDict(s)
	dict.Add(sdic.TypeFacMap, name, fac)
}

func ftToTypeFac[T0 any](ft T0, tlist []FType) T0 {
	return ft
}

func scRegisterType(s Scope, name string, ftype FType) {
	scRegisterTypeFac(s, name, (func(_r0 []FType) FType { return ftToTypeFac(ftype, _r0) }))
}

func scRegisterRecFac(s Scope, name string, fac RecordFactory) {
	sdic := SCSDict(s)
	dict.Add(sdic.RecFacMap, name, fac)
	dict.Add(sdic.TypeFacMap, name, (func(_r0 []FType) FType { return GenRecordFType(fac, _r0) }))
}

func scLookupRecFacCur(s Scope, fieldNames []string) frt.Tuple2[RecordFactory, bool] {
	sdic := SCSDict(s)
	return frt.Pipe(dict.Values(sdic.RecFacMap), (func(_r0 []RecordFactory) frt.Tuple2[RecordFactory, bool] {
		return slice.TryFind((func(_r0 RecordFactory) bool { return recFacMatch(fieldNames, _r0) }), _r0)
	}))
}

func emptyRecFac() RecordFactory {
	return frt.Empty[RecordFactory]()
}

func scLookupRecFac(s Scope, fieldNames []string) frt.Tuple2[RecordFactory, bool] {
	rfac, ok := frt.Destr2(scLookupRecFacCur(s, fieldNames))
	return frt.IfElse(ok, (func() frt.Tuple2[RecordFactory, bool] {
		return frt.NewTuple2(rfac, ok)
	}), (func() frt.Tuple2[RecordFactory, bool] {
		return frt.IfElse(SCHasParent(s), (func() frt.Tuple2[RecordFactory, bool] {
			return scLookupRecFac(SCParent(s), fieldNames)
		}), (func() frt.Tuple2[RecordFactory, bool] {
			return frt.NewTuple2(emptyRecFac(), false)
		}))
	}))
}

func scLookupRecFacByName(s Scope, name string) frt.Tuple2[RecordFactory, bool] {
	sd := SCSDict(s)
	rfac, ok := frt.Destr2(dict.TryFind(sd.RecFacMap, name))
	return frt.IfElse(ok, (func() frt.Tuple2[RecordFactory, bool] {
		return frt.NewTuple2(rfac, ok)
	}), (func() frt.Tuple2[RecordFactory, bool] {
		return frt.IfElse(SCHasParent(s), (func() frt.Tuple2[RecordFactory, bool] {
			return scLookupRecFacByName(SCParent(s), name)
		}), (func() frt.Tuple2[RecordFactory, bool] {
			return frt.NewTuple2(emptyRecFac(), false)
		}))
	}))
}

func scLookupTypeFac(s Scope, name string) frt.Tuple2[func([]FType) FType, bool] {
	sd := SCSDict(s)
	rec, ok := frt.Destr2(dict.TryFind(sd.TypeFacMap, name))
	return frt.IfElse(ok, (func() frt.Tuple2[func([]FType) FType, bool] {
		return frt.NewTuple2(rec, ok)
	}), (func() frt.Tuple2[func([]FType) FType, bool] {
		return frt.IfElse(SCHasParent(s), (func() frt.Tuple2[func([]FType) FType, bool] {
			return scLookupTypeFac(SCParent(s), name)
		}), (func() frt.Tuple2[func([]FType) FType, bool] {
			empty := (func(_r0 []FType) FType { return ftToTypeFac(New_FType_FUnit, _r0) })
			return frt.NewTuple2(empty, false)
		}))
	}))
}

type EquivSet struct {
	Dict dict.Dict[string, bool]
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

func MapL[T0 any, T1 any, T2 any](fn func(T0) T1, tup frt.Tuple2[T0, T2]) frt.Tuple2[T1, T2] {
	nl := frt.Pipe(frt.Fst(tup), fn)
	return frt.NewTuple2(nl, frt.Snd(tup))
}

func MapR[T0 any, T1 any, T2 any](fn func(T0) T1, tup frt.Tuple2[T2, T0]) frt.Tuple2[T2, T1] {
	nr := frt.Pipe(frt.Snd(tup), fn)
	return frt.NewTuple2(frt.Fst(tup), nr)
}

func PairL[T0 any, T1 any](l T0, r T1) frt.Tuple2[T0, T1] {
	return frt.NewTuple2(l, r)
}

func PairR[T0 any, T1 any](r T0, l T1) frt.Tuple2[T1, T0] {
	return frt.NewTuple2(l, r)
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
	scope := NewScope0()
	offside := ([]int{0})
	tva2 := NewTypeVarAllocator("_P")
	defdict := dict.New[string, FType]()
	adict := dict.New[string, string]()
	tvctx := newTypeVarCtx()
	tdctx := TypeDefCtx{tva: tva2, insideTD: false, defined: defdict, allocedDict: adict}
	return newParse(tkz, scope, offside, tvctx, tdctx)
}

func psPanic(ps ParseState, msg string) {
	tkzPanic(ps.tkz, msg)
}

func psForErrMsg(ps ParseState) {
	SetLastTkz(ps.tkz)
}

func psSetNewSrc(src string, ps ParseState) ParseState {
	tkz := newTkz(src)
	return psWithTkz(ps, tkz)
}

func psTypeVarGen(ps ParseState) func() TypeVar {
	return tvcToTypeVarGen(ps.tvc)
}

func psPushScope(org ParseState) ParseState {
	return frt.Pipe(NewScope(org.scope), (func(_r0 Scope) ParseState { return psWithScope(org, _r0) }))
}

func popScope(sc Scope) Scope {
	return SCParent(sc)
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
		psPanic(ps, "Overrun offside rule")
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

func psCurIsNot(expectTT TokenType, ps ParseState) bool {
	return frt.OpNotEqual(psCurrentTT(ps), expectTT)
}

func psNext(ps ParseState) ParseState {
	ntk := tkzNext(ps.tkz)
	return psWithTkz(ps, ntk)
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

func psExpectMsg(ttype TokenType, ps ParseState, msg string) {
	frt.IfOnly(frt.OpNot(psCurIs(ttype, ps)), (func() {
		psPanic(ps, msg)
	}))

}

func psUnexpect(ttype TokenType, ps ParseState, msg string) {
	frt.IfOnly(psCurIs(ttype, ps), (func() {
		psPanic(ps, msg)
	}))

}

func psExpect(ttype TokenType, ps ParseState) {
	psExpectMsg(ttype, ps, "non expected token")
}

func psConsume(ttype TokenType, ps ParseState) ParseState {
	psExpect(ttype, ps)
	return psNext(ps)
}

func psMulConsume(ttypes []TokenType, ps ParseState) ParseState {
	consOne := func(p ParseState, tt TokenType) ParseState {
		return psConsume(tt, p)
	}
	return slice.Fold(consOne, ps, ttypes)
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

func psIdentNameNxL(ps ParseState) frt.Tuple2[ParseState, string] {
	return frt.Pipe(psIdentNameNx(ps), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return MapL(psSkipEOL, _r0) }))
}

func psStringValNxL(ps ParseState) frt.Tuple2[ParseState, string] {
	return frt.Pipe(psStringValNx(ps), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return MapL(psSkipEOL, _r0) }))
}

func psResetTmpCtx(ps ParseState) ParseState {
	resetUniqueTmpCounter()
	return frt.Pipe(newTypeVarCtx(), (func(_r0 TypeVarCtx) ParseState { return psWithTVCtx(ps, _r0) }))
}

func psIsNeighborLT(ps ParseState) bool {
	return tkzIsNeighborLT(ps.tkz)
}

func ParseSepList[T0 any](one func(ParseState) frt.Tuple2[ParseState, T0], sep TokenType, ps ParseState) frt.Tuple2[ParseState, []T0] {
	endPred := (func(_r0 ParseState) bool { return psCurIsNot(sep, _r0) })
	next := (func(_r0 ParseState) ParseState { return psConsume(sep, _r0) })
	return ParseList2(one, endPred, next, ps)
}

func scRegFunFac(sc Scope, fname string, ff FuncFactory) {
	scRegisterVarFac(sc, fname, (func(_r0 []FType, _r1 func() TypeVar) VarRef { return GenFuncVar(fname, ff, _r0, _r1) }))
}

func udToUtOnly(ud UnionDef) UnionType {
	return UnionType{Name: ud.Name}
}

func csRegisterCtor(sc Scope, ud UnionDef, cas NameTypePair) {
	ctorName := csConstructorName(ud.Name, cas)
	frt.IfElseUnit(csIsVar(ud.Tparams, cas), (func() {
		ut := frt.Pipe(udToUtOnly(ud), New_FType_FUnion)
		frt.PipeUnit(Var{Name: ctorName, Ftype: ut}, (func(_r0 Var) { scDefVar(sc, cas.Name, _r0) }))
	}), (func() {
		targs := slice.Map(newTvf, ud.Tparams)
		uf := udToUniFac(ud)
		ut := frt.Pipe(GenUnionType(uf, targs), New_FType_FUnion)
		tps := ([]FType{cas.Ftype, ut})
		ffac := FuncFactory{Tparams: ud.Tparams, Targets: tps}
		scRegisterVarFac(sc, cas.Name, (func(_r0 []FType, _r1 func() TypeVar) VarRef { return GenFuncVar(ctorName, ffac, _r0, _r1) }))
	}))
}

func udRegisterCsCtors(sc Scope, ud UnionDef) {
	frt.PipeUnit(udCases(ud), (func(_r0 []NameTypePair) { slice.Iter((func(_r0 NameTypePair) { csRegisterCtor(sc, ud, _r0) }), _r0) }))
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

func scRegTFData(sc Scope, tname string, tfd TypeFactoryData) {
	scRegisterTypeFac(sc, tname, (func(_r0 []FType) FType { return GenType(tfd, _r0) }))
}

func piRegFF(pi PackageInfo, fname string, ff FuncFactory, ps ParseState) ParseState {
	dict.Add(pi.FuncInfo, fname, ff)
	scRegFunFac(ps.scope, fname, ff)
	return ps
}

func regFF(pi PackageInfo, sc Scope, sff frt.Tuple2[string, FuncFactory]) {
	ffname, ff := frt.Destr2(sff)
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

func rdToRecFac(rd RecordDef) RecordFactory {
	return RecordFactory(rd)
}

func psRegRecDefToTDCtx(rd RecordDef, ps ParseState) {
	rfac := rdToRecFac(rd)
	scRegisterRecFac(ps.scope, rd.Name, rfac)
	rtype, ok := frt.Destr2(tryRecFacToRecType(rfac))
	frt.IfOnly(ok, (func() {
		dict.Add(ps.tdctx.defined, rtype.Name, New_FType_FRecord(rtype))
	}))
}

func psRegUdToTDCtx(ud UnionDef, ps ParseState) {
	sc := ps.scope
	ufac := udToUniFac(ud)
	udRegisterCsCtors(sc, ud)
	scRegisterTypeFac(sc, ud.Name, (func(_r0 []FType) FType { return GenUnionFType(ufac, _r0) }))
	utype, ok := frt.Destr2(tryUniFacToUniType(ufac))
	frt.IfOnly(ok, (func() {
		dict.Add(ps.tdctx.defined, utype.Name, New_FType_FUnion(utype))
	}))
}

func transTVByTDCtx(tdctx TypeDefCtx, tv TypeVar) FType {
	rname, ok := frt.Destr2(dict.TryFind(tdctx.allocedDict, tv.Name))
	return frt.IfElse(ok, (func() FType {
		nt, ok2 := frt.Destr2(dict.TryFind(tdctx.defined, rname))
		return frt.IfElse(ok2, (func() FType {
			return nt
		}), (func() FType {
			frt.PipeUnit(frt.Sprintf1("Unresolved foward decl type: %s", rname), PanicNow)
			return nt
		}))
	}), (func() FType {
		return New_FType_FTypeVar(tv)
	}))
}

func resolveFwrdDecl(ps ParseState, md MultipleDefs) MultipleDefs {
	transTV := (func(_r0 TypeVar) FType { return transTVByTDCtx(ps.tdctx, _r0) })
	transD := (func(_r0 DefStmt) DefStmt { return transTVDefStmt(transTV, _r0) })
	ndefs := slice.Map(transD, md.Defs)
	return MultipleDefs{Defs: ndefs}
}

func scRegDefStmtType(sc Scope, df DefStmt) {
	switch _v1 := (df).(type) {
	case DefStmt_DRecordDef:
		rd := _v1.Value
		frt.PipeUnit(rdToRecFac(rd), (func(_r0 RecordFactory) { scRegisterRecFac(sc, rd.Name, _r0) }))
	case DefStmt_DUnionDef:
		ud := _v1.Value
		udRegisterCsCtors(sc, ud)
		ufac := udToUniFac(ud)
		scRegisterTypeFac(sc, ud.Name, (func(_r0 []FType) FType { return GenUnionFType(ufac, _r0) }))
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func psRegMdTypes(md MultipleDefs, ps ParseState) {
	slice.Iter((func(_r0 DefStmt) { scRegDefStmtType(ps.scope, _r0) }), md.Defs)
}
