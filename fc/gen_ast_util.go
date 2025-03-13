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

func umrMapBlock(pBlock func(Block) Block, umr UnionMatchRule) UnionMatchRule {
	nbody := pBlock(umr.Body)
	return UnionMatchRule{UnionPattern: umr.UnionPattern, Body: nbody}
}

func mrsMapBlock(pBlock func(Block) Block, mr MatchRules) MatchRules {
	one := (func(_r0 UnionMatchRule) UnionMatchRule { return umrMapBlock(pBlock, _r0) })
	switch _v3 := (mr).(type) {
	case MatchRules_Unions:
		us := _v3.Value
		return frt.Pipe(slice.Map(one, us), New_MatchRules_Unions)
	case MatchRules_UnionsWD:
		uds := _v3.Value
		nus := slice.Map(one, uds.Unions)
		nd := pBlock(uds.Default)
		return frt.Pipe(UnionMatchRulesWD{Unions: nus, Default: nd}, New_MatchRules_UnionsWD)
	default:
		panic("Union pattern fail. Never reached here.")
	}
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
	switch _v4 := (expr).(type) {
	case Expr_EVarRef:
		rv := _v4.Value
		switch _v5 := (rv).(type) {
		case VarRef_VRVar:
			v := _v5.Value
			return frt.Pipe(frt.Pipe(transV(v), New_VarRef_VRVar), New_Expr_EVarRef)
		case VarRef_VRSVar:
			sv := _v5.Value
			nv := transV(sv.Var)
			return frt.Pipe(frt.Pipe(SpecVar{Var: nv, SpecList: sv.SpecList}, New_VarRef_VRSVar), New_Expr_EVarRef)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_ESlice:
		es := _v4.Value
		return frt.Pipe(slice.Map(transE, es), New_Expr_ESlice)
	case Expr_EBinOpCall:
		bop := _v4.Value
		nlhs := transE(bop.Lhs)
		nrhs := transE(bop.Rhs)
		nret := transT(bop.Rtype)
		return frt.Pipe(BinOpCall{Op: bop.Op, Rtype: nret, Lhs: nlhs, Rhs: nrhs}, New_Expr_EBinOpCall)
	case Expr_ETupleExpr:
		es := _v4.Value
		return frt.Pipe(slice.Map(transE, es), New_Expr_ETupleExpr)
	case Expr_ELambda:
		le := _v4.Value
		nparams := slice.Map(transV, le.Params)
		nbody := transB(le.Body)
		return frt.Pipe(LambdaExpr{Params: nparams, Body: nbody}, New_Expr_ELambda)
	case Expr_ERecordGen:
		rg := _v4.Value
		newNV := slice.Map((func(_r0 NEPair) NEPair { return transExprNE(transE, _r0) }), rg.FieldsNV)
		nrec := transRecType(transT, rg.RecordType)
		return frt.Pipe(RecordGen{FieldsNV: newNV, RecordType: nrec}, New_Expr_ERecordGen)
	case Expr_ELazyBlock:
		lb := _v4.Value
		nbl := transB(lb.Block)
		return frt.Pipe(LazyBlock{Block: nbl}, New_Expr_ELazyBlock)
	case Expr_EReturnableExpr:
		re := _v4.Value
		switch _v6 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v6.Value
			return frt.Pipe(transBlock(transE, transS, bl), blockToExpr)
		case ReturnableExpr_RMatchExpr:
			me := _v6.Value
			ntarget := transE(me.Target)
			nrules := mrsMapBlock(transB, me.Rules)
			return frt.Pipe(frt.Pipe(MatchExpr{Target: ntarget, Rules: nrules}, New_ReturnableExpr_RMatchExpr), New_Expr_EReturnableExpr)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_EFunCall:
		fc := _v4.Value
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
		fa := _v4.Value
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
	_, ok := frt.Destr2(dict.TryFind(st.Dict, key))
	return ok
}

func SSetPut(st SSet, key string) {
	dict.Add(st.Dict, key, true)
}

func collectTVarFTypeWithSet(visited SSet, ft FType) []string {
	recurse := (func(_r0 FType) []string { return collectTVarFTypeWithSet(visited, _r0) })
	switch _v7 := (ft).(type) {
	case FType_FTypeVar:
		tv := _v7.Value
		return ([]string{tv.Name})
	case FType_FSlice:
		ts := _v7.Value
		return recurse(ts.ElemType)
	case FType_FTuple:
		ftup := _v7.Value
		return slice.Collect(recurse, ftup.ElemTypes)
	case FType_FFieldAccess:
		fa := _v7.Value
		return recurse(fa.RecType)
	case FType_FRecord:
		rt := _v7.Value
		ri := lookupRecInfo(rt)
		fres := frt.Pipe(frt.Pipe(ri.Fields, (func(_r0 []NameTypePair) []FType {
			return slice.Map(func(_v1 NameTypePair) FType {
				return _v1.Ftype
			}, _r0)
		})), (func(_r0 []FType) []string { return slice.Collect(recurse, _r0) }))
		tres := frt.Pipe(rt.Targs, (func(_r0 []FType) []string { return slice.Collect(recurse, _r0) }))
		return slice.Append(fres, tres)
	case FType_FUnion:
		ut := _v7.Value
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
		fnt := _v7.Value
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
	switch _v8 := (stmt).(type) {
	case Stmt_SLetVarDef:
		llvd := _v8.Value
		switch _v9 := (llvd).(type) {
		case LLetVarDef_LLOneVarDef:
			lvd := _v9.Value
			nvar := collectTVarFType(lvd.Lvar.Ftype)
			nrhs := collE(lvd.Rhs)
			return slice.Append(nvar, nrhs)
		case LLetVarDef_LLDestVarDef:
			ldvd := _v9.Value
			nvars := frt.Pipe(slice.Map(func(_v1 Var) FType {
				return _v1.Ftype
			}, ldvd.Lvars), (func(_r0 []FType) []string { return slice.Collect(collectTVarFType, _r0) }))
			nrhs := collE(ldvd.Rhs)
			return slice.Append(nvars, nrhs)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Stmt_SExprStmt:
		e := _v8.Value
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

func mrsToBlocks(rules MatchRules) []Block {
	switch _v10 := (rules).(type) {
	case MatchRules_Unions:
		us := _v10.Value
		return slice.Map(func(_v1 UnionMatchRule) Block {
			return _v1.Body
		}, us)
	case MatchRules_UnionsWD:
		uds := _v10.Value
		return frt.Pipe(slice.Map(func(_v2 UnionMatchRule) Block {
			return _v2.Body
		}, uds.Unions), (func(_r0 []Block) []Block { return slice.PushLast(uds.Default, _r0) }))
	case MatchRules_DefaultOnly:
		db := _v10.Value
		return ([]Block{db})
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func collectTVarExpr(expr Expr) []string {
	recurse := collectTVarExpr
	collS := (func(_r0 Stmt) []string { return collectTVarStmt(recurse, _r0) })
	collB := (func(_r0 Block) []string { return collectTVarBlock(recurse, collS, _r0) })
	switch _v11 := (expr).(type) {
	case Expr_EVarRef:
		vr := _v11.Value
		return frt.Pipe(varRefVarType(vr), collectTVarFType)
	case Expr_ESlice:
		es := _v11.Value
		return slice.Collect(recurse, es)
	case Expr_EBinOpCall:
		bop := _v11.Value
		lres := recurse(bop.Lhs)
		rres := recurse(bop.Rhs)
		return slice.Append(lres, rres)
	case Expr_ETupleExpr:
		es := _v11.Value
		return slice.Collect(recurse, es)
	case Expr_ELambda:
		le := _v11.Value
		pas := frt.Pipe(slice.Map(func(_v1 Var) FType {
			return _v1.Ftype
		}, le.Params), (func(_r0 []FType) []string { return slice.Collect(collectTVarFType, _r0) }))
		return frt.Pipe(collB(le.Body), (func(_r0 []string) []string { return slice.Append(pas, _r0) }))
	case Expr_ERecordGen:
		rg := _v11.Value
		return frt.Pipe(slice.Map(func(_v2 NEPair) Expr {
			return _v2.Expr
		}, rg.FieldsNV), (func(_r0 []Expr) []string { return slice.Collect(recurse, _r0) }))
	case Expr_ELazyBlock:
		lb := _v11.Value
		return collB(lb.Block)
	case Expr_EReturnableExpr:
		re := _v11.Value
		switch _v12 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v12.Value
			return collB(bl)
		case ReturnableExpr_RMatchExpr:
			me := _v12.Value
			return frt.Pipe(frt.Pipe(mrsToBlocks(me.Rules), (func(_r0 []Block) []string { return slice.Collect(collB, _r0) })), (func(_r0 []string) []string { return slice.Append(recurse(me.Target), _r0) }))
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_EFunCall:
		fc := _v11.Value
		colt := frt.Pipe(varRefVarType(fc.TargetFunc), collectTVarFType)
		return frt.Pipe(slice.Collect(recurse, fc.Args), (func(_r0 []string) []string { return slice.Append(colt, _r0) }))
	case Expr_EFieldAccess:
		fa := _v11.Value
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
	switch _v13 := (ftp).(type) {
	case FType_FTypeVar:
		tv := _v13.Value
		return transTV(tv)
	case FType_FSlice:
		ts := _v13.Value
		et := recurse(ts.ElemType)
		return New_FType_FSlice(SliceType{ElemType: et})
	case FType_FTuple:
		ftup := _v13.Value
		nts := slice.Map(recurse, ftup.ElemTypes)
		return frt.Pipe(TupleType{ElemTypes: nts}, New_FType_FTuple)
	case FType_FFieldAccess:
		fa := _v13.Value
		nrec := recurse(fa.RecType)
		return frt.Pipe(FieldAccessType{RecType: nrec, FieldName: fa.FieldName}, faResolve)
	case FType_FFunc:
		fnt := _v13.Value
		return frt.Pipe(slice.Map(recurse, fnt.Targets), newFFunc)
	case FType_FParamd:
		pt := _v13.Value
		nts := slice.Map(recurse, pt.Targs)
		return frt.Pipe(ParamdType{Name: pt.Name, Targs: nts}, New_FType_FParamd)
	case FType_FRecord:
		rt := _v13.Value
		return frt.Pipe(transRecType(recurse, rt), New_FType_FRecord)
	case FType_FUnion:
		ut := _v13.Value
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
			ntargs := slice.Map(recurse, ut.Targs)
			nut := UnionType{Name: ut.Name, Targs: ntargs}
			nui := UnionTypeInfo{Cases: ncases}
			updateUniInfo(nut, nui)
			return New_FType_FUnion(nut)
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
	switch _v14 := (df).(type) {
	case DefStmt_DRecordDef:
		rd := _v14.Value
		nfields := slice.Map((func(_r0 NameTypePair) NameTypePair { return transTVNTPair(transTV, _r0) }), rd.Fields)
		noTvFound := frt.Pipe(slice.Map(func(_v1 NameTypePair) FType {
			return _v1.Ftype
		}, nfields), (func(_r0 []FType) bool { return slice.Forall(noTDTVarInFType, _r0) }))
		frt.IfOnly(frt.OpNot(noTvFound), (func() {
			PanicNow("Unresolve type")
		}))
		return frt.Pipe(RecordDef{Name: rd.Name, Tparams: rd.Tparams, Fields: nfields}, New_DefStmt_DRecordDef)
	case DefStmt_DUnionDef:
		ud := _v14.Value
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
