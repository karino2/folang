package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func faToType(fa FieldAccess) FType {
	rt := fa.targetType
	field := frGetField(rt, fa.fieldName)
	return field.ftype
}

func lblockReturnType(toT func(Expr) FType, lb LazyBlock) FType {
	return toT(lb.block.finalExpr)
}

func lblockToType(toT func(Expr) FType, lb LazyBlock) FType {
	rtype := lblockReturnType(toT, lb)
	return New_FType_FFunc(FuncType{targets: ([]FType{New_FType_FUnit, rtype})})
}

func blockReturnType(toT func(Expr) FType, block Block) FType {
	return toT(block.finalExpr)
}

func blockToType(toT func(Expr) FType, b Block) FType {
	return blockReturnType(toT, b)
}

func blockToExpr(block Block) Expr {
	return New_Expr_EReturnableExpr(New_ReturnableExpr_RBlock(block))
}

func emptyBlock() Block {
	return Block{}
}

func exprToBlock(bexpr Expr) Block {
	switch _v35 := (bexpr).(type) {
	case Expr_EReturnableExpr:
		re := _v35.Value
		switch _v36 := (re).(type) {
		case ReturnableExpr_RBlock:
			b := _v36.Value
			return b
		default:
			frt.Panic("Not block, some ReturnableExpr.")
			return emptyBlock()
		}
	default:
		frt.Panic("Not block.")
		return emptyBlock()
	}
}

func meToType(toT func(Expr) FType, me MatchExpr) FType {
	frule := frt.Pipe(me.rules, slice.Head)
	return frt.Pipe(frt.Pipe(frule.body, blockToExpr), toT)
}

func returnableToType(toT func(Expr) FType, rexpr ReturnableExpr) FType {
	switch _v37 := (rexpr).(type) {
	case ReturnableExpr_RBlock:
		b := _v37.Value
		return blockToType(toT, b)
	case ReturnableExpr_RMatchExpr:
		me := _v37.Value
		return meToType(toT, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func fcToFuncType(fc FunCall) FuncType {
	tfv := fc.targetFunc
	ft := tfv.ftype
	switch _v38 := (ft).(type) {
	case FType_FFunc:
		ft := _v38.Value
		return ft
	default:
		return FuncType{}
	}
}

func fcToType(fc FunCall) FType {
	ft := fcToFuncType(fc)
	tlen := frt.Pipe(fargs(ft), slice.Length)
	alen := slice.Length(fc.args)
	return frt.IfElse(frt.OpEqual(alen, tlen), (func() FType {
		return freturn(ft)
	}), (func() FType {
		if alen > tlen {
			panic("too many arugments")
		}
		newts := slice.Skip(alen, ft.targets)
		return New_FType_FFunc(FuncType{targets: newts})
	}))
}

func ExprToType(expr Expr) FType {
	switch _v39 := (expr).(type) {
	case Expr_EGoEvalExpr:
		ge := _v39.Value
		return ge.typeArg
	case Expr_EStringLiteral:
		return New_FType_FString
	case Expr_EIntImm:
		return New_FType_FInt
	case Expr_EUnit:
		return New_FType_FUnit
	case Expr_EBoolLiteral:
		return New_FType_FBool
	case Expr_EFieldAccess:
		fa := _v39.Value
		return faToType(fa)
	case Expr_EVar:
		v := _v39.Value
		return v.ftype
	case Expr_ESlice:
		s := _v39.Value
		etp := frt.Pipe(slice.Head(s), ExprToType)
		st := SliceType{elemType: etp}
		return New_FType_FSlice(st)
	case Expr_ETupleExpr:
		es := _v39.Value
		ts := slice.Map(ExprToType, es)
		return frt.Pipe(TupleType{elemTypes: ts}, New_FType_FTuple)
	case Expr_ERecordGen:
		rg := _v39.Value
		return New_FType_FRecord(rg.recordType)
	case Expr_ELazyBlock:
		lb := _v39.Value
		return lblockToType(ExprToType, lb)
	case Expr_EReturnableExpr:
		re := _v39.Value
		return returnableToType(ExprToType, re)
	case Expr_EFunCall:
		fc := _v39.Value
		return fcToType(fc)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
