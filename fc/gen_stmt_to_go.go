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

func udUnionDef(ud UnionDef) string {
	b := buf.New()
	buf.Write(b, "type ")
	buf.Write(b, ud.name)
	buf.Write(b, " interface {\n")
	buf.Write(b, "  ")
	buf.Write(b, ud.name)
	buf.Write(b, "_Union()\n")
	buf.Write(b, "}\n")
	return buf.String(b)
}

func csToConformMethod(uname string, method string, cas NameTypePair) string {
	csname := unionCSName(uname, cas.name)
	return ((("func (" + csname) + ") ") + method)
}

func udCSConformMethods(ud UnionDef) string {
	method := (ud.name + "_Union(){}\n")
	return frt.Pipe(frt.Pipe(ud.cases, (func(_r0 []NameTypePair) []string {
		return slice.Map((func(_r0 NameTypePair) string { return csToConformMethod(ud.name, method, _r0) }), _r0)
	})), (func(_r0 []string) string { return strings.Concat("", _r0) }))
}

func udCSDef(ud UnionDef, cas NameTypePair) string {
	b := buf.New()
	buf.Write(b, "type ")
	frt.PipeUnit(unionCSName(ud.name, cas.name), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, " struct {\n")
	frt.IfOnly(frt.OpNotEqual(cas.ftype, New_FType_FUnit), (func() {
		buf.Write(b, "  Value ")
		frt.PipeUnit(FTypeToGo(cas.ftype), (func(_r0 string) { buf.Write(b, _r0) }))
		buf.Write(b, "\n")
	}))
	buf.Write(b, "}\n")
	return buf.String(b)
}

func csConstructorName(unionName string, cas NameTypePair) string {
	return ("New_" + unionCSName(unionName, cas.name))
}

func csConstructFunc(uname string, cas NameTypePair) string {
	b := buf.New()
	buf.Write(b, "func ")
	frt.PipeUnit(csConstructorName(uname, cas), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "(v ")
	frt.PipeUnit(FTypeToGo(cas.ftype), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ") ")
	buf.Write(b, uname)
	buf.Write(b, " { return ")
	buf.Write(b, unionCSName(uname, cas.name))
	buf.Write(b, "{v} }\n")
	return buf.String(b)
}

func csConstructVar(uname string, cas NameTypePair) string {
	b := buf.New()
	buf.Write(b, "var ")
	frt.PipeUnit(csConstructorName(uname, cas), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, " ")
	buf.Write(b, uname)
	buf.Write(b, " = ")
	buf.Write(b, unionCSName(uname, cas.name))
	buf.Write(b, "{}\n")
	return buf.String(b)
}

func csConstruct(uname string, cas NameTypePair) string {
	return frt.IfElse(frt.OpEqual(cas.ftype, New_FType_FUnit), (func() string {
		return csConstructVar(uname, cas)
	}), (func() string {
		return csConstructFunc(uname, cas)
	}))
}

func caseToGo(ud UnionDef, cas NameTypePair) string {
	sdf := udCSDef(ud, cas)
	csdf := csConstruct(ud.name, cas)
	return (((sdf + "\n") + csdf) + "\n")
}

func udfToGo(ud UnionDef) string {
	b := buf.New()
	frt.PipeUnit(udUnionDef(ud), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n")
	frt.PipeUnit(udCSConformMethods(ud), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n")
	frt.PipeUnit(frt.Pipe(frt.Pipe(ud.cases, (func(_r0 []NameTypePair) []string {
		return slice.Map((func(_r0 NameTypePair) string { return caseToGo(ud, _r0) }), _r0)
	})), (func(_r0 []string) string { return strings.Concat("", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	return buf.String(b)
}

func dsToGo(ds DefStmt) string {
	switch _v137 := (ds).(type) {
	case DefStmt_RecordDef:
		rd := _v137.Value
		return rdfToGo(rd)
	case DefStmt_UnionDef:
		ud := _v137.Value
		return udfToGo(ud)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func mdToGo(md MultipleDefs) string {
	return frt.Pipe(frt.Pipe(md.defs, (func(_r0 []DefStmt) []string { return slice.Map(dsToGo, _r0) })), (func(_r0 []string) string { return strings.Concat("\n", _r0) }))
}

func StmtToGo(stmt Stmt) string {
	eToGo := (func(_r0 Expr) string { return ExprToGo(StmtToGo, _r0) })
	reToGoRet := (func(_r0 ReturnableExpr) string { return reToGoReturn(StmtToGo, eToGo, _r0) })
	bToGoRet := (func(_r0 Block) string { return blockToGoReturn(StmtToGo, eToGo, reToGoRet, _r0) })
	switch _v138 := (stmt).(type) {
	case Stmt_Import:
		im := _v138.Value
		return imToGo(im)
	case Stmt_Package:
		pn := _v138.Value
		return pmToGo(pn)
	case Stmt_PackageInfo:
		return ""
	case Stmt_LetFuncDef:
		lfd := _v138.Value
		return lfdToGo(bToGoRet, lfd)
	case Stmt_LetVarDef:
		lvd := _v138.Value
		return lvdToGo(eToGo, lvd)
	case Stmt_ExprStmt:
		expr := _v138.Value
		return eToGo(expr)
	case Stmt_DefStmt:
		ds := _v138.Value
		return dsToGo(ds)
	case Stmt_MultipleDefs:
		md := _v138.Value
		return mdToGo(md)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
