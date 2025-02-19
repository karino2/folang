package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

type UniRel struct {
	SrcV string
	Dest FType
}

func emptyRels() []UniRel {
	return []UniRel{}
}

func tupApply(f func(FType, FType) frt.Tuple2[FType, []UniRel], tup frt.Tuple2[FType, FType]) frt.Tuple2[FType, []UniRel] {
	lhs, rhs := frt.Destr(tup)
	return f(lhs, rhs)
}

func withRels(rels []UniRel, tp FType) frt.Tuple2[FType, []UniRel] {
	return frt.NewTuple2(tp, rels)
}

func withTp(tp FType, rels []UniRel) frt.Tuple2[FType, []UniRel] {
	return frt.NewTuple2(tp, rels)
}

func compositeTpList(cOne func(FType, FType) frt.Tuple2[FType, []UniRel], lhs []FType, rhs []FType) frt.Tuple2[[]FType, []UniRel] {
	tups := frt.Pipe(slice.Zip(lhs, rhs), (func(_r0 []frt.Tuple2[FType, FType]) []frt.Tuple2[FType, []UniRel] {
		return slice.Map((func(_r0 frt.Tuple2[FType, FType]) frt.Tuple2[FType, []UniRel] { return tupApply(cOne, _r0) }), _r0)
	}))
	tps := frt.Pipe(tups, (func(_r0 []frt.Tuple2[FType, []UniRel]) []FType { return slice.Map(frt.Fst, _r0) }))
	rels := frt.Pipe(frt.Pipe(tups, (func(_r0 []frt.Tuple2[FType, []UniRel]) [][]UniRel { return slice.Map(frt.Snd, _r0) })), slice.Concat)
	return frt.NewTuple2(tps, rels)
}

func compositeTp(lhs FType, rhs FType) frt.Tuple2[FType, []UniRel] {
	switch _v412 := (lhs).(type) {
	case FType_FTypeVar:
		tv := _v412.Value
		switch _v413 := (rhs).(type) {
		case FType_FTypeVar:
			tv2 := _v413.Value
			return frt.IfElse(frt.OpEqual(tv.Name, tv2.Name), (func() frt.Tuple2[FType, []UniRel] {
				return frt.Pipe(emptyRels(), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
			}), (func() frt.Tuple2[FType, []UniRel] {
				return frt.IfElse((tv.Name > tv2.Name), (func() frt.Tuple2[FType, []UniRel] {
					return frt.Pipe(([]UniRel{UniRel{SrcV: tv.Name, Dest: rhs}}), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(rhs, _r0) }))
				}), (func() frt.Tuple2[FType, []UniRel] {
					return frt.Pipe(([]UniRel{UniRel{SrcV: tv2.Name, Dest: lhs}}), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
				}))
			}))
		default:
			return frt.Pipe(([]UniRel{UniRel{SrcV: tv.Name, Dest: rhs}}), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(rhs, _r0) }))
		}
	default:
		switch _v414 := (rhs).(type) {
		case FType_FTypeVar:
			tv2 := _v414.Value
			return frt.Pipe(([]UniRel{UniRel{SrcV: tv2.Name, Dest: lhs}}), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
		case FType_FSlice:
			ts2 := _v414.Value
			switch _v415 := (lhs).(type) {
			case FType_FSlice:
				ts1 := _v415.Value
				rtp, rels := frt.Destr(compositeTp(ts1.ElemType, ts2.ElemType))
				return frt.Pipe(frt.Pipe(SliceType{ElemType: rtp}, New_FType_FSlice), (func(_r0 FType) frt.Tuple2[FType, []UniRel] { return withRels(rels, _r0) }))
			case FType_FFieldAccess:
				return frt.Pipe(emptyRels(), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(rhs, _r0) }))
			default:
				frt.Panic("right is slice, left is neither slice nor field access.")
				return frt.Pipe(emptyRels(), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
			}
		case FType_FFieldAccess:
			fa2 := _v414.Value
			switch _v416 := (lhs).(type) {
			case FType_FFieldAccess:
				fa1 := _v416.Value
				rtp, rels := frt.Destr(compositeTp(fa1.RecType, fa2.RecType))
				return frt.Pipe(frt.Pipe(FieldAccessType{RecType: rtp, FieldName: fa1.FieldName}, faResolve), (func(_r0 FType) frt.Tuple2[FType, []UniRel] { return withRels(rels, _r0) }))
			case FType_FSlice:
				return frt.Pipe(emptyRels(), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
			default:
				frt.Panic("unknown case")
				return frt.Pipe(emptyRels(), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
			}
		case FType_FFunc:
			tf2 := _v414.Value
			tf1 := lhs.(FType_FFunc).Value
			tps, rels := frt.Destr(compositeTpList(compositeTp, tf1.Targets, tf2.Targets))
			return frt.Pipe(frt.Pipe(FuncType{Targets: tps}, New_FType_FFunc), (func(_r0 FType) frt.Tuple2[FType, []UniRel] { return withRels(rels, _r0) }))
		case FType_FTuple:
			tt2 := _v414.Value
			tt1 := lhs.(FType_FTuple).Value
			tps, rels := frt.Destr(compositeTpList(compositeTp, tt1.ElemTypes, tt2.ElemTypes))
			return frt.Pipe(frt.Pipe(TupleType{ElemTypes: tps}, New_FType_FTuple), (func(_r0 FType) frt.Tuple2[FType, []UniRel] { return withRels(rels, _r0) }))
		default:
			return frt.Pipe(emptyRels(), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
		}
	}
}

func unifyType(lhs FType, rhs FType) []UniRel {
	_, rels := frt.Destr(compositeTp(lhs, rhs))
	return rels
}

func unifyTupArg(tup frt.Tuple2[FType, FType]) []UniRel {
	lhs, rhs := frt.Destr(tup)
	return unifyType(lhs, rhs)
}

func unifyVETup(veTup frt.Tuple2[Var, FType]) []UniRel {
	v, ft := frt.Destr(veTup)
	return frt.IfElse(frt.OpEqual(v.Name, "_"), (func() []UniRel {
		return emptyRels()
	}), (func() []UniRel {
		return unifyType(v.Ftype, ft)
	}))
}

func vToT(v Var) FType {
	return v.Ftype
}

func varsToTupleType(vars []Var) FType {
	ets := slice.Map(vToT, vars)
	return frt.Pipe(TupleType{ElemTypes: ets}, New_FType_FTuple)
}

func collectStmtRel(ec func(Expr) []UniRel, stmt Stmt) []UniRel {
	switch _v417 := (stmt).(type) {
	case Stmt_SExprStmt:
		se := _v417.Value
		return ec(se)
	case Stmt_SLetVarDef:
		slvd := _v417.Value
		switch _v418 := (slvd).(type) {
		case LLetVarDef_LLOneVarDef:
			lvd := _v418.Value
			inside := ec(lvd.Rhs)
			return frt.Pipe(unifyType(lvd.Lvar.Ftype, ExprToType(lvd.Rhs)), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
		case LLetVarDef_LLDestVarDef:
			ldvd := _v418.Value
			inside := ec(ldvd.Rhs)
			rhtype := ExprToType(ldvd.Rhs)
			switch _v419 := (rhtype).(type) {
			case FType_FTuple:
				ft := _v419.Value
				return frt.Pipe(frt.Pipe(frt.Pipe(slice.Zip(ldvd.Lvars, ft.ElemTypes), (func(_r0 []frt.Tuple2[Var, FType]) [][]UniRel { return slice.Map(unifyVETup, _r0) })), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
			case FType_FTypeVar:
				lft := varsToTupleType(ldvd.Lvars)
				return frt.Pipe(unifyType(rhtype, lft), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
			default:
				frt.Panic("Destructuring of right is not tuple, NYI.")
				return inside
			}
		default:
			panic("Union pattern fail. Never reached here.")
		}
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func collectFunCall(fc FunCall) []UniRel {
	tftype := fc.TargetFunc.Ftype
	switch _v420 := (tftype).(type) {
	case FType_FFunc:
		fft := _v420.Value
		argTps := slice.Map(ExprToType, fc.Args)
		tpArgTps := frt.Pipe(fargs(fft), (func(_r0 []FType) []FType { return slice.Take(slice.Length(argTps), _r0) }))
		return frt.Pipe(frt.Pipe(slice.Zip(argTps, tpArgTps), (func(_r0 []frt.Tuple2[FType, FType]) [][]UniRel { return slice.Map(unifyTupArg, _r0) })), slice.Concat)
	default:
		frt.Panic("funcall with non func first arg, possibly TypeVar, NYI.")
		return emptyRels()
	}
}

func collectSlice(es []Expr) []UniRel {
	return frt.IfElse((slice.Length(es) <= 1), (func() []UniRel {
		return emptyRels()
	}), (func() []UniRel {
		headT := frt.Pipe(slice.Head(es), ExprToType)
		return frt.Pipe(frt.Pipe(frt.Pipe(slice.Tail(es), (func(_r0 []Expr) []FType { return slice.Map(ExprToType, _r0) })), (func(_r0 []FType) [][]UniRel {
			return slice.Map((func(_r0 FType) []UniRel { return unifyType(headT, _r0) }), _r0)
		})), slice.Concat)
	}))
}

func collectBlock(colE func(Expr) []UniRel, colS func(Stmt) []UniRel, block Block) []UniRel {
	return frt.Pipe(frt.Pipe(slice.Map(colS, block.Stmts), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(colE(block.FinalExpr), _r0) }))
}

func mrToBlock(mr MatchRule) Block {
	return mr.Body
}

func NEPToExpr(nep NEPair) Expr {
	return nep.Expr
}

func collectExprRel(expr Expr) []UniRel {
	colB := (func(_r0 Block) []UniRel {
		return collectBlock(collectExprRel, (func(_r0 Stmt) []UniRel { return collectStmtRel(collectExprRel, _r0) }), _r0)
	})
	switch _v421 := (expr).(type) {
	case Expr_EFunCall:
		fc := _v421.Value
		inside := frt.Pipe(slice.Map(collectExprRel, fc.Args), slice.Concat)
		return frt.Pipe(collectFunCall(fc), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
	case Expr_ESlice:
		es := _v421.Value
		inside := frt.Pipe(slice.Map(collectExprRel, es), slice.Concat)
		return frt.Pipe(collectSlice(es), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
	case Expr_ERecordGen:
		rg := _v421.Value
		return frt.Pipe(frt.Pipe(slice.Map(NEPToExpr, rg.FieldsNV), (func(_r0 []Expr) [][]UniRel { return slice.Map(collectExprRel, _r0) })), slice.Concat)
	case Expr_ELazyBlock:
		lb := _v421.Value
		return colB(lb.Block)
	case Expr_EReturnableExpr:
		re := _v421.Value
		switch _v422 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v422.Value
			return colB(bl)
		case ReturnableExpr_RMatchExpr:
			me := _v422.Value
			return frt.Pipe(frt.Pipe(frt.Pipe(slice.Map(mrToBlock, me.Rules), (func(_r0 []Block) [][]UniRel { return slice.Map(colB, _r0) })), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(collectExprRel(me.Target), _r0) }))
		default:
			panic("Union pattern fail. Never reached here.")
		}
	default:
		return emptyRels()
	}
}

func lfdRetType(lfd LetFuncDef) FType {
	switch _v423 := (lfd.Fvar.Ftype).(type) {
	case FType_FFunc:
		ft := _v423.Value
		return freturn(ft)
	default:
		frt.Panic("LetFuncDef's fvar is not FFunc type.")
		return New_FType_FUnit
	}
}

func collectLfdRels(lfd LetFuncDef) []UniRel {
	brels := frt.Pipe(blockToExpr(lfd.Body), collectExprRel)
	lastExprType := frt.Pipe(lfd.Body.FinalExpr, ExprToType)
	return frt.Pipe(unifyType(lfdRetType(lfd), lastExprType), (func(_r0 []UniRel) []UniRel { return slice.Append(brels, _r0) }))
}

type EquivInfo struct {
	eset    EquivSet
	resType FType
}

func eiUnion(e1 EquivInfo, e2 EquivInfo) frt.Tuple2[EquivInfo, []UniRel] {
	nset := eqsUnion(e1.eset, e2.eset)
	nres, rels := frt.Destr(compositeTp(e1.resType, e2.resType))
	nei := EquivInfo{eset: nset, resType: nres}
	return frt.NewTuple2(nei, rels)
}

func eiUpdateResT(e EquivInfo, tcan FType) frt.Tuple2[EquivInfo, []UniRel] {
	nres, rels := frt.Destr(compositeTp(e.resType, tcan))
	nei := EquivInfo{eset: e.eset, resType: nres}
	return frt.NewTuple2(nei, rels)
}

func eiInit(tv TypeVar) EquivInfo {
	es := NewEquivSet(tv)
	rtype := New_FType_FTypeVar(tv)
	return EquivInfo{eset: es, resType: rtype}
}

func rsLookupEI(res Resolver, tvname string) EquivInfo {
	ei, ok := frt.Destr(eidLookup(res.eid, tvname))
	return frt.IfElse(ok, (func() EquivInfo {
		return ei
	}), (func() EquivInfo {
		return eiInit(TypeVar{Name: tvname})
	}))
}

func rsRegisterTo(res Resolver, ei EquivInfo, key string) {
	eidPut(res.eid, key, ei)
}

func rsRegisterNewEI(res Resolver, ei EquivInfo) {
	frt.PipeUnit(eqsItems(ei.eset), (func(_r0 []string) { slice.Iter((func(_r0 string) { rsRegisterTo(res, ei, _r0) }), _r0) }))
}

func updateResOne(res Resolver, rel UniRel) []UniRel {
	ei1 := rsLookupEI(res, rel.SrcV)
	switch _v424 := (rel.Dest).(type) {
	case FType_FTypeVar:
		tvd := _v424.Value
		ei2 := rsLookupEI(res, tvd.Name)
		nei, rels := frt.Destr(eiUnion(ei1, ei2))
		rsRegisterNewEI(res, nei)
		return rels
	default:
		nei, rels := frt.Destr(eiUpdateResT(ei1, rel.Dest))
		return frt.IfElse(slice.IsEmpty(rels), (func() []UniRel {
			return emptyRels()
		}), (func() []UniRel {
			rsRegisterNewEI(res, nei)
			return rels
		}))
	}
}

func updateResolver(res Resolver, rels []UniRel) Resolver {
	nrels := frt.Pipe(frt.Pipe(rels, (func(_r0 []UniRel) [][]UniRel {
		return slice.Map((func(_r0 UniRel) []UniRel { return updateResOne(res, _r0) }), _r0)
	})), slice.Concat)
	return frt.IfElse(slice.IsEmpty(nrels), (func() Resolver {
		return res
	}), (func() Resolver {
		return updateResolver(res, nrels)
	}))
}

func rsResolveType(resT func(FType) FType, res Resolver, tvname string) FType {
	ei := rsLookupEI(res, tvname)
	rcand := ei.resType
	switch _v425 := (rcand).(type) {
	case FType_FTypeVar:
		tv := _v425.Value
		return frt.IfElse(frt.OpEqual(tv.Name, tvname), (func() FType {
			return rcand
		}), (func() FType {
			return resT(rcand)
		}))
	default:
		return resT(rcand)
	}
}

func transOneTypeVar(resT func(FType) FType, rsv Resolver, tv TypeVar) FType {
	return rsResolveType(resT, rsv, tv.Name)
}

func resolveType(rsv Resolver, ftp FType) FType {
	resT := (func(_r0 FType) FType { return resolveType(rsv, _r0) })
	return transTypeVarFType((func(_r0 TypeVar) FType { return transOneTypeVar(resT, rsv, _r0) }), ftp)
}

func transExprNE(cnv func(Expr) Expr, p NEPair) NEPair {
	return NEPair{Name: p.Name, Expr: cnv(p.Expr)}
}

func transVarStmt(transV func(Var) Var, transE func(Expr) Expr, stmt Stmt) Stmt {
	switch _v426 := (stmt).(type) {
	case Stmt_SLetVarDef:
		llvd := _v426.Value
		switch _v427 := (llvd).(type) {
		case LLetVarDef_LLOneVarDef:
			lvd := _v427.Value
			nvar := transV(lvd.Lvar)
			nrhs := transE(lvd.Rhs)
			return frt.Pipe(frt.Pipe(LetVarDef{Lvar: nvar, Rhs: nrhs}, New_LLetVarDef_LLOneVarDef), New_Stmt_SLetVarDef)
		case LLetVarDef_LLDestVarDef:
			ldvd := _v427.Value
			nvars := slice.Map(transV, ldvd.Lvars)
			nrhs := transE(ldvd.Rhs)
			return frt.Pipe(frt.Pipe(LetDestVarDef{Lvars: nvars, Rhs: nrhs}, New_LLetVarDef_LLDestVarDef), New_Stmt_SLetVarDef)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Stmt_SExprStmt:
		e := _v426.Value
		return frt.Pipe(transE(e), New_Stmt_SExprStmt)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func transExprMatchRule(pExpr func(Expr) Expr, mr MatchRule) MatchRule {
	nbody := frt.Pipe(frt.Pipe(blockToExpr(mr.Body), pExpr), exprToBlock)
	return MatchRule{Pattern: mr.Pattern, Body: nbody}
}

func transVarBlock(transE func(Expr) Expr, transS func(Stmt) Stmt, bl Block) Block {
	nss := frt.Pipe(bl.Stmts, (func(_r0 []Stmt) []Stmt { return slice.Map(transS, _r0) }))
	fexpr := transE(bl.FinalExpr)
	return Block{Stmts: nss, FinalExpr: fexpr}
}

func transVarExpr(transV func(Var) Var, expr Expr) Expr {
	transE := (func(_r0 Expr) Expr { return transVarExpr(transV, _r0) })
	transS := (func(_r0 Stmt) Stmt { return transVarStmt(transV, transE, _r0) })
	switch _v428 := (expr).(type) {
	case Expr_EVar:
		v := _v428.Value
		return frt.Pipe(transV(v), New_Expr_EVar)
	case Expr_ESlice:
		es := _v428.Value
		return frt.Pipe(slice.Map(transE, es), New_Expr_ESlice)
	case Expr_EBinOpCall:
		bop := _v428.Value
		nlhs := transE(bop.Lhs)
		nrhs := transE(bop.Rhs)
		return frt.Pipe(BinOpCall{Op: bop.Op, Rtype: bop.Rtype, Lhs: nlhs, Rhs: nrhs}, New_Expr_EBinOpCall)
	case Expr_ETupleExpr:
		es := _v428.Value
		return frt.Pipe(slice.Map((func(_r0 Expr) Expr { return transVarExpr(transV, _r0) }), es), New_Expr_ETupleExpr)
	case Expr_ERecordGen:
		rg := _v428.Value
		newNV := slice.Map((func(_r0 NEPair) NEPair { return transExprNE(transE, _r0) }), rg.FieldsNV)
		return frt.Pipe(RecordGen{FieldsNV: newNV, RecordType: rg.RecordType}, New_Expr_ERecordGen)
	case Expr_ELazyBlock:
		lb := _v428.Value
		nbl := transVarBlock(transE, transS, lb.Block)
		return frt.Pipe(LazyBlock{Block: nbl}, New_Expr_ELazyBlock)
	case Expr_EReturnableExpr:
		re := _v428.Value
		switch _v429 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v429.Value
			return frt.Pipe(transVarBlock(transE, transS, bl), blockToExpr)
		case ReturnableExpr_RMatchExpr:
			me := _v429.Value
			ntarget := transE(me.Target)
			nrules := slice.Map((func(_r0 MatchRule) MatchRule { return transExprMatchRule(transE, _r0) }), me.Rules)
			return frt.Pipe(frt.Pipe(MatchExpr{Target: ntarget, Rules: nrules}, New_ReturnableExpr_RMatchExpr), New_Expr_EReturnableExpr)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_EFunCall:
		fc := _v428.Value
		ntarget := transV(fc.TargetFunc)
		nargs := slice.Map(transE, fc.Args)
		return frt.Pipe(FunCall{TargetFunc: ntarget, Args: nargs}, New_Expr_EFunCall)
	case Expr_EBoolLiteral:
		return expr
	case Expr_EGoEvalExpr:
		return expr
	case Expr_EStringLiteral:
		return expr
	case Expr_EIntImm:
		return expr
	case Expr_EUnit:
		return expr
	case Expr_EFieldAccess:
		fa := _v428.Value
		ntarget := transE(fa.TargetExpr)
		return frt.Pipe(FieldAccess{TargetExpr: ntarget, FieldName: fa.FieldName}, New_Expr_EFieldAccess)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func resolveVarType(rsv Resolver, v Var) Var {
	return Var{Name: v.Name, Ftype: resolveType(rsv, v.Ftype)}
}

func resolveExprType(rsv Resolver, expr Expr) Expr {
	return transVarExpr((func(_r0 Var) Var { return resolveVarType(rsv, _r0) }), expr)
}

func resolveBlockType(rsv Resolver, bl Block) Block {
	return frt.Pipe(frt.Pipe(blockToExpr(bl), (func(_r0 Expr) Expr { return resolveExprType(rsv, _r0) })), exprToBlock)
}

func resolveLfd(rsv Resolver, lfd LetFuncDef) LetFuncDef {
	nfvar := resolveVarType(rsv, lfd.Fvar)
	nparams := slice.Map((func(_r0 Var) Var { return resolveVarType(rsv, _r0) }), lfd.Params)
	nbody := resolveBlockType(rsv, lfd.Body)
	return LetFuncDef{Fvar: nfvar, Params: nparams, Body: nbody}
}

func notFound(rsv Resolver, key string) bool {
	resT := (func(_r0 FType) FType { return resolveType(rsv, _r0) })
	rtype := rsResolveType(resT, rsv, key)
	switch _v430 := (rtype).(type) {
	case FType_FTypeVar:
		tv := _v430.Value
		return frt.OpEqual(tv.Name, key)
	default:
		return false
	}
}

func InferExpr(tvc TypeVarCtx, expr Expr) Expr {
	rels := collectExprRel(expr)
	updateResolver(tvc.resolver, rels)
	return resolveExprType(tvc.resolver, expr)
}

func InferLfd(tvc TypeVarCtx, lfd LetFuncDef) RootFuncDef {
	rels := collectLfdRels(lfd)
	updateResolver(tvc.resolver, rels)
	tps := []string{}
	nlfd := resolveLfd(tvc.resolver, lfd)
	return RootFuncDef{Tparams: tps, Lfd: nlfd}
}
