package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/slice"

type ParseState struct {
	tkz Tokenizer
}

func initParse(src string) ParseState {
	tkz := newTkz(src)
	return ParseState{tkz: tkz}
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
	return ParseState{tkz: ntk}
}

func psNextNOL(ps ParseState) ParseState {
	ntk := tkzNextNOL(ps.tkz)
	return ParseState{tkz: ntk}
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

func parsePackage(ps ParseState) frt.Tuple2[ParseState, Stmt] {
	ps2 := psConsume(New_TokenType_PACKAGE, ps)
	pname := psIdentName(ps2)
	ps3 := psNextNOL(ps2)
	pkg := New_Stmt_Package(pname)
	return frt.NewTuple2(ps3, pkg)
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
		tpair := parseType(ps3)
		ps4 := frt.Pipe(frt.Fst(tpair), (func(_r0 ParseState) ParseState { return psConsume(New_TokenType_RPAREN, _r0) }))
		tp := frt.Snd(tpair)
		v := Var{name: vname, ftype: tp}
		return frt.NewTuple2(ps4, New_Param_PVar(v))
	}
}

func parseParams(ps ParseState) frt.Tuple2[ParseState, []Param] {
	ftp := parseParam(ps)
	ps2 := frt.Fst(ftp)
	prm1 := frt.Snd(ftp)
	switch (prm1).(type) {
	case Param_PUnit:
		zero := []Param{}
		return frt.NewTuple2(ps2, zero)
	case Param_PVar:
		tt := psCurrentTT(ps2)
		switch (tt).(type) {
		case TokenType_LPAREN:
			ftp2 := parseParams(ps2)
			ps3 := frt.Fst(ftp2)
			prms2 := frt.Snd(ftp2)
			pas3 := slice.Append(prm1, prms2)
			return frt.NewTuple2(ps3, pas3)
		default:
			return frt.NewTuple2(ps2, ([]Param{prm1}))
		}
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
