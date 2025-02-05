package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/buf"

import "github.com/karino2/folang/pkg/strings"

func rgFVToGo(toGo func(Expr) string, fvPair frt.Tuple2[string, Expr]) string {
	fn := frt.Fst(fvPair)
	fv := frt.Snd(fvPair)
	fvGo := toGo(fv)
	return frt.OpPlus(frt.OpPlus(fn, ": "), fvGo)
}

func rgToGo(toGo func(Expr) string, rg RecordGen) string {
	rtype := rg.recordType
	b := buf.New()
	frt.PipeUnit(frStructName(rtype), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "{")
	fvGo := frt.Pipe(frt.Pipe(slice.Zip(rg.fieldNames, rg.fieldValues), (func(_r0 []frt.Tuple2[string, Expr]) []string {
		return slice.Map((func(_r0 frt.Tuple2[string, Expr]) string { return rgFVToGo(toGo, _r0) }), _r0)
	})), (func(_r0 []string) string { return strings.Concat(", ", _r0) }))
	buf.Write(b, fvGo)
	buf.Write(b, "}")
	return buf.String(b)
}

func buildReturn(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, stmts []Stmt, lastExpr Expr) string {
	stmtGos := frt.Pipe(slice.Map(sToGo, stmts), (func(_r0 []string) string { return strings.Concat("\n", _r0) }))
	lastGo := (func() string {
		switch _v80 := (lastExpr).(type) {
		case Expr_ReturnableExpr:
			re := _v80.Value
			return reToGoRet(re)
		default:
			mayReturn := frt.IfElse(frt.OpEqual(ExprToType(lastExpr), New_FType_FUnit), (func() string {
				return ""
			}), (func() string {
				return "return "
			}))
			lg := eToGo(lastExpr)
			return frt.OpPlus(mayReturn, lg)
		}
	})()
	return frt.OpPlus(frt.OpPlus(stmtGos, "\n"), lastGo)
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
	return frt.OpPlus(wf, "()")
}

func lbToGo(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, lb LazyBlock) string {
	returnBody := buildReturn(sToGo, eToGo, reToGoRet, lb.stmts, lb.finalExpr)
	rtype := ExprToType(lb.finalExpr)
	return wrapFunc(FTypeToGo, rtype, returnBody)
}

func blockToGoReturn(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, block Block) string {
	return buildReturn(sToGo, eToGo, reToGoRet, block.stmts, block.finalExpr)
}

func blockToGo(sToGo func(Stmt) string, eToGo func(Expr) string, reToGoRet func(ReturnableExpr) string, block Block) string {
	goRet := blockToGoReturn(sToGo, eToGo, reToGoRet, block)
	rtype := ExprToType(block.finalExpr)
	return wrapFunCall(FTypeToGo, rtype, goRet)
}

func mpToCaseHeader(uname string, mp MatchPattern, tmpVarName string) string {
	b := buf.New()
	buf.Write(b, "case ")
	frt.PipeUnit(unionCaseStructName(uname, mp.caseId), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ":\n")
	frt.IfOnly(frt.OpAnd(frt.OpNotEqual(mp.varName, "_"), frt.OpNotEqual(mp.varName, "")), (func() {
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

func mrHasCaseVar(mr MatchRule) bool {
	pat := mr.pattern
	return frt.OpAnd(frt.OpAnd(frt.OpNotEqual(pat.caseId, "_"), frt.OpNotEqual(pat.varName, "")), frt.OpNotEqual(pat.varName, "_"))
}

func meHasCaseVar(me MatchExpr) bool {
	return slice.Forall(mrHasCaseVar, me.rules)
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
	return frt.Pipe(New_ReturnableExpr_MatchExpr(me), New_Expr_ReturnableExpr)
}

func meToGo(toGo func(Expr) string, btogRet func(Block) string, me MatchExpr) string {
	goret := meToGoReturn(toGo, btogRet, me)
	rtype := ExprToType(meToExpr(me))
	return wrapFunCall(FTypeToGo, rtype, goret)
}

func reToGoReturn(sToGo func(Stmt) string, eToGo func(Expr) string, rexpr ReturnableExpr) string {
	rtgr := (func(_r0 ReturnableExpr) string { return reToGoReturn(sToGo, eToGo, _r0) })
	btogoRet := (func(_r0 Block) string { return blockToGoReturn(sToGo, eToGo, rtgr, _r0) })
	switch _v81 := (rexpr).(type) {
	case ReturnableExpr_Block:
		b := _v81.Value
		return blockToGoReturn(sToGo, eToGo, rtgr, b)
	case ReturnableExpr_MatchExpr:
		me := _v81.Value
		return meToGoReturn(eToGo, btogoRet, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func reToGo(sToGo func(Stmt) string, eToGo func(Expr) string, rexpr ReturnableExpr) string {
	rtgr := (func(_r0 ReturnableExpr) string { return reToGoReturn(sToGo, eToGo, _r0) })
	btogRet := (func(_r0 Block) string { return blockToGoReturn(sToGo, eToGo, rtgr, _r0) })
	switch _v82 := (rexpr).(type) {
	case ReturnableExpr_Block:
		b := _v82.Value
		return blockToGo(sToGo, eToGo, rtgr, b)
	case ReturnableExpr_MatchExpr:
		me := _v82.Value
		return meToGo(eToGo, btogRet, me)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func ExprToGo(sToGo func(Stmt) string, expr Expr) string {
	eToGo := (func(_r0 Expr) string { return ExprToGo(sToGo, _r0) })
	switch _v83 := (expr).(type) {
	case Expr_BoolLiteral:
		b := _v83.Value
		return frt.Sprintf1("%t", b)
	case Expr_GoEval:
		ge := _v83.Value
		return ge.goStmt
	case Expr_StringLiteral:
		s := _v83.Value
		return frt.Sprintf1("\\\"%s\\\"", s)
	case Expr_IntImm:
		i := _v83.Value
		return frt.Sprintf1("%d", i)
	case Expr_Unit:
		return ""
	case Expr_FieldAccess:
		fa := _v83.Value
		return frt.OpPlus(frt.OpPlus(fa.targetName, "."), fa.fieldName)
	case Expr_Var:
		v := _v83.Value
		return v.name
	case Expr_RecordGen:
		rg := _v83.Value
		return rgToGo(eToGo, rg)
	case Expr_ReturnableExpr:
		re := _v83.Value
		return reToGo(sToGo, eToGo, re)
	case Expr_FunCall:
		return "NYI"
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
