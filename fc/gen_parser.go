package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

func parsePackage(ps ParseState) frt.Tuple2[ParseState, RootStmt] {
	return frt.Pipe(frt.Pipe(psConsume(New_TokenType_PACKAGE, ps), psIdentNameNxL), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, RootStmt] {
		return CnvR(New_RootStmt_RSPackage, _r0)
	}))
}

func parseImport(ps ParseState) frt.Tuple2[ParseState, RootStmt] {
	ps2 := psConsume(New_TokenType_IMPORT, ps)
	return frt.IfElse(psCurIs(New_TokenType_IDENTIFIER, ps2), (func() frt.Tuple2[ParseState, RootStmt] {
		ps3, iname := frt.Destr(psIdentNameNxL(ps2))
		rstmt := frt.Pipe(frt.Sprintf1("github.com/karino2/folang/pkg/%s", iname), New_RootStmt_RSImport)
		return frt.NewTuple2(ps3, rstmt)
	}), (func() frt.Tuple2[ParseState, RootStmt] {
		return frt.Pipe(frt.Pipe(ps2, psStringValNxL), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, RootStmt] {
			return CnvR(New_RootStmt_RSImport, _r0)
		}))
	}))
}

func parseFullName(ps ParseState) frt.Tuple2[ParseState, string] {
	ps2, one := frt.Destr(psIdentNameNx(ps))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_DOT), (func() frt.Tuple2[ParseState, string] {
		ps3, rest := frt.Destr(frt.Pipe(psConsume(New_TokenType_DOT, ps2), parseFullName))
		return frt.NewTuple2(ps3, ((one + ".") + rest))
	}), (func() frt.Tuple2[ParseState, string] {
		return frt.NewTuple2(ps2, one)
	}))
}

func parseAtomType(pType func(ParseState) frt.Tuple2[ParseState, FType], pTerm func(ParseState) frt.Tuple2[ParseState, FType], ps ParseState) frt.Tuple2[ParseState, FType] {
	tk := psCurrent(ps)
	switch (tk.ttype).(type) {
	case TokenType_LPAREN:
		ps2 := psConsume(New_TokenType_LPAREN, ps)
		return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_RPAREN), (func() frt.Tuple2[ParseState, FType] {
			ps3 := psConsume(New_TokenType_RPAREN, ps2)
			return frt.NewTuple2(ps3, New_FType_FUnit)
		}), (func() frt.Tuple2[ParseState, FType] {
			return frt.Pipe(pType(ps2), (func(_r0 frt.Tuple2[ParseState, FType]) frt.Tuple2[ParseState, FType] {
				return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RPAREN, _r0) }), _r0)
			}))
		}))
	case TokenType_LSBRACKET:
		ps2, et := frt.Destr(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LSBRACKET, ps), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RSBRACKET, _r0) })), pTerm))
		return frt.Pipe(frt.Pipe(SliceType{elemType: et}, New_FType_FSlice), (func(_r0 FType) frt.Tuple2[ParseState, FType] { return withPs(ps2, _r0) }))
	case TokenType_IDENTIFIER:
		tname := tk.stringVal
		ps3 := psNext(ps)
		return frt.IfElse(frt.OpEqual(tname, "string"), (func() frt.Tuple2[ParseState, FType] {
			return frt.NewTuple2(ps3, New_FType_FString)
		}), (func() frt.Tuple2[ParseState, FType] {
			return frt.IfElse(frt.OpEqual(tname, "int"), (func() frt.Tuple2[ParseState, FType] {
				return frt.NewTuple2(ps3, New_FType_FInt)
			}), (func() frt.Tuple2[ParseState, FType] {
				return frt.IfElse(frt.OpEqual(tname, "bool"), (func() frt.Tuple2[ParseState, FType] {
					return frt.NewTuple2(ps3, New_FType_FBool)
				}), (func() frt.Tuple2[ParseState, FType] {
					ps4, fullName := frt.Destr(parseFullName(ps))
					res, ok := frt.Destr(scLookupType(ps3.scope, fullName))
					return frt.IfElse(ok, (func() frt.Tuple2[ParseState, FType] {
						return frt.NewTuple2(ps4, res)
					}), (func() frt.Tuple2[ParseState, FType] {
						frt.Panic("type not found.")
						return frt.NewTuple2(ps4, New_FType_FUnit)
					}))
				}))
			}))
		}))
	default:
		frt.Panic("Unknown type")
		return frt.NewTuple2(ps, New_FType_FUnit)
	}
}

func parseTermType(pType func(ParseState) frt.Tuple2[ParseState, FType], ps ParseState) frt.Tuple2[ParseState, FType] {
	pAtom := (func(_r0 ParseState) frt.Tuple2[ParseState, FType] {
		return parseAtomType(pType, (func(_r0 ParseState) frt.Tuple2[ParseState, FType] { return parseTermType(pType, _r0) }), _r0)
	})
	ps2, ft := frt.Destr(pAtom(ps))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_ASTER), (func() frt.Tuple2[ParseState, FType] {
		ps3, ft2 := frt.Destr(frt.Pipe(psConsume(New_TokenType_ASTER, ps2), pAtom))
		frt.IfOnly(frt.OpEqual(psCurrentTT(ps3), New_TokenType_ASTER), (func() {
			frt.Panic("More than three elem tuple, NYI")
		}))
		return frt.Pipe(frt.Pipe(TupleType{elemTypes: ([]FType{ft, ft2})}, New_FType_FTuple), (func(_r0 FType) frt.Tuple2[ParseState, FType] { return withPs(ps3, _r0) }))
	}), (func() frt.Tuple2[ParseState, FType] {
		return frt.NewTuple2(ps2, ft)
	}))
}

func parseTypeArrows(pType func(ParseState) frt.Tuple2[ParseState, FType], ps ParseState) frt.Tuple2[ParseState, []FType] {
	ps2, one := frt.Destr(parseTermType(pType, ps))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_RARROW), (func() frt.Tuple2[ParseState, []FType] {
		return frt.Pipe(frt.Pipe(psConsume(New_TokenType_RARROW, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, []FType] { return parseTypeArrows(pType, _r0) })), (func(_r0 frt.Tuple2[ParseState, []FType]) frt.Tuple2[ParseState, []FType] {
			return CnvR((func(_r0 []FType) []FType { return slice.PushHead(one, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []FType] {
		return frt.NewTuple2(ps2, ([]FType{one}))
	}))
}

func parseType(ps ParseState) frt.Tuple2[ParseState, FType] {
	ps2, tps := frt.Destr(parseTypeArrows(parseType, ps))
	return frt.IfElse(frt.OpEqual(slice.Length(tps), 1), (func() frt.Tuple2[ParseState, FType] {
		return frt.Pipe(slice.Head(tps), (func(_r0 FType) frt.Tuple2[ParseState, FType] { return withPs(ps2, _r0) }))
	}), (func() frt.Tuple2[ParseState, FType] {
		return frt.Pipe(New_FType_FFunc(FuncType{targets: tps}), (func(_r0 FType) frt.Tuple2[ParseState, FType] { return withPs(ps2, _r0) }))
	}))
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
	switch _v692 := (prm1).(type) {
	case Param_PUnit:
		zero := []Var{}
		return frt.NewTuple2(ps2, zero)
	case Param_PVar:
		v := _v692.Value
		tt := psCurrentTT(ps2)
		switch (tt).(type) {
		case TokenType_LPAREN:
			return frt.Pipe(parseParams(ps2), (func(_r0 frt.Tuple2[ParseState, []Var]) frt.Tuple2[ParseState, []Var] {
				return CnvR((func(_r0 []Var) []Var { return slice.PushHead(v, _r0) }), _r0)
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
			return CnvR((func(_r0 []NEPair) []NEPair { return slice.PushHead(nep, _r0) }), _r0)
		}))
	}))
}

func NEPToName(nvp NEPair) string {
	return nvp.name
}

func parseRecordGen(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2, neps := frt.Destr(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LBRACE, ps), (func(_r0 ParseState) frt.Tuple2[ParseState, []NEPair] { return parseFieldInitializers(parseE, _r0) })), (func(_r0 frt.Tuple2[ParseState, []NEPair]) frt.Tuple2[ParseState, []NEPair] {
		return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RBRACE, _r0) }), _r0)
	})))
	rtype, ok := frt.Destr(frt.Pipe(slice.Map(NEPToName, neps), (func(_r0 []string) frt.Tuple2[RecordType, bool] { return scLookupRecord(ps2.scope, _r0) })))
	return frt.IfElse(ok, (func() frt.Tuple2[ParseState, Expr] {
		return frt.Pipe(frt.Pipe(RecordGen{fieldsNV: neps, recordType: rtype}, New_Expr_ERecordGen), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(ps2, _r0) }))
	}), (func() frt.Tuple2[ParseState, Expr] {
		frt.Panic("record field name match to no record type.")
		return frt.NewTuple2(ps2, New_Expr_EUnit)
	}))
}

func refVar(vname string, ps ParseState) Expr {
	vfac, ok := frt.Destr(scLookupVarFac(ps.scope, vname))
	return frt.IfElse(ok, (func() Expr {
		return frt.Pipe(frt.Pipe(psTypeVarGen(ps), vfac), New_Expr_EVar)
	}), (func() Expr {
		frt.Panic("Unknown var ref")
		return New_Expr_EUnit
	}))
}

func parseFAAfterDot(ps ParseState, cur Expr) frt.Tuple2[ParseState, Expr] {
	ps2, fname := frt.Destr(frt.Pipe(psConsume(New_TokenType_DOT, ps), psIdentNameNx))
	switch _v695 := (ExprToType(cur)).(type) {
	case FType_FRecord:
		rtype := _v695.Value
		fexpr := frt.Pipe(FieldAccess{targetExpr: cur, targetType: rtype, fieldName: fname}, New_Expr_EFieldAccess)
		return frt.IfElse(psCurIs(New_TokenType_DOT, ps2), (func() frt.Tuple2[ParseState, Expr] {
			return parseFAAfterDot(ps2, fexpr)
		}), (func() frt.Tuple2[ParseState, Expr] {
			return frt.NewTuple2(ps2, fexpr)
		}))
	default:
		frt.Panic("non record type of field access")
		return frt.NewTuple2(ps2, New_Expr_EUnit)
	}
}

func parseVarRef(ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2, firstId := frt.Destr(psIdentNameNx(ps))
	return frt.IfElse(frt.OpNotEqual(psCurrentTT(ps2), New_TokenType_DOT), (func() frt.Tuple2[ParseState, Expr] {
		return frt.Pipe(refVar(firstId, ps2), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(ps2, _r0) }))
	}), (func() frt.Tuple2[ParseState, Expr] {
		vfac, ok := frt.Destr(scLookupVarFac(ps2.scope, firstId))
		return frt.IfElse(ok, (func() frt.Tuple2[ParseState, Expr] {
			return frt.Pipe(frt.Pipe(frt.Pipe(psTypeVarGen(ps2), vfac), New_Expr_EVar), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return parseFAAfterDot(ps2, _r0) }))
		}), (func() frt.Tuple2[ParseState, Expr] {
			ps3, fullName := frt.Destr(parseFullName(ps))
			return frt.Pipe(refVar(fullName, ps3), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(ps3, _r0) }))
		}))
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
	case TokenType_LPAREN:
		return frt.IfElse(frt.OpEqual(psCurrentTT(pn), New_TokenType_RPAREN), (func() frt.Tuple2[ParseState, Expr] {
			return frt.Pipe(frt.NewTuple2(pn, New_Expr_EUnit), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Expr] {
				return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RPAREN, _r0) }), _r0)
			}))
		}), (func() frt.Tuple2[ParseState, Expr] {
			ps2, e1 := frt.Destr(parseE(pn))
			return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_COMMA), (func() frt.Tuple2[ParseState, Expr] {
				ps3, e2 := frt.Destr(frt.Pipe(psConsume(New_TokenType_COMMA, ps2), parseE))
				frt.IfOnly(frt.OpEqual(psCurrentTT(ps3), New_TokenType_COMMA), (func() {
					frt.Panic("only pair is supported for tuple expr.")
				}))
				return frt.Pipe(frt.Pipe(([]Expr{e1, e2}), New_Expr_ETupleExpr), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(psConsume(New_TokenType_RPAREN, ps3), _r0) }))
			}), (func() frt.Tuple2[ParseState, Expr] {
				return frt.Pipe(frt.NewTuple2(ps2, e1), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Expr] {
					return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RPAREN, _r0) }), _r0)
				}))
			}))
		}))
	case TokenType_IDENTIFIER:
		return frt.IfElse(frt.OpEqual(cur.stringVal, "GoEval"), (func() frt.Tuple2[ParseState, Expr] {
			return parseGoEval(ps)
		}), (func() frt.Tuple2[ParseState, Expr] {
			return parseVarRef(ps)
		}))
	default:
		frt.Panic("Unown atom.")
		return frt.NewTuple2(ps, New_Expr_EUnit)
	}
}

func psCurIsBinOp(ps ParseState) bool {
	_, ok := frt.Destr(lookupBinOp(psCurrentTT(ps)))
	return ok
}

func psNextNonEOLIsBinOp(ps ParseState) bool {
	return frt.Pipe(psSkipEOL(ps), psCurIsBinOp)
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
		return psNextNonEOLIsBinOp(ps)
	}
}

func parseAtomList(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, []Expr] {
	ps2, one := frt.Destr(parseAtom(parseE, ps))
	return frt.IfElse(isEndOfTerm(ps2), (func() frt.Tuple2[ParseState, []Expr] {
		return frt.NewTuple2(ps2, ([]Expr{one}))
	}), (func() frt.Tuple2[ParseState, []Expr] {
		return frt.Pipe(parseAtomList(parseE, ps2), (func(_r0 frt.Tuple2[ParseState, []Expr]) frt.Tuple2[ParseState, []Expr] {
			return CnvR((func(_r0 []Expr) []Expr { return slice.PushHead(one, _r0) }), _r0)
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
			return CnvR((func(_r0 []MatchRule) []MatchRule { return slice.PushHead(one, _r0) }), _r0)
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
			return CnvR((func(_r0 []Expr) []Expr { return slice.PushHead(one, _r0) }), _r0)
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

func exprOnlyBlock(expr Expr) Block {
	emp := []Stmt{}
	return Block{stmts: emp, finalExpr: expr}
}

func parseInlineBlock(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Block] {
	ps2, expr := frt.Destr(pExpr(ps))
	return frt.Pipe(exprOnlyBlock(expr), (func(_r0 Block) frt.Tuple2[ParseState, Block] { return withPs(ps2, _r0) }))
}

func parseIfAfterIfExpr(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2, cond := frt.Destr(frt.Pipe(pExpr(ps), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Expr] {
		return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_THEN, _r0) }), _r0)
	})))
	tgen := psTypeVarGen(ps2)
	recurse := (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] { return parseIfAfterIfExpr(pExpr, pBlock, _r0) })
	return frt.IfElse(psCurIs(New_TokenType_EOL, ps2), (func() frt.Tuple2[ParseState, Expr] {
		ps3, tbody := frt.Destr(frt.Pipe(psSkipEOL(ps2), pBlock))
		ps4 := psSkipEOL(ps3)
		return frt.IfElse(psCurIs(New_TokenType_ELSE, ps4), (func() frt.Tuple2[ParseState, Expr] {
			pse2, fbody := frt.Destr(frt.Pipe(frt.Pipe(psConsume(New_TokenType_ELSE, ps4), psSkipEOL), pBlock))
			return frt.Pipe(newIfElseCall(tgen, cond, tbody, fbody), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(pse2, _r0) }))
		}), (func() frt.Tuple2[ParseState, Expr] {
			return frt.IfElse(psCurIs(New_TokenType_ELIF, ps4), (func() frt.Tuple2[ParseState, Expr] {
				ps5, elseExpr := frt.Destr(frt.Pipe(psConsume(New_TokenType_ELIF, ps4), recurse))
				ebody := exprOnlyBlock(elseExpr)
				return frt.Pipe(newIfElseCall(tgen, cond, tbody, ebody), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(ps5, _r0) }))
			}), (func() frt.Tuple2[ParseState, Expr] {
				return frt.Pipe(newIfOnlyCall(tgen, cond, tbody), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(ps3, _r0) }))
			}))
		}))
	}), (func() frt.Tuple2[ParseState, Expr] {
		psi2, tbody := frt.Destr(parseInlineBlock(pExpr, ps2))
		return frt.IfElse(psCurIs(New_TokenType_ELSE, psi2), (func() frt.Tuple2[ParseState, Expr] {
			psi3, fbody := frt.Destr(parseInlineBlock(pExpr, psi2))
			return frt.Pipe(newIfElseCall(tgen, cond, tbody, fbody), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(psi3, _r0) }))
		}), (func() frt.Tuple2[ParseState, Expr] {
			return frt.Pipe(newIfOnlyCall(tgen, cond, tbody), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(psi2, _r0) }))
		}))
	}))
}

func parseIfExpr(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, Expr] {
	return frt.Pipe(psConsume(New_TokenType_IF, ps), (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] { return parseIfAfterIfExpr(pExpr, pBlock, _r0) }))
}

func parseTerm(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, Expr] {
	switch (psCurrentTT(ps)).(type) {
	case TokenType_MATCH:
		return frt.Pipe(frt.Pipe(parseMatchExpr(pExpr, pBlock, ps), (func(_r0 frt.Tuple2[ParseState, MatchExpr]) frt.Tuple2[ParseState, ReturnableExpr] {
			return CnvR(New_ReturnableExpr_RMatchExpr, _r0)
		})), (func(_r0 frt.Tuple2[ParseState, ReturnableExpr]) frt.Tuple2[ParseState, Expr] {
			return CnvR(New_Expr_EReturnableExpr, _r0)
		}))
	case TokenType_LSBRACKET:
		return parseSliceExpr(pExpr, ps)
	case TokenType_IF:
		return parseIfExpr(pExpr, pBlock, ps)
	case TokenType_NOT:
		ps2, target := frt.Destr(frt.Pipe(psConsume(New_TokenType_NOT, ps), (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] { return parseTerm(pExpr, pBlock, _r0) })))
		return frt.Pipe(newUnaryNotCall(psTypeVarGen(ps2), target), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return withPs(ps2, _r0) }))
	default:
		ps2, es := frt.Destr(parseAtomList(pExpr, ps))
		return frt.IfElse(frt.OpEqual(slice.Length(es), 1), (func() frt.Tuple2[ParseState, Expr] {
			return frt.NewTuple2(ps2, slice.Head(es))
		}), (func() frt.Tuple2[ParseState, Expr] {
			head := slice.Head(es)
			tail := slice.Tail(es)
			switch _v701 := (head).(type) {
			case Expr_EVar:
				v := _v701.Value
				fc := FunCall{targetFunc: v, args: tail}
				return frt.NewTuple2(ps2, New_Expr_EFunCall(fc))
			default:
				frt.Panic("Funcall head is not var")
				return frt.NewTuple2(ps2, head)
			}
		}))
	}
}

func lookupBinOpNF(tk TokenType) BinOpInfo {
	res := lookupBinOp(tk)
	return frt.Fst(res)
}

func parseBinAfter(pEwithMinPrec func(int, ParseState) frt.Tuple2[ParseState, Expr], minPrec int, ps ParseState, cur Expr) frt.Tuple2[ParseState, Expr] {
	ps2 := psSkipEOL(ps)
	return frt.IfElse(psCurIsBinOp(ps2), (func() frt.Tuple2[ParseState, Expr] {
		btk := psCurrentTT(ps2)
		bop := lookupBinOpNF(btk)
		return frt.IfElse((bop.precedence < minPrec), (func() frt.Tuple2[ParseState, Expr] {
			return frt.NewTuple2(ps, cur)
		}), (func() frt.Tuple2[ParseState, Expr] {
			ps3, rhs := frt.Destr(frt.Pipe(psConsume(btk, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] { return pEwithMinPrec((bop.precedence + 1), _r0) })))
			tvgen := psTypeVarGen(ps3)
			return frt.Pipe(newBinOpCall(tvgen, btk, bop, cur, rhs), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return parseBinAfter(pEwithMinPrec, minPrec, ps3, _r0) }))
		}))
	}), (func() frt.Tuple2[ParseState, Expr] {
		return frt.NewTuple2(ps, cur)
	}))
}

func parseExprWithPrec(pBlock func(ParseState) frt.Tuple2[ParseState, Block], minPrec int, ps ParseState) frt.Tuple2[ParseState, Expr] {
	pExpr := (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] { return parseExprWithPrec(pBlock, 1, _r0) })
	ps2, expr := frt.Destr(parseTerm(pExpr, pBlock, ps))
	ps3 := psSkipEOL(ps2)
	return frt.IfElse(psCurIsBinOp(ps3), (func() frt.Tuple2[ParseState, Expr] {
		return parseBinAfter((func(_r0 int, _r1 ParseState) frt.Tuple2[ParseState, Expr] { return parseExprWithPrec(pBlock, _r0, _r1) }), minPrec, ps3, expr)
	}), (func() frt.Tuple2[ParseState, Expr] {
		return frt.NewTuple2(ps2, expr)
	}))
}

func parseExpr(pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, Expr] {
	return parseExprWithPrec(pBlock, 1, ps)
}

func parseStmt(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pLet func(ParseState) frt.Tuple2[ParseState, LetVarDef], ps ParseState) frt.Tuple2[ParseState, Stmt] {
	switch (psCurrentTT(ps)).(type) {
	case TokenType_LET:
		return frt.Pipe(pLet(ps), (func(_r0 frt.Tuple2[ParseState, LetVarDef]) frt.Tuple2[ParseState, Stmt] {
			return CnvR(New_Stmt_SLetVarDef, _r0)
		}))
	default:
		return frt.Pipe(pExpr(ps), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Stmt] {
			return CnvR(New_Stmt_SExprStmt, _r0)
		}))
	}
}

func isEndOfBlock(ps ParseState) bool {
	isOffside := (psCurCol(ps) < psCurOffside(ps))
	isEof := frt.OpEqual(psCurrentTT(ps), New_TokenType_EOF)
	return (isOffside || isEof)
}

func parseStmtList(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pLet func(ParseState) frt.Tuple2[ParseState, LetVarDef], ps ParseState) frt.Tuple2[ParseState, []Stmt] {
	ps2, one := frt.Destr(frt.Pipe(parseStmt(pExpr, pLet, ps), (func(_r0 frt.Tuple2[ParseState, Stmt]) frt.Tuple2[ParseState, Stmt] { return CnvL(psSkipEOL, _r0) })))
	return frt.IfElse(isEndOfBlock(ps2), (func() frt.Tuple2[ParseState, []Stmt] {
		return frt.NewTuple2(ps2, ([]Stmt{one}))
	}), (func() frt.Tuple2[ParseState, []Stmt] {
		return frt.Pipe(parseStmtList(pExpr, pLet, ps2), (func(_r0 frt.Tuple2[ParseState, []Stmt]) frt.Tuple2[ParseState, []Stmt] {
			return CnvR((func(_r0 []Stmt) []Stmt { return slice.PushHead(one, _r0) }), _r0)
		}))
	}))
}

func parseBlockAfterPushScope(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pLet func(ParseState) frt.Tuple2[ParseState, LetVarDef], ps ParseState) frt.Tuple2[ParseState, Block] {
	ps2, sls := frt.Destr(frt.Pipe(frt.Pipe(psPushOffside(ps), (func(_r0 ParseState) frt.Tuple2[ParseState, []Stmt] { return parseStmtList(pExpr, pLet, _r0) })), (func(_r0 frt.Tuple2[ParseState, []Stmt]) frt.Tuple2[ParseState, []Stmt] {
		return CnvL(psPopOffside, _r0)
	})))
	last := slice.Last(sls)
	stmts := slice.PopLast(sls)
	switch _v703 := (last).(type) {
	case Stmt_SExprStmt:
		e := _v703.Value
		return frt.NewTuple2(ps2, Block{stmts: stmts, finalExpr: e})
	default:
		frt.Panic("block of last is not expr")
		return frt.Pipe(emptyBlock(), (func(_r0 Block) frt.Tuple2[ParseState, Block] { return withPs(ps2, _r0) }))
	}
}

func parseBlock(pLet func(ParseState) frt.Tuple2[ParseState, LetVarDef], ps ParseState) frt.Tuple2[ParseState, Block] {
	pExpr := (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] {
		return parseExpr((func(_r0 ParseState) frt.Tuple2[ParseState, Block] { return parseBlock(pLet, _r0) }), _r0)
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
	return frt.Pipe(LetVarDef{lvar: v, rhs: rhs}, (func(_r0 LetVarDef) frt.Tuple2[ParseState, LetVarDef] { return withPs(ps3, _r0) }))
}

func vToT(v Var) FType {
	return v.ftype
}

func lfdToFuncVar(lfd LetFuncDef) Var {
	return lfd.fvar
}

func parseLetFuncDef(pLet func(ParseState) frt.Tuple2[ParseState, LetVarDef], ps ParseState) frt.Tuple2[ParseState, LetFuncDef] {
	ps2 := psConsume(New_TokenType_LET, ps)
	fname := psIdentName(ps2)
	ps3, params := frt.Destr(frt.Pipe(psNext(ps2), parseParams))
	ps4, block := frt.Destr(frt.Pipe(frt.Pipe(psConsume(New_TokenType_EQ, ps3), psSkipEOL), (func(_r0 ParseState) frt.Tuple2[ParseState, Block] { return parseBlock(pLet, _r0) })))
	rtype := frt.Pipe(blockToExpr(block), ExprToType)
	targets := frt.IfElse(frt.OpEqual(slice.Length(params), 0), (func() []FType {
		return ([]FType{New_FType_FUnit, rtype})
	}), (func() []FType {
		return frt.Pipe(slice.Map(vToT, params), (func(_r0 []FType) []FType { return slice.PushLast(rtype, _r0) }))
	}))
	ft := frt.Pipe(FuncType{targets: targets}, New_FType_FFunc)
	fnvar := Var{name: fname, ftype: ft}
	return frt.Pipe(LetFuncDef{fvar: fnvar, params: params, body: block}, (func(_r0 LetFuncDef) frt.Tuple2[ParseState, LetFuncDef] { return withPs(ps4, _r0) }))
}

func tvToN(tv TypeVar) string {
	return tv.name
}

func parseRootLetFuncDef(pLet func(ParseState) frt.Tuple2[ParseState, LetVarDef], ps ParseState) frt.Tuple2[ParseState, RootFuncDef] {
	ps2, lfd := frt.Destr(parseLetFuncDef(pLet, ps))
	frt.PipeUnit(lfdToFuncVar(lfd), (func(_r0 Var) { scDefVar(ps.scope, lfd.fvar.name, _r0) }))
	tnames := frt.Pipe(tvaListAlloced(ps.tva), (func(_r0 []TypeVar) []string { return slice.Map(tvToN, _r0) }))
	rfd := Infer(tnames, lfd)
	return frt.NewTuple2(ps2, rfd)
}

func parseRootLet(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps0 ParseState) frt.Tuple2[ParseState, RootStmt] {
	ps := psResetTmpCtx(ps0)
	pLet := (func(_r0 ParseState) frt.Tuple2[ParseState, LetVarDef] { return parseLetVarDef(pExpr, _r0) })
	psN := psNext(ps)
	switch (psCurrentTT(psN)).(type) {
	case TokenType_LPAREN:
		return frt.Pipe(parseRootLetFuncDef(pLet, ps), (func(_r0 frt.Tuple2[ParseState, RootFuncDef]) frt.Tuple2[ParseState, RootStmt] {
			return CnvR(New_RootStmt_RSRootFuncDef, _r0)
		}))
	default:
		psNN := psNext(psN)
		switch (psCurrentTT(psNN)).(type) {
		case TokenType_EQ:
			frt.Panic("Root let var def, NYI")
			return frt.NewTuple2(ps, New_RootStmt_RSImport("dummy"))
		default:
			return frt.Pipe(parseRootLetFuncDef(pLet, ps), (func(_r0 frt.Tuple2[ParseState, RootFuncDef]) frt.Tuple2[ParseState, RootStmt] {
				return CnvR(New_RootStmt_RSRootFuncDef, _r0)
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
	ps2, ntp := frt.Destr(frt.Pipe(psSkipEOL(ps), parseFieldDef))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_RBRACE), (func() frt.Tuple2[ParseState, []NameTypePair] {
		return frt.NewTuple2(ps2, ([]NameTypePair{ntp}))
	}), (func() frt.Tuple2[ParseState, []NameTypePair] {
		ps3 := frt.Pipe(psConsume(New_TokenType_SEMICOLON, ps2), psSkipEOL)
		return frt.IfElse(psCurIs(New_TokenType_RBRACE, ps3), (func() frt.Tuple2[ParseState, []NameTypePair] {
			return frt.NewTuple2(ps3, ([]NameTypePair{ntp}))
		}), (func() frt.Tuple2[ParseState, []NameTypePair] {
			return frt.Pipe(parseFieldDefs(ps3), (func(_r0 frt.Tuple2[ParseState, []NameTypePair]) frt.Tuple2[ParseState, []NameTypePair] {
				return CnvR((func(_r0 []NameTypePair) []NameTypePair { return slice.PushHead(ntp, _r0) }), _r0)
			}))
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
			return CnvR((func(_r0 []NameTypePair) []NameTypePair { return slice.PushHead(cs, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []NameTypePair] {
		return frt.NewTuple2(ps2, ([]NameTypePair{cs}))
	}))
}

func parseUnionDef(tname string, ps ParseState) frt.Tuple2[ParseState, UnionDef] {
	ps2, css := frt.Destr(parseCaseDefs(ps))
	ud := UnionDef{name: tname, cases: css}
	udRegisterToScope(ps2.scope, ud)
	return frt.NewTuple2(ps2, ud)
}

func parseTypeDef(ps ParseState) frt.Tuple2[ParseState, RootStmt] {
	ps2, tname := frt.Destr(frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_TYPE, ps), psIdentNameNxL), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] {
		return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_EQ, _r0) }), _r0)
	})), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return CnvL(psSkipEOL, _r0) })))
	switch (psCurrentTT(ps2)).(type) {
	case TokenType_LBRACE:
		return frt.Pipe(frt.Pipe(parseRecordDef(tname, ps2), (func(_r0 frt.Tuple2[ParseState, RecordDef]) frt.Tuple2[ParseState, DefStmt] {
			return CnvR(New_DefStmt_DRecordDef, _r0)
		})), (func(_r0 frt.Tuple2[ParseState, DefStmt]) frt.Tuple2[ParseState, RootStmt] {
			return CnvR(New_RootStmt_RSDefStmt, _r0)
		}))
	case TokenType_BAR:
		return frt.Pipe(frt.Pipe(parseUnionDef(tname, ps2), (func(_r0 frt.Tuple2[ParseState, UnionDef]) frt.Tuple2[ParseState, DefStmt] {
			return CnvR(New_DefStmt_DUnionDef, _r0)
		})), (func(_r0 frt.Tuple2[ParseState, DefStmt]) frt.Tuple2[ParseState, RootStmt] {
			return CnvR(New_RootStmt_RSDefStmt, _r0)
		}))
	default:
		frt.Panic("NYI")
		return frt.NewTuple2(ps2, New_RootStmt_RSImport("dummy"))
	}
}

func parseExtTypeDef(pi PackageInfo, ps ParseState) ParseState {
	ps2, tname := frt.Destr(frt.Pipe(psConsume(New_TokenType_TYPE, ps), psIdentNameNx))
	etype := piRegEType(pi, tname)
	scRegisterType(ps2.scope, tname, etype)
	return ps2
}

func parseTypeParams(ps ParseState) frt.Tuple2[ParseState, []string] {
	ps2, tname := frt.Destr(psIdentNameNx(ps))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_COMMA), (func() frt.Tuple2[ParseState, []string] {
		return frt.Pipe(frt.Pipe(psConsume(New_TokenType_COMMA, ps2), parseTypeParams), (func(_r0 frt.Tuple2[ParseState, []string]) frt.Tuple2[ParseState, []string] {
			return CnvR((func(_r0 []string) []string { return slice.PushHead(tname, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []string] {
		return frt.NewTuple2(ps2, ([]string{tname}))
	}))
}

func regTypeVar(ps ParseState, tname string) {
	frt.PipeUnit(New_FType_FTypeVar(TypeVar{name: tname}), (func(_r0 FType) { scRegisterType(ps.scope, tname, _r0) }))
}

func psRegTypeVars(ps ParseState, tnames []string) {
	slice.Iter((func(_r0 string) { regTypeVar(ps, _r0) }), tnames)
}

func parseExtFuncDef(pi PackageInfo, ps ParseState) ParseState {
	ps2, fname := frt.Destr(frt.Pipe(psConsume(New_TokenType_LET, ps), psIdentNameNx))
	ps3, tnames := frt.Destr(frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_LT), (func() frt.Tuple2[ParseState, []string] {
		return frt.Pipe(frt.Pipe(psConsume(New_TokenType_LT, ps2), parseTypeParams), (func(_r0 frt.Tuple2[ParseState, []string]) frt.Tuple2[ParseState, []string] {
			return CnvL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_GT, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []string] {
		return frt.Pipe([]string{}, (func(_r0 []string) frt.Tuple2[ParseState, []string] { return withPs(ps2, _r0) }))
	})))
	psRegTypeVars(ps3, tnames)
	ps4, fts := frt.Destr(frt.Pipe(psConsume(New_TokenType_COLON, ps3), (func(_r0 ParseState) frt.Tuple2[ParseState, []FType] { return parseTypeArrows(parseType, _r0) })))
	ff := FuncFactory{tparams: tnames, targets: fts}
	piRegFF(pi, fname, ff, ps4)
	return ps4
}

func parseExtDef(pi PackageInfo, ps ParseState) ParseState {
	switch (psCurrentTT(ps)).(type) {
	case TokenType_LET:
		return parseExtFuncDef(pi, ps)
	case TokenType_TYPE:
		return parseExtTypeDef(pi, ps)
	default:
		frt.Panic("Unknown pkginfo def")
		return ps
	}
}

func parseExtDefs(pi PackageInfo, ps ParseState) ParseState {
	ps2 := frt.Pipe(parseExtDef(pi, ps), psSkipEOL)
	return frt.IfElse(isEndOfBlock(ps2), (func() ParseState {
		return ps2
	}), (func() ParseState {
		return parseExtDefs(pi, ps2)
	}))
}

func parsePackageInfo(ps ParseState) frt.Tuple2[ParseState, RootStmt] {
	ps2 := psConsume(New_TokenType_PACKAGE_INFO, ps)
	ps3, pkgName := frt.Destr(frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_UNDER_SCORE), (func() frt.Tuple2[ParseState, string] {
		return frt.Pipe(frt.NewTuple2(ps2, "_"), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return CnvL(psNext, _r0) }))
	}), (func() frt.Tuple2[ParseState, string] {
		return psIdentNameNx(ps2)
	})))
	ps4 := frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_EQ, ps3), psSkipEOL), psPushOffside), psPushScope)
	pi := NewPackageInfo(pkgName)
	ps5 := frt.Pipe(frt.Pipe(parseExtDefs(pi, ps4), psPopScope), psPopOffside)
	piRegAll(pi, ps5.scope)
	return frt.NewTuple2(ps5, New_RootStmt_RSPackageInfo(pi))
}

func parseRootOneStmt(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, RootStmt] {
	switch (psCurrentTT(ps)).(type) {
	case TokenType_PACKAGE:
		return parsePackage(ps)
	case TokenType_IMPORT:
		return parseImport(ps)
	case TokenType_LET:
		return parseRootLet(pExpr, ps)
	case TokenType_TYPE:
		return parseTypeDef(ps)
	case TokenType_PACKAGE_INFO:
		return parsePackageInfo(ps)
	default:
		frt.Panic("Unknown stmt")
		return parsePackage(ps)
	}
}

func parseRootStmts(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, []RootStmt] {
	ps2 := psSkipEOL(ps)
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_EOF), (func() frt.Tuple2[ParseState, []RootStmt] {
		s := []RootStmt{}
		return frt.NewTuple2(ps2, s)
	}), (func() frt.Tuple2[ParseState, []RootStmt] {
		ps3, one := frt.Destr(frt.Pipe(parseRootOneStmt(pExpr, ps2), (func(_r0 frt.Tuple2[ParseState, RootStmt]) frt.Tuple2[ParseState, RootStmt] {
			return CnvL(psSkipEOL, _r0)
		})))
		ps4, rest := frt.Destr(parseRootStmts(pExpr, ps3))
		ss := slice.PushHead(one, rest)
		return frt.NewTuple2(ps4, ss)
	}))
}
