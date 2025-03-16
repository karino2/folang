package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/buf"

import "github.com/karino2/folang/pkg/strings"

func rgFVToGo(toGo func(Expr) string, fvPair NEPair) string {
	fn := fvPair.Name
	fv := fvPair.Expr
	fvGo := toGo(fv)
	return frt.SInterP("%s: %s", fn, fvGo)
}

func rgToGo(toGo func(Expr) string, rg RecordGen) string {
	rtype := rg.RecordType
	b := buf.New()
	frt.PipeUnit(frStructName(FTypeToGo, rtype), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "{")
	fvGo := frt.Pipe(frt.Pipe(rg.FieldsNV, (func(_r0 []NEPair) []string {
		return slice.Map((func(_r0 NEPair) string { return rgFVToGo(toGo, _r0) }), _r0)
	})), (func(_r0 []string) string { return strings.Concat(", ", _r0) }))
	buf.Write(b, fvGo)
	buf.Write(b, "}")
	return buf.String(b)
}

func buildReturn(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, stmts []Stmt, lastExpr Expr) string {
	stmtGos := frt.Pipe(slice.Map(sToGo, stmts), (func(_r0 []string) string { return strings.Concat("\n", _r0) }))
	lastGo := (func() string {
		switch _v1 := (lastExpr).(type) {
		case Expr_EReturnableExpr:
			re := _v1.Value
			return reToGoRet(re)
		default:
			mayReturn := frt.IfElse(frt.OpEqual(ExprToType(lastExpr), New_FType_FUnit), (func() string {
				return ""
			}), (func() string {
				return "return "
			}))
			lg := eToGo(lastExpr)
			return (mayReturn + lg)
		}
	})()
	return frt.IfElse(frt.OpEqual(stmtGos, ""), (func() string {
		return lastGo
	}), (func() string {
		return ((stmtGos + "\n") + lastGo)
	}))
}

func wrapFunc(toGo func(FType) string, rtype FType, goReturnBody string) string {
	b := buf.New()
	buf.Write(b, "(func () ")
	frt.PipeUnit(toGo(rtype), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, " {\n")
	buf.Write(b, goReturnBody)
	buf.Write(b, "})")
	return buf.String(b)
}

func wrapFunCall(toGo func(FType) string, rtype FType, goReturnBody string) string {
	wf := wrapFunc(toGo, rtype, goReturnBody)
	return (wf + "()")
}

func blockToGoReturn(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, block Block) string {
	return buildReturn(sToGo, eToGo, reToGoRet, block.Stmts, block.FinalExpr)
}

func blockToGo(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, block Block) string {
	goRet := blockToGoReturn(sToGo, eToGo, reToGoRet, block)
	rtype := ExprToType(block.FinalExpr)
	return wrapFunCall(FTypeToGo, rtype, goRet)
}

func lbToGo(bToRet func(Block) string, lb LazyBlock) string {
	returnBody := bToRet(lb.Block)
	rtype := ExprToType(lb.Block.FinalExpr)
	return wrapFunc(FTypeToGo, rtype, returnBody)
}

func umpToCaseHeader(uname string, ump UnionMatchPattern, tmpVarName string) string {
	b := buf.New()
	buf.Write(b, "case ")
	frt.PipeUnit(unionCSName(uname, ump.CaseId), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ":\n")
	frt.IfOnly((frt.OpNotEqual(ump.VarName, "_") && frt.OpNotEqual(ump.VarName, "")), (func() {
		buf.Write(b, ump.VarName)
		buf.Write(b, " := ")
		buf.Write(b, tmpVarName)
		buf.Write(b, ".Value")
		buf.Write(b, "\n")
	}))
	return buf.String(b)
}

func umrToCase(btogRet func(Block) string, uname string, tmpVarName string, umr UnionMatchRule) string {
	b := buf.New()
	mp := umr.UnionPattern
	frt.PipeUnit(umpToCaseHeader(uname, mp, tmpVarName), (func(_r0 string) { buf.Write(b, _r0) }))
	frt.PipeUnit(btogRet(umr.Body), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n")
	return buf.String(b)
}

func drToCase(btogRet func(Block) string, db Block) string {
	b := buf.New()
	buf.Write(b, "default:\n")
	frt.PipeUnit(btogRet(db), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n")
	return buf.String(b)
}

func umrHasNoCaseVar(umr UnionMatchRule) bool {
	pat := umr.UnionPattern
	return (frt.OpEqual(pat.VarName, "") || frt.OpEqual(pat.VarName, "_"))
}

func umrHasCaseVar(rules UnionMatchRules) bool {
	allNoCaseF := (func(_r0 []UnionMatchRule) bool { return slice.Forall(umrHasNoCaseVar, _r0) })
	switch _v2 := (rules).(type) {
	case UnionMatchRules_UCaseOnly:
		us := _v2.Value
		return frt.OpNot(allNoCaseF(us))
	case UnionMatchRules_UCaseWD:
		uds := _v2.Value
		return frt.OpNot(allNoCaseF(uds.Unions))
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func umrToGoReturn(toGo func(Expr) string, btogRet func(Block) string, target Expr, rules UnionMatchRules) string {
	ttype := ExprToType(target)
	uttype := CastNow[FType_FUnion](ttype).Value
	uname := utName(uttype)
	hasCaseVar := umrHasCaseVar(rules)
	tmpVarName := frt.IfElse(hasCaseVar, (func() string {
		return uniqueTmpVarName()
	}), (func() string {
		return ""
	}))
	b := buf.New()
	umrstocases := (func(_r0 []UnionMatchRule) []string {
		return slice.Map((func(_r0 UnionMatchRule) string { return umrToCase(btogRet, uname, tmpVarName, _r0) }), _r0)
	})
	writeUmrs := func(umrs []UnionMatchRule) {
		frt.PipeUnit(frt.Pipe(umrstocases(umrs), (func(_r0 []string) string { return strings.Concat("", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	}
	buf.Write(b, "switch ")
	frt.IfOnly(hasCaseVar, (func() {
		buf.Write(b, tmpVarName)
		buf.Write(b, " := ")
	}))
	buf.Write(b, "(")
	frt.PipeUnit(toGo(target), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ").(type){\n")
	(func() {
		switch _v3 := (rules).(type) {
		case UnionMatchRules_UCaseOnly:
			us := _v3.Value
			writeUmrs(us)
			buf.Write(b, "default:\npanic(\"Union pattern fail. Never reached here.\")\n")
		case UnionMatchRules_UCaseWD:
			uds := _v3.Value
			writeUmrs(uds.Unions)
			frt.PipeUnit(drToCase(btogRet, uds.Default), (func(_r0 string) { buf.Write(b, _r0) }))
		default:
			panic("Union pattern fail. Never reached here.")
		}
	})()
	buf.Write(b, "}")
	return buf.String(b)
}

func smrToCase(btogRet func(Block) string, smr StringMatchRule) string {
	b := buf.New()
	lp := smr.LiteralPattern
	buf.Write(b, frt.SInterP("case \"%s\":", lp))
	buf.Write(b, "\n")
	frt.PipeUnit(btogRet(smr.Body), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n")
	return buf.String(b)
}

func svrToCase(btogRet func(Block) string, svr StringVarMatchRule) string {
	b := buf.New()
	buf.Write(b, "default:\n")
	frt.PipeUnit(btogRet(svr.Body), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n")
	return buf.String(b)
}

func smrToGoReturn(toGo func(Expr) string, btogRet func(Block) string, target Expr, rules StringMatchRules) string {
	b := buf.New()
	smrstocases := (func(_r0 []StringMatchRule) []string {
		return slice.Map((func(_r0 StringMatchRule) string { return smrToCase(btogRet, _r0) }), _r0)
	})
	writeSmrs := func(smrs []StringMatchRule) {
		frt.PipeUnit(frt.Pipe(smrstocases(smrs), (func(_r0 []string) string { return strings.Concat("", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	}
	(func() {
		switch _v4 := (rules).(type) {
		case StringMatchRules_SCaseWV:
			swv := _v4.Value
			vname := swv.VarRule.VarName
			buf.Write(b, frt.SInterP("switch %s :=", vname))
			buf.Write(b, "(")
			frt.PipeUnit(toGo(target), (func(_r0 string) { buf.Write(b, _r0) }))
			buf.Write(b, frt.SInterP("); %s", vname))
			buf.Write(b, "{\n")
			writeSmrs(swv.Literals)
			svrToCase(btogRet, swv.VarRule)
			buf.Write(b, "}")
		case StringMatchRules_SCaseWD:
			sds := _v4.Value
			buf.Write(b, "switch (")
			frt.PipeUnit(toGo(target), (func(_r0 string) { buf.Write(b, _r0) }))
			buf.Write(b, "){\n")
			writeSmrs(sds.Literals)
			frt.PipeUnit(drToCase(btogRet, sds.Default), (func(_r0 string) { buf.Write(b, _r0) }))
			buf.Write(b, "}")
		default:
			panic("Union pattern fail. Never reached here.")
		}
	})()
	return buf.String(b)
}

func meToGoReturn(toGo func(Expr) string, btogRet func(Block) string, me MatchExpr) string {
	switch _v5 := (me.Rules).(type) {
	case MatchRules_RUnions:
		ru := _v5.Value
		return umrToGoReturn(toGo, btogRet, me.Target, ru)
	case MatchRules_RStrings:
		su := _v5.Value
		return smrToGoReturn(toGo, btogRet, me.Target, su)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func meToExpr(me MatchExpr) Expr {
	return frt.Pipe(New_ReturnableExpr_RMatchExpr(me), New_Expr_EReturnableExpr)
}

func meToGo(toGo func(Expr) string, btogRet func(Block) string, me MatchExpr) string {
	goret := meToGoReturn(toGo, btogRet, me)
	rtype := ExprToType(meToExpr(me))
	return wrapFunCall(FTypeToGo, rtype, goret)
}

func reToGoReturn(sToGo func(Stmt) string, eToGo func(Expr) string, rexpr ReturnableExpr) string {
	rtgr := (func(_r0 ReturnableExpr) string { return reToGoReturn(sToGo, eToGo, _r0) })
	btogoRet := (func(_r0 Block) string { return blockToGoReturn(sToGo, eToGo, rtgr, _r0) })
	switch _v6 := (rexpr).(type) {
	case ReturnableExpr_RBlock:
		b := _v6.Value
		return blockToGoReturn(sToGo, eToGo, rtgr, b)
	case ReturnableExpr_RMatchExpr:
		me := _v6.Value
		return meToGoReturn(eToGo, btogoRet, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func reToGo(sToGo func(Stmt) string, eToGo func(Expr) string, rexpr ReturnableExpr) string {
	rtgr := (func(_r0 ReturnableExpr) string { return reToGoReturn(sToGo, eToGo, _r0) })
	btogRet := (func(_r0 Block) string { return blockToGoReturn(sToGo, eToGo, rtgr, _r0) })
	switch _v7 := (rexpr).(type) {
	case ReturnableExpr_RBlock:
		b := _v7.Value
		return blockToGo(sToGo, eToGo, rtgr, b)
	case ReturnableExpr_RMatchExpr:
		me := _v7.Value
		return meToGo(eToGo, btogRet, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func ftiToParamName(i int, ft FType) string {
	return frt.SInterP("_r%s", i)
}

func ntpairToParam(tGo func(FType) string, ntp frt.Tuple2[string, FType]) string {
	tpgo := frt.Pipe(frt.Snd(ntp), tGo)
	name := frt.Fst(ntp)
	return frt.SInterP("%s %s", name, tpgo)
}

func varRefToGo(tGo func(FType) string, vr VarRef) string {
	switch _v8 := (vr).(type) {
	case VarRef_VRVar:
		v := _v8.Value
		return v.Name
	case VarRef_VRSVar:
		sv := _v8.Value
		return (sv.Var.Name + tArgsToGo(tGo, sv.SpecList))
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func fcPartialApplyGo(tGo func(FType) string, eGo func(Expr) string, fc FunCall) string {
	funcType := fcToFuncType(fc)
	fargTypes := fargs(funcType)
	argNum := slice.Length(fc.Args)
	restTypes := slice.Skip(argNum, fargTypes)
	restParamNames := slice.Mapi(ftiToParamName, restTypes)
	b := buf.New()
	buf.Write(b, "(func (")
	frt.PipeUnit(frt.Pipe(frt.Pipe(slice.Zip(restParamNames, restTypes), (func(_r0 []frt.Tuple2[string, FType]) []string {
		return slice.Map((func(_r0 frt.Tuple2[string, FType]) string { return ntpairToParam(tGo, _r0) }), _r0)
	})), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ") ")
	fret := freturn(funcType)
	frt.IfElseUnit(frt.OpEqual(fret, New_FType_FUnit), (func() {
		buf.Write(b, "{ ")
	}), (func() {
		frt.PipeUnit(tGo(fret), (func(_r0 string) { buf.Write(b, _r0) }))
		buf.Write(b, "{ return ")
	}))
	frt.PipeUnit(varRefToGo(tGo, fc.TargetFunc), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "(")
	frt.PipeUnit(frt.Pipe(slice.Map(eGo, fc.Args), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ", ")
	frt.PipeUnit(strings.Concat(", ", restParamNames), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ") })")
	return buf.String(b)
}

func fcUnitArgOnly(fc FunCall) bool {
	al := slice.Length(fc.Args)
	return frt.IfElse(frt.OpEqual(al, 1), (func() bool {
		return frt.OpEqual(New_Expr_EUnit, slice.Head(fc.Args))
	}), (func() bool {
		return false
	}))
}

func fcFullApplyGo(tGo func(FType) string, eGo func(Expr) string, fc FunCall) string {
	b := buf.New()
	frt.PipeUnit(varRefToGo(tGo, fc.TargetFunc), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "(")
	frt.IfOnly(frt.OpNot(fcUnitArgOnly(fc)), (func() {
		frt.PipeUnit(frt.Pipe(slice.Map(eGo, fc.Args), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	}))
	buf.Write(b, ")")
	return buf.String(b)
}

func fcToGo(tGo func(FType) string, eGo func(Expr) string, fc FunCall) string {
	funcType := fcToFuncType(fc)
	fargTypes := fargs(funcType)
	al := slice.Length(fc.Args)
	tal := slice.Length(fargTypes)
	frt.IfOnly((al > tal), (func() {
		panic("Too many argument")
	}))
	return frt.IfElse((al < tal), (func() string {
		return fcPartialApplyGo(tGo, eGo, fc)
	}), (func() string {
		return fcFullApplyGo(tGo, eGo, fc)
	}))
}

func sliceToGo(tGo func(FType) string, eGo func(Expr) string, exprs []Expr) string {
	b := buf.New()
	buf.Write(b, "(")
	frt.PipeUnit(frt.Pipe(frt.Pipe(New_Expr_ESlice(exprs), ExprToType), tGo), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "{")
	frt.PipeUnit(frt.Pipe(slice.Map(eGo, exprs), (func(_r0 []string) string { return strings.Concat(",", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "}")
	buf.Write(b, ")")
	return buf.String(b)
}

func tupleToGo(eGo func(Expr) string, exprs []Expr) string {
	b := buf.New()
	len := slice.Length(exprs)
	buf.Write(b, frt.SInterP("frt.NewTuple%s(", len))
	frt.PipeUnit(frt.Pipe(slice.Map(eGo, exprs), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ")")
	return buf.String(b)
}

func binOpToGo(eGo func(Expr) string, binOp BinOpCall) string {
	b := buf.New()
	buf.Write(b, "(")
	frt.PipeUnit(frt.Pipe(frt.Pipe(([]Expr{binOp.Lhs, binOp.Rhs}), (func(_r0 []Expr) []string { return slice.Map(eGo, _r0) })), (func(_r0 []string) string { return strings.Concat(binOp.Op, _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ")")
	return buf.String(b)
}

func faToGo(eGo func(Expr) string, fa FieldAccess) string {
	target := eGo(fa.TargetExpr)
	return frt.SInterP("%s.%s", target, fa.FieldName)
}

func paramsToGo(pm Var) string {
	ts := FTypeToGo(pm.Ftype)
	return frt.SInterP("%s %s", pm.Name, ts)
}

func lambdaToGo(bToGoRet func(Block) string, le LambdaExpr) string {
	b := buf.New()
	buf.Write(b, "func (")
	frt.PipeUnit(frt.Pipe(slice.Map(paramsToGo, le.Params), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ")")
	frt.PipeUnit(frt.Pipe(blockToType(ExprToType, le.Body), FTypeToGo), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "{\n")
	frt.PipeUnit(bToGoRet(le.Body), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n}")
	return buf.String(b)
}

func sinterpToGo(s string) string {
	fm, vs := frt.Destr2(ParseSInterP(s))
	b := buf.New()
	frt.PipeUnit(frt.Sprintf1("frt.SInterP(\"%s\", ", fm), (func(_r0 string) { buf.Write(b, _r0) }))
	frt.PipeUnit(frt.Pipe(vs, (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ")")
	return buf.String(b)
}

func ExprToGo(sToGo func(Stmt) string, expr Expr) string {
	eToGo := (func(_r0 Expr) string { return ExprToGo(sToGo, _r0) })
	reToGoRet := (func(_r0 ReturnableExpr) string { return reToGoReturn(sToGo, eToGo, _r0) })
	bToGoRet := (func(_r0 Block) string { return blockToGoReturn(sToGo, eToGo, reToGoRet, _r0) })
	switch _v9 := (expr).(type) {
	case Expr_EBoolLiteral:
		b := _v9.Value
		return frt.Sprintf1("%t", b)
	case Expr_EGoEvalExpr:
		ge := _v9.Value
		return reinterpretEscape(ge.GoStmt)
	case Expr_EStringLiteral:
		s := _v9.Value
		return frt.Sprintf1("\"%s\"", s)
	case Expr_ESInterP:
		sp := _v9.Value
		return sinterpToGo(sp)
	case Expr_EIntImm:
		i := _v9.Value
		return frt.Sprintf1("%d", i)
	case Expr_EUnit:
		return ""
	case Expr_EFieldAccess:
		fa := _v9.Value
		return faToGo(eToGo, fa)
	case Expr_EVarRef:
		vr := _v9.Value
		return varRefName(vr)
	case Expr_ESlice:
		es := _v9.Value
		return sliceToGo(FTypeToGo, eToGo, es)
	case Expr_ETupleExpr:
		es := _v9.Value
		return tupleToGo(eToGo, es)
	case Expr_ELambda:
		le := _v9.Value
		return lambdaToGo(bToGoRet, le)
	case Expr_EBinOpCall:
		bop := _v9.Value
		return binOpToGo(eToGo, bop)
	case Expr_ERecordGen:
		rg := _v9.Value
		return rgToGo(eToGo, rg)
	case Expr_EReturnableExpr:
		re := _v9.Value
		return reToGo(sToGo, eToGo, re)
	case Expr_EFunCall:
		fc := _v9.Value
		return fcToGo(FTypeToGo, eToGo, fc)
	case Expr_ELazyBlock:
		lb := _v9.Value
		return lbToGo(bToGoRet, lb)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
