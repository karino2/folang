package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func faToType(eToT func(Expr) FType, fa FieldAccess) FType {
	tart := eToT(fa.TargetExpr)
	switch _v1 := (tart).(type) {
	case FType_FRecord:
		rt := _v1.Value
		field := frGetField(rt, fa.FieldName)
		return field.Ftype
	default:
		return frt.Pipe(FieldAccessType{RecType: tart, FieldName: fa.FieldName}, New_FType_FFieldAccess)
	}
}

func lblockReturnType(toT func(Expr) FType, lb LazyBlock) FType {
	return toT(lb.Block.FinalExpr)
}

func lblockToType(toT func(Expr) FType, lb LazyBlock) FType {
	rtype := lblockReturnType(toT, lb)
	return New_FType_FFunc(FuncType{Targets: ([]FType{New_FType_FUnit, rtype})})
}

func blockReturnType(toT func(Expr) FType, block Block) FType {
	return toT(block.FinalExpr)
}

func blockToType(toT func(Expr) FType, b Block) FType {
	return blockReturnType(toT, b)
}

func blockToExpr(block Block) Expr {
	return New_Expr_EReturnableExpr(New_ReturnableExpr_RBlock(block))
}

func exprToBlock(bexpr Expr) Block {
	switch _v2 := (bexpr).(type) {
	case Expr_EReturnableExpr:
		re := _v2.Value
		switch _v3 := (re).(type) {
		case ReturnableExpr_RBlock:
			b := _v3.Value
			return b
		default:
			PanicNow("Not block, some ReturnableExpr.")
			return frt.Empty[Block]()
		}
	default:
		PanicNow("Not block.")
		return frt.Empty[Block]()
	}
}

func meToType(toT func(Expr) FType, me MatchExpr) FType {
	bToT := func(b Block) FType {
		return frt.Pipe(blockToExpr(b), toT)
	}
	switch _v4 := (me.Rules).(type) {
	case MatchRules_Unions:
		us := _v4.Value
		return frt.Pipe(frt.Pipe(slice.Head(us), func(_v1 UnionMatchRule) Block {
			return _v1.Body
		}), bToT)
	case MatchRules_UnionsWD:
		uds := _v4.Value
		return frt.Pipe(frt.Pipe(slice.Head(uds.Unions), func(_v2 UnionMatchRule) Block {
			return _v2.Body
		}), bToT)
	case MatchRules_DefaultOnly:
		db := _v4.Value
		return bToT(db)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func returnableToType(toT func(Expr) FType, rexpr ReturnableExpr) FType {
	switch _v5 := (rexpr).(type) {
	case ReturnableExpr_RBlock:
		b := _v5.Value
		return blockToType(toT, b)
	case ReturnableExpr_RMatchExpr:
		me := _v5.Value
		return meToType(toT, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func fcToFuncType(fc FunCall) FuncType {
	tfv := fc.TargetFunc
	ft := varRefVarType(tfv)
	switch _v6 := (ft).(type) {
	case FType_FFunc:
		ft := _v6.Value
		return ft
	default:
		return FuncType{}
	}
}

func fcToType(fc FunCall) FType {
	firstTp := varRefVarType(fc.TargetFunc)
	switch (firstTp).(type) {
	case FType_FTypeVar:
		return firstTp
	default:
		ft := fcToFuncType(fc)
		tlen := frt.Pipe(fargs(ft), slice.Length)
		alen := slice.Length(fc.Args)
		return frt.IfElse(frt.OpEqual(alen, tlen), (func() FType {
			return freturn(ft)
		}), (func() FType {
			if alen > tlen {
				panic("too many arugments")
			}
			newts := slice.Skip(alen, ft.Targets)
			return New_FType_FFunc(FuncType{Targets: newts})
		}))
	}
}

func lambdaToType(bToT func(Block) FType, le LambdaExpr) FType {
	return frt.Pipe(frt.Pipe(slice.Map(func(_v1 Var) FType {
		return _v1.Ftype
	}, le.Params), (func(_r0 []FType) []FType { return slice.PushLast(bToT(le.Body), _r0) })), newFFunc)
}

func ExprToType(expr Expr) FType {
	switch _v7 := (expr).(type) {
	case Expr_EGoEvalExpr:
		ge := _v7.Value
		return ge.TypeArg
	case Expr_EStringLiteral:
		return New_FType_FString
	case Expr_ESInterP:
		return New_FType_FString
	case Expr_EIntImm:
		return New_FType_FInt
	case Expr_EUnit:
		return New_FType_FUnit
	case Expr_EBoolLiteral:
		return New_FType_FBool
	case Expr_EFieldAccess:
		fa := _v7.Value
		return faToType(ExprToType, fa)
	case Expr_ELambda:
		le := _v7.Value
		return lambdaToType((func(_r0 Block) FType { return blockToType(ExprToType, _r0) }), le)
	case Expr_EVarRef:
		vr := _v7.Value
		v := varRefVar(vr)
		return v.Ftype
	case Expr_ESlice:
		s := _v7.Value
		etp := frt.Pipe(slice.Head(s), ExprToType)
		st := SliceType{ElemType: etp}
		return New_FType_FSlice(st)
	case Expr_ETupleExpr:
		es := _v7.Value
		ts := slice.Map(ExprToType, es)
		return frt.Pipe(TupleType{ElemTypes: ts}, New_FType_FTuple)
	case Expr_ERecordGen:
		rg := _v7.Value
		return New_FType_FRecord(rg.RecordType)
	case Expr_ELazyBlock:
		lb := _v7.Value
		return lblockToType(ExprToType, lb)
	case Expr_EReturnableExpr:
		re := _v7.Value
		return returnableToType(ExprToType, re)
	case Expr_EFunCall:
		fc := _v7.Value
		return fcToType(fc)
	case Expr_EBinOpCall:
		bc := _v7.Value
		return bc.Rtype
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
