package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

type ParseState struct {
	tkz        Tokenizer
	scope      Scope
	offsideCol []int
}

func newParse(tkz Tokenizer, scope Scope, offCols []int) ParseState {
	return ParseState{tkz: tkz, scope: scope, offsideCol: offCols}
}

func psWithTkz(org ParseState, tkz Tokenizer) ParseState {
	return newParse(tkz, org.scope, org.offsideCol)
}

func psWithScope(org ParseState, nsc Scope) ParseState {
	return newParse(org.tkz, nsc, org.offsideCol)
}

func psWithOffside(org ParseState, offs []int) ParseState {
	return newParse(org.tkz, org.scope, offs)
}

func psPushScope(org ParseState) ParseState {
	return frt.Pipe(newScope(org.scope), (func(_r0 Scope) ParseState { return psWithScope(org, _r0) }))
}

func psPopScope(org ParseState) ParseState {
	return frt.Pipe(popScope(org.scope), (func(_r0 Scope) ParseState { return psWithScope(org, _r0) }))
}

func psCurOffside(ps ParseState) int {
	return slice.Last(ps.offsideCol)
}

func psCurCol(ps ParseState) int {
	return ps.tkz.col
}

func psPushOffside(ps ParseState) ParseState {
	curCol := psCurCol(ps)
	frt.IfOnly((psCurOffside(ps) >= curCol), (func() {
		frt.Panic("Overrun offside rule")
	}))
	return frt.Pipe(slice.Append(curCol, ps.offsideCol), (func(_r0 []int) ParseState { return psWithOffside(ps, _r0) }))
}

func psPopOffside(ps ParseState) ParseState {
	return frt.Pipe(slice.PopLast(ps.offsideCol), (func(_r0 []int) ParseState { return psWithOffside(ps, _r0) }))
}

func initParse(src string) ParseState {
	tkz := newTkz(src)
	scope := NewScope()
	offside := ([]int{0})
	return newParse(tkz, scope, offside)
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

func psStrNx(f func(ParseState) string, ps ParseState) frt.Tuple2[ParseState, string] {
	s := f(ps)
	ps2 := psNext(ps)
	return frt.NewTuple2(ps2, s)
}

func psIdentNameNx(ps ParseState) frt.Tuple2[ParseState, string] {
	return psStrNx(psIdentName, ps)
}

func psStringValNx(ps ParseState) frt.Tuple2[ParseState, string] {
	return psStrNx(psStringVal, ps)
}

func psCurrentNx(ps ParseState) frt.Tuple2[ParseState, Token] {
	tk := psCurrent(ps)
	ps2 := psNext(ps)
	return frt.NewTuple2(ps2, tk)
}

func psCurrentTTNx(ps ParseState) frt.Tuple2[ParseState, TokenType] {
	tt := psCurrentTT(ps)
	ps2 := psNext(ps)
	return frt.NewTuple2(ps2, tt)
}

func psIdentNameNxL(ps ParseState) frt.Tuple2[ParseState, string] {
	return frt.Pipe(psIdentNameNx(ps), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return CnvL(psSkipEOL, _r0) }))
}

func psStringValNxL(ps ParseState) frt.Tuple2[ParseState, string] {
	return frt.Pipe(psStringValNx(ps), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return CnvL(psSkipEOL, _r0) }))
}

func psCurrentNxL(ps ParseState) frt.Tuple2[ParseState, Token] {
	return frt.Pipe(psCurrentNx(ps), (func(_r0 frt.Tuple2[ParseState, Token]) frt.Tuple2[ParseState, Token] { return CnvL(psSkipEOL, _r0) }))
}

func psCurrentTTNxL(ps ParseState) frt.Tuple2[ParseState, TokenType] {
	return frt.Pipe(psCurrentTTNx(ps), (func(_r0 frt.Tuple2[ParseState, TokenType]) frt.Tuple2[ParseState, TokenType] {
		return CnvL(psSkipEOL, _r0)
	}))
}

func parsePackage(ps ParseState) frt.Tuple2[ParseState, Stmt] {
	return frt.Pipe(frt.Pipe(psConsume(New_TokenType_PACKAGE, ps), psIdentNameNxL), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, Stmt] {
		return CnvR(New_Stmt_SPackage, _r0)
	}))
}

func parseImport(ps ParseState) frt.Tuple2[ParseState, Stmt] {
	return frt.Pipe(frt.Pipe(psConsume(New_TokenType_IMPORT, ps), psStringValNxL), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, Stmt] {
		return CnvR(New_Stmt_SImport, _r0)
	}))
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
					res, ok := frt.Destr(scLookupType(ps3.scope, tname))
					return frt.IfElse(ok, (func() FType {
						return res
					}), (func() FType {
						frt.Panic("type not found.")
						return New_FType_FUnit
					}))
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
		ps3, tp := frt.Destr(frt.Pipe(frt.Pipe(frt.Pipe(psNext(ps2), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_COLON, _r0) })), parseType), (func(_r0 frt.Tuple2[ParseState, FType]) frt.Tuple2[ParseState, FType] {
			return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RPAREN, _r0) }), _r0)
		})))
		v := Var{name: vname, ftype: tp}
		scDefVar(ps3.scope, vname, v)
		return frt.NewTuple2(ps3, New_Param_PVar(v))
	}
}

func parseParams(ps ParseState) frt.Tuple2[ParseState, []Var] {
	ps2, prm1 := frt.Destr(parseParam(ps))
	switch _v380 := (prm1).(type) {
	case Param_PUnit:
		zero := []Var{}
		return frt.NewTuple2(ps2, zero)
	case Param_PVar:
		v := _v380.Value
		tt := psCurrentTT(ps2)
		switch (tt).(type) {
		case TokenType_LPAREN:
			return frt.Pipe(parseParams(ps2), (func(_r0 frt.Tuple2[ParseState, []Var]) frt.Tuple2[ParseState, []Var] {
				return CnvR((func(_r0 []Var) []Var { return slice.Append(v, _r0) }), _r0)
			}))
		default:
			return frt.NewTuple2(ps2, ([]Var{v}))
		}
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func parseGoEval(ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2 := psNext(ps)
	switch (psCurrentTT(ps2)).(type) {
	case TokenType_LT:
		ps3, ft := frt.Destr(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LT, ps2), parseType), (func(_r0 frt.Tuple2[ParseState, FType]) frt.Tuple2[ParseState, FType] {
			return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_GT, _r0) }), _r0)
		})))
		ps4, s := frt.Destr(psStringValNx(ps3))
		ge := GoEvalExpr{goStmt: s, typeArg: ft}
		return frt.NewTuple2(ps4, New_Expr_EGoEvalExpr(ge))
	case TokenType_STRING:
		ps3, s := frt.Destr(psStringValNx(ps2))
		ge := GoEvalExpr{goStmt: s, typeArg: New_FType_FUnit}
		return frt.NewTuple2(ps3, New_Expr_EGoEvalExpr(ge))
	default:
		frt.Panic("Wrong arg for GoEval")
		return frt.NewTuple2(ps2, New_Expr_EUnit)
	}
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
		return frt.Pipe(frt.Pipe(psConsume(New_TokenType_SEMICOLON, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, []NEPair] { return parseFieldInitializers(parseE, _r0) })), (func(_r0 frt.Tuple2[ParseState, []NEPair]) frt.Tuple2[ParseState, []NEPair] {
			return CnvR((func(_r0 []NEPair) []NEPair { return slice.Prepend(nep, _r0) }), _r0)
		}))
	}))
}

func NVPToName(nvp NEPair) string {
	return nvp.name
}

func parseRecordGen(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2, neps := frt.Destr(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LBRACE, ps), (func(_r0 ParseState) frt.Tuple2[ParseState, []NEPair] { return parseFieldInitializers(parseE, _r0) })), (func(_r0 frt.Tuple2[ParseState, []NEPair]) frt.Tuple2[ParseState, []NEPair] {
		return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RBRACE, _r0) }), _r0)
	})))
	rtype, ok := frt.Destr(frt.Pipe(slice.Map(NVPToName, neps), (func(_r0 []string) frt.Tuple2[RecordType, bool] { return scLookupRecord(ps2.scope, _r0) })))
	return frt.IfElse(ok, (func() frt.Tuple2[ParseState, Expr] {
		return frt.Pipe(frt.Pipe(RecordGen{fieldsNV: neps, recordType: rtype}, New_Expr_ERecordGen), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(ps2, _r0) }))
	}), (func() frt.Tuple2[ParseState, Expr] {
		frt.Panic("record field name match to no record type.")
		return frt.NewTuple2(ps2, New_Expr_EUnit)
	}))
}

func parseAtom(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Expr] {
	cur := psCurrent(ps)
	pn := psNext(ps)
	switch (cur.ttype).(type) {
	case TokenType_STRING:
		return frt.Pipe(New_Expr_EStringLiteral(cur.stringVal), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(pn, _r0) }))
	case TokenType_INT_IMM:
		return frt.Pipe(New_Expr_EIntImm(cur.intVal), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(pn, _r0) }))
	case TokenType_TRUE:
		return frt.Pipe(New_Expr_EBoolLiteral(true), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(pn, _r0) }))
	case TokenType_FALSE:
		return frt.Pipe(New_Expr_EBoolLiteral(false), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(pn, _r0) }))
	case TokenType_LBRACE:
		return parseRecordGen(parseE, ps)
	case TokenType_IDENTIFIER:
		return frt.IfElse(frt.OpEqual(cur.stringVal, "GoEval"), (func() frt.Tuple2[ParseState, Expr] {
			return parseGoEval(ps)
		}), (func() frt.Tuple2[ParseState, Expr] {
			vfac, ok := frt.Destr(scLookupVarFac(ps.scope, cur.stringVal))
			return frt.IfElse(ok, (func() frt.Tuple2[ParseState, Expr] {
				return frt.Pipe(frt.Pipe(vfac(), New_Expr_EVar), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(pn, _r0) }))
			}), (func() frt.Tuple2[ParseState, Expr] {
				frt.Panic("Unknown var ref")
				return frt.NewTuple2(ps, New_Expr_EUnit)
			}))
		}))
	default:
		return frt.NewTuple2(ps, New_Expr_EUnit)
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
	case TokenType_RPAREN:
		return true
	case TokenType_RSBRACKET:
		return true
	case TokenType_WITH:
		return true
	case TokenType_THEN:
		return true
	case TokenType_ELSE:
		return true
	case TokenType_COMMA:
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
		return frt.Pipe(parseAtomList(parseE, ps2), (func(_r0 frt.Tuple2[ParseState, []Expr]) frt.Tuple2[ParseState, []Expr] {
			return CnvR((func(_r0 []Expr) []Expr { return slice.Prepend(one, _r0) }), _r0)
		}))
	}))
}

func parseMatchRule(pBlock func(ParseState) frt.Tuple2[ParseState, Block], target Expr, ps ParseState) frt.Tuple2[ParseState, MatchRule] {
	ps2 := psConsume(New_TokenType_BAR, ps)
	switch (psCurrentTT(ps2)).(type) {
	case TokenType_UNDER_SCORE:
		ps3, block := frt.Destr(frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_UNDER_SCORE, ps2), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RARROW, _r0) })), psSkipEOL), pBlock))
		mp := MatchPattern{caseId: "_", varName: ""}
		return frt.Pipe(MatchRule{pattern: mp, body: block}, (func(_r0 MatchRule) frt.Tuple2[ParseState, MatchRule] { return withPs(ps3, _r0) }))
	default:
		ps3, cname := frt.Destr(psIdentNameNx(ps2))
		ps4, vname := frt.Destr((func() frt.Tuple2[ParseState, string] {
			switch (psCurrentTT(ps3)).(type) {
			case TokenType_RARROW:
				return frt.NewTuple2(ps3, "")
			case TokenType_UNDER_SCORE:
				return frt.Pipe(frt.NewTuple2(ps3, "_"), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return CnvL(psNext, _r0) }))
			default:
				return psIdentNameNx(ps3)
			}
		})())
		ps5 := frt.Pipe(frt.Pipe(psConsume(New_TokenType_RARROW, ps4), psSkipEOL), psPushScope)
		frt.IfOnly((frt.OpNotEqual(vname, "") && frt.OpNotEqual(vname, "_")), (func() {
			tt := ExprToType(target)
			fu := tt.(FType_FUnion).Value
			cp := lookupCase(fu, cname)
			scDefVar(ps5.scope, vname, Var{name: vname, ftype: cp.ftype})
		}))
		ps6, block := frt.Destr(pBlock(ps5))
		mp := MatchPattern{caseId: cname, varName: vname}
		return frt.Pipe(MatchRule{pattern: mp, body: block}, (func(_r0 MatchRule) frt.Tuple2[ParseState, MatchRule] { return withPs(ps6, _r0) }))
	}
}

func parseMatchRules(pBlock func(ParseState) frt.Tuple2[ParseState, Block], target Expr, ps ParseState) frt.Tuple2[ParseState, []MatchRule] {
	ps2, one := frt.Destr(frt.Pipe(parseMatchRule(pBlock, target, ps), (func(_r0 frt.Tuple2[ParseState, MatchRule]) frt.Tuple2[ParseState, MatchRule] {
		return CnvL(psSkipEOL, _r0)
	})))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_BAR), (func() frt.Tuple2[ParseState, []MatchRule] {
		return frt.Pipe(parseMatchRules(pBlock, target, ps2), (func(_r0 frt.Tuple2[ParseState, []MatchRule]) frt.Tuple2[ParseState, []MatchRule] {
			return CnvR((func(_r0 []MatchRule) []MatchRule { return slice.Prepend(one, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []MatchRule] {
		return frt.NewTuple2(ps2, ([]MatchRule{one}))
	}))
}

func parseMatchExpr(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, MatchExpr] {
	ps2, target := frt.Destr(frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_MATCH, ps), pExpr), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Expr] {
		return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_WITH, _r0) }), _r0)
	})), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Expr] { return CnvL(psSkipEOL, _r0) })))
	ps3, rules := frt.Destr(parseMatchRules(pBlock, target, ps2))
	return frt.Pipe(MatchExpr{target: target, rules: rules}, (func(_r0 MatchExpr) frt.Tuple2[ParseState, MatchExpr] { return withPs(ps3, _r0) }))
}

func parseSemiExprs(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, []Expr] {
	ps2, one := frt.Destr(pExpr(ps))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_SEMICOLON), (func() frt.Tuple2[ParseState, []Expr] {
		return frt.Pipe(frt.Pipe(psConsume(New_TokenType_SEMICOLON, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, []Expr] { return parseSemiExprs(pExpr, _r0) })), (func(_r0 frt.Tuple2[ParseState, []Expr]) frt.Tuple2[ParseState, []Expr] {
			return CnvR((func(_r0 []Expr) []Expr { return slice.Prepend(one, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []Expr] {
		return frt.NewTuple2(ps2, ([]Expr{one}))
	}))
}

func parseSliceExpr(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Expr] {
	return frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LSBRACKET, ps), (func(_r0 ParseState) frt.Tuple2[ParseState, []Expr] { return parseSemiExprs(pExpr, _r0) })), (func(_r0 frt.Tuple2[ParseState, []Expr]) frt.Tuple2[ParseState, []Expr] {
		return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RSBRACKET, _r0) }), _r0)
	})), (func(_r0 frt.Tuple2[ParseState, []Expr]) frt.Tuple2[ParseState, Expr] {
		return CnvR(New_Expr_ESlice, _r0)
	}))
}

func parseTerm(pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, Expr] {
	pExpr := (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] { return parseTerm(pBlock, _r0) })
	switch (psCurrentTT(ps)).(type) {
	case TokenType_MATCH:
		return frt.Pipe(frt.Pipe(parseMatchExpr(pExpr, pBlock, ps), (func(_r0 frt.Tuple2[ParseState, MatchExpr]) frt.Tuple2[ParseState, ReturnableExpr] {
			return CnvR(New_ReturnableExpr_RMatchExpr, _r0)
		})), (func(_r0 frt.Tuple2[ParseState, ReturnableExpr]) frt.Tuple2[ParseState, Expr] {
			return CnvR(New_Expr_EReturnableExpr, _r0)
		}))
	case TokenType_LSBRACKET:
		return parseSliceExpr(pExpr, ps)
	default:
		ps2, es := frt.Destr(parseAtomList(pExpr, ps))
		return frt.IfElse(frt.OpEqual(slice.Length(es), 1), (func() frt.Tuple2[ParseState, Expr] {
			return frt.NewTuple2(ps2, slice.Head(es))
		}), (func() frt.Tuple2[ParseState, Expr] {
			head := slice.Head(es)
			tail := slice.Tail(es)
			switch _v388 := (head).(type) {
			case Expr_EVar:
				v := _v388.Value
				fc := FunCall{targetFunc: v, args: tail}
				return frt.NewTuple2(ps2, New_Expr_EFunCall(fc))
			default:
				frt.Panic("Funcall head is not var")
				return frt.NewTuple2(ps2, head)
			}
		}))
	}
}

type StmtLike interface {
	StmtLike_Union()
}

func (StmtLike_SLExpr) StmtLike_Union()    {}
func (StmtLike_SLLetStmt) StmtLike_Union() {}

type StmtLike_SLExpr struct {
	Value Expr
}

func New_StmtLike_SLExpr(v Expr) StmtLike { return StmtLike_SLExpr{v} }

type StmtLike_SLLetStmt struct {
	Value LetVarDef
}

func New_StmtLike_SLLetStmt(v LetVarDef) StmtLike { return StmtLike_SLLetStmt{v} }

func slToStmt(sl StmtLike) Stmt {
	switch _v389 := (sl).(type) {
	case StmtLike_SLExpr:
		e := _v389.Value
		return New_Stmt_SExprStmt(e)
	case StmtLike_SLLetStmt:
		l := _v389.Value
		return New_Stmt_SLetVarDef(l)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func parseStmtLike(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pLet func(ParseState) frt.Tuple2[ParseState, LetVarDef], ps ParseState) frt.Tuple2[ParseState, StmtLike] {
	switch (psCurrentTT(ps)).(type) {
	case TokenType_LET:
		return frt.Pipe(pLet(ps), (func(_r0 frt.Tuple2[ParseState, LetVarDef]) frt.Tuple2[ParseState, StmtLike] {
			return CnvR(New_StmtLike_SLLetStmt, _r0)
		}))
	default:
		return frt.Pipe(pExpr(ps), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, StmtLike] {
			return CnvR(New_StmtLike_SLExpr, _r0)
		}))
	}
}

func isEndOfBlock(ps ParseState) bool {
	isOffside := (psCurCol(ps) < psCurOffside(ps))
	isEof := frt.OpEqual(psCurrentTT(ps), New_TokenType_EOF)
	return (isOffside || isEof)
}

func parseStmtLikeList(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pLet func(ParseState) frt.Tuple2[ParseState, LetVarDef], ps ParseState) frt.Tuple2[ParseState, []StmtLike] {
	ps2, one := frt.Destr(frt.Pipe(parseStmtLike(pExpr, pLet, ps), (func(_r0 frt.Tuple2[ParseState, StmtLike]) frt.Tuple2[ParseState, StmtLike] {
		return CnvL(psSkipEOL, _r0)
	})))
	return frt.IfElse(isEndOfBlock(ps2), (func() frt.Tuple2[ParseState, []StmtLike] {
		return frt.NewTuple2(ps2, ([]StmtLike{one}))
	}), (func() frt.Tuple2[ParseState, []StmtLike] {
		return frt.Pipe(parseStmtLikeList(pExpr, pLet, ps2), (func(_r0 frt.Tuple2[ParseState, []StmtLike]) frt.Tuple2[ParseState, []StmtLike] {
			return CnvR((func(_r0 []StmtLike) []StmtLike { return slice.Prepend(one, _r0) }), _r0)
		}))
	}))
}

func emptyBlock() Block {
	return Block{}
}

func parseBlockAfterPushScope(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pLet func(ParseState) frt.Tuple2[ParseState, LetVarDef], ps ParseState) frt.Tuple2[ParseState, Block] {
	ps2, sls := frt.Destr(frt.Pipe(frt.Pipe(psPushOffside(ps), (func(_r0 ParseState) frt.Tuple2[ParseState, []StmtLike] { return parseStmtLikeList(pExpr, pLet, _r0) })), (func(_r0 frt.Tuple2[ParseState, []StmtLike]) frt.Tuple2[ParseState, []StmtLike] {
		return CnvL(psPopOffside, _r0)
	})))
	last := slice.Last(sls)
	stmts := frt.Pipe(slice.PopLast(sls), (func(_r0 []StmtLike) []Stmt { return slice.Map(slToStmt, _r0) }))
	switch _v391 := (last).(type) {
	case StmtLike_SLExpr:
		e := _v391.Value
		return frt.NewTuple2(ps2, Block{stmts: stmts, finalExpr: e})
	default:
		frt.Panic("block of last is not expr")
		return frt.Pipe(emptyBlock(), (func(_r0 Block) frt.Tuple2[ParseState, Block] { return withPs(ps2, _r0) }))
	}
}

func parseBlock(pLet func(ParseState) frt.Tuple2[ParseState, LetVarDef], ps ParseState) frt.Tuple2[ParseState, Block] {
	pExpr := (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] {
		return parseTerm((func(_r0 ParseState) frt.Tuple2[ParseState, Block] { return parseBlock(pLet, _r0) }), _r0)
	})
	return frt.Pipe(psPushScope(ps), (func(_r0 ParseState) frt.Tuple2[ParseState, Block] { return parseBlockAfterPushScope(pExpr, pLet, _r0) }))
}

func parseLetVarDef(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, LetVarDef] {
	ps2, vname := frt.Destr(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LET, ps), psIdentNameNx), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] {
		return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_EQ, _r0) }), _r0)
	})))
	ps3, rhs := frt.Destr(pExpr(ps2))
	v := Var{name: vname, ftype: ExprToType(rhs)}
	scDefVar(ps3.scope, vname, v)
	return frt.Pipe(LetVarDef{name: vname, rhs: rhs}, (func(_r0 LetVarDef) frt.Tuple2[ParseState, LetVarDef] { return withPs(ps3, _r0) }))
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

func parseLetFuncDef(pLet func(ParseState) frt.Tuple2[ParseState, LetVarDef], ps ParseState) frt.Tuple2[ParseState, LetFuncDef] {
	ps2 := psConsume(New_TokenType_LET, ps)
	fname := psIdentName(ps2)
	ps3, params := frt.Destr(frt.Pipe(psNext(ps2), parseParams))
	ps4, block := frt.Destr(frt.Pipe(frt.Pipe(psConsume(New_TokenType_EQ, ps3), psSkipEOL), (func(_r0 ParseState) frt.Tuple2[ParseState, Block] { return parseBlock(pLet, _r0) })))
	lfd := LetFuncDef{name: fname, params: params, body: block}
	frt.PipeUnit(lfdToFuncVar(lfd), (func(_r0 Var) { scDefVar(ps4.scope, fname, _r0) }))
	return frt.NewTuple2(ps4, lfd)
}

func parseLet(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Stmt] {
	pLet := (func(_r0 ParseState) frt.Tuple2[ParseState, LetVarDef] { return parseLetVarDef(pExpr, _r0) })
	psN := psNext(ps)
	switch (psCurrentTT(psN)).(type) {
	case TokenType_LPAREN:
		return frt.Pipe(parseLetFuncDef(pLet, ps), (func(_r0 frt.Tuple2[ParseState, LetFuncDef]) frt.Tuple2[ParseState, Stmt] {
			return CnvR(New_Stmt_SLetFuncDef, _r0)
		}))
	default:
		psNN := psNext(psN)
		switch (psCurrentTT(psNN)).(type) {
		case TokenType_EQ:
			return frt.Pipe(pLet(ps), (func(_r0 frt.Tuple2[ParseState, LetVarDef]) frt.Tuple2[ParseState, Stmt] {
				return CnvR(New_Stmt_SLetVarDef, _r0)
			}))
		default:
			return frt.Pipe(parseLetFuncDef(pLet, ps), (func(_r0 frt.Tuple2[ParseState, LetFuncDef]) frt.Tuple2[ParseState, Stmt] {
				return CnvR(New_Stmt_SLetFuncDef, _r0)
			}))
		}
	}
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
		return frt.Pipe(frt.Pipe(psConsume(New_TokenType_SEMICOLON, ps2), parseFieldDefs), (func(_r0 frt.Tuple2[ParseState, []NameTypePair]) frt.Tuple2[ParseState, []NameTypePair] {
			return CnvR((func(_r0 []NameTypePair) []NameTypePair { return slice.Prepend(ntp, _r0) }), _r0)
		}))
	}))
}

func rdToRecType(rd RecordDef) RecordType {
	return RecordType(rd)
}

func parseRecordDef(tname string, ps ParseState) frt.Tuple2[ParseState, RecordDef] {
	ps2, ntps := frt.Destr(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LBRACE, ps), parseFieldDefs), (func(_r0 frt.Tuple2[ParseState, []NameTypePair]) frt.Tuple2[ParseState, []NameTypePair] {
		return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RBRACE, _r0) }), _r0)
	})))
	rd := RecordDef{name: tname, fields: ntps}
	frt.PipeUnit(rdToRecType(rd), (func(_r0 RecordType) { scRegisterRecType(ps2.scope, _r0) }))
	return frt.NewTuple2(ps2, rd)
}

func parseOneCaseDef(ps ParseState) frt.Tuple2[ParseState, NameTypePair] {
	ps2, cname := frt.Destr(frt.Pipe(psConsume(New_TokenType_BAR, ps), psIdentNameNx))
	switch (psCurrentTT(ps2)).(type) {
	case TokenType_OF:
		ps3, tp := frt.Destr(frt.Pipe(psConsume(New_TokenType_OF, ps2), parseType))
		cs := NameTypePair{name: cname, ftype: tp}
		return frt.NewTuple2(ps3, cs)
	default:
		ps3 := psConsume(New_TokenType_EOL, ps2)
		cs := NameTypePair{name: cname, ftype: New_FType_FUnit}
		return frt.NewTuple2(ps3, cs)
	}
}

func parseCaseDefs(ps ParseState) frt.Tuple2[ParseState, []NameTypePair] {
	ps2, cs := frt.Destr(parseOneCaseDef(ps))
	ps3 := psSkipEOL(ps2)
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps3), New_TokenType_BAR), (func() frt.Tuple2[ParseState, []NameTypePair] {
		return frt.Pipe(parseCaseDefs(ps3), (func(_r0 frt.Tuple2[ParseState, []NameTypePair]) frt.Tuple2[ParseState, []NameTypePair] {
			return CnvR((func(_r0 []NameTypePair) []NameTypePair { return slice.Prepend(cs, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []NameTypePair] {
		return frt.NewTuple2(ps2, ([]NameTypePair{cs}))
	}))
}

func udToUt(ud UnionDef) UnionType {
	return UnionType(ud)
}

func udToFUt(ud UnionDef) FType {
	return frt.Pipe(udToUt(ud), New_FType_FUnion)
}

func csRegisterCtor(sc Scope, ud UnionDef, cas NameTypePair) int {
	ctorName := csConstructorName(ud.name, cas)
	ut := udToFUt(ud)
	v := (func() Var {
		switch (cas.ftype).(type) {
		case FType_FUnit:
			return Var{name: ctorName, ftype: ut}
		default:
			tps := ([]FType{cas.ftype, ut})
			funcTp := New_FType_FFunc(FuncType{targets: tps})
			return Var{name: ctorName, ftype: funcTp}
		}
	})()
	scDefVar(sc, cas.name, v)
	return 1
}

func udRegisterCsCtors(sc Scope, ud UnionDef) int {
	frt.Pipe(ud.cases, (func(_r0 []NameTypePair) []int {
		return slice.Map((func(_r0 NameTypePair) int { return csRegisterCtor(sc, ud, _r0) }), _r0)
	}))
	return 1
}

func udRegisterToScope(sc Scope, ud UnionDef) {
	udRegisterCsCtors(sc, ud)
	frt.PipeUnit(udToFUt(ud), (func(_r0 FType) { scRegisterType(sc, ud.name, _r0) }))
}

func parseUnionDef(tname string, ps ParseState) frt.Tuple2[ParseState, UnionDef] {
	ps2, css := frt.Destr(parseCaseDefs(ps))
	ud := UnionDef{name: tname, cases: css}
	udRegisterToScope(ps2.scope, ud)
	return frt.NewTuple2(ps2, ud)
}

func parseTypeDef(ps ParseState) frt.Tuple2[ParseState, Stmt] {
	ps2, tname := frt.Destr(frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_TYPE, ps), psIdentNameNxL), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] {
		return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_EQ, _r0) }), _r0)
	})), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return CnvL(psSkipEOL, _r0) })))
	switch (psCurrentTT(ps2)).(type) {
	case TokenType_LBRACE:
		return frt.Pipe(frt.Pipe(parseRecordDef(tname, ps2), (func(_r0 frt.Tuple2[ParseState, RecordDef]) frt.Tuple2[ParseState, DefStmt] {
			return CnvR(New_DefStmt_DRecordDef, _r0)
		})), (func(_r0 frt.Tuple2[ParseState, DefStmt]) frt.Tuple2[ParseState, Stmt] {
			return CnvR(New_Stmt_SDefStmt, _r0)
		}))
	case TokenType_BAR:
		return frt.Pipe(frt.Pipe(parseUnionDef(tname, ps2), (func(_r0 frt.Tuple2[ParseState, UnionDef]) frt.Tuple2[ParseState, DefStmt] {
			return CnvR(New_DefStmt_DUnionDef, _r0)
		})), (func(_r0 frt.Tuple2[ParseState, DefStmt]) frt.Tuple2[ParseState, Stmt] {
			return CnvR(New_Stmt_SDefStmt, _r0)
		}))
	default:
		frt.Panic("NYI")
		return frt.NewTuple2(ps2, New_Stmt_SExprStmt(New_Expr_EUnit))
	}
}

func dummyPBlock(ps ParseState) frt.Tuple2[ParseState, Block] {
	frt.Panic("Should never called.")
	return frt.Pipe(emptyBlock(), (func(_r0 Block) frt.Tuple2[ParseState, Block] { return withPs(ps, _r0) }))
}

func parseExprNoBlock(ps ParseState) frt.Tuple2[ParseState, Expr] {
	return parseTerm(dummyPBlock, ps)
}

func parseStmt(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Stmt] {
	switch (psCurrentTT(ps)).(type) {
	case TokenType_PACKAGE:
		return parsePackage(ps)
	case TokenType_IMPORT:
		return parseImport(ps)
	case TokenType_LET:
		return parseLet(pExpr, ps)
	case TokenType_TYPE:
		return parseTypeDef(ps)
	default:
		frt.Panic("Unknown stmt")
		return parsePackage(ps)
	}
}

func parseStmts(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, []Stmt] {
	ps2 := psSkipEOL(ps)
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_EOF), (func() frt.Tuple2[ParseState, []Stmt] {
		s := []Stmt{}
		return frt.NewTuple2(ps2, s)
	}), (func() frt.Tuple2[ParseState, []Stmt] {
		ps3, one := frt.Destr(frt.Pipe(parseStmt(pExpr, ps), (func(_r0 frt.Tuple2[ParseState, Stmt]) frt.Tuple2[ParseState, Stmt] { return CnvL(psSkipEOL, _r0) })))
		ps4, rest := frt.Destr(parseStmts(pExpr, ps3))
		ss := slice.Prepend(one, rest)
		return frt.NewTuple2(ps4, ss)
	}))
}
