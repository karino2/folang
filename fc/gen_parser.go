package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

type ParseState struct {
	tkz   Tokenizer
	scope Scope
}

func newParse(tkz Tokenizer, scope Scope) ParseState {
	return ParseState{tkz: tkz, scope: scope}
}

func psWithTkz(org ParseState, tkz Tokenizer) ParseState {
	return ParseState{tkz: tkz, scope: org.scope}
}

func psWithScope(org ParseState, nsc Scope) ParseState {
	return ParseState{tkz: org.tkz, scope: nsc}
}

func initParse(src string) ParseState {
	tkz := newTkz(src)
	scope := NewScope()
	return newParse(tkz, scope)
}

func psCurrent(ps ParseState) Token {
	return ps.tkz.current
}

func psCurrentTT(ps ParseState) TokenType {
	tk := psCurrent(ps)
	return tk.ttype
}

func psNext(ps ParseState) ParseState {
	ntk := tkzNext(ps.tkz)
	return psWithTkz(ps, ntk)
}

func psNextNOL(ps ParseState) ParseState {
	ntk := tkzNextNOL(ps.tkz)
	return psWithTkz(ps, ntk)
}

func psSkipEOL(ps ParseState) ParseState {
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps), New_TokenType_EOL), (func() ParseState {
		return psNextNOL(ps)
	}), (func() ParseState {
		return ps
	}))
}

func psExpect(ttype TokenType, ps ParseState) {
	cur := psCurrent(ps)
	frt.IfOnly(frt.OpNotEqual(cur.ttype, ttype), (func() {
		frt.Panic("non expected token")
	}))

}

func psConsume(ttype TokenType, ps ParseState) ParseState {
	psExpect(ttype, ps)
	return psNext(ps)
}

func psIdentName(ps ParseState) string {
	psExpect(New_TokenType_IDENTIFIER, ps)
	cur := psCurrent(ps)
	return cur.stringVal
}

func psStringVal(ps ParseState) string {
	psExpect(New_TokenType_STRING, ps)
	cur := psCurrent(ps)
	return cur.stringVal
}

func parsePackage(ps ParseState) frt.Tuple2[ParseState, Stmt] {
	ps2 := psConsume(New_TokenType_PACKAGE, ps)
	pname := psIdentName(ps2)
	ps3 := psNextNOL(ps2)
	pkg := New_Stmt_Package(pname)
	return frt.NewTuple2(ps3, pkg)
}

func parseImport(ps ParseState) frt.Tuple2[ParseState, Stmt] {
	ps2 := psConsume(New_TokenType_IMPORT, ps)
	pname := psStringVal(ps2)
	ps3 := psNextNOL(ps2)
	imp := New_Stmt_Import(pname)
	return frt.NewTuple2(ps3, imp)
}

func parseType(ps ParseState) frt.Tuple2[ParseState, FType] {
	tk := psCurrent(ps)
	switch (tk.ttype).(type) {
	case TokenType_LPAREN:
		ps2 := frt.Pipe(frt.Pipe(ps, (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_LPAREN, _r0) })), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RPAREN, _r0) }))
		return frt.NewTuple2(ps2, New_FType_FUnit)
	case TokenType_IDENTIFIER:
		tname := tk.stringVal
		ps3 := psNext(ps)
		rtype := frt.IfElse(frt.OpEqual(tname, "string"), (func() FType {
			return New_FType_FString
		}), (func() FType {
			return frt.IfElse(frt.OpEqual(tname, "int"), (func() FType {
				return New_FType_FInt
			}), (func() FType {
				return frt.IfElse(frt.OpEqual(tname, "bool"), (func() FType {
					return New_FType_FBool
				}), (func() FType {
					frt.Panic("NYI")
					return New_FType_FUnit
				}))
			}))
		}))
		return frt.NewTuple2(ps3, rtype)
	default:
		frt.Panic("Unknown type")
		return frt.NewTuple2(ps, New_FType_FUnit)
	}
}

type Param interface {
	Param_Union()
}

func (Param_PVar) Param_Union()  {}
func (Param_PUnit) Param_Union() {}

type Param_PVar struct {
	Value Var
}

func New_Param_PVar(v Var) Param { return Param_PVar{v} }

type Param_PUnit struct {
}

var New_Param_PUnit Param = Param_PUnit{}

func parseParam(ps ParseState) frt.Tuple2[ParseState, Param] {
	ps2 := psConsume(New_TokenType_LPAREN, ps)
	tk := psCurrent(ps2)
	switch (tk.ttype).(type) {
	case TokenType_RPAREN:
		ps3 := psConsume(New_TokenType_RPAREN, ps2)
		return frt.NewTuple2(ps3, New_Param_PUnit)
	default:
		vname := psIdentName(ps2)
		ps3 := frt.Pipe(psNext(ps2), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_COLON, _r0) }))
		ps4, tp := frt.Destr(parseType(ps3))
		ps5 := psConsume(New_TokenType_RPAREN, ps4)
		v := Var{name: vname, ftype: tp}
		return frt.NewTuple2(ps5, New_Param_PVar(v))
	}
}

func parseParams(ps ParseState) frt.Tuple2[ParseState, []Var] {
	ps2, prm1 := frt.Destr(parseParam(ps))
	switch _v218 := (prm1).(type) {
	case Param_PUnit:
		zero := []Var{}
		return frt.NewTuple2(ps2, zero)
	case Param_PVar:
		v := _v218.Value
		tt := psCurrentTT(ps2)
		switch (tt).(type) {
		case TokenType_LPAREN:
			ftp2 := parseParams(ps2)
			ps3 := frt.Fst(ftp2)
			prms2 := frt.Snd(ftp2)
			pas3 := slice.Append(v, prms2)
			return frt.NewTuple2(ps3, pas3)
		default:
			return frt.NewTuple2(ps2, ([]Var{v}))
		}
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func parseGoEval(ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2 := psNext(ps)
	cur := psCurrent(ps2)
	ge := GoEvalExpr{goStmt: cur.stringVal, typeArg: New_FType_FUnit}
	ps3 := psNext(ps2)
	return frt.NewTuple2(ps3, New_Expr_GoEvalExpr(ge))
}

func parseFiIni(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, NEPair] {
	fname := psIdentName(ps)
	ps2, expr := frt.Destr(frt.Pipe(frt.Pipe(frt.Pipe(psNextNOL(ps), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_EQ, _r0) })), psSkipEOL), parseE))
	return frt.NewTuple2(ps2, NEPair{name: fname, expr: expr})
}

func parseFieldInitializers(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, []NEPair] {
	ps2, nep := frt.Destr(parseFiIni(parseE, ps))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_RBRACE), (func() frt.Tuple2[ParseState, []NEPair] {
		return frt.NewTuple2(ps2, ([]NEPair{nep}))
	}), (func() frt.Tuple2[ParseState, []NEPair] {
		ps3, rest := frt.Destr(frt.Pipe(psConsume(New_TokenType_SEMICOLON, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, []NEPair] { return parseFieldInitializers(parseE, _r0) })))
		return frt.Pipe(slice.Prepend(nep, rest), (func(_r0 []NEPair) frt.Tuple2[ParseState, []NEPair] { return withPs(ps3, _r0) }))
	}))
}

func NVPToName(nvp NEPair) string {
	return nvp.name
}

func parseRecordGen(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2, neps := frt.Destr(frt.Pipe(psConsume(New_TokenType_LBRACE, ps), (func(_r0 ParseState) frt.Tuple2[ParseState, []NEPair] { return parseFieldInitializers(parseE, _r0) })))
	ps3 := psConsume(New_TokenType_RBRACE, ps2)
	rtype, ok := frt.Destr(frt.Pipe(slice.Map(NVPToName, neps), (func(_r0 []string) frt.Tuple2[RecordType, bool] { return scLookupRecord(ps3.scope, _r0) })))
	return frt.IfElse(ok, (func() frt.Tuple2[ParseState, Expr] {
		return frt.Pipe(frt.Pipe(RecordGen{fieldsNV: neps, recordType: rtype}, New_Expr_RecordGen), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(ps3, _r0) }))
	}), (func() frt.Tuple2[ParseState, Expr] {
		frt.Panic("record field name match to no record type.")
		return frt.NewTuple2(ps3, New_Expr_Unit)
	}))
}

func parseAtom(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Expr] {
	cur := psCurrent(ps)
	pn := psNext(ps)
	switch (cur.ttype).(type) {
	case TokenType_STRING:
		return frt.Pipe(New_Expr_StringLiteral(cur.stringVal), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(pn, _r0) }))
	case TokenType_INT_IMM:
		return frt.Pipe(New_Expr_IntImm(cur.intVal), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(pn, _r0) }))
	case TokenType_LBRACE:
		return parseRecordGen(parseE, ps)
	case TokenType_IDENTIFIER:
		return frt.IfElse(frt.OpEqual(cur.stringVal, "GoEval"), (func() frt.Tuple2[ParseState, Expr] {
			return parseGoEval(ps)
		}), (func() frt.Tuple2[ParseState, Expr] {
			vfac, ok := frt.Destr(scLookupVarFac(ps.scope, cur.stringVal))
			return frt.IfElse(ok, (func() frt.Tuple2[ParseState, Expr] {
				return frt.Pipe(frt.Pipe(vfac(), New_Expr_Var), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(pn, _r0) }))
			}), (func() frt.Tuple2[ParseState, Expr] {
				frt.Panic("Unkonw var ref")
				return frt.NewTuple2(ps, New_Expr_Unit)
			}))
		}))
	default:
		return frt.NewTuple2(ps, New_Expr_Unit)
	}
}

func isEndOfTerm(ps ParseState) bool {
	switch (psCurrentTT(ps)).(type) {
	case TokenType_EOF:
		return true
	case TokenType_EOL:
		return true
	case TokenType_SEMICOLON:
		return true
	case TokenType_RBRACE:
		return true
	default:
		return false
	}
}

func parseAtomList(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, []Expr] {
	ps2, one := frt.Destr(parseAtom(parseE, ps))
	return frt.IfElse(isEndOfTerm(ps2), (func() frt.Tuple2[ParseState, []Expr] {
		return frt.NewTuple2(ps2, ([]Expr{one}))
	}), (func() frt.Tuple2[ParseState, []Expr] {
		ps3, rest := frt.Destr(parseAtomList(parseE, ps2))
		al := slice.Prepend(one, rest)
		return frt.NewTuple2(ps3, al)
	}))
}

func parseTerm(ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2, es := frt.Destr(parseAtomList(parseTerm, ps))
	return frt.IfElse(frt.OpEqual(slice.Length(es), 1), (func() frt.Tuple2[ParseState, Expr] {
		return frt.NewTuple2(ps2, slice.Head(es))
	}), (func() frt.Tuple2[ParseState, Expr] {
		head := slice.Head(es)
		tail := slice.Tail(es)
		switch _v222 := (head).(type) {
		case Expr_Var:
			v := _v222.Value
			fc := FunCall{targetFunc: v, args: tail}
			return frt.NewTuple2(ps2, New_Expr_FunCall(fc))
		default:
			frt.Panic("Funcall head is not var")
			return frt.NewTuple2(ps2, head)
		}
	}))
}

func parseBlock(ps ParseState) frt.Tuple2[ParseState, Block] {
	ps2, expr := frt.Destr(parseTerm(ps))
	block := Block{[]Stmt{}, expr}
	ps3 := psSkipEOL(ps2)
	return frt.NewTuple2(ps3, block)
}

func vToT(v Var) FType {
	return v.ftype
}

func lfdToFuncType(lfd LetFuncDef) FuncType {
	rtype := frt.Pipe(blockToExpr(lfd.body), ExprToType)
	targets := frt.Pipe(slice.Map(vToT, lfd.params), (func(_r0 []FType) []FType { return slice.Append(rtype, _r0) }))
	return FuncType{targets: targets}
}

func lfdToFuncVar(lfd LetFuncDef) Var {
	ft := frt.Pipe(lfdToFuncType(lfd), New_FType_FFunc)
	return Var{name: lfd.name, ftype: ft}
}

func parseLetFuncDef(ps ParseState) frt.Tuple2[ParseState, Stmt] {
	ps2 := psConsume(New_TokenType_LET, ps)
	fname := psIdentName(ps2)
	ps3, params := frt.Destr(frt.Pipe(psNext(ps2), parseParams))
	ps4, block := frt.Destr(frt.Pipe(frt.Pipe(psConsume(New_TokenType_EQ, ps3), psSkipEOL), parseBlock))
	lfd := LetFuncDef{name: fname, params: params, body: block}
	frt.PipeUnit(lfdToFuncVar(lfd), (func(_r0 Var) { scDefVar(ps4.scope, fname, _r0) }))
	stmt := New_Stmt_LetFuncDef(lfd)
	return frt.NewTuple2(ps4, stmt)
}

func parseFieldDef(ps ParseState) frt.Tuple2[ParseState, NameTypePair] {
	fname := psIdentName(ps)
	ps2, tp := frt.Destr(frt.Pipe(frt.Pipe(psNextNOL(ps), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_COLON, _r0) })), parseType))
	ntp := NameTypePair{name: fname, ftype: tp}
	return frt.NewTuple2(ps2, ntp)
}

func parseFieldDefs(ps ParseState) frt.Tuple2[ParseState, []NameTypePair] {
	ps2, ntp := frt.Destr(parseFieldDef(ps))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_RBRACE), (func() frt.Tuple2[ParseState, []NameTypePair] {
		return frt.NewTuple2(ps2, ([]NameTypePair{ntp}))
	}), (func() frt.Tuple2[ParseState, []NameTypePair] {
		ps3, rest := frt.Destr(frt.Pipe(psConsume(New_TokenType_SEMICOLON, ps2), parseFieldDefs))
		result := slice.Prepend(ntp, rest)
		return frt.NewTuple2(ps3, result)
	}))
}

func rdToRecType(rd RecordDef) RecordType {
	return RecordType(rd)
}

func parseRecordDef(tname string, ps ParseState) frt.Tuple2[ParseState, RecordDef] {
	ps2, ntps := frt.Destr(frt.Pipe(psConsume(New_TokenType_LBRACE, ps), parseFieldDefs))
	ps3 := psConsume(New_TokenType_RBRACE, ps2)
	rd := RecordDef{name: tname, fields: ntps}
	recType := rdToRecType(rd)
	scRegisterRecType(ps3.scope, recType)
	return frt.NewTuple2(ps3, rd)
}

func parseTypeDef(ps ParseState) frt.Tuple2[ParseState, Stmt] {
	ps2 := psConsume(New_TokenType_TYPE, ps)
	tname := psIdentName(ps2)
	ps3, rd := frt.Destr(frt.Pipe(frt.Pipe(psNextNOL(ps2), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_EQ, _r0) })), (func(_r0 ParseState) frt.Tuple2[ParseState, RecordDef] { return parseRecordDef(tname, _r0) })))
	rdstmt := frt.Pipe(New_DefStmt_RecordDef(rd), New_Stmt_DefStmt)
	return frt.NewTuple2(ps3, rdstmt)
}

func parseStmt(ps ParseState) frt.Tuple2[ParseState, Stmt] {
	switch (psCurrentTT(ps)).(type) {
	case TokenType_PACKAGE:
		return parsePackage(ps)
	case TokenType_IMPORT:
		return parseImport(ps)
	case TokenType_LET:
		return parseLetFuncDef(ps)
	case TokenType_TYPE:
		return parseTypeDef(ps)
	default:
		frt.Panic("Unknown stmt")
		return parsePackage(ps)
	}
}

func parseStmts(ps ParseState) frt.Tuple2[ParseState, []Stmt] {
	ps2 := psSkipEOL(ps)
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_EOF), (func() frt.Tuple2[ParseState, []Stmt] {
		s := []Stmt{}
		return frt.NewTuple2(ps2, s)
	}), (func() frt.Tuple2[ParseState, []Stmt] {
		ps3, one := frt.Destr(parseStmt(ps))
		ps4, rest := frt.Destr(parseStmts(psSkipEOL(ps3)))
		ss := slice.Prepend(one, rest)
		return frt.NewTuple2(ps4, ss)
	}))
}
