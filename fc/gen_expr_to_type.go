package main

import (
	"github.com/karino2/folang/pkg/frt"
	"github.com/karino2/folang/pkg/slice"
)

func faToType(fa FieldAccess) FType {
	rt := fa.targetType
	field := frGetField(rt, fa.fieldName)
	return field.ftype
}

func lblockReturnType(toT func(Expr) FType, lb LazyBlock) FType {
	return toT(lb.finalExpr)
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
	return New_Expr_ReturnableExpr(New_ReturnableExpr_Block(block))
}

func meToType(toT func(Expr) FType, me MatchExpr) FType {
	frule := frt.Pipe(me.rules, slice.Head)
	return frt.Pipe(frt.Pipe(frule.body, blockToExpr), toT)
}

func returnableToType(toT func(Expr) FType, rexpr ReturnableExpr) FType {
	switch _v36 := (rexpr).(type) {
	case ReturnableExpr_Block:
		b := _v36.Value
		return blockToType(toT, b)
	case ReturnableExpr_MatchExpr:
		me := _v36.Value
		return meToType(toT, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func fcToFuncType(fc FunCall) FuncType {
	tfv := fc.targetFunc
	ft := tfv.ftype
	switch _v37 := (ft).(type) {
	case FType_FFunc:
		ft := _v37.Value
		return ft
	default:
		return FuncType{}
	}
}

func fcArgTypes(fc FunCall) []FType {
	return frt.Pipe(fcToFuncType(fc), fargs)
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
	switch _v38 := (expr).(type) {
	case Expr_GoEval:
		ge := _v38.Value
		return ge.typeArg
	case Expr_StringLiteral:
		return New_FType_FString
	case Expr_IntImm:
		return New_FType_FInt
	case Expr_Unit:
		return New_FType_FUnit
	case Expr_BoolLiteral:
		return New_FType_FBool
	case Expr_FieldAccess:
		fa := _v38.Value
		return faToType(fa)
	case Expr_Var:
		v := _v38.Value
		return v.ftype
	case Expr_RecordGen:
		rg := _v38.Value
		return New_FType_FRecord(rg.recordType)
	case Expr_LazyBlock:
		lb := _v38.Value
		return lblockToType(ExprToType, lb)
	case Expr_ReturnableExpr:
		re := _v38.Value
		return returnableToType(ExprToType, re)
	case Expr_FunCall:
		fc := _v38.Value
		return fcToType(fc)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
