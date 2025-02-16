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
	switch _v281 := (lhs).(type) {
	case FType_FTypeVar:
		tv := _v281.Value
		switch _v282 := (rhs).(type) {
		case FType_FTypeVar:
			tv2 := _v282.Value
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
		switch _v283 := (rhs).(type) {
		case FType_FTypeVar:
			tv2 := _v283.Value
			return ([]UniRel{UniRel{srcV: tv2.name, dest: lhs}})
		case FType_FSlice:
			ts2 := _v283.Value
			ts1 := lhs.(FType_FSlice).Value
			return unifyType(ts1.elemType, ts2.elemType)
		case FType_FFunc:
			tf2 := _v283.Value
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
	switch _v284 := (stmt).(type) {
	case Stmt_SExprStmt:
		se := _v284.Value
		return ec(se)
	case Stmt_SLetVarDef:
		lvd := _v284.Value
		inside := ec(lvd.rhs)
		one := unifyType(lvd.lvar.ftype, ExprToType(lvd.rhs))
		return slice.Append(one, inside)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func collectFunCall(fc FunCall) []UniRel {
	tftype := fc.targetFunc.ftype
	switch _v285 := (tftype).(type) {
	case FType_FFunc:
		fft := _v285.Value
		return frt.Pipe(frt.Pipe(frt.Pipe(slice.Map(ExprToType, fc.args), (func(_r0 []FType) []frt.Tuple2[FType, FType] { return slice.Zip(fargs(fft), _r0) })), (func(_r0 []frt.Tuple2[FType, FType]) [][]UniRel { return slice.Map(unifyTupArg, _r0) })), slice.Concat)
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
	switch _v286 := (expr).(type) {
	case Expr_EFunCall:
		fc := _v286.Value
		inside := frt.Pipe(slice.Map(collectExprRel, fc.args), slice.Concat)
		return frt.Pipe(collectFunCall(fc), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
	case Expr_ESlice:
		es := _v286.Value
		inside := frt.Pipe(slice.Map(collectExprRel, es), slice.Concat)
		return frt.Pipe(collectSlice(es), (func(_r0 []UniRel) []UniRel { return slice.Append(inside, _r0) }))
	case Expr_ERecordGen:
		rg := _v286.Value
		return frt.Pipe(frt.Pipe(slice.Map(NEPToExpr, rg.fieldsNV), (func(_r0 []Expr) [][]UniRel { return slice.Map(collectExprRel, _r0) })), slice.Concat)
	case Expr_ELazyBlock:
		lb := _v286.Value
		return colB(lb.block)
	case Expr_EReturnableExpr:
		re := _v286.Value
		switch _v287 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v287.Value
			return colB(bl)
		case ReturnableExpr_RMatchExpr:
			me := _v287.Value
			return frt.Pipe(frt.Pipe(frt.Pipe(slice.Map(mrToBlock, me.rules), (func(_r0 []Block) [][]UniRel { return slice.Map(colB, _r0) })), slice.Concat), (func(_r0 []UniRel) []UniRel { return slice.Append(collectExprRel(me.target), _r0) }))
		default:
			panic("Union pattern fail. Never reached here.")
		}
	default:
		return emptyRels()
	}
}

func lfdRetType(lfd LetFuncDef) FType {
	switch _v288 := (lfd.fvar.ftype).(type) {
	case FType_FFunc:
		ft := _v288.Value
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

func transTypeVarFType(transTV func(TypeVar) FType, ftp FType) FType {
	recurse := (func(_r0 FType) FType { return transTypeVarFType(transTV, _r0) })
	switch _v289 := (ftp).(type) {
	case FType_FTypeVar:
		tv := _v289.Value
		return transTV(tv)
	case FType_FSlice:
		ts := _v289.Value
		et := recurse(ts.elemType)
		return New_FType_FSlice(SliceType{elemType: et})
	case FType_FFunc:
		fnt := _v289.Value
		nts := slice.Map(recurse, fnt.targets)
		return frt.Pipe(FuncType{targets: nts}, New_FType_FFunc)
	default:
		return ftp
	}
}

func transOneTypeVar(tdict TypeDict, tv TypeVar) FType {
	return tdLookupNF(tdict, tv.name)
}

func resolveType(tdict TypeDict, ftp FType) FType {
	return transTypeVarFType((func(_r0 TypeVar) FType { return transOneTypeVar(tdict, _r0) }), ftp)
}

func transExprNE(cnv func(Expr) Expr, p NEPair) NEPair {
	return NEPair{name: p.name, expr: cnv(p.expr)}
}

func transVarStmt(transV func(Var) Var, transE func(Expr) Expr, stmt Stmt) Stmt {
	switch _v290 := (stmt).(type) {
	case Stmt_SLetVarDef:
		lvd := _v290.Value
		nvar := transV(lvd.lvar)
		nrhs := transE(lvd.rhs)
		return frt.Pipe(LetVarDef{lvar: nvar, rhs: nrhs}, New_Stmt_SLetVarDef)
	case Stmt_SExprStmt:
		e := _v290.Value
		return frt.Pipe(transE(e), New_Stmt_SExprStmt)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func transExprMatchRule(pExpr func(Expr) Expr, mr MatchRule) MatchRule {
	nbody := frt.Pipe(frt.Pipe(blockToExpr(mr.body), pExpr), exprToBlock)
	return MatchRule{pattern: mr.pattern, body: nbody}
}

func transVarExpr(transV func(Var) Var, expr Expr) Expr {
	transE := (func(_r0 Expr) Expr { return transVarExpr(transV, _r0) })
	switch _v291 := (expr).(type) {
	case Expr_EVar:
		v := _v291.Value
		return frt.Pipe(transV(v), New_Expr_EVar)
	case Expr_ESlice:
		es := _v291.Value
		return frt.Pipe(slice.Map((func(_r0 Expr) Expr { return transVarExpr(transV, _r0) }), es), New_Expr_ESlice)
	case Expr_ERecordGen:
		rg := _v291.Value
		newNV := slice.Map((func(_r0 NEPair) NEPair { return transExprNE(transE, _r0) }), rg.fieldsNV)
		return frt.Pipe(RecordGen{fieldsNV: newNV, recordType: rg.recordType}, New_Expr_ERecordGen)
	case Expr_EReturnableExpr:
		re := _v291.Value
		switch _v292 := (re).(type) {
		case ReturnableExpr_RBlock:
			bl := _v292.Value
			nss := frt.Pipe(bl.stmts, (func(_r0 []Stmt) []Stmt {
				return slice.Map((func(_r0 Stmt) Stmt { return transVarStmt(transV, transE, _r0) }), _r0)
			}))
			fexpr := transE(bl.finalExpr)
			return frt.Pipe(Block{stmts: nss, finalExpr: fexpr}, blockToExpr)
		case ReturnableExpr_RMatchExpr:
			me := _v292.Value
			ntarget := transE(me.target)
			nrules := slice.Map((func(_r0 MatchRule) MatchRule { return transExprMatchRule(transE, _r0) }), me.rules)
			return frt.Pipe(frt.Pipe(MatchExpr{target: ntarget, rules: nrules}, New_ReturnableExpr_RMatchExpr), New_Expr_EReturnableExpr)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Expr_EFunCall:
		fc := _v291.Value
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

func resolveVarType(tdict TypeDict, v Var) Var {
	return Var{name: v.name, ftype: resolveType(tdict, v.ftype)}
}

func resolveExprType(tdict TypeDict, expr Expr) Expr {
	return transVarExpr((func(_r0 Var) Var { return resolveVarType(tdict, _r0) }), expr)
}

func resolveBlockType(tdict TypeDict, bl Block) Block {
	return frt.Pipe(frt.Pipe(blockToExpr(bl), (func(_r0 Expr) Expr { return resolveExprType(tdict, _r0) })), exprToBlock)
}

func resolveLfd(tdict TypeDict, lfd LetFuncDef) LetFuncDef {
	nfvar := resolveVarType(tdict, lfd.fvar)
	nparams := slice.Map((func(_r0 Var) Var { return resolveVarType(tdict, _r0) }), lfd.params)
	nbody := resolveBlockType(tdict, lfd.body)
	return LetFuncDef{fvar: nfvar, params: nparams, body: nbody}
}

func notFound(dic TypeDict, key string) bool {
	_, ok := frt.Destr(tdLookup(dic, key))
	return frt.OpNot(ok)
}

func Infer(tvnames []string, lfd LetFuncDef) RootFuncDef {
	rels := collectLfdRels(lfd)
	rdict := relsToTDict(rels)
	tps := []string{}
	nlfd := resolveLfd(rdict, lfd)
	return RootFuncDef{tparams: tps, lfd: nlfd}
}
