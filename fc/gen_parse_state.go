package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/dict"

import "github.com/karino2/folang/pkg/strings"

func transExprNE(cnv func(Expr) Expr, p NEPair) NEPair {
	return NEPair{Name: p.Name, Expr: cnv(p.Expr)}
}

func transStmt(transV func(Var) Var, transE func(Expr) Expr, stmt Stmt) Stmt {
	switch _v1 := (stmt).(type) {
	case Stmt_SLetVarDef:
		llvd := _v1.Value
		switch _v2 := (llvd).(type) {
		case LLetVarDef_LLOneVarDef:
			lvd := _v2.Value
			nvar := transV(lvd.Lvar)
			nrhs := transE(lvd.Rhs)
			return frt.Pipe(frt.Pipe(LetVarDef{Lvar: nvar, Rhs: nrhs}, New_LLetVarDef_LLOneVarDef), New_Stmt_SLetVarDef)
		case LLetVarDef_LLDestVarDef:
			ldvd := _v2.Value
			nvars := slice.Map(transV, ldvd.Lvars)
			nrhs := transE(ldvd.Rhs)
			return frt.Pipe(frt.Pipe(LetDestVarDef{Lvars: nvars, Rhs: nrhs}, New_LLetVarDef_LLDestVarDef), New_Stmt_SLetVarDef)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Stmt_SExprStmt:
		e := _v1.Value
		return frt.Pipe(transE(e), New_Stmt_SExprStmt)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func transExprMatchRule(pExpr func(Expr) Expr, mr MatchRule) MatchRule {
	nbody := frt.Pipe(frt.Pipe(blockToExpr(mr.Body), pExpr), exprToBlock)
	return MatchRule{Pattern: mr.Pattern, Body: nbody}
}

func transBlock(transE func(Expr) Expr, transS func(Stmt) Stmt, bl Block) Block {
	nss := frt.Pipe(bl.Stmts, (func(_r0 []Stmt) []Stmt { return slice.Map(transS, _r0) }))
	fexpr := transE(bl.FinalExpr)
	return Block{Stmts: nss, FinalExpr: fexpr}
}

func transRecType(transT func(FType) FType, rt RecordType) RecordType {
	ri := lookupRecInfo(rt)
	ntps := frt.Pipe(slice.Map(func(_v1 NameTypePair) FType {
		return _v1.Ftype
	}, ri.Fields), (func(_r0 []FType) []FType { return slice.Map(transT, _r0) }))
	names := slice.Map(func(_v2 NameTypePair) string {
		return _v2.Name
	}, ri.Fields)
	nfields := frt.Pipe(slice.Zip(names, ntps), (func(_r0 []frt.Tuple2[string, FType]) []NameTypePair {
		return slice.Map(func(tp frt.Tuple2[string, FType]) NameTypePair {
			return newNTPair(frt.Fst(tp), frt.Snd(tp))
		}, _r0)
	}))
	ntargs := slice.Map(transT, rt.Targs)
	nri := RecordTypeInfo{Fields: nfields}
	nrt := RecordType{Name: rt.Name, Targs: ntargs}
	updateRecInfo(nrt, nri)
	return nrt
}

func transExpr(transT func(FType) FType, transV func(Var) Var, transS func(Stmt) Stmt, transB func(Block) Block, expr Expr) Expr {
	transE := (func(_r0 Expr) Expr { return transExpr(transT, transV, transS, transB, _r0) })
	switch _v3 := (expr).(type) {
	case Expr_EVarRef:
		rv := _v3.Value
		switch _v4 := (rv).(type) {
		case VarRef_VRVar:
			v := _v4.Value
			return frt.Pipe(frt.Pipe(transV(v), New_VarRef_VRVar), New_Expr_EVarRef)
		case VarRef_VRSVar:
			sv := _v4.Value
			nv := transV(sv.Var)
			return frt.Pipe(frt.Pipe(SpecVar{Var: nv, SpecList: sv.SpecList}, New_VarRef_VRSVar), New_Expr_EVarRef)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_ESlice:
		es := _v3.Value
		return frt.Pipe(slice.Map(transE, es), New_Expr_ESlice)
	case Expr_EBinOpCall:
		bop := _v3.Value
		nlhs := transE(bop.Lhs)
		nrhs := transE(bop.Rhs)
		nret := transT(bop.Rtype)
		return frt.Pipe(BinOpCall{Op: bop.Op, Rtype: nret, Lhs: nlhs, Rhs: nrhs}, New_Expr_EBinOpCall)
	case Expr_ETupleExpr:
		es := _v3.Value
		return frt.Pipe(slice.Map(transE, es), New_Expr_ETupleExpr)
	case Expr_ELambda:
		le := _v3.Value
		nparams := slice.Map(transV, le.Params)
		nbody := transB(le.Body)
		return frt.Pipe(LambdaExpr{Params: nparams, Body: nbody}, New_Expr_ELambda)
	case Expr_ERecordGen:
		rg := _v3.Value
		newNV := slice.Map((func(_r0 NEPair) NEPair { return transExprNE(transE, _r0) }), rg.FieldsNV)
		nrec := transRecType(transT, rg.RecordType)
		return frt.Pipe(RecordGen{FieldsNV: newNV, RecordType: nrec}, New_Expr_ERecordGen)
	case Expr_ELazyBlock:
		lb := _v3.Value
		nbl := transB(lb.Block)
		return frt.Pipe(LazyBlock{Block: nbl}, New_Expr_ELazyBlock)
	case Expr_EReturnableExpr:
		re := _v3.Value
		switch _v5 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v5.Value
			return frt.Pipe(transBlock(transE, transS, bl), blockToExpr)
		case ReturnableExpr_RMatchExpr:
			me := _v5.Value
			ntarget := transE(me.Target)
			nrules := slice.Map((func(_r0 MatchRule) MatchRule { return transExprMatchRule(transE, _r0) }), me.Rules)
			return frt.Pipe(frt.Pipe(MatchExpr{Target: ntarget, Rules: nrules}, New_ReturnableExpr_RMatchExpr), New_Expr_EReturnableExpr)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_EFunCall:
		fc := _v3.Value
		ntarget := transVarVR(transV, fc.TargetFunc)
		nargs := slice.Map(transE, fc.Args)
		return frt.Pipe(FunCall{TargetFunc: ntarget, Args: nargs}, New_Expr_EFunCall)
	case Expr_EBoolLiteral:
		return expr
	case Expr_EGoEvalExpr:
		return expr
	case Expr_EStringLiteral:
		return expr
	case Expr_ESInterP:
		return expr
	case Expr_EIntImm:
		return expr
	case Expr_EUnit:
		return expr
	case Expr_EFieldAccess:
		fa := _v3.Value
		ntarget := transE(fa.TargetExpr)
		return frt.Pipe(FieldAccess{TargetExpr: ntarget, FieldName: fa.FieldName}, New_Expr_EFieldAccess)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

type SSet struct {
	Dict dict.Dict[string, bool]
}

func NewSSet() SSet {
	return SSet{Dict: dict.New[string, bool]()}
}

func SSetHasKey(st SSet, key string) bool {
	_, ok := frt.Destr(dict.TryFind(st.Dict, key))
	return ok
}

func SSetPut(st SSet, key string) {
	dict.Add(st.Dict, key, true)
}

func collectTVarFTypeWithSet(visited SSet, ft FType) []string {
	recurse := (func(_r0 FType) []string { return collectTVarFTypeWithSet(visited, _r0) })
	switch _v6 := (ft).(type) {
	case FType_FTypeVar:
		tv := _v6.Value
		return ([]string{tv.Name})
	case FType_FSlice:
		ts := _v6.Value
		return recurse(ts.ElemType)
	case FType_FTuple:
		ftup := _v6.Value
		return slice.Collect(recurse, ftup.ElemTypes)
	case FType_FFieldAccess:
		fa := _v6.Value
		return recurse(fa.RecType)
	case FType_FRecord:
		rt := _v6.Value
		ri := lookupRecInfo(rt)
		fres := frt.Pipe(frt.Pipe(ri.Fields, (func(_r0 []NameTypePair) []FType {
			return slice.Map(func(_v1 NameTypePair) FType {
				return _v1.Ftype
			}, _r0)
		})), (func(_r0 []FType) []string { return slice.Collect(recurse, _r0) }))
		tres := frt.Pipe(rt.Targs, (func(_r0 []FType) []string { return slice.Collect(recurse, _r0) }))
		return slice.Append(fres, tres)
	case FType_FUnion:
		ut := _v6.Value
		uname := utName(ut)
		return frt.IfElse(SSetHasKey(visited, uname), (func() []string {
			return slice.New[string]()
		}), (func() []string {
			SSetPut(visited, uname)
			return frt.Pipe(frt.Pipe(utCases(ut), (func(_r0 []NameTypePair) []FType {
				return slice.Map(func(_v2 NameTypePair) FType {
					return _v2.Ftype
				}, _r0)
			})), (func(_r0 []FType) []string { return slice.Collect(recurse, _r0) }))
		}))
	case FType_FFunc:
		fnt := _v6.Value
		return slice.Collect(recurse, fnt.Targets)
	default:
		return slice.New[string]()
	}
}

func collectTVarFType(ft FType) []string {
	visited := NewSSet()
	return collectTVarFTypeWithSet(visited, ft)
}

func collectTVarStmt(collE func(Expr) []string, stmt Stmt) []string {
	switch _v7 := (stmt).(type) {
	case Stmt_SLetVarDef:
		llvd := _v7.Value
		switch _v8 := (llvd).(type) {
		case LLetVarDef_LLOneVarDef:
			lvd := _v8.Value
			nvar := collectTVarFType(lvd.Lvar.Ftype)
			nrhs := collE(lvd.Rhs)
			return slice.Append(nvar, nrhs)
		case LLetVarDef_LLDestVarDef:
			ldvd := _v8.Value
			nvars := frt.Pipe(slice.Map(func(_v1 Var) FType {
				return _v1.Ftype
			}, ldvd.Lvars), (func(_r0 []FType) []string { return slice.Collect(collectTVarFType, _r0) }))
			nrhs := collE(ldvd.Rhs)
			return slice.Append(nvars, nrhs)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Stmt_SExprStmt:
		e := _v7.Value
		return collE(e)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func collectTVarBlock(collE func(Expr) []string, collS func(Stmt) []string, bl Block) []string {
	nss := frt.Pipe(bl.Stmts, (func(_r0 []Stmt) []string { return slice.Collect(collS, _r0) }))
	fexpr := collE(bl.FinalExpr)
	return slice.Append(nss, fexpr)
}

func collectTVarExpr(expr Expr) []string {
	recurse := collectTVarExpr
	collS := (func(_r0 Stmt) []string { return collectTVarStmt(recurse, _r0) })
	collB := (func(_r0 Block) []string { return collectTVarBlock(recurse, collS, _r0) })
	switch _v9 := (expr).(type) {
	case Expr_EVarRef:
		vr := _v9.Value
		return frt.Pipe(varRefVarType(vr), collectTVarFType)
	case Expr_ESlice:
		es := _v9.Value
		return slice.Collect(recurse, es)
	case Expr_EBinOpCall:
		bop := _v9.Value
		lres := recurse(bop.Lhs)
		rres := recurse(bop.Rhs)
		return slice.Append(lres, rres)
	case Expr_ETupleExpr:
		es := _v9.Value
		return slice.Collect(recurse, es)
	case Expr_ELambda:
		le := _v9.Value
		pas := frt.Pipe(slice.Map(func(_v1 Var) FType {
			return _v1.Ftype
		}, le.Params), (func(_r0 []FType) []string { return slice.Collect(collectTVarFType, _r0) }))
		return frt.Pipe(collB(le.Body), (func(_r0 []string) []string { return slice.Append(pas, _r0) }))
	case Expr_ERecordGen:
		rg := _v9.Value
		return frt.Pipe(slice.Map(func(_v2 NEPair) Expr {
			return _v2.Expr
		}, rg.FieldsNV), (func(_r0 []Expr) []string { return slice.Collect(recurse, _r0) }))
	case Expr_ELazyBlock:
		lb := _v9.Value
		return collB(lb.Block)
	case Expr_EReturnableExpr:
		re := _v9.Value
		switch _v10 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v10.Value
			return collB(bl)
		case ReturnableExpr_RMatchExpr:
			me := _v10.Value
			return frt.Pipe(frt.Pipe(slice.Map(func(_v3 MatchRule) Block {
				return _v3.Body
			}, me.Rules), (func(_r0 []Block) []string { return slice.Collect(collB, _r0) })), (func(_r0 []string) []string { return slice.Append(recurse(me.Target), _r0) }))
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_EFunCall:
		fc := _v9.Value
		colt := frt.Pipe(varRefVarType(fc.TargetFunc), collectTVarFType)
		return frt.Pipe(slice.Collect(recurse, fc.Args), (func(_r0 []string) []string { return slice.Append(colt, _r0) }))
	case Expr_EFieldAccess:
		fa := _v9.Value
		return recurse(fa.TargetExpr)
	default:
		return []string{}
	}
}

func collectTVarBlockFacade(b Block) []string {
	collE := collectTVarExpr
	collS := (func(_r0 Stmt) []string { return collectTVarStmt(collE, _r0) })
	return collectTVarBlock(collE, collS, b)
}

func transTVFTypeWithSet(visited SSet, transTV func(TypeVar) FType, ftp FType) FType {
	recurse := (func(_r0 FType) FType { return transTVFTypeWithSet(visited, transTV, _r0) })
	switch _v11 := (ftp).(type) {
	case FType_FTypeVar:
		tv := _v11.Value
		return transTV(tv)
	case FType_FSlice:
		ts := _v11.Value
		et := recurse(ts.ElemType)
		return New_FType_FSlice(SliceType{ElemType: et})
	case FType_FTuple:
		ftup := _v11.Value
		nts := slice.Map(recurse, ftup.ElemTypes)
		return frt.Pipe(TupleType{ElemTypes: nts}, New_FType_FTuple)
	case FType_FFieldAccess:
		fa := _v11.Value
		nrec := recurse(fa.RecType)
		return frt.Pipe(FieldAccessType{RecType: nrec, FieldName: fa.FieldName}, faResolve)
	case FType_FFunc:
		fnt := _v11.Value
		return frt.Pipe(slice.Map(recurse, fnt.Targets), newFFunc)
	case FType_FParamd:
		pt := _v11.Value
		nts := slice.Map(recurse, pt.Targs)
		return frt.Pipe(ParamdType{Name: pt.Name, Targs: nts}, New_FType_FParamd)
	case FType_FRecord:
		rt := _v11.Value
		return frt.Pipe(transRecType(recurse, rt), New_FType_FRecord)
	case FType_FUnion:
		ut := _v11.Value
		uname := utName(ut)
		return frt.IfElse(SSetHasKey(visited, uname), (func() FType {
			return ftp
		}), (func() FType {
			SSetPut(visited, uname)
			cases := utCases(ut)
			ntps := frt.Pipe(slice.Map(func(_v1 NameTypePair) FType {
				return _v1.Ftype
			}, cases), (func(_r0 []FType) []FType { return slice.Map(recurse, _r0) }))
			names := slice.Map(func(_v2 NameTypePair) string {
				return _v2.Name
			}, cases)
			ncases := frt.Pipe(slice.Zip(names, ntps), (func(_r0 []frt.Tuple2[string, FType]) []NameTypePair {
				return slice.Map(func(tp frt.Tuple2[string, FType]) NameTypePair {
					return newNTPair(frt.Fst(tp), frt.Snd(tp))
				}, _r0)
			}))
			utUpdateCases(ut, ncases)
			return New_FType_FUnion(ut)
		}))
	default:
		return ftp
	}
}

func transTVFType(transTV func(TypeVar) FType, ftp FType) FType {
	visited := NewSSet()
	return transTVFTypeWithSet(visited, transTV, ftp)
}

func transTVVar(transTV func(TypeVar) FType, v Var) Var {
	ntp := transTVFType(transTV, v.Ftype)
	return Var{Name: v.Name, Ftype: ntp}
}

func isTDTVar(tvname string) bool {
	return strings.HasPrefix("_P", tvname)
}

func noTDTVarInFType(ft FType) bool {
	return frt.Pipe(frt.Pipe(collectTVarFType(ft), (func(_r0 []string) []string { return slice.Filter(isTDTVar, _r0) })), slice.IsEmpty)
}

func transTRecurse(transT func(FType) FType, count int, ft FType) FType {
	frt.IfOnly((count > 1000), (func() {
		PanicNow("Too deep recurse fwddecl, maybe cyclic, give up")
	}))
	nt := transT(ft)
	noTvFound := noTDTVarInFType(nt)
	return frt.IfElse(noTvFound, (func() FType {
		return nt
	}), (func() FType {
		return transTRecurse(transT, (count + 1), nt)
	}))
}

func transTVNTPair(transV func(TypeVar) FType, ntp NameTypePair) NameTypePair {
	transTOne := (func(_r0 FType) FType { return transTVFType(transV, _r0) })
	nt := transTRecurse(transTOne, 0, ntp.Ftype)
	return NameTypePair{Name: ntp.Name, Ftype: nt}
}

func transTVDefStmt(transTV func(TypeVar) FType, df DefStmt) DefStmt {
	switch _v12 := (df).(type) {
	case DefStmt_DRecordDef:
		rd := _v12.Value
		nfields := slice.Map((func(_r0 NameTypePair) NameTypePair { return transTVNTPair(transTV, _r0) }), rd.Fields)
		noTvFound := frt.Pipe(slice.Map(func(_v1 NameTypePair) FType {
			return _v1.Ftype
		}, nfields), (func(_r0 []FType) bool { return slice.Forall(noTDTVarInFType, _r0) }))
		frt.IfOnly(frt.OpNot(noTvFound), (func() {
			PanicNow("Unresolve type")
		}))
		return frt.Pipe(RecordDef{Name: rd.Name, Tparams: rd.Tparams, Fields: nfields}, New_DefStmt_DRecordDef)
	case DefStmt_DUnionDef:
		ud := _v12.Value
		ncases := frt.Pipe(udCases(ud), (func(_r0 []NameTypePair) []NameTypePair {
			return slice.Map((func(_r0 NameTypePair) NameTypePair { return transTVNTPair(transTV, _r0) }), _r0)
		}))
		noTvFound := frt.Pipe(slice.Map(func(_v2 NameTypePair) FType {
			return _v2.Ftype
		}, ncases), (func(_r0 []FType) bool { return slice.Forall(noTDTVarInFType, _r0) }))
		frt.IfOnly(frt.OpNot(noTvFound), (func() {
			PanicNow("Unresolve type2")
		}))
		return frt.Pipe(udUpdate(ud, ncases), New_DefStmt_DUnionDef)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func transTVExpr(transTV func(TypeVar) FType, expr Expr) Expr {
	transE := (func(_r0 Expr) Expr { return transTVExpr(transTV, _r0) })
	transT := (func(_r0 FType) FType { return transTVFType(transTV, _r0) })
	transV := (func(_r0 Var) Var { return transTVVar(transTV, _r0) })
	transS := (func(_r0 Stmt) Stmt { return transStmt(transV, transE, _r0) })
	transB := (func(_r0 Block) Block { return transBlock(transE, transS, _r0) })
	return transExpr(transT, transV, transS, transB, expr)
}

func transTVBlock(transTV func(TypeVar) FType, block Block) Block {
	transE := (func(_r0 Expr) Expr { return transTVExpr(transTV, _r0) })
	transV := (func(_r0 Var) Var { return transTVVar(transTV, _r0) })
	transS := (func(_r0 Stmt) Stmt { return transStmt(transV, transE, _r0) })
	return transBlock(transE, transS, block)
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
	vfac, ok := frt.Destr(dict.TryFind(sd.VarFacMap, name))
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

func emptyRec() RecordType {
	return frt.Empty[RecordType]()
}

func emptyRecFac() RecordFactory {
	return frt.Empty[RecordFactory]()
}

func scLookupRecFac(s Scope, fieldNames []string) frt.Tuple2[RecordFactory, bool] {
	rfac, ok := frt.Destr(scLookupRecFacCur(s, fieldNames))
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
	rfac, ok := frt.Destr(dict.TryFind(sd.RecFacMap, name))
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
	rec, ok := frt.Destr(dict.TryFind(sd.TypeFacMap, name))
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

func CnvL[T0 any, T1 any, T2 any](fn func(T0) T1, tup frt.Tuple2[T0, T2]) frt.Tuple2[T1, T2] {
	nl := frt.Pipe(frt.Fst(tup), fn)
	return frt.NewTuple2(nl, frt.Snd(tup))
}

func CnvR[T0 any, T1 any, T2 any](fn func(T0) T1, tup frt.Tuple2[T2, T0]) frt.Tuple2[T2, T1] {
	nr := frt.Pipe(frt.Snd(tup), fn)
	return frt.NewTuple2(frt.Fst(tup), nr)
}

func withPs[T0 any](ps ParseState, v T0) frt.Tuple2[ParseState, T0] {
	return frt.NewTuple2(ps, v)
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
		psPanic(ps, "non expected token")
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
	ut := UnionType{Name: ud.Name}
	ui := UnionTypeInfo{Cases: ud.Cases}
	updateUniInfo(ut, ui)
	return ut
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
	switch _v13 := (rht).(type) {
	case FType_FFunc:
		ft := _v13.Value
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

func rdToRecFac(rd RecordDef) RecordFactory {
	return RecordFactory(rd)
}

func psRegRecDefToTDCtx(rd RecordDef, ps ParseState) {
	rfac := rdToRecFac(rd)
	scRegisterRecFac(ps.scope, rd.Name, rfac)
	rtype, ok := frt.Destr(tryRecFacToRecType(rfac))
	frt.IfOnly(ok, (func() {
		dict.Add(ps.tdctx.defined, rtype.Name, New_FType_FRecord(rtype))
	}))
}

func psRegUdToTDCtx(ud UnionDef, ps ParseState) {
	sc := ps.scope
	udRegisterCsCtors(sc, ud)
	fut := udToFUt(ud)
	scRegisterType(sc, ud.Name, fut)
	dict.Add(ps.tdctx.defined, ud.Name, fut)
}

func transTVByTDCtx(tdctx TypeDefCtx, tv TypeVar) FType {
	rname, ok := frt.Destr(dict.TryFind(tdctx.allocedDict, tv.Name))
	return frt.IfElse(ok, (func() FType {
		nt, ok2 := frt.Destr(dict.TryFind(tdctx.defined, rname))
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
	switch _v14 := (df).(type) {
	case DefStmt_DRecordDef:
		rd := _v14.Value
		frt.PipeUnit(rdToRecFac(rd), (func(_r0 RecordFactory) { scRegisterRecFac(sc, rd.Name, _r0) }))
	case DefStmt_DUnionDef:
		ud := _v14.Value
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
