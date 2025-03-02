package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/buf"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/strings"

func imToGo(pn string) string {
	return frt.SInterP("import \"%s\"", pn)
}

func pmToGo(pn string) string {
	return frt.SInterP("package %s", pn)
}

func lfdParamsToGo(lfd LetFuncDef) string {
	return frt.Pipe(frt.Pipe(lfd.Params, (func(_r0 []Var) []string { return slice.Map(paramsToGo, _r0) })), (func(_r0 []string) string { return strings.Concat(", ", _r0) }))
}

func lfdToGo(bToGoRet func(Block) string, lfd LetFuncDef) string {
	b := buf.New()
	buf.Write(b, "func ")
	buf.Write(b, lfd.Fvar.Name)
	buf.Write(b, "(")
	frt.PipeUnit(lfdParamsToGo(lfd), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ") ")
	frt.PipeUnit(frt.Pipe(blockToType(ExprToType, lfd.Body), FTypeToGo), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "{\n")
	frt.PipeUnit(bToGoRet(lfd.Body), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n}")
	return buf.String(b)
}

func pany(s string) string {
	return frt.Sprintf1("%s any", s)
}

func writeTParamsIfAny(b buf.Buffer, tparams []string) {
	frt.IfOnly(frt.OpNot(slice.IsEmpty(tparams)), (func() {
		buf.Write(b, "[")
		frt.PipeUnit(frt.Pipe(frt.Pipe(tparams, (func(_r0 []string) []string { return slice.Map(pany, _r0) })), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
		buf.Write(b, "]")
	}))
}

func rfdToGo(bToGoRet func(Block) string, rfd RootFuncDef) string {
	lfd := rfd.Lfd
	b := buf.New()
	buf.Write(b, "func ")
	buf.Write(b, lfd.Fvar.Name)
	writeTParamsIfAny(b, rfd.Tparams)
	buf.Write(b, "(")
	frt.PipeUnit(lfdParamsToGo(lfd), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ") ")
	frt.PipeUnit(frt.Pipe(blockToType(ExprToType, lfd.Body), FTypeToGo), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "{\n")
	frt.PipeUnit(bToGoRet(lfd.Body), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n}")
	return buf.String(b)
}

func lvdToGo(eToGo func(Expr) string, lvd LetVarDef) string {
	rhs := eToGo(lvd.Rhs)
	return ((lvd.Lvar.Name + " := ") + rhs)
}

func ldvdToGo(eToGo func(Expr) string, ldvd LetDestVarDef) string {
	b := buf.New()
	frt.PipeUnit(frt.Pipe(slice.Map(func(_v1 Var) string {
		return _v1.Name
	}, ldvd.Lvars), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, " := frt.Destr(")
	frt.PipeUnit(eToGo(ldvd.Rhs), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ")")
	return buf.String(b)
}

func rdffieldToGo(field NameTypePair) string {
	return ((("  " + field.Name) + " ") + FTypeToGo(field.Ftype))
}

func rdfToGo(rdf RecordDef) string {
	b := buf.New()
	buf.Write(b, frt.SInterP("type %s", rdf.Name))
	writeTParamsIfAny(b, rdf.Tparams)
	buf.Write(b, " struct {\n")
	frt.PipeUnit(frt.Pipe(frt.Pipe(rdf.Fields, (func(_r0 []NameTypePair) []string { return slice.Map(rdffieldToGo, _r0) })), (func(_r0 []string) string { return strings.Concat("\n", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n}")
	return buf.String(b)
}

func udUnionDef(ud UnionDef) string {
	b := buf.New()
	buf.Write(b, frt.SInterP("type %s interface ", ud.Name))
	buf.Write(b, "{\n")
	buf.Write(b, frt.SInterP("  %s_Union()\n", ud.Name))
	buf.Write(b, "}\n")
	return buf.String(b)
}

func csToConformMethod(uname string, method string, cas NameTypePair) string {
	csname := unionCSName(uname, cas.Name)
	return frt.SInterP("func (%s) %s", csname, method)
}

func udCSConformMethods(ud UnionDef) string {
	method := (frt.SInterP("%s_Union()", ud.Name) + "{}\n")
	return frt.Pipe(frt.Pipe(udCases(ud), (func(_r0 []NameTypePair) []string {
		return slice.Map((func(_r0 NameTypePair) string { return csToConformMethod(ud.Name, method, _r0) }), _r0)
	})), (func(_r0 []string) string { return strings.Concat("", _r0) }))
}

func udCSDef(ud UnionDef, cas NameTypePair) string {
	b := buf.New()
	buf.Write(b, "type ")
	frt.PipeUnit(unionCSName(ud.Name, cas.Name), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, " struct {\n")
	frt.IfOnly(frt.OpNotEqual(cas.Ftype, New_FType_FUnit), (func() {
		buf.Write(b, "  Value ")
		frt.PipeUnit(FTypeToGo(cas.Ftype), (func(_r0 string) { buf.Write(b, _r0) }))
		buf.Write(b, "\n")
	}))
	buf.Write(b, "}\n")
	return buf.String(b)
}

func csConstructorName(unionName string, cas NameTypePair) string {
	return ("New_" + unionCSName(unionName, cas.Name))
}

func csConstructFunc(uname string, cas NameTypePair) string {
	b := buf.New()
	buf.Write(b, "func ")
	frt.PipeUnit(csConstructorName(uname, cas), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "(v ")
	frt.PipeUnit(FTypeToGo(cas.Ftype), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, ") ")
	buf.Write(b, uname)
	buf.Write(b, " { return ")
	buf.Write(b, unionCSName(uname, cas.Name))
	buf.Write(b, "{v} }\n")
	return buf.String(b)
}

func csConstructVar(uname string, cas NameTypePair) string {
	b := buf.New()
	buf.Write(b, "var ")
	frt.PipeUnit(csConstructorName(uname, cas), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, frt.SInterP(" %s = ", uname))
	buf.Write(b, unionCSName(uname, cas.Name))
	buf.Write(b, "{}\n")
	return buf.String(b)
}

func csConstruct(uname string, cas NameTypePair) string {
	return frt.IfElse(frt.OpEqual(cas.Ftype, New_FType_FUnit), (func() string {
		return csConstructVar(uname, cas)
	}), (func() string {
		return csConstructFunc(uname, cas)
	}))
}

func caseToGo(ud UnionDef, cas NameTypePair) string {
	sdf := udCSDef(ud, cas)
	csdf := csConstruct(ud.Name, cas)
	return (((sdf + "\n") + csdf) + "\n")
}

func udfToGo(ud UnionDef) string {
	b := buf.New()
	frt.PipeUnit(udUnionDef(ud), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n")
	frt.PipeUnit(udCSConformMethods(ud), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n")
	frt.PipeUnit(frt.Pipe(frt.Pipe(udCases(ud), (func(_r0 []NameTypePair) []string {
		return slice.Map((func(_r0 NameTypePair) string { return caseToGo(ud, _r0) }), _r0)
	})), (func(_r0 []string) string { return strings.Concat("", _r0) })), (func(_r0 string) { buf.Write(b, _r0) }))
	return buf.String(b)
}

func dsToGo(ds DefStmt) string {
	switch _v1 := (ds).(type) {
	case DefStmt_DRecordDef:
		rd := _v1.Value
		return rdfToGo(rd)
	case DefStmt_DUnionDef:
		ud := _v1.Value
		return udfToGo(ud)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func mdToGo(md MultipleDefs) string {
	return frt.Pipe(frt.Pipe(md.Defs, (func(_r0 []DefStmt) []string { return slice.Map(dsToGo, _r0) })), (func(_r0 []string) string { return strings.Concat("\n", _r0) }))
}

func StmtToGo(stmt Stmt) string {
	eToGo := (func(_r0 Expr) string { return ExprToGo(StmtToGo, _r0) })
	switch _v2 := (stmt).(type) {
	case Stmt_SLetVarDef:
		llvd := _v2.Value
		switch _v3 := (llvd).(type) {
		case LLetVarDef_LLOneVarDef:
			lvd := _v3.Value
			return lvdToGo(eToGo, lvd)
		case LLetVarDef_LLDestVarDef:
			ldvd := _v3.Value
			return ldvdToGo(eToGo, ldvd)
		default:
			panic("Union pattern fail. Never reached here.")
		}
	case Stmt_SExprStmt:
		expr := _v2.Value
		return eToGo(expr)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func rootVarDefToGo(eToGo func(Expr) string, rvd RootVarDef) string {
	vname := rvd.Vdef.Lvar.Name
	rhs := rvd.Vdef.Rhs
	b := buf.New()
	buf.Write(b, frt.SInterP("var %s = ", vname))
	frt.PipeUnit(eToGo(rhs), (func(_r0 string) { buf.Write(b, _r0) }))
	return buf.String(b)
}

func RootStmtToGo(rstmt RootStmt) string {
	eToGo := (func(_r0 Expr) string { return ExprToGo(StmtToGo, _r0) })
	reToGoRet := (func(_r0 ReturnableExpr) string { return reToGoReturn(StmtToGo, eToGo, _r0) })
	bToGoRet := (func(_r0 Block) string { return blockToGoReturn(StmtToGo, eToGo, reToGoRet, _r0) })
	switch _v4 := (rstmt).(type) {
	case RootStmt_RSImport:
		im := _v4.Value
		return imToGo(im)
	case RootStmt_RSPackage:
		pn := _v4.Value
		return pmToGo(pn)
	case RootStmt_RSPackageInfo:
		return ""
	case RootStmt_RSRootFuncDef:
		rfd := _v4.Value
		return rfdToGo(bToGoRet, rfd)
	case RootStmt_RSRootVarDef:
		rvd := _v4.Value
		return rootVarDefToGo(eToGo, rvd)
	case RootStmt_RSDefStmt:
		ds := _v4.Value
		return dsToGo(ds)
	case RootStmt_RSMultipleDefs:
		md := _v4.Value
		return mdToGo(md)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func RootStmtsToGo(rstmts []RootStmt) string {
	return frt.Pipe(frt.Pipe(slice.Map(RootStmtToGo, rstmts), (func(_r0 []string) string { return strings.Concat("\n\n", _r0) })), (func(_r0 string) string { return strings.AppendTail("\n", _r0) }))
}
