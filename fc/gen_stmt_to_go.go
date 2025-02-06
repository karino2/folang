package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/buf"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/strings"

func imToGo(pn string) string {
	return frt.Sprintf1("import \"%s\"", pn)
}

func pmToGo(pn string) string {
	return frt.Sprintf1("package %s", pn)
}

func paramsToGo(pm Var) string {
	ts := FTypeToGo(pm.ftype)
	return ((pm.name + " ") + ts)
}

func lfdParamsToGo(lfd LetFuncDef) string {
	return frt.Pipe(frt.Pipe(lfd.params, (func(_r0 []Var) []string { return slice.Map(paramsToGo, _r0) })), (func(_r0 []string) string { return strings.Concat(", ", _r0) }))
}

func lfdToGo(bToGoRet func(Block) string, lfd LetFuncDef) string {
	b := buf.New()
	buf.Write(b, "func ")
	buf.Write(b, lfd.name)
	buf.Write(b, "(")
	frt.PipeUnit(lfdParamsToGo(lfd), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ")")
	frt.PipeUnit(frt.Pipe(blockToType(ExprToType, lfd.body), FTypeToGo), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "{\n")
	frt.PipeUnit(bToGoRet(lfd.body), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n}")
	return buf.String(b)
}

func lvdToGo(eToGo func(Expr) string, lvd LetVarDef) string {
	rhs := eToGo(lvd.rhs)
	return ((lvd.name + " := ") + rhs)
}

func rdffieldToGo(field NameTypePair) string {
	return ((("  " + field.name) + " ") + FTypeToGo(field.ftype))
}

func rdfToGo(rdf RecordDef) string {
	b := buf.New()
	buf.Write(b, "type ")
	buf.Write(b, rdf.name)
	buf.Write(b, " struct {\n")
	frt.PipeUnit(frt.Pipe(frt.Pipe(rdf.fields, (func(_r0 []NameTypePair) []string { return slice.Map(rdffieldToGo, _r0) })), (func(_r0 []string) string { return strings.Concat("\n", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n}")
	return buf.String(b)
}

func StmtToGo(stmt Stmt) string {
	eToGo := (func(_r0 Expr) string { return ExprToGo(StmtToGo, _r0) })
	reToGoRet := (func(_r0 ReturnableExpr) string { return reToGoReturn(StmtToGo, eToGo, _r0) })
	bToGoRet := (func(_r0 Block) string { return blockToGoReturn(StmtToGo, eToGo, reToGoRet, _r0) })
	switch _v107 := (stmt).(type) {
	case Stmt_Import:
		im := _v107.Value
		return imToGo(im)
	case Stmt_Package:
		pn := _v107.Value
		return pmToGo(pn)
	case Stmt_PackageInfo:
		return ""
	case Stmt_LetFuncDef:
		lfd := _v107.Value
		return lfdToGo(bToGoRet, lfd)
	case Stmt_LetVarDef:
		lvd := _v107.Value
		return lvdToGo(eToGo, lvd)
	case Stmt_ExprStmt:
		expr := _v107.Value
		return eToGo(expr)
	case Stmt_DefStmt:
		return "NYI"
	case Stmt_MultipleDefs:
		return "NYI"
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
