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

func relToTp(rel UniRel) frt.Tuple2[string, FType] {
	return frt.NewTuple2(rel.srcV, rel.dest)
}

func relsToTDict(rels []UniRel) TypeDict {
	return frt.Pipe(slice.Map(relToTp, rels), toTDict)
}

func tupApply(f func(FType, FType) []UniRel, tup frt.Tuple2[FType, FType]) []UniRel {
	lhs, rhs := frt.Destr(tup)
	return f(lhs, rhs)
}

func unifyType(lhs FType, rhs FType) []UniRel {
	switch _v253 := (lhs).(type) {
	case FType_FTypeVar:
		tv := _v253.Value
		switch _v254 := (rhs).(type) {
		case FType_FTypeVar:
			tv2 := _v254.Value
			return frt.IfElse(frt.OpEqual(tv.name, tv2.name), (func() []UniRel {
				return emptyRels()
			}), (func() []UniRel {
				return frt.IfElse((tv.name > tv2.name), (func() []UniRel {
					return ([]UniRel{UniRel{srcV: tv.name, dest: rhs}})
				}), (func() []UniRel {
					return ([]UniRel{UniRel{srcV: tv2.name, dest: lhs}})
				}))
			}))
		default:
			return ([]UniRel{UniRel{srcV: tv.name, dest: rhs}})
		}
	default:
		switch _v255 := (rhs).(type) {
		case FType_FTypeVar:
			tv2 := _v255.Value
			return ([]UniRel{UniRel{srcV: tv2.name, dest: lhs}})
		case FType_FSlice:
			ts2 := _v255.Value
			ts1 := lhs.(FType_FSlice).Value
			return unifyType(ts1.elemType, ts2.elemType)
		case FType_FFunc:
			tf2 := _v255.Value
			tf1 := lhs.(FType_FFunc).Value
			return frt.Pipe(frt.Pipe(slice.Zip(tf1.targets, tf2.targets), (func(_r0 []frt.Tuple2[FType, FType]) [][]UniRel {
				return slice.Map((func(_r0 frt.Tuple2[FType, FType]) []UniRel { return tupApply(unifyType, _r0) }), _r0)
			})), slice.Concat)
		default:
			return emptyRels()
		}
	}
}

func unifyTupArg(tup frt.Tuple2[FType, FType]) []UniRel {
	lhs, rhs := frt.Destr(tup)
	return unifyType(lhs, rhs)
}

func collectStmtRel(ec func(Expr) []UniRel, stmt Stmt) []UniRel {
	switch _v256 := (stmt).(type) {
	case Stmt_SExprStmt:
		se := _v256.Value
		return ec(se)
	case Stmt_SLetVarDef:
		lvd := _v256.Value
		inside := ec(lvd.rhs)
		one := unifyType(lvd.lvar.ftype, ExprToType(lvd.rhs))
		return slice.Append(one, inside)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func collectExprRel(expr Expr) []UniRel {
	switch _v257 := (expr).(type) {
	case Expr_EFunCall:
		fc := _v257.Value
		inside := frt.Pipe(slice.Map(collectExprRel, fc.args), slice.Concat)
		tftype := fc.targetFunc.ftype
		switch _v258 := (tftype).(type) {
		case FType_FFunc:
			fft := _v258.Value
			return frt.Pipe(frt.Pipe(frt.Pipe(frt.Pipe(slice.Map(ExprToType, fc.args), (func(_r0 []FType) []frt.Tuple2[FType, FType] { return slice.Zip(fargs(fft), _r0) })), (func(_r0 []frt.Tuple2[FType, FType]) [][]UniRel { return slice.Map(unifyTupArg, _r0) })), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
		default:
			frt.Panic("funcall with non func first arg, possibly TypeVar, NYI.")
			return emptyRels()
		}
	case Expr_ESlice:
		es := _v257.Value
		inside := frt.Pipe(slice.Map(collectExprRel, es), slice.Concat)
		return frt.IfElse((slice.Length(es) <= 1), (func() []UniRel {
			return inside
		}), (func() []UniRel {
			headT := frt.Pipe(slice.Head(es), ExprToType)
			return frt.Pipe(frt.Pipe(frt.Pipe(frt.Pipe(slice.Tail(es), (func(_r0 []Expr) []FType { return slice.Map(ExprToType, _r0) })), (func(_r0 []FType) [][]UniRel {
				return slice.Map((func(_r0 FType) []UniRel { return unifyType(headT, _r0) }), _r0)
			})), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
		}))
	case Expr_ERecordGen:
		return emptyRels()
	case Expr_ELazyBlock:
		lb := _v257.Value
		return frt.Pipe(frt.Pipe(slice.Map((func(_r0 Stmt) []UniRel { return collectStmtRel(collectExprRel, _r0) }), lb.stmts), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(collectExprRel(lb.finalExpr), _r0) }))
	case Expr_EReturnableExpr:
		re := _v257.Value
		switch _v259 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v259.Value
			return frt.Pipe(frt.Pipe(slice.Map((func(_r0 Stmt) []UniRel { return collectStmtRel(collectExprRel, _r0) }), bl.stmts), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(collectExprRel(bl.finalExpr), _r0) }))
		case ReturnableExpr_RMatchExpr:
			return emptyRels()
		default:
			panic("Union pattern fail. Never reached here.")
		}
	default:
		return emptyRels()
	}
}

func lfdRetType(lfd LetFuncDef) FType {
	switch _v260 := (lfd.fvar.ftype).(type) {
	case FType_FFunc:
		ft := _v260.Value
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

func notFound(dic TypeDict, key string) bool {
	_, ok := frt.Destr(tdLookup(dic, key))
	return frt.OpNot(ok)
}

func resolveType(tdict TypeDict, ftp FType) FType {
	switch _v261 := (ftp).(type) {
	case FType_FTypeVar:
		tv := _v261.Value
		return tdLookupNF(tdict, tv.name)
	case FType_FSlice:
		ts := _v261.Value
		et := resolveType(tdict, ts.elemType)
		return New_FType_FSlice(SliceType{elemType: et})
	case FType_FFunc:
		fnt := _v261.Value
		nts := slice.Map((func(_r0 FType) FType { return resolveType(tdict, _r0) }), fnt.targets)
		return frt.Pipe(FuncType{targets: nts}, New_FType_FFunc)
	default:
		return ftp
	}
}

func resolveVarType(tdict TypeDict, v Var) Var {
	return Var{name: v.name, ftype: resolveType(tdict, v.ftype)}
}

func resolveStmtType(resExpr func(Expr) Expr, tdict TypeDict, stmt Stmt) Stmt {
	switch _v262 := (stmt).(type) {
	case Stmt_SLetVarDef:
		lvd := _v262.Value
		nvar := resolveVarType(tdict, lvd.lvar)
		nrhs := resExpr(lvd.rhs)
		return frt.Pipe(LetVarDef{lvar: nvar, rhs: nrhs}, New_Stmt_SLetVarDef)
	case Stmt_SExprStmt:
		e := _v262.Value
		return frt.Pipe(resExpr(e), New_Stmt_SExprStmt)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func resolveBlockType(resExpr func(Expr) Expr, tdict TypeDict, bl Block) Block {
	nstmts := slice.Map((func(_r0 Stmt) Stmt { return resolveStmtType(resExpr, tdict, _r0) }), bl.stmts)
	nfexpr := resExpr(bl.finalExpr)
	return Block{stmts: nstmts, finalExpr: nfexpr}
}

func resolveExprType(tdict TypeDict, expr Expr) Expr {
	switch _v263 := (expr).(type) {
	case Expr_EVar:
		v := _v263.Value
		return frt.Pipe(resolveVarType(tdict, v), New_Expr_EVar)
	case Expr_EFunCall:
		fc := _v263.Value
		ntf := resolveVarType(tdict, fc.targetFunc)
		nargs := slice.Map((func(_r0 Expr) Expr { return resolveExprType(tdict, _r0) }), fc.args)
		return New_Expr_EFunCall(FunCall{targetFunc: ntf, args: nargs})
	case Expr_EReturnableExpr:
		re := _v263.Value
		switch _v264 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v264.Value
			return frt.Pipe(resolveBlockType((func(_r0 Expr) Expr { return resolveExprType(tdict, _r0) }), tdict, bl), blockToExpr)
		default:
			return expr
		}
	default:
		return expr
	}
}

func resolveLfd(tdict TypeDict, lfd LetFuncDef) LetFuncDef {
	nfvar := resolveVarType(tdict, lfd.fvar)
	nparams := slice.Map((func(_r0 Var) Var { return resolveVarType(tdict, _r0) }), lfd.params)
	nbody := resolveBlockType((func(_r0 Expr) Expr { return resolveExprType(tdict, _r0) }), tdict, lfd.body)
	return LetFuncDef{fvar: nfvar, params: nparams, body: nbody}
}

func Infer(tvnames []string, lfd LetFuncDef) RootFuncDef {
	rels := collectLfdRels(lfd)
	rdict := relsToTDict(rels)
	tps := []string{}
	nlfd := resolveLfd(rdict, lfd)
	return RootFuncDef{tparams: tps, lfd: nlfd}
}
