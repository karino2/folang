package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/buf"

import "github.com/karino2/folang/pkg/strings"

func rgFVToGo(toGo func(Expr) string, fvPair NEPair) string {
	fn := fvPair.name
	fv := fvPair.expr
	fvGo := toGo(fv)
	return ((fn + ": ") + fvGo)
}

func rgToGo(toGo func(Expr) string, rg RecordGen) string {
	rtype := rg.recordType
	b := buf.New()
	frt.PipeUnit(frStructName(rtype), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "{")
	fvGo := frt.Pipe(frt.Pipe(rg.fieldsNV, (func(_r0 []NEPair) []string {
		return slice.Map((func(_r0 NEPair) string { return rgFVToGo(toGo, _r0) }), _r0)
	})), (func(_r0 []string) string { return strings.Concat(", ", _r0) }))
	buf.Write(b, fvGo)
	buf.Write(b, "}")
	return buf.String(b)
}

func buildReturn(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, stmts []Stmt, lastExpr Expr) string {
	stmtGos := frt.Pipe(slice.Map(sToGo, stmts), (func(_r0 []string) string { return strings.Concat("\n", _r0) }))
	lastGo := (func() string {
		switch _v114 := (lastExpr).(type) {
		case Expr_EReturnableExpr:
			re := _v114.Value
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
	return buildReturn(sToGo, eToGo, reToGoRet, block.stmts, block.finalExpr)
}

func blockToGo(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, block Block) string {
	goRet := blockToGoReturn(sToGo, eToGo, reToGoRet, block)
	rtype := ExprToType(block.finalExpr)
	return wrapFunCall(FTypeToGo, rtype, goRet)
}

func lbToGo(bToRet func(Block) string, lb LazyBlock) string {
	returnBody := bToRet(lb.block)
	rtype := ExprToType(lb.block.finalExpr)
	return wrapFunc(FTypeToGo, rtype, returnBody)
}

func mpToCaseHeader(uname string, mp MatchPattern, tmpVarName string) string {
	b := buf.New()
	buf.Write(b, "case ")
	frt.PipeUnit(unionCSName(uname, mp.caseId), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ":\n")
	frt.IfOnly((frt.OpNotEqual(mp.varName, "_") && frt.OpNotEqual(mp.varName, "")), (func() {
		buf.Write(b, mp.varName)
		buf.Write(b, " := ")
		buf.Write(b, tmpVarName)
		buf.Write(b, ".Value")
		buf.Write(b, "\n")
	}))
	return buf.String(b)
}

func mrToCase(btogRet func(Block) string, uname string, tmpVarName string, mr MatchRule) string {
	b := buf.New()
	mp := mr.pattern
	cheader := frt.IfElse(frt.OpEqual(mp.caseId, "_"), (func() string {
		return "default:\n"
	}), (func() string {
		return mpToCaseHeader(uname, mp, tmpVarName)
	}))
	buf.Write(b, cheader)
	frt.PipeUnit(btogRet(mr.body), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n")
	return buf.String(b)
}

func mrHasNoCaseVar(mr MatchRule) bool {
	pat := mr.pattern
	return ((frt.OpEqual(pat.caseId, "_") || frt.OpEqual(pat.varName, "")) || frt.OpEqual(pat.varName, "_"))
}

func meHasCaseVar(me MatchExpr) bool {
	return frt.OpNot(slice.Forall(mrHasNoCaseVar, me.rules))
}

func mrIsDefault(mr MatchRule) bool {
	pat := mr.pattern
	return frt.OpEqual(pat.caseId, "_")
}

func meToGoReturn(toGo func(Expr) string, btogRet func(Block) string, me MatchExpr) string {
	ttype := ExprToType(me.target)
	uttype := ttype.(FType_FUnion).Value
	hasCaseVar := meHasCaseVar(me)
	hasDefault := slice.Forany(mrIsDefault, me.rules)
	tmpVarName := frt.IfElse(hasCaseVar, (func() string {
		return uniqueTmpVarName()
	}), (func() string {
		return ""
	}))
	mrtocase := (func(_r0 MatchRule) string { return mrToCase(btogRet, uttype.name, tmpVarName, _r0) })
	b := buf.New()
	buf.Write(b, "switch ")
	frt.IfOnly(hasCaseVar, (func() {
		buf.Write(b, tmpVarName)
		buf.Write(b, " := ")
	}))
	buf.Write(b, "(")
	frt.PipeUnit(toGo(me.target), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ").(type){\n")
	frt.PipeUnit(frt.Pipe(slice.Map(mrtocase, me.rules), (func(_r0 []string) string { return strings.Concat("", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
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
	switch _v115 := (rexpr).(type) {
	case ReturnableExpr_RBlock:
		b := _v115.Value
		return blockToGoReturn(sToGo, eToGo, rtgr, b)
	case ReturnableExpr_RMatchExpr:
		me := _v115.Value
		return meToGoReturn(eToGo, btogoRet, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func reToGo(sToGo func(Stmt) string, eToGo func(Expr) string, rexpr ReturnableExpr) string {
	rtgr := (func(_r0 ReturnableExpr) string { return reToGoReturn(sToGo, eToGo, _r0) })
	btogRet := (func(_r0 Block) string { return blockToGoReturn(sToGo, eToGo, rtgr, _r0) })
	switch _v116 := (rexpr).(type) {
	case ReturnableExpr_RBlock:
		b := _v116.Value
		return blockToGo(sToGo, eToGo, rtgr, b)
	case ReturnableExpr_RMatchExpr:
		me := _v116.Value
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

func fcPartialApplyGo(tGo func(FType) string, eGo func(Expr) string, fc FunCall) string {
	funcType := fcToFuncType(fc)
	fargTypes := fargs(funcType)
	argNum := slice.Length(fc.args)
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
	buf.Write(b, fc.targetFunc.name)
	buf.Write(b, "(")
	frt.PipeUnit(frt.Pipe(slice.Map(eGo, fc.args), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ", ")
	frt.PipeUnit(strings.Concat(", ", restParamNames), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ") })")
	return buf.String(b)
}

func fcUnitArgOnly(fc FunCall) bool {
	al := slice.Length(fc.args)
	return frt.IfElse(frt.OpEqual(al, 1), (func() bool {
		return frt.OpEqual(New_Expr_EUnit, slice.Head(fc.args))
	}), (func() bool {
		return false
	}))
}

func fcFullApplyGo(eGo func(Expr) string, fc FunCall) string {
	b := buf.New()
	buf.Write(b, fc.targetFunc.name)
	buf.Write(b, "(")
	frt.IfOnly(frt.OpNot(fcUnitArgOnly(fc)), (func() {
		frt.PipeUnit(frt.Pipe(slice.Map(eGo, fc.args), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	}))
	buf.Write(b, ")")
	return buf.String(b)
}

func fcToGo(tGo func(FType) string, eGo func(Expr) string, fc FunCall) string {
	funcType := fcToFuncType(fc)
	fargTypes := fargs(funcType)
	al := slice.Length(fc.args)
	tal := slice.Length(fargTypes)
	frt.IfOnly((al > tal), (func() {
		panic("Too many argument")
	}))
	return frt.IfElse((al < tal), (func() string {
		return fcPartialApplyGo(tGo, eGo, fc)
	}), (func() string {
		return fcFullApplyGo(eGo, fc)
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
	frt.PipeUnit(frt.Pipe(frt.Pipe(([]Expr{binOp.lhs, binOp.rhs}), (func(_r0 []Expr) []string { return slice.Map(eGo, _r0) })), (func(_r0 []string) string { return strings.Concat(binOp.op, _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ")")
	return buf.String(b)
}

func faToGo(eGo func(Expr) string, fa FieldAccess) string {
	target := eGo(fa.targetExpr)
	return ((target + ".") + fa.fieldName)
}

func ExprToGo(sToGo func(Stmt) string, expr Expr) string {
	eToGo := (func(_r0 Expr) string { return ExprToGo(sToGo, _r0) })
	reToGoRet := (func(_r0 ReturnableExpr) string { return reToGoReturn(sToGo, eToGo, _r0) })
	switch _v117 := (expr).(type) {
	case Expr_EBoolLiteral:
		b := _v117.Value
		return frt.Sprintf1("%t", b)
	case Expr_EGoEvalExpr:
		ge := _v117.Value
		return reinterpretEscape(ge.goStmt)
	case Expr_EStringLiteral:
		s := _v117.Value
		return frt.Sprintf1("\"%s\"", s)
	case Expr_EIntImm:
		i := _v117.Value
		return frt.Sprintf1("%d", i)
	case Expr_EUnit:
		return ""
	case Expr_EFieldAccess:
		fa := _v117.Value
		return faToGo(eToGo, fa)
	case Expr_EVar:
		v := _v117.Value
		return v.name
	case Expr_ESlice:
		es := _v117.Value
		return sliceToGo(FTypeToGo, eToGo, es)
	case Expr_ETupleExpr:
		es := _v117.Value
		return tupleToGo(eToGo, es)
	case Expr_EBinOpCall:
		bop := _v117.Value
		return binOpToGo(eToGo, bop)
	case Expr_ERecordGen:
		rg := _v117.Value
		return rgToGo(eToGo, rg)
	case Expr_EReturnableExpr:
		re := _v117.Value
		return reToGo(sToGo, eToGo, re)
	case Expr_EFunCall:
		fc := _v117.Value
		return fcToGo(FTypeToGo, eToGo, fc)
	case Expr_ELazyBlock:
		lb := _v117.Value
		return lbToGo((func(_r0 Block) string { return blockToGoReturn(sToGo, eToGo, reToGoRet, _r0) }), lb)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
