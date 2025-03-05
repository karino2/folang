package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/dict"

func tpReplaceOne(tdic dict.Dict[string, FType], tv TypeVar) FType {
	tp, ok := frt.Destr(dict.TryFind(tdic, tv.Name))
	return frt.IfElse(ok, (func() FType {
		return tp
	}), (func() FType {
		return New_FType_FTypeVar(tv)
	}))
}

func tpreplace(tdic dict.Dict[string, FType], ft FType) FType {
	return transTVFType((func(_r0 TypeVar) FType { return tpReplaceOne(tdic, _r0) }), ft)
}

func emptyFtps() []FType {
	return slice.New[FType]()
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
		PanicNow("Too many type specified.")
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
		PanicNow("wrong type param num for instantiate.")
	}))
	return frt.Pipe(ParamdType{Name: tfd.Name, Targs: targs}, New_FType_FParamd)
}

type RecordFactory struct {
	Name    string
	Tparams []string
	Fields  []NameTypePair
}

func GenRecordType(rf RecordFactory, stlist []FType) RecordType {
	frt.IfOnly(frt.OpNotEqual(slice.Len(stlist), slice.Len(rf.Tparams)), (func() {
		PanicNow("wrong type param num for instantiate.")
	}))
	tdic := frt.Pipe(slice.Zip(rf.Tparams, stlist), dict.ToDict)
	nftypes := frt.Pipe(slice.Map(func(_v1 NameTypePair) FType {
		return _v1.Ftype
	}, rf.Fields), (func(_r0 []FType) []FType {
		return slice.Map((func(_r0 FType) FType { return tpreplace(tdic, _r0) }), _r0)
	}))
	fnames := slice.Map(func(_v2 NameTypePair) string {
		return _v2.Name
	}, rf.Fields)
	nfields := frt.Pipe(slice.Zip(fnames, nftypes), (func(_r0 []frt.Tuple2[string, FType]) []NameTypePair { return slice.Map(tupToNTPair, _r0) }))
	rt := RecordType{Name: rf.Name, Targs: stlist}
	ri := RecordTypeInfo{Fields: nfields}
	updateRecInfo(rt, ri)
	return rt
}

func GenRecordFType(rf RecordFactory, stlist []FType) FType {
	return frt.Pipe(GenRecordType(rf, stlist), New_FType_FRecord)
}

func GenRecordTypeByTgen(rf RecordFactory, tvgen func() TypeVar) RecordType {
	ftvs := slice.Map(func(x string) FType {
		return frt.Pipe(tvgen(), New_FType_FTypeVar)
	}, rf.Tparams)
	return GenRecordType(rf, ftvs)
}

func recFacMatch(fieldNames []string, rf RecordFactory) bool {
	return frt.IfElse(frt.OpNotEqual(slice.Length(fieldNames), slice.Length(rf.Fields)), (func() bool {
		return false
	}), (func() bool {
		sortedInput := frt.Pipe(fieldNames, slice.Sort)
		sortedFName := frt.Pipe(slice.Map(func(_v1 NameTypePair) string {
			return _v1.Name
		}, rf.Fields), slice.Sort)
		return frt.OpEqual(sortedInput, sortedFName)
	}))
}

func tryRecFacToRecType(rf RecordFactory) frt.Tuple2[RecordType, bool] {
	return frt.IfElse(slice.IsEmpty(rf.Tparams), (func() frt.Tuple2[RecordType, bool] {
		rt := RecordType{Name: rf.Name}
		ri := RecordTypeInfo{Fields: rf.Fields}
		updateRecInfo(rt, ri)
		return frt.NewTuple2(rt, true)
	}), (func() frt.Tuple2[RecordType, bool] {
		return frt.NewTuple2(frt.Empty[RecordType](), false)
	}))
}

type UnionFactory struct {
	Name    string
	Tparams []string
	Cases   []NameTypePair
}

func ufCases(uf UnionFactory) []NameTypePair {
	return uf.Cases
}

func GenUnionType(uf UnionFactory, stlist []FType) UnionType {
	frt.IfOnly(frt.OpNotEqual(slice.Len(stlist), slice.Len(uf.Tparams)), (func() {
		PanicNow("wrong type param num for instantiate.")
	}))
	tdic := frt.Pipe(slice.Zip(uf.Tparams, stlist), dict.ToDict)
	cases := ufCases(uf)
	nftypes := frt.Pipe(slice.Map(func(_v1 NameTypePair) FType {
		return _v1.Ftype
	}, cases), (func(_r0 []FType) []FType {
		return slice.Map((func(_r0 FType) FType { return tpreplace(tdic, _r0) }), _r0)
	}))
	cnames := slice.Map(func(_v2 NameTypePair) string {
		return _v2.Name
	}, cases)
	ncases := frt.Pipe(slice.Zip(cnames, nftypes), (func(_r0 []frt.Tuple2[string, FType]) []NameTypePair { return slice.Map(tupToNTPair, _r0) }))
	ui := UnionTypeInfo{Cases: ncases}
	ut := UnionType{Name: uf.Name, Targs: stlist}
	updateUniInfo(ut, ui)
	return ut
}

func GenUnionFType(ufac UnionFactory, stlist []FType) FType {
	return frt.Pipe(GenUnionType(ufac, stlist), New_FType_FUnion)
}

func tryUniFacToUniType(uf UnionFactory) frt.Tuple2[UnionType, bool] {
	return frt.IfElse(slice.IsEmpty(uf.Tparams), (func() frt.Tuple2[UnionType, bool] {
		ut := UnionType{Name: uf.Name}
		ui := UnionTypeInfo{Cases: uf.Cases}
		updateUniInfo(ut, ui)
		return frt.NewTuple2(ut, true)
	}), (func() frt.Tuple2[UnionType, bool] {
		return frt.NewTuple2(frt.Empty[UnionType](), false)
	}))
}

func udToUniFac(ud UnionDef) UnionFactory {
	return UnionFactory(ud)
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
	secFncT := newFFunc(([]FType{t1type, t2type}))
	names := ([]string{t1name, t2name})
	tps := ([]FType{t1type, secFncT, t2type})
	args := ([]Expr{lhs, rhs})
	return genBuiltinFunCall(tvgen, "frt.Pipe", names, tps, args)
}

func newPipeCallUnit(tvgen func() TypeVar, lhs Expr, rhs Expr) Expr {
	t1name := "T1"
	t1type := newTvf(t1name)
	secFncT := newFFunc(([]FType{t1type, New_FType_FUnit}))
	names := ([]string{t1name})
	tps := ([]FType{t1type, secFncT, New_FType_FUnit})
	args := ([]Expr{lhs, rhs})
	return genBuiltinFunCall(tvgen, "frt.PipeUnit", names, tps, args)
}

func newPipeCall(tvgen func() TypeVar, lhs Expr, rhs Expr) Expr {
	rht := ExprToType(rhs)
	switch _v1 := (rht).(type) {
	case FType_FFunc:
		ft := _v1.Value
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
	return frt.Pipe(([]FType{argType, retType}), newFFunc)
}

func emptySS() []string {
	return slice.New[string]()
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
