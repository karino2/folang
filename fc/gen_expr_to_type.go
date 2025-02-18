package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func faToType(eToT func(Expr) FType, fa FieldAccess) FType {
	tart := eToT(fa.TargetExpr)
	switch _v40 := (tart).(type) {
	case FType_FRecord:
		rt := _v40.Value
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

func emptyBlock() Block {
	return Block{}
}

func exprToBlock(bexpr Expr) Block {
	switch _v41 := (bexpr).(type) {
	case Expr_EReturnableExpr:
		re := _v41.Value
		switch _v42 := (re).(type) {
		case ReturnableExpr_RBlock:
			b := _v42.Value
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
	frule := frt.Pipe(me.Rules, slice.Head)
	return frt.Pipe(frt.Pipe(frule.Body, blockToExpr), toT)
}

func returnableToType(toT func(Expr) FType, rexpr ReturnableExpr) FType {
	switch _v43 := (rexpr).(type) {
	case ReturnableExpr_RBlock:
		b := _v43.Value
		return blockToType(toT, b)
	case ReturnableExpr_RMatchExpr:
		me := _v43.Value
		return meToType(toT, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func fcToFuncType(fc FunCall) FuncType {
	tfv := fc.TargetFunc
	ft := tfv.Ftype
	switch _v44 := (ft).(type) {
	case FType_FFunc:
		ft := _v44.Value
		return ft
	default:
		return FuncType{}
	}
}

func fcToType(fc FunCall) FType {
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

func ExprToType(expr Expr) FType {
	switch _v45 := (expr).(type) {
	case Expr_EGoEvalExpr:
		ge := _v45.Value
		return ge.TypeArg
	case Expr_EStringLiteral:
		return New_FType_FString
	case Expr_EIntImm:
		return New_FType_FInt
	case Expr_EUnit:
		return New_FType_FUnit
	case Expr_EBoolLiteral:
		return New_FType_FBool
	case Expr_EFieldAccess:
		fa := _v45.Value
		return faToType(ExprToType, fa)
	case Expr_EVar:
		v := _v45.Value
		return v.Ftype
	case Expr_ESlice:
		s := _v45.Value
		etp := frt.Pipe(slice.Head(s), ExprToType)
		st := SliceType{ElemType: etp}
		return New_FType_FSlice(st)
	case Expr_ETupleExpr:
		es := _v45.Value
		ts := slice.Map(ExprToType, es)
		return frt.Pipe(TupleType{ElemTypes: ts}, New_FType_FTuple)
	case Expr_ERecordGen:
		rg := _v45.Value
		return New_FType_FRecord(rg.RecordType)
	case Expr_ELazyBlock:
		lb := _v45.Value
		return lblockToType(ExprToType, lb)
	case Expr_EReturnableExpr:
		re := _v45.Value
		return returnableToType(ExprToType, re)
	case Expr_EFunCall:
		fc := _v45.Value
		return fcToType(fc)
	case Expr_EBinOpCall:
		bc := _v45.Value
		return bc.Rtype
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
