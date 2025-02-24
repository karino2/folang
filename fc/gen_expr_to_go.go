package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/buf"

import "github.com/karino2/folang/pkg/strings"

func rgFVToGo(toGo func(Expr) string, fvPair NEPair) string {
	fn := fvPair.Name
	fv := fvPair.Expr
	fvGo := toGo(fv)
	return ((fn + ": ") + fvGo)
}

func rgToGo(toGo func(Expr) string, rg RecordGen) string {
	rtype := rg.RecordType
	b := buf.New()
	frt.PipeUnit(frStructName(rtype), (func(_r0 string) { buf.Write(b, _r0) }))
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

func mpToCaseHeader(uname string, mp MatchPattern, tmpVarName string) string {
	b := buf.New()
	buf.Write(b, "case ")
	frt.PipeUnit(unionCSName(uname, mp.CaseId), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ":\n")
	frt.IfOnly((frt.OpNotEqual(mp.VarName, "_") && frt.OpNotEqual(mp.VarName, "")), (func() {
		buf.Write(b, mp.VarName)
		buf.Write(b, " := ")
		buf.Write(b, tmpVarName)
		buf.Write(b, ".Value")
		buf.Write(b, "\n")
	}))
	return buf.String(b)
}

func mrToCase(btogRet func(Block) string, uname string, tmpVarName string, mr MatchRule) string {
	b := buf.New()
	mp := mr.Pattern
	cheader := frt.IfElse(frt.OpEqual(mp.CaseId, "_"), (func() string {
		return "default:\n"
	}), (func() string {
		return mpToCaseHeader(uname, mp, tmpVarName)
	}))
	buf.Write(b, cheader)
	frt.PipeUnit(btogRet(mr.Body), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n")
	return buf.String(b)
}

func mrHasNoCaseVar(mr MatchRule) bool {
	pat := mr.Pattern
	return ((frt.OpEqual(pat.CaseId, "_") || frt.OpEqual(pat.VarName, "")) || frt.OpEqual(pat.VarName, "_"))
}

func meHasCaseVar(me MatchExpr) bool {
	return frt.OpNot(slice.Forall(mrHasNoCaseVar, me.Rules))
}

func mrIsDefault(mr MatchRule) bool {
	pat := mr.Pattern
	return frt.OpEqual(pat.CaseId, "_")
}

func meToGoReturn(toGo func(Expr) string, btogRet func(Block) string, me MatchExpr) string {
	ttype := ExprToType(me.Target)
	uttype := ttype.(FType_FUnion).Value
	uname := utName(uttype)
	hasCaseVar := meHasCaseVar(me)
	hasDefault := slice.Forany(mrIsDefault, me.Rules)
	tmpVarName := frt.IfElse(hasCaseVar, (func() string {
		return uniqueTmpVarName()
	}), (func() string {
		return ""
	}))
	mrtocase := (func(_r0 MatchRule) string { return mrToCase(btogRet, uname, tmpVarName, _r0) })
	b := buf.New()
	buf.Write(b, "switch ")
	frt.IfOnly(hasCaseVar, (func() {
		buf.Write(b, tmpVarName)
		buf.Write(b, " := ")
	}))
	buf.Write(b, "(")
	frt.PipeUnit(toGo(me.Target), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ").(type){\n")
	frt.PipeUnit(frt.Pipe(slice.Map(mrtocase, me.Rules), (func(_r0 []string) string { return strings.Concat("", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	frt.IfOnly(frt.OpNot(hasDefault), (func() {
		buf.Write(b, "default:\npanic(\"Union pattern fail. Never reached here.\")\n")
	}))
	buf.Write(b, "}")
	return buf.String(b)
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
	switch _v2 := (rexpr).(type) {
	case ReturnableExpr_RBlock:
		b := _v2.Value
		return blockToGoReturn(sToGo, eToGo, rtgr, b)
	case ReturnableExpr_RMatchExpr:
		me := _v2.Value
		return meToGoReturn(eToGo, btogoRet, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func reToGo(sToGo func(Stmt) string, eToGo func(Expr) string, rexpr ReturnableExpr) string {
	rtgr := (func(_r0 ReturnableExpr) string { return reToGoReturn(sToGo, eToGo, _r0) })
	btogRet := (func(_r0 Block) string { return blockToGoReturn(sToGo, eToGo, rtgr, _r0) })
	switch _v3 := (rexpr).(type) {
	case ReturnableExpr_RBlock:
		b := _v3.Value
		return blockToGo(sToGo, eToGo, rtgr, b)
	case ReturnableExpr_RMatchExpr:
		me := _v3.Value
		return meToGo(eToGo, btogRet, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func ftiToParamName(i int, ft FType) string {
	return frt.Sprintf1("_r%d", i)
}

func ntpairToParam(tGo func(FType) string, ntp frt.Tuple2[string, FType]) string {
	tpgo := frt.Pipe(frt.Snd(ntp), tGo)
	name := frt.Fst(ntp)
	return ((name + " ") + tpgo)
}

func varRefToGo(tGo func(FType) string, vr VarRef) string {
	switch _v4 := (vr).(type) {
	case VarRef_VRVar:
		v := _v4.Value
		return v.Name
	case VarRef_VRSVar:
		sv := _v4.Value
		tlis := frt.Pipe(frt.Pipe(sv.SpecList, (func(_r0 []FType) []string { return slice.Map(tGo, _r0) })), (func(_r0 []string) string { return strings.Concat(", ", _r0) }))
		return (((sv.Var.Name + "[") + tlis) + "]")
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
	buf.Write(b, "frt.NewTuple2(")
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
	return ((target + ".") + fa.FieldName)
}

func paramsToGo(pm Var) string {
	ts := FTypeToGo(pm.Ftype)
	return ((pm.Name + " ") + ts)
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

func ExprToGo(sToGo func(Stmt) string, expr Expr) string {
	eToGo := (func(_r0 Expr) string { return ExprToGo(sToGo, _r0) })
	reToGoRet := (func(_r0 ReturnableExpr) string { return reToGoReturn(sToGo, eToGo, _r0) })
	bToGoRet := (func(_r0 Block) string { return blockToGoReturn(sToGo, eToGo, reToGoRet, _r0) })
	switch _v5 := (expr).(type) {
	case Expr_EBoolLiteral:
		b := _v5.Value
		return frt.Sprintf1("%t", b)
	case Expr_EGoEvalExpr:
		ge := _v5.Value
		return reinterpretEscape(ge.GoStmt)
	case Expr_EStringLiteral:
		s := _v5.Value
		return frt.Sprintf1("\"%s\"", s)
	case Expr_EIntImm:
		i := _v5.Value
		return frt.Sprintf1("%d", i)
	case Expr_EUnit:
		return ""
	case Expr_EFieldAccess:
		fa := _v5.Value
		return faToGo(eToGo, fa)
	case Expr_EVarRef:
		vr := _v5.Value
		return varRefName(vr)
	case Expr_ESlice:
		es := _v5.Value
		return sliceToGo(FTypeToGo, eToGo, es)
	case Expr_ETupleExpr:
		es := _v5.Value
		return tupleToGo(eToGo, es)
	case Expr_ELambda:
		le := _v5.Value
		return lambdaToGo(bToGoRet, le)
	case Expr_EBinOpCall:
		bop := _v5.Value
		return binOpToGo(eToGo, bop)
	case Expr_ERecordGen:
		rg := _v5.Value
		return rgToGo(eToGo, rg)
	case Expr_EReturnableExpr:
		re := _v5.Value
		return reToGo(sToGo, eToGo, re)
	case Expr_EFunCall:
		fc := _v5.Value
		return fcToGo(FTypeToGo, eToGo, fc)
	case Expr_ELazyBlock:
		lb := _v5.Value
		return lbToGo(bToGoRet, lb)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
