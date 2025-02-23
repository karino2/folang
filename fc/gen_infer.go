package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/dict"

type UniRel struct {
	SrcV string
	Dest FType
}

func emptyRels() []UniRel {
	return frt.Empty[[]UniRel]()
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
	switch _v1 := (lhs).(type) {
	case FType_FTypeVar:
		tv := _v1.Value
		switch _v2 := (rhs).(type) {
		case FType_FTypeVar:
			tv2 := _v2.Value
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
		switch _v3 := (rhs).(type) {
		case FType_FTypeVar:
			tv2 := _v3.Value
			return frt.Pipe(([]UniRel{UniRel{SrcV: tv2.Name, Dest: lhs}}), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
		case FType_FSlice:
			ts2 := _v3.Value
			switch _v4 := (lhs).(type) {
			case FType_FSlice:
				ts1 := _v4.Value
				rtp, rels := frt.Destr(compositeTp(ts1.ElemType, ts2.ElemType))
				return frt.Pipe(frt.Pipe(SliceType{ElemType: rtp}, New_FType_FSlice), (func(_r0 FType) frt.Tuple2[FType, []UniRel] { return withRels(rels, _r0) }))
			case FType_FFieldAccess:
				return frt.Pipe(emptyRels(), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(rhs, _r0) }))
			default:
				frt.Panic("right is slice, left is neither slice nor field access.")
				return frt.Pipe(emptyRels(), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
			}
		case FType_FFieldAccess:
			fa2 := _v3.Value
			fa22 := faResolve(fa2)
			switch (fa22).(type) {
			case FType_FFieldAccess:
				switch _v5 := (lhs).(type) {
				case FType_FFieldAccess:
					fa1 := _v5.Value
					fa12 := faResolve(fa1)
					switch (fa12).(type) {
					case FType_FFieldAccess:
						rtp, rels := frt.Destr(compositeTp(fa1.RecType, fa2.RecType))
						return frt.Pipe(frt.Pipe(FieldAccessType{RecType: rtp, FieldName: fa1.FieldName}, faResolve), (func(_r0 FType) frt.Tuple2[FType, []UniRel] { return withRels(rels, _r0) }))
					default:
						return compositeTp(fa12, rhs)
					}
				case FType_FSlice:
					return frt.Pipe(emptyRels(), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
				default:
					return frt.Pipe(emptyRels(), (func(_r0 []UniRel) frt.Tuple2[FType, []UniRel] { return withTp(lhs, _r0) }))
				}
			default:
				return compositeTp(lhs, fa22)
			}
		case FType_FFunc:
			tf2 := _v3.Value
			tf1 := lhs.(FType_FFunc).Value
			tps, rels := frt.Destr(compositeTpList(compositeTp, tf1.Targets, tf2.Targets))
			return frt.Pipe(newFFunc(tps), (func(_r0 FType) frt.Tuple2[FType, []UniRel] { return withRels(rels, _r0) }))
		case FType_FParamd:
			pt2 := _v3.Value
			pt1 := lhs.(FType_FParamd).Value
			tps, rels := frt.Destr(compositeTpList(compositeTp, pt1.Targs, pt2.Targs))
			return frt.Pipe(frt.Pipe(ParamdType{Name: pt1.Name, Targs: tps}, New_FType_FParamd), (func(_r0 FType) frt.Tuple2[FType, []UniRel] { return withRels(rels, _r0) }))
		case FType_FTuple:
			tt2 := _v3.Value
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

func varsToTupleType(vars []Var) FType {
	ets := slice.Map(func(_v1 Var) FType {
		return _v1.Ftype
	}, vars)
	return frt.Pipe(TupleType{ElemTypes: ets}, New_FType_FTuple)
}

func collectStmtRel(ec func(Expr) []UniRel, stmt Stmt) []UniRel {
	switch _v6 := (stmt).(type) {
	case Stmt_SExprStmt:
		se := _v6.Value
		return ec(se)
	case Stmt_SLetVarDef:
		slvd := _v6.Value
		switch _v7 := (slvd).(type) {
		case LLetVarDef_LLOneVarDef:
			lvd := _v7.Value
			inside := ec(lvd.Rhs)
			return frt.Pipe(unifyType(lvd.Lvar.Ftype, ExprToType(lvd.Rhs)), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
		case LLetVarDef_LLDestVarDef:
			ldvd := _v7.Value
			inside := ec(ldvd.Rhs)
			rhtype := ExprToType(ldvd.Rhs)
			switch _v8 := (rhtype).(type) {
			case FType_FTuple:
				ft := _v8.Value
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
	tftype := varRefVarType(fc.TargetFunc)
	switch _v9 := (tftype).(type) {
	case FType_FFunc:
		fft := _v9.Value
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

func NEPToNT(nep NEPair) frt.Tuple2[string, FType] {
	return frt.NewTuple2(nep.Name, ExprToType(nep.Expr))
}

func recNTUnify(rec RecordType, ntp frt.Tuple2[string, FType]) []UniRel {
	name, ftp := frt.Destr(ntp)
	rpair := frGetField(rec, name)
	return unifyType(ftp, rpair.Ftype)
}

func collectExprRel(expr Expr) []UniRel {
	colE := collectExprRel
	colB := (func(_r0 Block) []UniRel {
		return collectBlock(colE, (func(_r0 Stmt) []UniRel { return collectStmtRel(colE, _r0) }), _r0)
	})
	switch _v10 := (expr).(type) {
	case Expr_EFunCall:
		fc := _v10.Value
		inside := frt.Pipe(slice.Map(colE, fc.Args), slice.Concat)
		return frt.Pipe(collectFunCall(fc), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
	case Expr_EBinOpCall:
		bop := _v10.Value
		insideL := colE(bop.Lhs)
		insideR := colE(bop.Rhs)
		lft := ExprToType(bop.Lhs)
		rft := ExprToType(bop.Rhs)
		teq := unifyType(lft, rft)
		retEq := (func() []UniRel {
			switch (bop.Rtype).(type) {
			case FType_FBool:
				return emptyRels()
			default:
				return unifyType(bop.Rtype, lft)
			}
		})()
		all := ([][]UniRel{insideL, insideR, teq, retEq})
		return slice.Concat(all)
	case Expr_ETupleExpr:
		tes := _v10.Value
		return frt.Pipe(slice.Map(colE, tes), slice.Concat)
	case Expr_ELambda:
		le := _v10.Value
		return colB(le.Body)
	case Expr_ESlice:
		es := _v10.Value
		inside := frt.Pipe(slice.Map(colE, es), slice.Concat)
		return frt.Pipe(collectSlice(es), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
	case Expr_ERecordGen:
		rg := _v10.Value
		fieldValEs := slice.Map(func(_v1 NEPair) Expr {
			return _v1.Expr
		}, rg.FieldsNV)
		inside := frt.Pipe(frt.Pipe(fieldValEs, (func(_r0 []Expr) [][]UniRel { return slice.Map(colE, _r0) })), slice.Concat)
		return frt.Pipe(frt.Pipe(frt.Pipe(slice.Map(NEPToNT, rg.FieldsNV), (func(_r0 []frt.Tuple2[string, FType]) [][]UniRel {
			return slice.Map((func(_r0 frt.Tuple2[string, FType]) []UniRel { return recNTUnify(rg.RecordType, _r0) }), _r0)
		})), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
	case Expr_ELazyBlock:
		lb := _v10.Value
		return colB(lb.Block)
	case Expr_EReturnableExpr:
		re := _v10.Value
		switch _v11 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v11.Value
			return colB(bl)
		case ReturnableExpr_RMatchExpr:
			me := _v11.Value
			return frt.Pipe(frt.Pipe(frt.Pipe(slice.Map(func(_v2 MatchRule) Block {
				return _v2.Body
			}, me.Rules), (func(_r0 []Block) [][]UniRel { return slice.Map(colB, _r0) })), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(colE(me.Target), _r0) }))
		default:
			panic("Union pattern fail. Never reached here.")
		}
	default:
		return emptyRels()
	}
}

func lfdRetType(lfd LetFuncDef) FType {
	switch _v12 := (lfd.Fvar.Ftype).(type) {
	case FType_FFunc:
		ft := _v12.Value
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

func newEquivSet0() EquivSet {
	dic := dict.New[string, bool]()
	return EquivSet{Dict: dic}
}

func NewEquivSet(tv TypeVar) EquivSet {
	es := newEquivSet0()
	dict.Add(es.Dict, tv.Name, true)
	return es
}

func eqsItems(es EquivSet) []string {
	return dict.Keys(es.Dict)
}

func setAddKeys(d dict.Dict[string, bool], k string) {
	dict.Add(d, k, true)
}

func eqsUnion(es1 EquivSet, es2 EquivSet) EquivSet {
	e3 := newEquivSet0()
	frt.PipeUnit(dict.Keys(es1.Dict), (func(_r0 []string) { slice.Iter((func(_r0 string) { setAddKeys(e3.Dict, _r0) }), _r0) }))
	frt.PipeUnit(dict.Keys(es2.Dict), (func(_r0 []string) { slice.Iter((func(_r0 string) { setAddKeys(e3.Dict, _r0) }), _r0) }))
	return e3
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
	ei, ok := frt.Destr(dict.TryFind(res.eid, tvname))
	return frt.IfElse(ok, (func() EquivInfo {
		return ei
	}), (func() EquivInfo {
		return eiInit(TypeVar{Name: tvname})
	}))
}

func rsRegisterTo(res Resolver, ei EquivInfo, key string) {
	dict.Add(res.eid, key, ei)
}

func rsRegisterNewEI(res Resolver, ei EquivInfo) {
	frt.PipeUnit(eqsItems(ei.eset), (func(_r0 []string) { slice.Iter((func(_r0 string) { rsRegisterTo(res, ei, _r0) }), _r0) }))
}

func updateResOne(res Resolver, rel UniRel) []UniRel {
	ei1 := rsLookupEI(res, rel.SrcV)
	switch _v13 := (rel.Dest).(type) {
	case FType_FTypeVar:
		tvd := _v13.Value
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

func transExprNE(cnv func(Expr) Expr, p NEPair) NEPair {
	return NEPair{Name: p.Name, Expr: cnv(p.Expr)}
}

func transTypeStmt(transTV func(TypeVar) FType, transE func(Expr) Expr, stmt Stmt) Stmt {
	transV := (func(_r0 Var) Var { return transOneVar(transTV, _r0) })
	switch _v14 := (stmt).(type) {
	case Stmt_SLetVarDef:
		llvd := _v14.Value
		switch _v15 := (llvd).(type) {
		case LLetVarDef_LLOneVarDef:
			lvd := _v15.Value
			nvar := transV(lvd.Lvar)
			nrhs := transE(lvd.Rhs)
			return frt.Pipe(frt.Pipe(LetVarDef{Lvar: nvar, Rhs: nrhs}, New_LLetVarDef_LLOneVarDef), New_Stmt_SLetVarDef)
		case LLetVarDef_LLDestVarDef:
			ldvd := _v15.Value
			nvars := slice.Map(transV, ldvd.Lvars)
			nrhs := transE(ldvd.Rhs)
			return frt.Pipe(frt.Pipe(LetDestVarDef{Lvars: nvars, Rhs: nrhs}, New_LLetVarDef_LLDestVarDef), New_Stmt_SLetVarDef)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Stmt_SExprStmt:
		e := _v14.Value
		return frt.Pipe(transE(e), New_Stmt_SExprStmt)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func transExprMatchRule(pExpr func(Expr) Expr, mr MatchRule) MatchRule {
	nbody := frt.Pipe(frt.Pipe(blockToExpr(mr.Body), pExpr), exprToBlock)
	return MatchRule{Pattern: mr.Pattern, Body: nbody}
}

func transTypeBlock(transE func(Expr) Expr, transS func(Stmt) Stmt, bl Block) Block {
	nss := frt.Pipe(bl.Stmts, (func(_r0 []Stmt) []Stmt { return slice.Map(transS, _r0) }))
	fexpr := transE(bl.FinalExpr)
	return Block{Stmts: nss, FinalExpr: fexpr}
}

func transTypeExpr(transTV func(TypeVar) FType, expr Expr) Expr {
	transV := (func(_r0 Var) Var { return transOneVar(transTV, _r0) })
	transT := (func(_r0 FType) FType { return transTypeVarFType(transTV, _r0) })
	transE := (func(_r0 Expr) Expr { return transTypeExpr(transTV, _r0) })
	transS := (func(_r0 Stmt) Stmt { return transTypeStmt(transTV, transE, _r0) })
	transB := (func(_r0 Block) Block { return transTypeBlock(transE, transS, _r0) })
	switch _v16 := (expr).(type) {
	case Expr_EVarRef:
		rv := _v16.Value
		switch _v17 := (rv).(type) {
		case VarRef_VRVar:
			v := _v17.Value
			return frt.Pipe(frt.Pipe(transV(v), New_VarRef_VRVar), New_Expr_EVarRef)
		case VarRef_VRSVar:
			sv := _v17.Value
			nv := transV(sv.Var)
			return frt.Pipe(frt.Pipe(SpecVar{Var: nv, SpecList: sv.SpecList}, New_VarRef_VRSVar), New_Expr_EVarRef)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_ESlice:
		es := _v16.Value
		return frt.Pipe(slice.Map(transE, es), New_Expr_ESlice)
	case Expr_EBinOpCall:
		bop := _v16.Value
		nlhs := transE(bop.Lhs)
		nrhs := transE(bop.Rhs)
		nret := transT(bop.Rtype)
		return frt.Pipe(BinOpCall{Op: bop.Op, Rtype: nret, Lhs: nlhs, Rhs: nrhs}, New_Expr_EBinOpCall)
	case Expr_ETupleExpr:
		es := _v16.Value
		return frt.Pipe(slice.Map(transE, es), New_Expr_ETupleExpr)
	case Expr_ELambda:
		le := _v16.Value
		nparams := slice.Map(transV, le.Params)
		nbody := transB(le.Body)
		return frt.Pipe(LambdaExpr{Params: nparams, Body: nbody}, New_Expr_ELambda)
	case Expr_ERecordGen:
		rg := _v16.Value
		newNV := slice.Map((func(_r0 NEPair) NEPair { return transExprNE(transE, _r0) }), rg.FieldsNV)
		return frt.Pipe(RecordGen{FieldsNV: newNV, RecordType: rg.RecordType}, New_Expr_ERecordGen)
	case Expr_ELazyBlock:
		lb := _v16.Value
		nbl := transB(lb.Block)
		return frt.Pipe(LazyBlock{Block: nbl}, New_Expr_ELazyBlock)
	case Expr_EReturnableExpr:
		re := _v16.Value
		switch _v18 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v18.Value
			return frt.Pipe(transTypeBlock(transE, transS, bl), blockToExpr)
		case ReturnableExpr_RMatchExpr:
			me := _v18.Value
			ntarget := transE(me.Target)
			nrules := slice.Map((func(_r0 MatchRule) MatchRule { return transExprMatchRule(transE, _r0) }), me.Rules)
			return frt.Pipe(frt.Pipe(MatchExpr{Target: ntarget, Rules: nrules}, New_ReturnableExpr_RMatchExpr), New_Expr_EReturnableExpr)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_EFunCall:
		fc := _v16.Value
		ntarget := transVarVR(transV, fc.TargetFunc)
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
		fa := _v16.Value
		ntarget := transE(fa.TargetExpr)
		return frt.Pipe(FieldAccess{TargetExpr: ntarget, FieldName: fa.FieldName}, New_Expr_EFieldAccess)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func transTypeBlockFacade(transTV func(TypeVar) FType, block Block) Block {
	transE := (func(_r0 Expr) Expr { return transTypeExpr(transTV, _r0) })
	transS := (func(_r0 Stmt) Stmt { return transTypeStmt(transTV, transE, _r0) })
	return transTypeBlock(transE, transS, block)
}

func transTypeLfd(transTV func(TypeVar) FType, lfd LetFuncDef) LetFuncDef {
	transV := (func(_r0 Var) Var { return transOneVar(transTV, _r0) })
	nfvar := transV(lfd.Fvar)
	nparams := slice.Map(transV, lfd.Params)
	nbody := transTypeBlockFacade(transTV, lfd.Body)
	return LetFuncDef{Fvar: nfvar, Params: nparams, Body: nbody}
}

func collectTVarStmt(collE func(Expr) []string, stmt Stmt) []string {
	switch _v19 := (stmt).(type) {
	case Stmt_SLetVarDef:
		llvd := _v19.Value
		switch _v20 := (llvd).(type) {
		case LLetVarDef_LLOneVarDef:
			lvd := _v20.Value
			nvar := collectTVarFType(lvd.Lvar.Ftype)
			nrhs := collE(lvd.Rhs)
			return slice.Append(nvar, nrhs)
		case LLetVarDef_LLDestVarDef:
			ldvd := _v20.Value
			nvars := frt.Pipe(slice.Map(func(_v1 Var) FType {
				return _v1.Ftype
			}, ldvd.Lvars), (func(_r0 []FType) []string { return slice.Collect(collectTVarFType, _r0) }))
			nrhs := collE(ldvd.Rhs)
			return slice.Append(nvars, nrhs)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Stmt_SExprStmt:
		e := _v19.Value
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
	switch _v21 := (expr).(type) {
	case Expr_EVarRef:
		vr := _v21.Value
		return frt.Pipe(varRefVarType(vr), collectTVarFType)
	case Expr_ESlice:
		es := _v21.Value
		return slice.Collect(recurse, es)
	case Expr_EBinOpCall:
		bop := _v21.Value
		lres := recurse(bop.Lhs)
		rres := recurse(bop.Rhs)
		return slice.Append(lres, rres)
	case Expr_ETupleExpr:
		es := _v21.Value
		return slice.Collect(recurse, es)
	case Expr_ELambda:
		le := _v21.Value
		pas := frt.Pipe(slice.Map(func(_v1 Var) FType {
			return _v1.Ftype
		}, le.Params), (func(_r0 []FType) []string { return slice.Collect(collectTVarFType, _r0) }))
		return frt.Pipe(collB(le.Body), (func(_r0 []string) []string { return slice.Append(pas, _r0) }))
	case Expr_ERecordGen:
		rg := _v21.Value
		return frt.Pipe(slice.Map(func(_v2 NEPair) Expr {
			return _v2.Expr
		}, rg.FieldsNV), (func(_r0 []Expr) []string { return slice.Collect(recurse, _r0) }))
	case Expr_ELazyBlock:
		lb := _v21.Value
		return collB(lb.Block)
	case Expr_EReturnableExpr:
		re := _v21.Value
		switch _v22 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v22.Value
			return collB(bl)
		case ReturnableExpr_RMatchExpr:
			me := _v22.Value
			return frt.Pipe(frt.Pipe(slice.Map(func(_v3 MatchRule) Block {
				return _v3.Body
			}, me.Rules), (func(_r0 []Block) []string { return slice.Collect(collB, _r0) })), (func(_r0 []string) []string { return slice.Append(recurse(me.Target), _r0) }))
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_EFunCall:
		fc := _v21.Value
		colt := frt.Pipe(varRefVarType(fc.TargetFunc), collectTVarFType)
		return frt.Pipe(slice.Collect(recurse, fc.Args), (func(_r0 []string) []string { return slice.Append(colt, _r0) }))
	case Expr_EFieldAccess:
		fa := _v21.Value
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

func resolveOneTypeVar(rsv Resolver, tv TypeVar) FType {
	recurse := (func(_r0 TypeVar) FType { return resolveOneTypeVar(rsv, _r0) })
	ei := rsLookupEI(rsv, tv.Name)
	rcand := ei.resType
	switch _v23 := (rcand).(type) {
	case FType_FTypeVar:
		tv2 := _v23.Value
		return frt.IfElse(frt.OpEqual(tv2.Name, tv.Name), (func() FType {
			return rcand
		}), (func() FType {
			return transTypeVarFType(recurse, rcand)
		}))
	default:
		return transTypeVarFType(recurse, rcand)
	}
}

func resolveType(rsv Resolver, ftp FType) FType {
	return transTypeVarFType((func(_r0 TypeVar) FType { return resolveOneTypeVar(rsv, _r0) }), ftp)
}

func resolveExprType(rsv Resolver, expr Expr) Expr {
	return transTypeExpr((func(_r0 TypeVar) FType { return resolveOneTypeVar(rsv, _r0) }), expr)
}

func resolveLfd(rsv Resolver, lfd LetFuncDef) LetFuncDef {
	return transTypeLfd((func(_r0 TypeVar) FType { return resolveOneTypeVar(rsv, _r0) }), lfd)
}

func InferExpr(tvc TypeVarCtx, expr Expr) Expr {
	rels := collectExprRel(expr)
	updateResolver(tvc.resolver, rels)
	return resolveExprType(tvc.resolver, expr)
}

func collectTVarLfd(lfd LetFuncDef) []string {
	vres := collectTVarFType(lfd.Fvar.Ftype)
	pres := frt.Pipe(slice.Map(func(_v1 Var) FType {
		return _v1.Ftype
	}, lfd.Params), (func(_r0 []FType) []string { return slice.Collect(collectTVarFType, _r0) }))
	bres := collectTVarBlockFacade(lfd.Body)
	res := ([][]string{vres, pres, bres})
	return slice.Concat(res)
}

func newTName(i int, n string) string {
	return frt.Sprintf1("T%d", i)
}

func replaceSDict(ttdict dict.Dict[string, string], tv TypeVar) FType {
	nname := dict.Item(ttdict, tv.Name)
	return frt.Pipe(TypeVar{Name: nname}, New_FType_FTypeVar)
}

func hoistTVar(unresT []string, lfd LetFuncDef) frt.Tuple2[[]string, LetFuncDef] {
	newTs := slice.Mapi(newTName, unresT)
	ttdict := frt.Pipe(slice.Zip(unresT, newTs), dict.ToDict)
	transTV := (func(_r0 TypeVar) FType { return replaceSDict(ttdict, _r0) })
	nlfd := transTypeLfd(transTV, lfd)
	return frt.NewTuple2(newTs, nlfd)
}

func InferLfd(tvc TypeVarCtx, lfd LetFuncDef) RootFuncDef {
	rels := collectLfdRels(lfd)
	updateResolver(tvc.resolver, rels)
	nlfd := resolveLfd(tvc.resolver, lfd)
	unresTvs := frt.Pipe(collectTVarLfd(nlfd), slice.Distinct)
	newTvs, nlfd2 := frt.Destr(hoistTVar(unresTvs, nlfd))
	return RootFuncDef{Tparams: newTvs, Lfd: nlfd2}
}
