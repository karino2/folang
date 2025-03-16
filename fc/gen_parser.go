package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/dict"

func parsePackage(ps ParseState) frt.Tuple2[ParseState, RootStmt] {
	return frt.Pipe(frt.Pipe(psConsume(New_TokenType_PACKAGE, ps), psIdentNameNxL), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, RootStmt] {
		return MapR(New_RootStmt_RSPackage, _r0)
	}))
}

func parseImport(ps ParseState) frt.Tuple2[ParseState, RootStmt] {
	ps2 := psConsume(New_TokenType_IMPORT, ps)
	return frt.IfElse(psCurIs(New_TokenType_IDENTIFIER, ps2), (func() frt.Tuple2[ParseState, RootStmt] {
		ps3, iname := frt.Destr2(psIdentNameNxL(ps2))
		rstmt := frt.Pipe(frt.Sprintf1("github.com/karino2/folang/pkg/%s", iname), New_RootStmt_RSImport)
		return frt.NewTuple2(ps3, rstmt)
	}), (func() frt.Tuple2[ParseState, RootStmt] {
		return frt.Pipe(frt.Pipe(ps2, psStringValNxL), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, RootStmt] {
			return MapR(New_RootStmt_RSImport, _r0)
		}))
	}))
}

func parseFullName(ps ParseState) frt.Tuple2[ParseState, string] {
	ps2, one := frt.Destr2(psIdentNameNx(ps))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_DOT), (func() frt.Tuple2[ParseState, string] {
		ps3, rest := frt.Destr2(frt.Pipe(psConsume(New_TokenType_DOT, ps2), parseFullName))
		return frt.NewTuple2(ps3, ((one + ".") + rest))
	}), (func() frt.Tuple2[ParseState, string] {
		return frt.NewTuple2(ps2, one)
	}))
}

func parseTypeList(pType func(ParseState) frt.Tuple2[ParseState, FType], ps ParseState) frt.Tuple2[ParseState, []FType] {
	ps2, one := frt.Destr2(pType(ps))
	return frt.IfElse(psCurIs(New_TokenType_COMMA, ps2), (func() frt.Tuple2[ParseState, []FType] {
		return frt.Pipe(frt.Pipe(psConsume(New_TokenType_COMMA, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, []FType] { return parseTypeList(pType, _r0) })), (func(_r0 frt.Tuple2[ParseState, []FType]) frt.Tuple2[ParseState, []FType] {
			return MapR((func(_r0 []FType) []FType { return slice.PushHead(one, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []FType] {
		return frt.NewTuple2(ps2, ([]FType{one}))
	}))
}

func mightParseSpecifiedTypeList(pType func(ParseState) frt.Tuple2[ParseState, FType], ps ParseState) frt.Tuple2[ParseState, []FType] {
	return frt.IfElse(psCurIs(New_TokenType_LT, ps), (func() frt.Tuple2[ParseState, []FType] {
		return frt.Pipe(frt.Pipe(psConsume(New_TokenType_LT, ps), (func(_r0 ParseState) frt.Tuple2[ParseState, []FType] { return parseTypeList(pType, _r0) })), (func(_r0 frt.Tuple2[ParseState, []FType]) frt.Tuple2[ParseState, []FType] {
			return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_GT, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []FType] {
		return frt.Pipe(emptyFtps(), (func(_r0 []FType) frt.Tuple2[ParseState, []FType] { return PairL(ps, _r0) }))
	}))
}

func parseAtomType(pType func(ParseState) frt.Tuple2[ParseState, FType], ps ParseState) frt.Tuple2[ParseState, FType] {
	tk := psCurrent(ps)
	switch (tk.ttype).(type) {
	case TokenType_LPAREN:
		ps2 := psConsume(New_TokenType_LPAREN, ps)
		return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_RPAREN), (func() frt.Tuple2[ParseState, FType] {
			return frt.Pipe(psConsume(New_TokenType_RPAREN, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, FType] { return PairR(New_FType_FUnit, _r0) }))
		}), (func() frt.Tuple2[ParseState, FType] {
			return frt.Pipe(pType(ps2), (func(_r0 frt.Tuple2[ParseState, FType]) frt.Tuple2[ParseState, FType] {
				return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RPAREN, _r0) }), _r0)
			}))
		}))
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
					return frt.IfElse(frt.OpEqual(tname, "float"), (func() frt.Tuple2[ParseState, FType] {
						return frt.NewTuple2(ps3, New_FType_FFloat)
					}), (func() frt.Tuple2[ParseState, FType] {
						return frt.IfElse(frt.OpEqual(tname, "any"), (func() frt.Tuple2[ParseState, FType] {
							return frt.NewTuple2(ps3, New_FType_FAny)
						}), (func() frt.Tuple2[ParseState, FType] {
							ps4, fullName := frt.Destr2(parseFullName(ps))
							tfac, ok := frt.Destr2(scLookupTypeFac(ps3.scope, fullName))
							return frt.IfElse(ok, (func() frt.Tuple2[ParseState, FType] {
								return frt.Pipe(mightParseSpecifiedTypeList(pType, ps4), (func(_r0 frt.Tuple2[ParseState, []FType]) frt.Tuple2[ParseState, FType] { return MapR(tfac, _r0) }))
							}), (func() frt.Tuple2[ParseState, FType] {
								return frt.IfElse(psInsideTypeDef(ps4), (func() frt.Tuple2[ParseState, FType] {
									tvarf := tdctxTVFAlloc(ps4.tdctx, fullName)
									return frt.NewTuple2(ps4, tvarf)
								}), (func() frt.Tuple2[ParseState, FType] {
									frt.Panicf1("type not found: %s.", fullName)
									return frt.NewTuple2(ps4, New_FType_FUnit)
								}))
							}))
						}))
					}))
				}))
			}))
		}))
	default:
		psPanic(ps, "Unknown type")
		return frt.NewTuple2(ps, New_FType_FUnit)
	}
}

func parseTermType(pType func(ParseState) frt.Tuple2[ParseState, FType], ps ParseState) frt.Tuple2[ParseState, FType] {
	recurse := (func(_r0 ParseState) frt.Tuple2[ParseState, FType] { return parseTermType(pType, _r0) })
	return frt.IfElse(psCurIs(New_TokenType_LSBRACKET, ps), (func() frt.Tuple2[ParseState, FType] {
		ps2, et := frt.Destr2(frt.Pipe(psMulConsume(([]TokenType{New_TokenType_LSBRACKET, New_TokenType_RSBRACKET}), ps), recurse))
		return frt.Pipe(frt.Pipe(SliceType{ElemType: et}, New_FType_FSlice), (func(_r0 FType) frt.Tuple2[ParseState, FType] { return PairL(ps2, _r0) }))
	}), (func() frt.Tuple2[ParseState, FType] {
		return parseAtomType(pType, ps)
	}))
}

func parseElemType(pType func(ParseState) frt.Tuple2[ParseState, FType], ps ParseState) frt.Tuple2[ParseState, FType] {
	pTerm := (func(_r0 ParseState) frt.Tuple2[ParseState, FType] { return parseTermType(pType, _r0) })
	ps2, fts := frt.Destr2(ParseSepList(pTerm, New_TokenType_ASTER, ps))
	return frt.IfElse(frt.OpEqual(slice.Length(fts), 1), (func() frt.Tuple2[ParseState, FType] {
		return frt.NewTuple2(ps2, slice.Head(fts))
	}), (func() frt.Tuple2[ParseState, FType] {
		return frt.Pipe(frt.Pipe(TupleType{ElemTypes: fts}, New_FType_FTuple), (func(_r0 FType) frt.Tuple2[ParseState, FType] { return PairL(ps2, _r0) }))
	}))
}

func parseTypeArrows(pType func(ParseState) frt.Tuple2[ParseState, FType], ps ParseState) frt.Tuple2[ParseState, []FType] {
	ps2, one := frt.Destr2(parseElemType(pType, ps))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_RARROW), (func() frt.Tuple2[ParseState, []FType] {
		return frt.Pipe(frt.Pipe(psConsume(New_TokenType_RARROW, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, []FType] { return parseTypeArrows(pType, _r0) })), (func(_r0 frt.Tuple2[ParseState, []FType]) frt.Tuple2[ParseState, []FType] {
			return MapR((func(_r0 []FType) []FType { return slice.PushHead(one, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []FType] {
		return frt.NewTuple2(ps2, ([]FType{one}))
	}))
}

func parseType(ps ParseState) frt.Tuple2[ParseState, FType] {
	ps2, tps := frt.Destr2(parseTypeArrows(parseType, ps))
	return frt.IfElse(frt.OpEqual(slice.Length(tps), 1), (func() frt.Tuple2[ParseState, FType] {
		return frt.Pipe(slice.Head(tps), (func(_r0 FType) frt.Tuple2[ParseState, FType] { return PairL(ps2, _r0) }))
	}), (func() frt.Tuple2[ParseState, FType] {
		return frt.Pipe(newFFunc(tps), (func(_r0 FType) frt.Tuple2[ParseState, FType] { return PairL(ps2, _r0) }))
	}))
}

type Param interface {
	Param_Union()
}

func (Param_PVar) Param_Union()  {}
func (Param_PUnit) Param_Union() {}

func (v Param_PVar) String() string  { return frt.Sprintf1("(PVar: %v)", v.Value) }
func (v Param_PUnit) String() string { return "(PUnit)" }

type Param_PVar struct {
	Value Var
}

func New_Param_PVar(v Var) Param { return Param_PVar{v} }

type Param_PUnit struct {
}

var New_Param_PUnit Param = Param_PUnit{}

func psNewTypeVar(ps ParseState) TypeVar {
	tgen := psTypeVarGen(ps)
	return tgen()
}

func parseParam(ps ParseState) frt.Tuple2[ParseState, Param] {
	switch (psCurrentTT(ps)).(type) {
	case TokenType_LPAREN:
		ps2 := psConsume(New_TokenType_LPAREN, ps)
		tk := psCurrent(ps2)
		switch (tk.ttype).(type) {
		case TokenType_RPAREN:
			return frt.Pipe(psConsume(New_TokenType_RPAREN, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, Param] { return PairR(New_Param_PUnit, _r0) }))
		default:
			vname := psIdentName(ps2)
			ps3, tp := frt.Destr2(frt.Pipe(frt.Pipe(frt.Pipe(psNext(ps2), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_COLON, _r0) })), parseType), (func(_r0 frt.Tuple2[ParseState, FType]) frt.Tuple2[ParseState, FType] {
				return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RPAREN, _r0) }), _r0)
			})))
			v := Var{Name: vname, Ftype: tp}
			scDefVar(ps3.scope, vname, v)
			return frt.NewTuple2(ps3, New_Param_PVar(v))
		}
	case TokenType_IDENTIFIER:
		ps2, vname := frt.Destr2(psIdentNameNx(ps))
		ftv := frt.Pipe(psNewTypeVar(ps2), New_FType_FTypeVar)
		v := Var{Name: vname, Ftype: ftv}
		scDefVar(ps2.scope, vname, v)
		return frt.NewTuple2(ps2, New_Param_PVar(v))
	default:
		psPanic(ps, "Unexpected token in parameter")
		return frt.NewTuple2(ps, New_Param_PUnit)
	}
}

func parseParams(ps ParseState) frt.Tuple2[ParseState, []Var] {
	ps2, prm1 := frt.Destr2(parseParam(ps))
	switch _v1 := (prm1).(type) {
	case Param_PUnit:
		zero := []Var{}
		return frt.NewTuple2(ps2, zero)
	case Param_PVar:
		v := _v1.Value
		tt := psCurrentTT(ps2)
		switch (tt).(type) {
		case TokenType_LPAREN:
			return frt.Pipe(parseParams(ps2), (func(_r0 frt.Tuple2[ParseState, []Var]) frt.Tuple2[ParseState, []Var] {
				return MapR((func(_r0 []Var) []Var { return slice.PushHead(v, _r0) }), _r0)
			}))
		case TokenType_IDENTIFIER:
			return frt.Pipe(parseParams(ps2), (func(_r0 frt.Tuple2[ParseState, []Var]) frt.Tuple2[ParseState, []Var] {
				return MapR((func(_r0 []Var) []Var { return slice.PushHead(v, _r0) }), _r0)
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
		ps3, ft := frt.Destr2(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LT, ps2), parseType), (func(_r0 frt.Tuple2[ParseState, FType]) frt.Tuple2[ParseState, FType] {
			return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_GT, _r0) }), _r0)
		})))
		ps4, s := frt.Destr2(psStringValNx(ps3))
		ge := GoEvalExpr{GoStmt: s, TypeArg: ft}
		return frt.NewTuple2(ps4, New_Expr_EGoEvalExpr(ge))
	case TokenType_STRING:
		ps3, s := frt.Destr2(psStringValNx(ps2))
		ge := GoEvalExpr{GoStmt: s, TypeArg: New_FType_FUnit}
		return frt.NewTuple2(ps3, New_Expr_EGoEvalExpr(ge))
	default:
		psPanic(ps2, "Wrong arg for GoEval")
		return frt.NewTuple2(ps2, New_Expr_EUnit)
	}
}

type fiInfo struct {
	RecName string
	NePair  NEPair
}

func parseFiIni(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, fiInfo] {
	ps2, fname := frt.Destr2(psIdentNameNxL(ps))
	return frt.IfElse(psCurIs(New_TokenType_DOT, ps2), (func() frt.Tuple2[ParseState, fiInfo] {
		ps3, fname2 := frt.Destr2(frt.Pipe(psConsume(New_TokenType_DOT, ps2), psIdentNameNxL))
		ps4, expr := frt.Destr2(frt.Pipe(frt.Pipe(psConsume(New_TokenType_EQ, ps3), psSkipEOL), parseE))
		fi := NEPair{Name: fname2, Expr: expr}
		return frt.NewTuple2(ps4, fiInfo{RecName: fname, NePair: fi})
	}), (func() frt.Tuple2[ParseState, fiInfo] {
		ps3, expr := frt.Destr2(frt.Pipe(frt.Pipe(frt.Pipe(ps2, (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_EQ, _r0) })), psSkipEOL), parseE))
		fi := NEPair{Name: fname, Expr: expr}
		return frt.NewTuple2(ps3, fiInfo{RecName: "", NePair: fi})
	}))
}

type fiListInfo struct {
	RecName string
	NePairs []NEPair
}

func parseFieldInitializers(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, fiListInfo] {
	ps2, fii := frt.Destr2(parseFiIni(parseE, ps))
	nep := fii.NePair
	recN := fii.RecName
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_RBRACE), (func() frt.Tuple2[ParseState, fiListInfo] {
		return frt.NewTuple2(ps2, fiListInfo{RecName: recN, NePairs: ([]NEPair{nep})})
	}), (func() frt.Tuple2[ParseState, fiListInfo] {
		ps3, fiInfos := frt.Destr2(frt.Pipe(psConsume(New_TokenType_SEMICOLON, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, fiListInfo] { return parseFieldInitializers(parseE, _r0) })))
		recN2 := fiInfos.RecName
		neps := fiInfos.NePairs
		neps2 := slice.PushHead(nep, neps)
		recN3 := frt.IfElse(frt.OpNotEqual(recN, ""), (func() string {
			return recN
		}), (func() string {
			return recN2
		}))
		return frt.NewTuple2(ps3, fiListInfo{RecName: recN3, NePairs: neps2})
	}))
}

func retRecordGen(ok bool, rfac RecordFactory, neps []NEPair, ps ParseState) frt.Tuple2[ParseState, Expr] {
	return frt.IfElse(ok, (func() frt.Tuple2[ParseState, Expr] {
		rtype := frt.Pipe(psTypeVarGen(ps), (func(_r0 func() TypeVar) RecordType { return GenRecordTypeByTgen(rfac, _r0) }))
		return frt.Pipe(frt.Pipe(RecordGen{FieldsNV: neps, RecordType: rtype}, New_Expr_ERecordGen), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(ps, _r0) }))
	}), (func() frt.Tuple2[ParseState, Expr] {
		psPanic(ps, "can't find record type.")
		return frt.NewTuple2(ps, New_Expr_EUnit)
	}))
}

func parseRecordGen(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2, fiInfos := frt.Destr2(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LBRACE, ps), (func(_r0 ParseState) frt.Tuple2[ParseState, fiListInfo] { return parseFieldInitializers(parseE, _r0) })), (func(_r0 frt.Tuple2[ParseState, fiListInfo]) frt.Tuple2[ParseState, fiListInfo] {
		return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RBRACE, _r0) }), _r0)
	})))
	neps := fiInfos.NePairs
	recName := fiInfos.RecName
	return frt.IfElse(frt.OpEqual(recName, ""), (func() frt.Tuple2[ParseState, Expr] {
		rfac, ok := frt.Destr2(frt.Pipe(slice.Map(func(_v1 NEPair) string {
			return _v1.Name
		}, neps), (func(_r0 []string) frt.Tuple2[RecordFactory, bool] { return scLookupRecFac(ps2.scope, _r0) })))
		return retRecordGen(ok, rfac, neps, ps2)
	}), (func() frt.Tuple2[ParseState, Expr] {
		rfac, ok := frt.Destr2(scLookupRecFacByName(ps2.scope, recName))
		return retRecordGen(ok, rfac, neps, ps2)
	}))
}

func refVar(vname string, stlist []FType, ps ParseState) Expr {
	vfac, ok := frt.Destr2(scLookupVarFac(ps.scope, vname))
	return frt.IfElse(ok, (func() Expr {
		return frt.Pipe(frt.Pipe(psTypeVarGen(ps), (func(_r0 func() TypeVar) VarRef { return vfac(stlist, _r0) })), New_Expr_EVarRef)
	}), (func() Expr {
		frt.PipeUnit(frt.Sprintf1("Unknown var ref: %s", vname), (func(_r0 string) { psPanic(ps, _r0) }))
		return New_Expr_EUnit
	}))
}

func parseFAAfterDot(ps ParseState, cur Expr) frt.Tuple2[ParseState, Expr] {
	ps2, fname := frt.Destr2(frt.Pipe(psConsume(New_TokenType_DOT, ps), psIdentNameNx))
	fexpr := frt.Pipe(FieldAccess{TargetExpr: cur, FieldName: fname}, New_Expr_EFieldAccess)
	return frt.IfElse(psCurIs(New_TokenType_DOT, ps2), (func() frt.Tuple2[ParseState, Expr] {
		return parseFAAfterDot(ps2, fexpr)
	}), (func() frt.Tuple2[ParseState, Expr] {
		return frt.NewTuple2(ps2, fexpr)
	}))
}

func parseVarRef(ps ParseState) frt.Tuple2[ParseState, Expr] {
	firstId := psIdentName(ps)
	ps2, stlist := frt.Destr2(frt.IfElse(psIsNeighborLT(ps), (func() frt.Tuple2[ParseState, []FType] {
		return frt.Pipe(psNext(ps), (func(_r0 ParseState) frt.Tuple2[ParseState, []FType] {
			return mightParseSpecifiedTypeList(parseType, _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []FType] {
		return frt.Pipe(emptyFtps(), (func(_r0 []FType) frt.Tuple2[ParseState, []FType] { return PairL(psNext(ps), _r0) }))
	})))
	return frt.IfElse(frt.OpNotEqual(psCurrentTT(ps2), New_TokenType_DOT), (func() frt.Tuple2[ParseState, Expr] {
		return frt.Pipe(refVar(firstId, stlist, ps2), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(ps2, _r0) }))
	}), (func() frt.Tuple2[ParseState, Expr] {
		vfac, ok := frt.Destr2(scLookupVarFac(ps2.scope, firstId))
		return frt.IfElse(ok, (func() frt.Tuple2[ParseState, Expr] {
			return frt.Pipe(frt.Pipe(frt.Pipe(psTypeVarGen(ps2), (func(_r0 func() TypeVar) VarRef { return vfac(emptyFtps(), _r0) })), New_Expr_EVarRef), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return parseFAAfterDot(ps2, _r0) }))
		}), (func() frt.Tuple2[ParseState, Expr] {
			ps3, fullName := frt.Destr2(parseFullName(ps))
			ps4, stlist := frt.Destr2(mightParseSpecifiedTypeList(parseType, ps3))
			return frt.Pipe(refVar(fullName, stlist, ps4), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(ps4, _r0) }))
		}))
	}))
}

func parseSemiExprs(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, []Expr] {
	return ParseSepList(pExpr, New_TokenType_SEMICOLON, ps)
}

func parseSliceExpr(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Expr] {
	return frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LSBRACKET, ps), (func(_r0 ParseState) frt.Tuple2[ParseState, []Expr] { return parseSemiExprs(pExpr, _r0) })), (func(_r0 frt.Tuple2[ParseState, []Expr]) frt.Tuple2[ParseState, []Expr] {
		return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RSBRACKET, _r0) }), _r0)
	})), (func(_r0 frt.Tuple2[ParseState, []Expr]) frt.Tuple2[ParseState, Expr] {
		return MapR(New_Expr_ESlice, _r0)
	}))
}

func exprOnlyBlock(expr Expr) Block {
	emp := slice.New[Stmt]()
	return Block{Stmts: emp, FinalExpr: expr}
}

func parseUSPropAcc(ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2, fname := frt.Destr2(frt.Pipe(psMulConsume(([]TokenType{New_TokenType_UNDER_SCORE, New_TokenType_DOT}), ps), psIdentNameNx))
	tname := uniqueTmpVarName()
	ttype := frt.Pipe(psNewTypeVar(ps2), New_FType_FTypeVar)
	tmpVar := Var{Name: tname, Ftype: ttype}
	vexpr := varToExpr(tmpVar)
	faExpr := frt.Pipe(FieldAccess{TargetExpr: vexpr, FieldName: fname}, New_Expr_EFieldAccess)
	body := exprOnlyBlock(faExpr)
	return frt.Pipe(frt.Pipe(LambdaExpr{Params: ([]Var{tmpVar}), Body: body}, New_Expr_ELambda), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(ps2, _r0) }))
}

func parseAtom(parseE func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Expr] {
	cur := psCurrent(ps)
	pn := psNext(ps)
	switch (cur.ttype).(type) {
	case TokenType_STRING:
		return frt.Pipe(New_Expr_EStringLiteral(cur.stringVal), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(pn, _r0) }))
	case TokenType_SINTERP:
		return frt.Pipe(New_Expr_ESInterP(cur.stringVal), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(pn, _r0) }))
	case TokenType_INT_IMM:
		return frt.Pipe(New_Expr_EIntImm(cur.intVal), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(pn, _r0) }))
	case TokenType_TRUE:
		return frt.Pipe(New_Expr_EBoolLiteral(true), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(pn, _r0) }))
	case TokenType_FALSE:
		return frt.Pipe(New_Expr_EBoolLiteral(false), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(pn, _r0) }))
	case TokenType_LBRACE:
		return parseRecordGen(parseE, ps)
	case TokenType_LSBRACKET:
		return parseSliceExpr(parseE, ps)
	case TokenType_LPAREN:
		return frt.IfElse(frt.OpEqual(psCurrentTT(pn), New_TokenType_RPAREN), (func() frt.Tuple2[ParseState, Expr] {
			return frt.Pipe(frt.NewTuple2(pn, New_Expr_EUnit), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Expr] {
				return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RPAREN, _r0) }), _r0)
			}))
		}), (func() frt.Tuple2[ParseState, Expr] {
			ps2, e1 := frt.Destr2(parseE(pn))
			return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_COMMA), (func() frt.Tuple2[ParseState, Expr] {
				ps3, elist := frt.Destr2(frt.Pipe(frt.Pipe(psConsume(New_TokenType_COMMA, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, []Expr] {
					return ParseSepList(parseE, New_TokenType_COMMA, _r0)
				})), (func(_r0 frt.Tuple2[ParseState, []Expr]) frt.Tuple2[ParseState, []Expr] {
					return MapR((func(_r0 []Expr) []Expr { return slice.PushHead(e1, _r0) }), _r0)
				})))
				frt.IfOnly((slice.Length(elist) > 3), (func() {
					psPanic(ps3, "More then 3 elem tuple, NYI.")
				}))
				return frt.Pipe(frt.Pipe(elist, New_Expr_ETupleExpr), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(psConsume(New_TokenType_RPAREN, ps3), _r0) }))
			}), (func() frt.Tuple2[ParseState, Expr] {
				return frt.Pipe(frt.NewTuple2(ps2, e1), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Expr] {
					return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RPAREN, _r0) }), _r0)
				}))
			}))
		}))
	case TokenType_UNDER_SCORE:
		return parseUSPropAcc(ps)
	case TokenType_IDENTIFIER:
		return frt.IfElse(frt.OpEqual(cur.stringVal, "GoEval"), (func() frt.Tuple2[ParseState, Expr] {
			return parseGoEval(ps)
		}), (func() frt.Tuple2[ParseState, Expr] {
			return parseVarRef(ps)
		}))
	default:
		psPanic(ps, "Unown atom.")
		return frt.NewTuple2(ps, New_Expr_EUnit)
	}
}

func psCurIsBinOp(ps ParseState) bool {
	_, ok := frt.Destr2(lookupBinOp(psCurrentTT(ps)))
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
	ps2, one := frt.Destr2(parseAtom(parseE, ps))
	return frt.IfElse(isEndOfTerm(ps2), (func() frt.Tuple2[ParseState, []Expr] {
		return frt.NewTuple2(ps2, ([]Expr{one}))
	}), (func() frt.Tuple2[ParseState, []Expr] {
		return frt.Pipe(parseAtomList(parseE, ps2), (func(_r0 frt.Tuple2[ParseState, []Expr]) frt.Tuple2[ParseState, []Expr] {
			return MapR((func(_r0 []Expr) []Expr { return slice.PushHead(one, _r0) }), _r0)
		}))
	}))
}

func isDefaultMR(ps ParseState) bool {
	return frt.IfElse(psCurIsNot(New_TokenType_BAR, ps), (func() bool {
		return false
	}), (func() bool {
		ps2 := psConsume(New_TokenType_BAR, ps)
		return frt.OpEqual(psCurrentTT(ps2), New_TokenType_UNDER_SCORE)
	}))
}

func parseDefaultMatchRule(pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, Block] {
	return frt.Pipe(frt.Pipe(psMulConsume(([]TokenType{New_TokenType_BAR, New_TokenType_UNDER_SCORE, New_TokenType_RARROW}), ps), psSkipEOL), pBlock)
}

func parseUnionMatchRule(pBlock func(ParseState) frt.Tuple2[ParseState, Block], target Expr, ps ParseState) frt.Tuple2[ParseState, UnionMatchRule] {
	ps3, cname := frt.Destr2(frt.Pipe(psConsume(New_TokenType_BAR, ps), psIdentNameNx))
	ps4, vname := frt.Destr2((func() frt.Tuple2[ParseState, string] {
		switch (psCurrentTT(ps3)).(type) {
		case TokenType_RARROW:
			return frt.NewTuple2(ps3, "")
		case TokenType_UNDER_SCORE:
			return frt.Pipe(frt.NewTuple2(ps3, "_"), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return MapL(psNext, _r0) }))
		default:
			return psIdentNameNx(ps3)
		}
	})())
	ps5 := frt.Pipe(frt.Pipe(psConsume(New_TokenType_RARROW, ps4), psSkipEOL), psPushScope)
	frt.IfOnly((frt.OpNotEqual(vname, "") && frt.OpNotEqual(vname, "_")), (func() {
		tt := ExprToType(target)
		fu := Cast[FType_FUnion](tt, ps5).Value
		cp := lookupCase(fu, cname)
		scDefVar(ps5.scope, vname, Var{Name: vname, Ftype: cp.Ftype})
	}))
	ps6, block := frt.Destr2(frt.Pipe(pBlock(ps5), (func(_r0 frt.Tuple2[ParseState, Block]) frt.Tuple2[ParseState, Block] { return MapL(psPopScope, _r0) })))
	ump := UnionMatchPattern{CaseId: cname, VarName: vname}
	return frt.Pipe(UnionMatchRule{UnionPattern: ump, Body: block}, (func(_r0 UnionMatchRule) frt.Tuple2[ParseState, UnionMatchRule] { return PairL(ps6, _r0) }))
}

func insideOffside(ps ParseState) bool {
	curCol := psCurCol(ps)
	curOff := psCurOffside(ps)
	return (curCol >= curOff)
}

func parseUnionMatchRules(pBlock func(ParseState) frt.Tuple2[ParseState, Block], target Expr, ps ParseState) frt.Tuple2[ParseState, []UnionMatchRule] {
	nextIsUMR := func(tps ParseState) bool {
		return ((insideOffside(tps) && psCurIs(New_TokenType_BAR, tps)) && frt.OpNot(isDefaultMR(tps)))
	}
	endPred := func(tps ParseState) bool {
		return frt.OpNot(nextIsUMR(tps))
	}
	parseOne := (func(_r0 ParseState) frt.Tuple2[ParseState, UnionMatchRule] {
		return parseUnionMatchRule(pBlock, target, _r0)
	})
	return ParseList2(parseOne, endPred, psSkipEOL, ps)
}

func exaustiveCheck(ttype FType, ucases []UnionMatchRule, ps ParseState) {
	switch _v2 := (ttype).(type) {
	case FType_FUnion:
		fu := _v2.Value
		ui := lookupUniInfo(fu)
		cmap := frt.Pipe(frt.Pipe(ui.Cases, (func(_r0 []NameTypePair) []frt.Tuple2[string, bool] {
			return slice.Map(func(ntp NameTypePair) frt.Tuple2[string, bool] {
				return frt.NewTuple2(ntp.Name, false)
			}, _r0)
		})), dict.ToDict)
		cases := frt.Pipe(frt.Pipe(ucases, (func(_r0 []UnionMatchRule) []UnionMatchPattern {
			return slice.Map(func(_v1 UnionMatchRule) UnionMatchPattern {
				return _v1.UnionPattern
			}, _r0)
		})), (func(_r0 []UnionMatchPattern) []string {
			return slice.Map(func(_v2 UnionMatchPattern) string {
				return _v2.CaseId
			}, _r0)
		}))
		folder := func(dic dict.Dict[string, bool], one string) dict.Dict[string, bool] {
			dict.Add(dic, one, true)
			return dic
		}
		frt.Pipe(cases, (func(_r0 []string) dict.Dict[string, bool] { return slice.Fold(folder, cmap, _r0) }))
		notFounds := frt.Pipe(dict.KVs(cmap), (func(_r0 []frt.Tuple2[string, bool]) []frt.Tuple2[string, bool] {
			return slice.Filter(func(p frt.Tuple2[string, bool]) bool {
				return frt.OpNot(frt.Snd(p))
			}, _r0)
		}))
		frt.IfOnly(slice.IsNotEmpty(notFounds), (func() {
			name, _ := frt.Destr2(slice.Head(notFounds))
			psPanic(ps, frt.SInterP("match does not cover all cases. Can't find case: %s.", name))
		}))
	default:

	}
}

func parseURules(pBlock func(ParseState) frt.Tuple2[ParseState, Block], target Expr, ps ParseState) frt.Tuple2[ParseState, UnionMatchRules] {
	ps2, us := frt.Destr2(parseUnionMatchRules(pBlock, target, ps))
	return frt.IfElse((insideOffside(ps2) && isDefaultMR(ps2)), (func() frt.Tuple2[ParseState, UnionMatchRules] {
		ps3, db := frt.Destr2(parseDefaultMatchRule(pBlock, ps2))
		return frt.Pipe(frt.Pipe(UnionMatchRulesWD{Unions: us, Default: db}, New_UnionMatchRules_UCaseWD), (func(_r0 UnionMatchRules) frt.Tuple2[ParseState, UnionMatchRules] { return PairL(ps3, _r0) }))
	}), (func() frt.Tuple2[ParseState, UnionMatchRules] {
		exaustiveCheck(ExprToType(target), us, ps)
		return frt.NewTuple2(ps2, New_UnionMatchRules_UCaseOnly(us))
	}))
}

func parseStringMatchRule(pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, StringMatchRule] {
	ps2, pat := frt.Destr2(frt.Pipe(psConsume(New_TokenType_BAR, ps), psStringValNx))
	ps3 := frt.Pipe(frt.Pipe(frt.Pipe(ps2, (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RARROW, _r0) })), psSkipEOL), psPushScope)
	ps4, block := frt.Destr2(frt.Pipe(pBlock(ps3), (func(_r0 frt.Tuple2[ParseState, Block]) frt.Tuple2[ParseState, Block] { return MapL(psPopScope, _r0) })))
	return frt.Pipe(StringMatchRule{LiteralPattern: pat, Body: block}, (func(_r0 StringMatchRule) frt.Tuple2[ParseState, StringMatchRule] { return PairL(ps4, _r0) }))
}

func parseStringVarRule(pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, StringVarMatchRule] {
	ps2, vname := frt.Destr2(frt.Pipe(psConsume(New_TokenType_BAR, ps), psIdentNameNx))
	ps3 := frt.Pipe(frt.Pipe(frt.Pipe(ps2, (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RARROW, _r0) })), psSkipEOL), psPushScope)
	scDefVar(ps3.scope, vname, Var{Name: vname, Ftype: New_FType_FString})
	ps4, block := frt.Destr2(frt.Pipe(pBlock(ps3), (func(_r0 frt.Tuple2[ParseState, Block]) frt.Tuple2[ParseState, Block] { return MapL(psPopScope, _r0) })))
	return frt.Pipe(StringVarMatchRule{VarName: vname, Body: block}, (func(_r0 StringVarMatchRule) frt.Tuple2[ParseState, StringVarMatchRule] { return PairL(ps4, _r0) }))
}

func isSLitRule(ps ParseState) bool {
	return frt.IfElse(psCurIsNot(New_TokenType_BAR, ps), (func() bool {
		return false
	}), (func() bool {
		return frt.IfElse(psNextIs(New_TokenType_STRING, ps), (func() bool {
			return true
		}), (func() bool {
			return false
		}))
	}))
}

func parseSMRules(pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, []StringMatchRule] {
	one := (func(_r0 ParseState) frt.Tuple2[ParseState, StringMatchRule] { return parseStringMatchRule(pBlock, _r0) })
	endPred := func(tps ParseState) bool {
		return frt.OpNot(isSLitRule(tps))
	}
	next := func(tps ParseState) ParseState {
		return tps
	}
	return ParseList2(one, endPred, next, ps)
}

func parseSRules(pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, StringMatchRules] {
	ps2, ss := frt.Destr2(parseSMRules(pBlock, ps))
	return frt.IfElse((insideOffside(ps2) && isDefaultMR(ps2)), (func() frt.Tuple2[ParseState, StringMatchRules] {
		ps3, db := frt.Destr2(parseDefaultMatchRule(pBlock, ps2))
		return frt.Pipe(frt.Pipe(StringMatchRulesWD{Literals: ss, Default: db}, New_StringMatchRules_SCaseWD), (func(_r0 StringMatchRules) frt.Tuple2[ParseState, StringMatchRules] { return PairL(ps3, _r0) }))
	}), (func() frt.Tuple2[ParseState, StringMatchRules] {
		ps3, vr := frt.Destr2(parseStringVarRule(pBlock, ps2))
		return frt.Pipe(frt.Pipe(StringMatchRulesWV{Literals: ss, VarRule: vr}, New_StringMatchRules_SCaseWV), (func(_r0 StringMatchRules) frt.Tuple2[ParseState, StringMatchRules] { return PairL(ps3, _r0) }))
	}))
}

func isUnionMatchRules(target Expr, ps0 ParseState) bool {
	ps := psConsume(New_TokenType_BAR, ps0)
	switch (ExprToType(target)).(type) {
	case FType_FUnion:
		return true
	case FType_FString:
		return false
	default:
		switch (psCurrentTT(ps)).(type) {
		case TokenType_STRING:
			return false
		case TokenType_IDENTIFIER:
			switch (psNextTT(ps)).(type) {
			case TokenType_IDENTIFIER:
				return true
			case TokenType_RARROW:
				psPanic(ps, "Can't distinguish String var pattern or union case only pattern. Syntax error for a while.")
				return false
			default:
				psPanic(ps, "Unknown case rule of match expr(2)")
				return false
			}
		default:
			psPanic(ps, "Unknown case rule of match expr")
			return false
		}
	}
}

func isStringMatchRules(target Expr, ps0 ParseState) bool {
	ps := psConsume(New_TokenType_BAR, ps0)
	switch (ExprToType(target)).(type) {
	case FType_FUnion:
		return false
	case FType_FString:
		return true
	default:
		switch (psCurrentTT(ps)).(type) {
		case TokenType_STRING:
			return true
		case TokenType_IDENTIFIER:
			switch (psNextTT(ps)).(type) {
			case TokenType_IDENTIFIER:
				return false
			case TokenType_RARROW:
				psPanic(ps, "Can't distinguish String var pattern or union case only pattern. Syntax error for a while.")
				return false
			default:
				psPanic(ps, "Unknown case rule of match expr(2)")
				return false
			}
		default:
			psPanic(ps, "Unknown case rule of match expr")
			return false
		}
	}
}

func parseMatchRules(pBlock func(ParseState) frt.Tuple2[ParseState, Block], target Expr, ps ParseState) frt.Tuple2[ParseState, MatchRules] {
	dummy := frt.Pipe(frt.Empty[MatchRules](), (func(_r0 MatchRules) frt.Tuple2[ParseState, MatchRules] { return PairL(ps, _r0) }))
	return frt.IfElse(isDefaultMR(ps), (func() frt.Tuple2[ParseState, MatchRules] {
		psPanic(ps, "Only default case, illegal.")
		return dummy
	}), (func() frt.Tuple2[ParseState, MatchRules] {
		return frt.IfElse(isUnionMatchRules(target, ps), (func() frt.Tuple2[ParseState, MatchRules] {
			return frt.Pipe(parseURules(pBlock, target, ps), (func(_r0 frt.Tuple2[ParseState, UnionMatchRules]) frt.Tuple2[ParseState, MatchRules] {
				return MapR(New_MatchRules_RUnions, _r0)
			}))
		}), (func() frt.Tuple2[ParseState, MatchRules] {
			return frt.IfElse(isStringMatchRules(target, ps), (func() frt.Tuple2[ParseState, MatchRules] {
				return frt.Pipe(parseSRules(pBlock, ps), (func(_r0 frt.Tuple2[ParseState, StringMatchRules]) frt.Tuple2[ParseState, MatchRules] {
					return MapR(New_MatchRules_RStrings, _r0)
				}))
			}), (func() frt.Tuple2[ParseState, MatchRules] {
				psPanic(ps, "Unknown match case, illegal.")
				return dummy
			}))
		}))
	}))
}

func parseMatchExpr(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, MatchExpr] {
	ps2, target := frt.Destr2(frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_MATCH, ps), pExpr), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Expr] {
		return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_WITH, _r0) }), _r0)
	})), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Expr] { return MapL(psSkipEOL, _r0) })))
	ps3, rules := frt.Destr2(parseMatchRules(pBlock, target, ps2))
	return frt.Pipe(MatchExpr{Target: target, Rules: rules}, (func(_r0 MatchExpr) frt.Tuple2[ParseState, MatchExpr] { return PairL(ps3, _r0) }))
}

func parseInlineBlock(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, Block] {
	ps2, expr := frt.Destr2(pExpr(ps))
	return frt.Pipe(exprOnlyBlock(expr), (func(_r0 Block) frt.Tuple2[ParseState, Block] { return PairL(ps2, _r0) }))
}

func parseIfAfterIfExpr(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2, cond := frt.Destr2(frt.Pipe(pExpr(ps), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Expr] {
		return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_THEN, _r0) }), _r0)
	})))
	tgen := psTypeVarGen(ps2)
	recurse := (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] { return parseIfAfterIfExpr(pExpr, pBlock, _r0) })
	return frt.IfElse(psCurIs(New_TokenType_EOL, ps2), (func() frt.Tuple2[ParseState, Expr] {
		ps3, tbody := frt.Destr2(frt.Pipe(psSkipEOL(ps2), pBlock))
		ps4 := psSkipEOL(ps3)
		return frt.IfElse(psCurIs(New_TokenType_ELSE, ps4), (func() frt.Tuple2[ParseState, Expr] {
			pse2, fbody := frt.Destr2(frt.Pipe(frt.Pipe(psConsume(New_TokenType_ELSE, ps4), psSkipEOL), pBlock))
			return frt.Pipe(newIfElseCall(tgen, cond, tbody, fbody), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(pse2, _r0) }))
		}), (func() frt.Tuple2[ParseState, Expr] {
			return frt.IfElse(psCurIs(New_TokenType_ELIF, ps4), (func() frt.Tuple2[ParseState, Expr] {
				ps5, elseExpr := frt.Destr2(frt.Pipe(psConsume(New_TokenType_ELIF, ps4), recurse))
				ebody := exprOnlyBlock(elseExpr)
				return frt.Pipe(newIfElseCall(tgen, cond, tbody, ebody), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(ps5, _r0) }))
			}), (func() frt.Tuple2[ParseState, Expr] {
				return frt.Pipe(newIfOnlyCall(tgen, cond, tbody), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(ps3, _r0) }))
			}))
		}))
	}), (func() frt.Tuple2[ParseState, Expr] {
		psi2, tbody := frt.Destr2(parseInlineBlock(pExpr, ps2))
		return frt.IfElse(psCurIs(New_TokenType_ELSE, psi2), (func() frt.Tuple2[ParseState, Expr] {
			psi3, fbody := frt.Destr2(frt.Pipe(psConsume(New_TokenType_ELSE, psi2), (func(_r0 ParseState) frt.Tuple2[ParseState, Block] { return parseInlineBlock(pExpr, _r0) })))
			return frt.Pipe(newIfElseCall(tgen, cond, tbody, fbody), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(psi3, _r0) }))
		}), (func() frt.Tuple2[ParseState, Expr] {
			return frt.Pipe(newIfOnlyCall(tgen, cond, tbody), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(psi2, _r0) }))
		}))
	}))
}

func parseIfExpr(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, Expr] {
	return frt.Pipe(psConsume(New_TokenType_IF, ps), (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] { return parseIfAfterIfExpr(pExpr, pBlock, _r0) }))
}

func tvgen2ftvgen(tgen func() TypeVar) FType {
	return frt.Pipe(tgen(), New_FType_FTypeVar)
}

func updateFunCallFunType(tgen func() TypeVar, res Resolver, vr VarRef, args []Expr) VarRef {
	switch _v3 := (vr).(type) {
	case VarRef_VRVar:
		v := _v3.Value
		switch _v4 := (v.Ftype).(type) {
		case FType_FFunc:
			return vr
		case FType_FTypeVar:
			tv := _v4.Value
			ntv := frt.Pipe(tgen(), New_FType_FTypeVar)
			nftype := frt.Pipe(frt.Pipe(slice.Map(ExprToType, args), (func(_r0 []FType) []FType { return slice.PushLast(ntv, _r0) })), newFFunc)
			frt.Pipe(UniRel{SrcV: tv.Name, Dest: nftype}, (func(_r0 UniRel) []UniRel { return updateResOne(res, _r0) }))
			return frt.Pipe(Var{Name: v.Name, Ftype: nftype}, New_VarRef_VRVar)
		default:
			PanicNow("Unknown funcall first arg type.")
			return vr
		}
	case VarRef_VRSVar:
		return vr
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func parseFunExpr(pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, Expr] {
	ps2, params := frt.Destr2(frt.Pipe(frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_FUN, ps), psPushScope), parseParams), (func(_r0 frt.Tuple2[ParseState, []Var]) frt.Tuple2[ParseState, []Var] {
		return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RARROW, _r0) }), _r0)
	})), (func(_r0 frt.Tuple2[ParseState, []Var]) frt.Tuple2[ParseState, []Var] { return MapL(psSkipEOL, _r0) })))
	ps3, body := frt.Destr2(frt.Pipe(pBlock(ps2), (func(_r0 frt.Tuple2[ParseState, Block]) frt.Tuple2[ParseState, Block] { return MapL(psPopScope, _r0) })))
	return frt.NewTuple2(ps3, New_Expr_ELambda(LambdaExpr{Params: params, Body: body}))
}

func parseTerm(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pBlock func(ParseState) frt.Tuple2[ParseState, Block], ps ParseState) frt.Tuple2[ParseState, Expr] {
	switch (psCurrentTT(ps)).(type) {
	case TokenType_MATCH:
		return frt.Pipe(frt.Pipe(parseMatchExpr(pExpr, pBlock, ps), (func(_r0 frt.Tuple2[ParseState, MatchExpr]) frt.Tuple2[ParseState, ReturnableExpr] {
			return MapR(New_ReturnableExpr_RMatchExpr, _r0)
		})), (func(_r0 frt.Tuple2[ParseState, ReturnableExpr]) frt.Tuple2[ParseState, Expr] {
			return MapR(New_Expr_EReturnableExpr, _r0)
		}))
	case TokenType_LSBRACKET:
		return parseSliceExpr(pExpr, ps)
	case TokenType_FUN:
		return parseFunExpr(pBlock, ps)
	case TokenType_IF:
		return parseIfExpr(pExpr, pBlock, ps)
	case TokenType_NOT:
		ps2, target := frt.Destr2(frt.Pipe(psConsume(New_TokenType_NOT, ps), (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] { return parseTerm(pExpr, pBlock, _r0) })))
		return frt.Pipe(newUnaryNotCall(psTypeVarGen(ps2), target), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return PairL(ps2, _r0) }))
	default:
		ps2, es := frt.Destr2(parseAtomList(pExpr, ps))
		return frt.IfElse(frt.OpEqual(slice.Length(es), 1), (func() frt.Tuple2[ParseState, Expr] {
			return frt.NewTuple2(ps2, slice.Head(es))
		}), (func() frt.Tuple2[ParseState, Expr] {
			head := slice.Head(es)
			tail := slice.Tail(es)
			switch _v5 := (head).(type) {
			case Expr_EVarRef:
				vr := _v5.Value
				tgen := psTypeVarGen(ps2)
				nvr := updateFunCallFunType(tgen, ps2.tvc.resolver, vr, tail)
				fc := FunCall{TargetFunc: nvr, Args: tail}
				return frt.NewTuple2(ps2, New_Expr_EFunCall(fc))
			default:
				psPanic(ps2, "Funcall head is not var")
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
		return frt.IfElse((bop.Precedence < minPrec), (func() frt.Tuple2[ParseState, Expr] {
			return frt.NewTuple2(ps, cur)
		}), (func() frt.Tuple2[ParseState, Expr] {
			ps3, rhs := frt.Destr2(frt.Pipe(psConsume(btk, ps2), (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] { return pEwithMinPrec((bop.Precedence + 1), _r0) })))
			tvgen := psTypeVarGen(ps3)
			return frt.Pipe(newBinOpCall(tvgen, btk, bop, cur, rhs), (func(_r0 Expr) frt.Tuple2[ParseState, Expr] { return parseBinAfter(pEwithMinPrec, minPrec, ps3, _r0) }))
		}))
	}), (func() frt.Tuple2[ParseState, Expr] {
		return frt.NewTuple2(ps, cur)
	}))
}

func parseExprWithPrec(pBlock func(ParseState) frt.Tuple2[ParseState, Block], minPrec int, ps ParseState) frt.Tuple2[ParseState, Expr] {
	pExpr := (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] { return parseExprWithPrec(pBlock, 1, _r0) })
	ps2, expr := frt.Destr2(parseTerm(pExpr, pBlock, ps))
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

func parseStmt(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pLet func(ParseState) frt.Tuple2[ParseState, LLetVarDef], ps ParseState) frt.Tuple2[ParseState, Stmt] {
	switch (psCurrentTT(ps)).(type) {
	case TokenType_LET:
		return frt.Pipe(pLet(ps), (func(_r0 frt.Tuple2[ParseState, LLetVarDef]) frt.Tuple2[ParseState, Stmt] {
			return MapR(New_Stmt_SLetVarDef, _r0)
		}))
	default:
		return frt.Pipe(pExpr(ps), (func(_r0 frt.Tuple2[ParseState, Expr]) frt.Tuple2[ParseState, Stmt] {
			return MapR(New_Stmt_SExprStmt, _r0)
		}))
	}
}

func isEndOfBlock(ps ParseState) bool {
	isOffside := (psCurCol(ps) < psCurOffside(ps))
	return ((isOffside || psCurIs(New_TokenType_EOF, ps)) || psCurIs(New_TokenType_RPAREN, ps))
}

func psStmtOneForList(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pLet func(ParseState) frt.Tuple2[ParseState, LLetVarDef], ps ParseState) frt.Tuple2[ParseState, Stmt] {
	psForErrMsg(ps)
	return frt.Pipe(parseStmt(pExpr, pLet, ps), (func(_r0 frt.Tuple2[ParseState, Stmt]) frt.Tuple2[ParseState, Stmt] { return MapL(psSkipEOL, _r0) }))
}

func parseStmtList(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pLet func(ParseState) frt.Tuple2[ParseState, LLetVarDef], ps ParseState) frt.Tuple2[ParseState, []Stmt] {
	pOne := (func(_r0 ParseState) frt.Tuple2[ParseState, Stmt] { return psStmtOneForList(pExpr, pLet, _r0) })
	return ParseList2(pOne, isEndOfBlock, func(p ParseState) ParseState {
		return p
	}, ps)
}

func parseBlockAfterPushScope(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], pLet func(ParseState) frt.Tuple2[ParseState, LLetVarDef], ps ParseState) frt.Tuple2[ParseState, Block] {
	ps2, sls := frt.Destr2(frt.Pipe(frt.Pipe(frt.Pipe(psPushOffside(ps), (func(_r0 ParseState) frt.Tuple2[ParseState, []Stmt] { return parseStmtList(pExpr, pLet, _r0) })), (func(_r0 frt.Tuple2[ParseState, []Stmt]) frt.Tuple2[ParseState, []Stmt] {
		return MapL(psPopOffside, _r0)
	})), (func(_r0 frt.Tuple2[ParseState, []Stmt]) frt.Tuple2[ParseState, []Stmt] { return MapL(psPopScope, _r0) })))
	last := slice.Last(sls)
	stmts := slice.PopLast(sls)
	switch _v6 := (last).(type) {
	case Stmt_SExprStmt:
		e := _v6.Value
		return frt.NewTuple2(ps2, Block{Stmts: stmts, FinalExpr: e})
	default:
		psPanic(ps2, "block of last is not expr")
		return frt.Pipe(frt.Empty[Block](), (func(_r0 Block) frt.Tuple2[ParseState, Block] { return PairL(ps2, _r0) }))
	}
}

func parseBlock(pLet func(ParseState) frt.Tuple2[ParseState, LLetVarDef], ps ParseState) frt.Tuple2[ParseState, Block] {
	pExpr := (func(_r0 ParseState) frt.Tuple2[ParseState, Expr] {
		return parseExpr((func(_r0 ParseState) frt.Tuple2[ParseState, Block] { return parseBlock(pLet, _r0) }), _r0)
	})
	return frt.Pipe(psPushScope(ps), (func(_r0 ParseState) frt.Tuple2[ParseState, Block] { return parseBlockAfterPushScope(pExpr, pLet, _r0) }))
}

func parseLetOneVarDef(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, LetVarDef] {
	ps2, vname := frt.Destr2(frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LET, ps), psIdentNameNx), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] {
		return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_EQ, _r0) }), _r0)
	})), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return MapL(psSkipEOL, _r0) })))
	ps3, rhs0 := frt.Destr2(pExpr(ps2))
	rhs := InferExpr(ps3.tvc, rhs0)
	v := Var{Name: vname, Ftype: ExprToType(rhs)}
	scDefVar(ps3.scope, vname, v)
	return frt.Pipe(LetVarDef{Lvar: v, Rhs: rhs}, (func(_r0 LetVarDef) frt.Tuple2[ParseState, LetVarDef] { return PairL(ps3, _r0) }))
}

func psIdentOrUSNameNx(ps ParseState) frt.Tuple2[ParseState, string] {
	return frt.IfElse(psCurIs(New_TokenType_IDENTIFIER, ps), (func() frt.Tuple2[ParseState, string] {
		return psIdentNameNx(ps)
	}), (func() frt.Tuple2[ParseState, string] {
		ps2 := psNext(ps)
		return frt.NewTuple2(ps2, "_")
	}))
}

func defVarIfNecessary(sc Scope, v Var) {
	frt.IfOnly(frt.OpNotEqual(v.Name, "_"), (func() {
		scDefVar(sc, v.Name, v)
	}))
}

func parseLetDestVarDef(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, LetDestVarDef] {
	ps2 := psMulConsume(([]TokenType{New_TokenType_LET, New_TokenType_LPAREN}), ps)
	ps3, vnames := frt.Destr2(frt.Pipe(frt.Pipe(ParseSepList(psIdentOrUSNameNx, New_TokenType_COMMA, ps2), (func(_r0 frt.Tuple2[ParseState, []string]) frt.Tuple2[ParseState, []string] {
		return MapL((func(_r0 ParseState) ParseState {
			return psMulConsume(([]TokenType{New_TokenType_RPAREN, New_TokenType_EQ}), _r0)
		}), _r0)
	})), (func(_r0 frt.Tuple2[ParseState, []string]) frt.Tuple2[ParseState, []string] {
		return MapL(psSkipEOL, _r0)
	})))
	frt.IfOnly((slice.Length(vnames) > 3), (func() {
		psPanic(ps3, "More than 3 let destructuring, NYI")
	}))
	ps4, rhs0 := frt.Destr2(pExpr(ps3))
	rhs := InferExpr(ps4.tvc, rhs0)
	rtype := ExprToType(rhs)
	switch _v7 := (rtype).(type) {
	case FType_FTuple:
		tup := _v7.Value
		vars := frt.Pipe(slice.Zip(vnames, tup.ElemTypes), (func(_r0 []frt.Tuple2[string, FType]) []Var {
			return slice.Map(func(t frt.Tuple2[string, FType]) Var {
				return newVar(frt.Fst(t), frt.Snd(t))
			}, _r0)
		}))
		slice.Iter((func(_r0 Var) { defVarIfNecessary(ps4.scope, _r0) }), vars)
		return frt.Pipe(LetDestVarDef{Lvars: vars, Rhs: rhs}, (func(_r0 LetDestVarDef) frt.Tuple2[ParseState, LetDestVarDef] { return PairL(ps4, _r0) }))
	case FType_FTypeVar:
		tpgen := psTypeVarGen(ps4)
		name2tp := func(name string) FType {
			return frt.IfElse(frt.OpEqual(name, "_"), (func() FType {
				return New_FType_FUnit
			}), (func() FType {
				return frt.Pipe(tpgen(), New_FType_FTypeVar)
			}))
		}
		vtypes := slice.Map(name2tp, vnames)
		vars := frt.Pipe(slice.Zip(vnames, vtypes), (func(_r0 []frt.Tuple2[string, FType]) []Var {
			return slice.Map(func(t frt.Tuple2[string, FType]) Var {
				return newVar(frt.Fst(t), frt.Snd(t))
			}, _r0)
		}))
		slice.Iter((func(_r0 Var) { defVarIfNecessary(ps4.scope, _r0) }), vars)
		return frt.Pipe(LetDestVarDef{Lvars: vars, Rhs: rhs}, (func(_r0 LetDestVarDef) frt.Tuple2[ParseState, LetDestVarDef] { return PairL(ps4, _r0) }))
	default:
		psPanic(ps2, "Destructuring let, but rhs is not tuple. NYI.")
		dummy := frt.Empty[[]Var]()
		return frt.NewTuple2(ps2, LetDestVarDef{Lvars: dummy, Rhs: New_Expr_EUnit})
	}
}

func parseLetFuncDef(pLet func(ParseState) frt.Tuple2[ParseState, LLetVarDef], ps ParseState) frt.Tuple2[ParseState, LetFuncDef] {
	ps2 := frt.Pipe(psConsume(New_TokenType_LET, ps), psPushScope)
	fname := psIdentName(ps2)
	ps3, params := frt.Destr2(frt.Pipe(psNext(ps2), parseParams))
	ps4, rtypeDef := frt.Destr2(frt.IfElse(psCurIs(New_TokenType_COLON, ps3), (func() frt.Tuple2[ParseState, FType] {
		return frt.Pipe(psConsume(New_TokenType_COLON, ps3), parseType)
	}), (func() frt.Tuple2[ParseState, FType] {
		tvgen := psTypeVarGen(ps3)
		tvf := frt.Pipe(tvgen(), New_FType_FTypeVar)
		return frt.NewTuple2(ps3, tvf)
	})))
	paramTypes := slice.Map(func(_v1 Var) FType {
		return _v1.Ftype
	}, params)
	defTargets := slice.PushLast(rtypeDef, paramTypes)
	defFt := newFFunc(defTargets)
	defVar := Var{Name: fname, Ftype: defFt}
	scDefVar(ps4.scope, fname, defVar)
	ps5, block := frt.Destr2(frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_EQ, ps4), psSkipEOL), (func(_r0 ParseState) frt.Tuple2[ParseState, Block] { return parseBlock(pLet, _r0) })), (func(_r0 frt.Tuple2[ParseState, Block]) frt.Tuple2[ParseState, Block] { return MapL(psPopScope, _r0) })))
	rtype := (func() FType {
		switch (rtypeDef).(type) {
		case FType_FTypeVar:
			return frt.Pipe(blockToExpr(block), ExprToType)
		default:
			return rtypeDef
		}
	})()
	targets := frt.IfElse(frt.OpEqual(slice.Length(params), 0), (func() []FType {
		return ([]FType{New_FType_FUnit, rtype})
	}), (func() []FType {
		return frt.Pipe(paramTypes, (func(_r0 []FType) []FType { return slice.PushLast(rtype, _r0) }))
	}))
	ft := newFFunc(targets)
	fnvar := Var{Name: fname, Ftype: ft}
	return frt.Pipe(LetFuncDef{Fvar: fnvar, Params: params, Body: block}, (func(_r0 LetFuncDef) frt.Tuple2[ParseState, LetFuncDef] { return PairL(ps5, _r0) }))
}

func rfdToFuncFactory(rfd RootFuncDef) FuncFactory {
	targets := (func() []FType {
		switch _v8 := (rfd.Lfd.Fvar.Ftype).(type) {
		case FType_FFunc:
			ft := _v8.Value
			return ft.Targets
		default:
			PanicNow("root func def let with non func var, bug")
			return []FType{}
		}
	})()
	return FuncFactory{Tparams: rfd.Tparams, Targets: targets}
}

func parseRootLetFuncDef(pLet func(ParseState) frt.Tuple2[ParseState, LLetVarDef], ps ParseState) frt.Tuple2[ParseState, RootFuncDef] {
	psForErrMsg(ps)
	ps2, lfd := frt.Destr2(parseLetFuncDef(pLet, ps))
	psForErrMsg(ps)
	rfd := InferLfd(ps2.tvc, lfd)
	frt.PipeUnit(rfdToFuncFactory(rfd), (func(_r0 FuncFactory) { scRegFunFac(ps2.scope, rfd.Lfd.Fvar.Name, _r0) }))
	return frt.NewTuple2(ps2, rfd)
}

func parseLetVarDef(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, LLetVarDef] {
	peekPs := psConsume(New_TokenType_LET, ps)
	return frt.IfElse(psCurIs(New_TokenType_LPAREN, peekPs), (func() frt.Tuple2[ParseState, LLetVarDef] {
		return frt.Pipe(parseLetDestVarDef(pExpr, ps), (func(_r0 frt.Tuple2[ParseState, LetDestVarDef]) frt.Tuple2[ParseState, LLetVarDef] {
			return MapR(New_LLetVarDef_LLDestVarDef, _r0)
		}))
	}), (func() frt.Tuple2[ParseState, LLetVarDef] {
		return frt.Pipe(parseLetOneVarDef(pExpr, ps), (func(_r0 frt.Tuple2[ParseState, LetVarDef]) frt.Tuple2[ParseState, LLetVarDef] {
			return MapR(New_LLetVarDef_LLOneVarDef, _r0)
		}))
	}))
}

func lfdToLetVar(lfd LetFuncDef) LetVarDef {
	elambda := frt.Pipe(LambdaExpr{Params: lfd.Params, Body: lfd.Body}, New_Expr_ELambda)
	return LetVarDef{Lvar: lfd.Fvar, Rhs: elambda}
}

func rawLetToLetVar(rawLet RawLetDef, ps ParseState) LLetVarDef {
	switch _v9 := (rawLet).(type) {
	case RawLetDef_RLetOneVar:
		one := _v9.Value
		return New_LLetVarDef_LLOneVarDef(one)
	case RawLetDef_RLetDestVar:
		dest := _v9.Value
		return New_LLetVarDef_LLDestVarDef(dest)
	case RawLetDef_RLetFunc:
		lfd := _v9.Value
		lvd := lfdToLetVar(lfd)
		nrhs := InferExpr(ps.tvc, lvd.Rhs)
		nv := Var{Name: lvd.Lvar.Name, Ftype: ExprToType(nrhs)}
		scDefVar(ps.scope, nv.Name, nv)
		return frt.Pipe(LetVarDef{Lvar: nv, Rhs: nrhs}, New_LLetVarDef_LLOneVarDef)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func parseRawLet(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, RawLetDef] {
	pLet := func(tps ParseState) frt.Tuple2[ParseState, LLetVarDef] {
		tps2, rl := frt.Destr2(parseRawLet(pExpr, tps))
		return frt.Pipe(rawLetToLetVar(rl, tps2), (func(_r0 LLetVarDef) frt.Tuple2[ParseState, LLetVarDef] { return PairL(tps2, _r0) }))
	}
	psN := psNext(ps)
	psNN := psNext(psN)
	switch (psCurrentTT(psN)).(type) {
	case TokenType_LPAREN:
		return frt.Pipe(parseLetDestVarDef(pExpr, ps), (func(_r0 frt.Tuple2[ParseState, LetDestVarDef]) frt.Tuple2[ParseState, RawLetDef] {
			return MapR(New_RawLetDef_RLetDestVar, _r0)
		}))
	default:
		switch (psCurrentTT(psNN)).(type) {
		case TokenType_EQ:
			return frt.Pipe(parseLetOneVarDef(pExpr, ps), (func(_r0 frt.Tuple2[ParseState, LetVarDef]) frt.Tuple2[ParseState, RawLetDef] {
				return MapR(New_RawLetDef_RLetOneVar, _r0)
			}))
		default:
			return frt.Pipe(parseLetFuncDef(pLet, ps), (func(_r0 frt.Tuple2[ParseState, LetFuncDef]) frt.Tuple2[ParseState, RawLetDef] {
				return MapR(New_RawLetDef_RLetFunc, _r0)
			}))
		}
	}
}

func parseRootLet(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps0 ParseState) frt.Tuple2[ParseState, RootStmt] {
	ps := psResetTmpCtx(ps0)
	psForErrMsg(ps)
	ps2, rlet := frt.Destr2(parseRawLet(pExpr, ps))
	switch _v10 := (rlet).(type) {
	case RawLetDef_RLetOneVar:
		one := _v10.Value
		lv := one.Lvar
		scDefVar(ps2.scope, lv.Name, lv)
		rootVd := frt.Pipe(RootVarDef{Vdef: one}, New_RootStmt_RSRootVarDef)
		return frt.NewTuple2(ps2, rootVd)
	case RawLetDef_RLetDestVar:
		psPanic(ps, "Root destructuring let, NYI.")
		return frt.NewTuple2(ps2, frt.Empty[RootStmt]())
	case RawLetDef_RLetFunc:
		lfd := _v10.Value
		psForErrMsg(ps)
		rfd := InferLfd(ps2.tvc, lfd)
		frt.PipeUnit(rfdToFuncFactory(rfd), (func(_r0 FuncFactory) { scRegFunFac(ps2.scope, rfd.Lfd.Fvar.Name, _r0) }))
		return frt.NewTuple2(ps2, New_RootStmt_RSRootFuncDef(rfd))
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func parseFieldDef(ps ParseState) frt.Tuple2[ParseState, NameTypePair] {
	fname := psIdentName(ps)
	ps2, tp := frt.Destr2(frt.Pipe(frt.Pipe(psNextNOL(ps), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_COLON, _r0) })), parseType))
	ntp := NameTypePair{Name: fname, Ftype: tp}
	return frt.NewTuple2(ps2, ntp)
}

func parseFieldDefs(ps ParseState) frt.Tuple2[ParseState, []NameTypePair] {
	ps2, ntp := frt.Destr2(frt.Pipe(psSkipEOL(ps), parseFieldDef))
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_RBRACE), (func() frt.Tuple2[ParseState, []NameTypePair] {
		return frt.NewTuple2(ps2, ([]NameTypePair{ntp}))
	}), (func() frt.Tuple2[ParseState, []NameTypePair] {
		ps3 := frt.Pipe(psConsume(New_TokenType_SEMICOLON, ps2), psSkipEOL)
		return frt.IfElse(psCurIs(New_TokenType_RBRACE, ps3), (func() frt.Tuple2[ParseState, []NameTypePair] {
			return frt.NewTuple2(ps3, ([]NameTypePair{ntp}))
		}), (func() frt.Tuple2[ParseState, []NameTypePair] {
			return frt.Pipe(parseFieldDefs(ps3), (func(_r0 frt.Tuple2[ParseState, []NameTypePair]) frt.Tuple2[ParseState, []NameTypePair] {
				return MapR((func(_r0 []NameTypePair) []NameTypePair { return slice.PushHead(ntp, _r0) }), _r0)
			}))
		}))
	}))
}

func regTypeVar(ps ParseState, tname string) {
	frt.PipeUnit(New_FType_FTypeVar(TypeVar{Name: tname}), (func(_r0 FType) { scRegisterType(ps.scope, tname, _r0) }))
}

func psRegTypeVars(ps ParseState, tnames []string) {
	slice.Iter((func(_r0 string) { regTypeVar(ps, _r0) }), tnames)
}

func parseRecordDef(tname string, pnames []string, ps0 ParseState) frt.Tuple2[ParseState, RecordDef] {
	ps := psPushScope(ps0)
	psRegTypeVars(ps, pnames)
	ps2, ntps := frt.Destr2(frt.Pipe(frt.Pipe(frt.Pipe(psConsume(New_TokenType_LBRACE, ps), parseFieldDefs), (func(_r0 frt.Tuple2[ParseState, []NameTypePair]) frt.Tuple2[ParseState, []NameTypePair] {
		return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RBRACE, _r0) }), _r0)
	})), (func(_r0 frt.Tuple2[ParseState, []NameTypePair]) frt.Tuple2[ParseState, []NameTypePair] {
		return MapL(psPopScope, _r0)
	})))
	rd := RecordDef{Name: tname, Tparams: pnames, Fields: ntps}
	psRegRecDefToTDCtx(rd, ps2)
	return frt.NewTuple2(ps2, rd)
}

func parseOneCaseDef(ps ParseState) frt.Tuple2[ParseState, NameTypePair] {
	ps2, cname := frt.Destr2(frt.Pipe(psConsume(New_TokenType_BAR, ps), psIdentNameNx))
	switch (psCurrentTT(ps2)).(type) {
	case TokenType_OF:
		ps3, tp := frt.Destr2(frt.Pipe(psConsume(New_TokenType_OF, ps2), parseType))
		cs := NameTypePair{Name: cname, Ftype: tp}
		return frt.NewTuple2(ps3, cs)
	default:
		ps3 := psConsume(New_TokenType_EOL, ps2)
		cs := NameTypePair{Name: cname, Ftype: New_FType_FUnit}
		return frt.NewTuple2(ps3, cs)
	}
}

func parseCaseDefs(ps ParseState) frt.Tuple2[ParseState, []NameTypePair] {
	ps2, cs := frt.Destr2(parseOneCaseDef(ps))
	ps3 := psSkipEOL(ps2)
	return frt.IfElse(frt.OpEqual(psCurrentTT(ps3), New_TokenType_BAR), (func() frt.Tuple2[ParseState, []NameTypePair] {
		return frt.Pipe(parseCaseDefs(ps3), (func(_r0 frt.Tuple2[ParseState, []NameTypePair]) frt.Tuple2[ParseState, []NameTypePair] {
			return MapR((func(_r0 []NameTypePair) []NameTypePair { return slice.PushHead(cs, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []NameTypePair] {
		return frt.NewTuple2(ps2, ([]NameTypePair{cs}))
	}))
}

func parseUnionDef(tname string, pnames []string, ps0 ParseState) frt.Tuple2[ParseState, UnionDef] {
	ps := psPushScope(ps0)
	psRegTypeVars(ps, pnames)
	ps2, css := frt.Destr2(frt.Pipe(parseCaseDefs(ps), (func(_r0 frt.Tuple2[ParseState, []NameTypePair]) frt.Tuple2[ParseState, []NameTypePair] {
		return MapL(psPopScope, _r0)
	})))
	ud := NewUnionDef(tname, pnames, css)
	psRegUdToTDCtx(ud, ps2)
	return frt.NewTuple2(ps2, ud)
}

func emptyDefStmt() DefStmt {
	return frt.Pipe(frt.Empty[RecordDef](), New_DefStmt_DRecordDef)
}

func parseIdList(ps ParseState) frt.Tuple2[ParseState, []string] {
	return ParseSepList(psIdentNameNx, New_TokenType_COMMA, ps)
}

func mightParseIdList(ps ParseState) frt.Tuple2[ParseState, []string] {
	return frt.IfElse(psCurIs(New_TokenType_LT, ps), (func() frt.Tuple2[ParseState, []string] {
		return frt.Pipe(frt.Pipe(psConsume(New_TokenType_LT, ps), parseIdList), (func(_r0 frt.Tuple2[ParseState, []string]) frt.Tuple2[ParseState, []string] {
			return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_GT, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []string] {
		return frt.Pipe([]string{}, (func(_r0 []string) frt.Tuple2[ParseState, []string] { return PairL(ps, _r0) }))
	}))
}

func parseTypeDefBody(ps ParseState) frt.Tuple2[ParseState, DefStmt] {
	ps2, tname := frt.Destr2(psIdentNameNxL(ps))
	ps3, pnames := frt.Destr2(frt.Pipe(frt.Pipe(mightParseIdList(ps2), (func(_r0 frt.Tuple2[ParseState, []string]) frt.Tuple2[ParseState, []string] {
		return MapL((func(_r0 ParseState) ParseState { return psConsume(New_TokenType_EQ, _r0) }), _r0)
	})), (func(_r0 frt.Tuple2[ParseState, []string]) frt.Tuple2[ParseState, []string] {
		return MapL(psSkipEOL, _r0)
	})))
	switch (psCurrentTT(ps3)).(type) {
	case TokenType_LBRACE:
		return frt.Pipe(parseRecordDef(tname, pnames, ps3), (func(_r0 frt.Tuple2[ParseState, RecordDef]) frt.Tuple2[ParseState, DefStmt] {
			return MapR(New_DefStmt_DRecordDef, _r0)
		}))
	case TokenType_BAR:
		return frt.Pipe(parseUnionDef(tname, pnames, ps3), (func(_r0 frt.Tuple2[ParseState, UnionDef]) frt.Tuple2[ParseState, DefStmt] {
			return MapR(New_DefStmt_DUnionDef, _r0)
		}))
	default:
		psPanic(ps3, "NYI")
		return frt.NewTuple2(ps3, emptyDefStmt())
	}
}

func parseTypeDefBodyList(ps ParseState) frt.Tuple2[ParseState, []DefStmt] {
	ps2, df := frt.Destr2(frt.Pipe(parseTypeDefBody(ps), (func(_r0 frt.Tuple2[ParseState, DefStmt]) frt.Tuple2[ParseState, DefStmt] { return MapL(psSkipEOL, _r0) })))
	return frt.IfElse(psCurIs(New_TokenType_AND, ps2), (func() frt.Tuple2[ParseState, []DefStmt] {
		return frt.Pipe(frt.Pipe(psConsume(New_TokenType_AND, ps2), parseTypeDefBodyList), (func(_r0 frt.Tuple2[ParseState, []DefStmt]) frt.Tuple2[ParseState, []DefStmt] {
			return MapR((func(_r0 []DefStmt) []DefStmt { return slice.PushHead(df, _r0) }), _r0)
		}))
	}), (func() frt.Tuple2[ParseState, []DefStmt] {
		return frt.NewTuple2(ps2, ([]DefStmt{df}))
	}))
}

func parseTypeDef(ps ParseState) frt.Tuple2[ParseState, RootStmt] {
	ps2, defList := frt.Destr2(frt.Pipe(frt.Pipe(frt.Pipe(frt.Pipe(frt.Pipe(psEnterTypeDef(ps), psPushScope), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_TYPE, _r0) })), parseTypeDefBodyList), (func(_r0 frt.Tuple2[ParseState, []DefStmt]) frt.Tuple2[ParseState, []DefStmt] {
		return MapL(psPopScope, _r0)
	})), (func(_r0 frt.Tuple2[ParseState, []DefStmt]) frt.Tuple2[ParseState, []DefStmt] {
		return MapL(psLeaveTypeDef, _r0)
	})))
	nmdefs := frt.Pipe(MultipleDefs{Defs: defList}, (func(_r0 MultipleDefs) MultipleDefs { return resolveFwrdDecl(ps2, _r0) }))
	psRegMdTypes(nmdefs, ps2)
	return frt.Pipe(frt.Pipe(nmdefs, New_RootStmt_RSMultipleDefs), (func(_r0 RootStmt) frt.Tuple2[ParseState, RootStmt] { return PairL(ps2, _r0) }))
}

func parseExtTypeDef(pi PackageInfo, ps ParseState) ParseState {
	ps2, tname := frt.Destr2(frt.Pipe(psConsume(New_TokenType_TYPE, ps), psIdentNameNx))
	ps3, pnames := frt.Destr2(mightParseIdList(ps2))
	tfd := piRegEType(pi, tname, pnames)
	scRegTFData(ps3.scope, tname, tfd)
	return ps3
}

func parseExtFuncDef(pi PackageInfo, ps ParseState) ParseState {
	ps2, fname := frt.Destr2(frt.Pipe(psConsume(New_TokenType_LET, ps), psIdentNameNx))
	ps3, tnames := frt.Destr2(mightParseIdList(ps2))
	psRegTypeVars(ps3, tnames)
	ps4, fts := frt.Destr2(frt.Pipe(psConsume(New_TokenType_COLON, ps3), (func(_r0 ParseState) frt.Tuple2[ParseState, []FType] { return parseTypeArrows(parseType, _r0) })))
	ff := FuncFactory{Tparams: tnames, Targets: fts}
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
		psPanic(ps, "Unknown pkginfo def")
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
	ps3, pkgName := frt.Destr2(frt.IfElse(frt.OpEqual(psCurrentTT(ps2), New_TokenType_UNDER_SCORE), (func() frt.Tuple2[ParseState, string] {
		return frt.Pipe(frt.NewTuple2(ps2, "_"), (func(_r0 frt.Tuple2[ParseState, string]) frt.Tuple2[ParseState, string] { return MapL(psNext, _r0) }))
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
	psForErrMsg(ps)
	frt.IfOnly((SCLen(ps.scope) > 1), (func() {
		psPanic(ps, "Scope is not property poped.")
	}))
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
		psPanic(ps, "Unknown stmt")
		return parsePackage(ps)
	}
}

func parseRootOneStmtSk(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, RootStmt] {
	return frt.Pipe(parseRootOneStmt(pExpr, ps), (func(_r0 frt.Tuple2[ParseState, RootStmt]) frt.Tuple2[ParseState, RootStmt] {
		return MapL(psSkipEOL, _r0)
	}))
}

func psIsRootStmtsEnd(ps ParseState) bool {
	return frt.OpEqual(psCurrentTT(ps), New_TokenType_EOF)
}

func parseRootStmts(pExpr func(ParseState) frt.Tuple2[ParseState, Expr], ps ParseState) frt.Tuple2[ParseState, []RootStmt] {
	return frt.Pipe(psSkipEOL(ps), (func(_r0 ParseState) frt.Tuple2[ParseState, []RootStmt] {
		return ParseList((func(_r0 ParseState) frt.Tuple2[ParseState, RootStmt] { return parseRootOneStmtSk(pExpr, _r0) }), psIsRootStmtsEnd, _r0)
	}))
}
