package main

import "github.com/karino2/folang/pkg/frt"

type ParseState struct {
	tkz Tokenizer
}

func initParse(src string) ParseState {
	tkz := newTkz(src)
	return ParseState{tkz: tkz}
}

func currentTk(ps ParseState) Token {
	return ps.tkz.current
}

func psNext(ps ParseState) ParseState {
	ntk := tkzNext(ps.tkz)
	return ParseState{tkz: ntk}
}

func psNextNOL(ps ParseState) ParseState {
	ntk := tkzNextNOL(ps.tkz)
	return ParseState{tkz: ntk}
}

func psExpect(ps ParseState, ttype TokenType) {
	cur := currentTk(ps)
	frt.IfOnly(frt.OpNotEqual(cur.ttype, ttype), (func() {
		frt.Panic("non expected token")
	}))

}

func psConsume(ps ParseState, ttype TokenType) ParseState {
	psExpect(ps, ttype)
	return psNext(ps)
}

func psIdentName(ps ParseState) string {
	psExpect(ps, New_TokenType_IDENTIFIER)
	cur := currentTk(ps)
	return cur.stringVal
}

func parsePackage(ps ParseState) frt.Tuple2[ParseState, Stmt] {
	nps := psConsume(ps, New_TokenType_PACKAGE)
	pname := psIdentName(nps)
	nnps := psNextNOL(nps)
	pkg := New_Stmt_Package(pname)
	return frt.NewTuple2(nnps, pkg)
}
