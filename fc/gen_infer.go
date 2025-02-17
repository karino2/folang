package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

type UniRel struct {
	srcV string
	dest FType
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
	switch _v376 := (lhs).(type) {
	case FType_FTypeVar:
		tv := _v376.Value
		switch _v377 := (rhs).(type) {
		case FType_FTypeVar:
			tv2 := _v377.Value
			return frt.IfElse(frt.OpEqual(tv.name, tv2.name), (func() frt.Tuple2[FType, []UniRel] {
				return frt.Pipe(emptyRels(), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
			}), (func() frt.Tuple2[FType, []UniRel] {
				return frt.IfElse((tv.name > tv2.name), (func() frt.Tuple2[FType, []UniRel] {
					return frt.Pipe(([]UniRel{UniRel{srcV: tv.name, dest: rhs}}), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(rhs, _r0) }))
				}), (func() frt.Tuple2[FType, []UniRel] {
					return frt.Pipe(([]UniRel{UniRel{srcV: tv2.name, dest: lhs}}), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
				}))
			}))
		default:
			return frt.Pipe(([]UniRel{UniRel{srcV: tv.name, dest: rhs}}), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(rhs, _r0) }))
		}
	default:
		switch _v378 := (rhs).(type) {
		case FType_FTypeVar:
			tv2 := _v378.Value
			return frt.Pipe(([]UniRel{UniRel{srcV: tv2.name, dest: lhs}}), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
		case FType_FSlice:
			ts2 := _v378.Value
			ts1 := lhs.(FType_FSlice).Value
			rtp, rels := frt.Destr(compositeTp(ts1.elemType, ts2.elemType))
			return frt.Pipe(frt.Pipe(SliceType{elemType: rtp}, New_FType_FSlice), (func(_r0 FType) frt.Tuple2[FType, []UniRel] { return withRels(rels, _r0) }))
		case FType_FFunc:
			tf2 := _v378.Value
			tf1 := lhs.(FType_FFunc).Value
			tps, rels := frt.Destr(compositeTpList(compositeTp, tf1.targets, tf2.targets))
			return frt.Pipe(frt.Pipe(FuncType{targets: tps}, New_FType_FFunc), (func(_r0 FType) frt.Tuple2[FType, []UniRel] { return withRels(rels, _r0) }))
		case FType_FTuple:
			tt2 := _v378.Value
			tt1 := lhs.(FType_FTuple).Value
			tps, rels := frt.Destr(compositeTpList(compositeTp, tt1.elemTypes, tt2.elemTypes))
			return frt.Pipe(frt.Pipe(TupleType{elemTypes: tps}, New_FType_FTuple), (func(_r0 FType) frt.Tuple2[FType, []UniRel] { return withRels(rels, _r0) }))
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
	return frt.IfElse(frt.OpEqual(v.name, "_"), (func() []UniRel {
		return emptyRels()
	}), (func() []UniRel {
		return unifyType(v.ftype, ft)
	}))
}

func collectStmtRel(ec func(Expr) []UniRel, stmt Stmt) []UniRel {
	switch _v379 := (stmt).(type) {
	case Stmt_SExprStmt:
		se := _v379.Value
		return ec(se)
	case Stmt_SLetVarDef:
		slvd := _v379.Value
		switch _v380 := (slvd).(type) {
		case LLetVarDef_LLOneVarDef:
			lvd := _v380.Value
			inside := ec(lvd.rhs)
			return frt.Pipe(unifyType(lvd.lvar.ftype, ExprToType(lvd.rhs)), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
		case LLetVarDef_LLDestVarDef:
			ldvd := _v380.Value
			inside := ec(ldvd.rhs)
			rhtype := ExprToType(ldvd.rhs)
			switch _v381 := (rhtype).(type) {
			case FType_FTuple:
				ft := _v381.Value
				return frt.Pipe(frt.Pipe(frt.Pipe(slice.Zip(ldvd.lvars, ft.elemTypes), (func(_r0 []frt.Tuple2[Var, FType]) [][]UniRel { return slice.Map(unifyVETup, _r0) })), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
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
	tftype := fc.targetFunc.ftype
	switch _v382 := (tftype).(type) {
	case FType_FFunc:
		fft := _v382.Value
		argTps := slice.Map(ExprToType, fc.args)
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
	return frt.Pipe(frt.Pipe(slice.Map(colS, block.stmts), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(colE(block.finalExpr), _r0) }))
}

func mrToBlock(mr MatchRule) Block {
	return mr.body
}

func NEPToExpr(nep NEPair) Expr {
	return nep.expr
}

func collectExprRel(expr Expr) []UniRel {
	colB := (func(_r0 Block) []UniRel {
		return collectBlock(collectExprRel, (func(_r0 Stmt) []UniRel { return collectStmtRel(collectExprRel, _r0) }), _r0)
	})
	switch _v383 := (expr).(type) {
	case Expr_EFunCall:
		fc := _v383.Value
		inside := frt.Pipe(slice.Map(collectExprRel, fc.args), slice.Concat)
		return frt.Pipe(collectFunCall(fc), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
	case Expr_ESlice:
		es := _v383.Value
		inside := frt.Pipe(slice.Map(collectExprRel, es), slice.Concat)
		return frt.Pipe(collectSlice(es), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
	case Expr_ERecordGen:
		rg := _v383.Value
		return frt.Pipe(frt.Pipe(slice.Map(NEPToExpr, rg.fieldsNV), (func(_r0 []Expr) [][]UniRel { return slice.Map(collectExprRel, _r0) })), slice.Concat)
	case Expr_ELazyBlock:
		lb := _v383.Value
		return colB(lb.block)
	case Expr_EReturnableExpr:
		re := _v383.Value
		switch _v384 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v384.Value
			return colB(bl)
		case ReturnableExpr_RMatchExpr:
			me := _v384.Value
			return frt.Pipe(frt.Pipe(frt.Pipe(slice.Map(mrToBlock, me.rules), (func(_r0 []Block) [][]UniRel { return slice.Map(colB, _r0) })), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(collectExprRel(me.target), _r0) }))
		default:
			panic("Union pattern fail. Never reached here.")
		}
	default:
		return emptyRels()
	}
}

func lfdRetType(lfd LetFuncDef) FType {
	switch _v385 := (lfd.fvar.ftype).(type) {
	case FType_FFunc:
		ft := _v385.Value
		return freturn(ft)
	default:
		frt.Panic("LetFuncDef's fvar is not FFunc type.")
		return New_FType_FUnit
	}
}

func collectLfdRels(lfd LetFuncDef) []UniRel {
	brels := frt.Pipe(blockToExpr(lfd.body), collectExprRel)
	lastExprType := frt.Pipe(lfd.body.finalExpr, ExprToType)
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

type Resolver struct {
	eid EquivInfoDict
}

func newResolver() Resolver {
	neid := NewEquivInfoDict()
	return Resolver{eid: neid}
}

func rsLookupEI(res Resolver, tvname string) EquivInfo {
	ei, ok := frt.Destr(eidLookup(res.eid, tvname))
	return frt.IfElse(ok, (func() EquivInfo {
		return ei
	}), (func() EquivInfo {
		return eiInit(TypeVar{name: tvname})
	}))
}

func rsRegisterTo(res Resolver, ei EquivInfo, key string) {
	eidPut(res.eid, key, ei)
}

func rsRegisterNewEI(res Resolver, ei EquivInfo) {
	frt.PipeUnit(eqsItems(ei.eset), (func(_r0 []string) { slice.Iter((func(_r0 string) { rsRegisterTo(res, ei, _r0) }), _r0) }))
}

func updateResOne(res Resolver, rel UniRel) []UniRel {
	ei1 := rsLookupEI(res, rel.srcV)
	switch _v386 := (rel.dest).(type) {
	case FType_FTypeVar:
		tvd := _v386.Value
		ei2 := rsLookupEI(res, tvd.name)
		nei, rels := frt.Destr(eiUnion(ei1, ei2))
		rsRegisterNewEI(res, nei)
		return rels
	default:
		nei, rels := frt.Destr(eiUpdateResT(ei1, rel.dest))
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

func buildResolver(rels []UniRel) Resolver {
	res := newResolver()
	updateResolver(res, rels)
	return res
}

func rsResolveType(resT func(FType) FType, res Resolver, tvname string) FType {
	ei := rsLookupEI(res, tvname)
	rcand := ei.resType
	switch _v387 := (rcand).(type) {
	case FType_FTypeVar:
		tv := _v387.Value
		return frt.IfElse(frt.OpEqual(tv.name, tvname), (func() FType {
			return rcand
		}), (func() FType {
			return resT(rcand)
		}))
	default:
		return resT(rcand)
	}
}

func transOneTypeVar(resT func(FType) FType, rsv Resolver, tv TypeVar) FType {
	return rsResolveType(resT, rsv, tv.name)
}

func resolveType(rsv Resolver, ftp FType) FType {
	resT := (func(_r0 FType) FType { return resolveType(rsv, _r0) })
	return transTypeVarFType((func(_r0 TypeVar) FType { return transOneTypeVar(resT, rsv, _r0) }), ftp)
}

func transExprNE(cnv func(Expr) Expr, p NEPair) NEPair {
	return NEPair{name: p.name, expr: cnv(p.expr)}
}

func transVarStmt(transV func(Var) Var, transE func(Expr) Expr, stmt Stmt) Stmt {
	switch _v388 := (stmt).(type) {
	case Stmt_SLetVarDef:
		llvd := _v388.Value
		switch _v389 := (llvd).(type) {
		case LLetVarDef_LLOneVarDef:
			lvd := _v389.Value
			nvar := transV(lvd.lvar)
			nrhs := transE(lvd.rhs)
			return frt.Pipe(frt.Pipe(LetVarDef{lvar: nvar, rhs: nrhs}, New_LLetVarDef_LLOneVarDef), New_Stmt_SLetVarDef)
		case LLetVarDef_LLDestVarDef:
			ldvd := _v389.Value
			nvars := slice.Map(transV, ldvd.lvars)
			nrhs := transE(ldvd.rhs)
			return frt.Pipe(frt.Pipe(LetDestVarDef{lvars: nvars, rhs: nrhs}, New_LLetVarDef_LLDestVarDef), New_Stmt_SLetVarDef)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Stmt_SExprStmt:
		e := _v388.Value
		return frt.Pipe(transE(e), New_Stmt_SExprStmt)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func transExprMatchRule(pExpr func(Expr) Expr, mr MatchRule) MatchRule {
	nbody := frt.Pipe(frt.Pipe(blockToExpr(mr.body), pExpr), exprToBlock)
	return MatchRule{pattern: mr.pattern, body: nbody}
}

func transVarBlock(transE func(Expr) Expr, transS func(Stmt) Stmt, bl Block) Block {
	nss := frt.Pipe(bl.stmts, (func(_r0 []Stmt) []Stmt { return slice.Map(transS, _r0) }))
	fexpr := transE(bl.finalExpr)
	return Block{stmts: nss, finalExpr: fexpr}
}

func transVarExpr(transV func(Var) Var, expr Expr) Expr {
	transE := (func(_r0 Expr) Expr { return transVarExpr(transV, _r0) })
	transS := (func(_r0 Stmt) Stmt { return transVarStmt(transV, transE, _r0) })
	switch _v390 := (expr).(type) {
	case Expr_EVar:
		v := _v390.Value
		return frt.Pipe(transV(v), New_Expr_EVar)
	case Expr_ESlice:
		es := _v390.Value
		return frt.Pipe(slice.Map(transE, es), New_Expr_ESlice)
	case Expr_EBinOpCall:
		bop := _v390.Value
		nlhs := transE(bop.lhs)
		nrhs := transE(bop.rhs)
		return frt.Pipe(BinOpCall{op: bop.op, rtype: bop.rtype, lhs: nlhs, rhs: nrhs}, New_Expr_EBinOpCall)
	case Expr_ETupleExpr:
		es := _v390.Value
		return frt.Pipe(slice.Map((func(_r0 Expr) Expr { return transVarExpr(transV, _r0) }), es), New_Expr_ETupleExpr)
	case Expr_ERecordGen:
		rg := _v390.Value
		newNV := slice.Map((func(_r0 NEPair) NEPair { return transExprNE(transE, _r0) }), rg.fieldsNV)
		return frt.Pipe(RecordGen{fieldsNV: newNV, recordType: rg.recordType}, New_Expr_ERecordGen)
	case Expr_ELazyBlock:
		lb := _v390.Value
		nbl := transVarBlock(transE, transS, lb.block)
		return frt.Pipe(LazyBlock{block: nbl}, New_Expr_ELazyBlock)
	case Expr_EReturnableExpr:
		re := _v390.Value
		switch _v391 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v391.Value
			return frt.Pipe(transVarBlock(transE, transS, bl), blockToExpr)
		case ReturnableExpr_RMatchExpr:
			me := _v391.Value
			ntarget := transE(me.target)
			nrules := slice.Map((func(_r0 MatchRule) MatchRule { return transExprMatchRule(transE, _r0) }), me.rules)
			return frt.Pipe(frt.Pipe(MatchExpr{target: ntarget, rules: nrules}, New_ReturnableExpr_RMatchExpr), New_Expr_EReturnableExpr)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_EFunCall:
		fc := _v390.Value
		ntarget := transV(fc.targetFunc)
		nargs := slice.Map(transE, fc.args)
		return frt.Pipe(FunCall{targetFunc: ntarget, args: nargs}, New_Expr_EFunCall)
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
		return expr
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func resolveVarType(rsv Resolver, v Var) Var {
	return Var{name: v.name, ftype: resolveType(rsv, v.ftype)}
}

func resolveExprType(rsv Resolver, expr Expr) Expr {
	return transVarExpr((func(_r0 Var) Var { return resolveVarType(rsv, _r0) }), expr)
}

func resolveBlockType(rsv Resolver, bl Block) Block {
	return frt.Pipe(frt.Pipe(blockToExpr(bl), (func(_r0 Expr) Expr { return resolveExprType(rsv, _r0) })), exprToBlock)
}

func resolveLfd(rsv Resolver, lfd LetFuncDef) LetFuncDef {
	nfvar := resolveVarType(rsv, lfd.fvar)
	nparams := slice.Map((func(_r0 Var) Var { return resolveVarType(rsv, _r0) }), lfd.params)
	nbody := resolveBlockType(rsv, lfd.body)
	return LetFuncDef{fvar: nfvar, params: nparams, body: nbody}
}

func notFound(rsv Resolver, key string) bool {
	resT := (func(_r0 FType) FType { return resolveType(rsv, _r0) })
	rtype := rsResolveType(resT, rsv, key)
	switch _v392 := (rtype).(type) {
	case FType_FTypeVar:
		tv := _v392.Value
		return frt.OpEqual(tv.name, key)
	default:
		return false
	}
}

func Infer(tvnames []string, lfd LetFuncDef) RootFuncDef {
	rels := collectLfdRels(lfd)
	rsv := buildResolver(rels)
	tps := []string{}
	nlfd := resolveLfd(rsv, lfd)
	return RootFuncDef{tparams: tps, lfd: nlfd}
}
