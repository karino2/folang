package main

import "github.com/karino2/folang/pkg/frt"

type TokenType interface {
	TokenType_Union()
}

func (TokenType_ILLEGAL) TokenType_Union()      {}
func (TokenType_EOF) TokenType_Union()          {}
func (TokenType_SPACE) TokenType_Union()        {}
func (TokenType_IDENTIFIER) TokenType_Union()   {}
func (TokenType_EQ) TokenType_Union()           {}
func (TokenType_LET) TokenType_Union()          {}
func (TokenType_FUN) TokenType_Union()          {}
func (TokenType_TYPE) TokenType_Union()         {}
func (TokenType_EOL) TokenType_Union()          {}
func (TokenType_PACKAGE) TokenType_Union()      {}
func (TokenType_IMPORT) TokenType_Union()       {}
func (TokenType_LPAREN) TokenType_Union()       {}
func (TokenType_RPAREN) TokenType_Union()       {}
func (TokenType_LBRACE) TokenType_Union()       {}
func (TokenType_RBRACE) TokenType_Union()       {}
func (TokenType_LSBRACKET) TokenType_Union()    {}
func (TokenType_RSBRACKET) TokenType_Union()    {}
func (TokenType_LT) TokenType_Union()           {}
func (TokenType_GT) TokenType_Union()           {}
func (TokenType_LE) TokenType_Union()           {}
func (TokenType_GE) TokenType_Union()           {}
func (TokenType_BRACKET) TokenType_Union()      {}
func (TokenType_PIPE) TokenType_Union()         {}
func (TokenType_STRING) TokenType_Union()       {}
func (TokenType_COLON) TokenType_Union()        {}
func (TokenType_COMMA) TokenType_Union()        {}
func (TokenType_SEMICOLON) TokenType_Union()    {}
func (TokenType_INT_IMM) TokenType_Union()      {}
func (TokenType_OF) TokenType_Union()           {}
func (TokenType_BAR) TokenType_Union()          {}
func (TokenType_BARBAR) TokenType_Union()       {}
func (TokenType_RARROW) TokenType_Union()       {}
func (TokenType_UNDER_SCORE) TokenType_Union()  {}
func (TokenType_MATCH) TokenType_Union()        {}
func (TokenType_WITH) TokenType_Union()         {}
func (TokenType_TRUE) TokenType_Union()         {}
func (TokenType_FALSE) TokenType_Union()        {}
func (TokenType_PACKAGE_INFO) TokenType_Union() {}
func (TokenType_DOT) TokenType_Union()          {}
func (TokenType_AND) TokenType_Union()          {}
func (TokenType_AMP) TokenType_Union()          {}
func (TokenType_AMPAMP) TokenType_Union()       {}
func (TokenType_PLUS) TokenType_Union()         {}
func (TokenType_MINUS) TokenType_Union()        {}
func (TokenType_ASTER) TokenType_Union()        {}
func (TokenType_SLASH) TokenType_Union()        {}
func (TokenType_IF) TokenType_Union()           {}
func (TokenType_THEN) TokenType_Union()         {}
func (TokenType_ELSE) TokenType_Union()         {}
func (TokenType_ELIF) TokenType_Union()         {}
func (TokenType_NOT) TokenType_Union()          {}

type TokenType_ILLEGAL struct {
}

var New_TokenType_ILLEGAL TokenType = TokenType_ILLEGAL{}

type TokenType_EOF struct {
}

var New_TokenType_EOF TokenType = TokenType_EOF{}

type TokenType_SPACE struct {
}

var New_TokenType_SPACE TokenType = TokenType_SPACE{}

type TokenType_IDENTIFIER struct {
}

var New_TokenType_IDENTIFIER TokenType = TokenType_IDENTIFIER{}

type TokenType_EQ struct {
}

var New_TokenType_EQ TokenType = TokenType_EQ{}

type TokenType_LET struct {
}

var New_TokenType_LET TokenType = TokenType_LET{}

type TokenType_FUN struct {
}

var New_TokenType_FUN TokenType = TokenType_FUN{}

type TokenType_TYPE struct {
}

var New_TokenType_TYPE TokenType = TokenType_TYPE{}

type TokenType_EOL struct {
}

var New_TokenType_EOL TokenType = TokenType_EOL{}

type TokenType_PACKAGE struct {
}

var New_TokenType_PACKAGE TokenType = TokenType_PACKAGE{}

type TokenType_IMPORT struct {
}

var New_TokenType_IMPORT TokenType = TokenType_IMPORT{}

type TokenType_LPAREN struct {
}

var New_TokenType_LPAREN TokenType = TokenType_LPAREN{}

type TokenType_RPAREN struct {
}

var New_TokenType_RPAREN TokenType = TokenType_RPAREN{}

type TokenType_LBRACE struct {
}

var New_TokenType_LBRACE TokenType = TokenType_LBRACE{}

type TokenType_RBRACE struct {
}

var New_TokenType_RBRACE TokenType = TokenType_RBRACE{}

type TokenType_LSBRACKET struct {
}

var New_TokenType_LSBRACKET TokenType = TokenType_LSBRACKET{}

type TokenType_RSBRACKET struct {
}

var New_TokenType_RSBRACKET TokenType = TokenType_RSBRACKET{}

type TokenType_LT struct {
}

var New_TokenType_LT TokenType = TokenType_LT{}

type TokenType_GT struct {
}

var New_TokenType_GT TokenType = TokenType_GT{}

type TokenType_LE struct {
}

var New_TokenType_LE TokenType = TokenType_LE{}

type TokenType_GE struct {
}

var New_TokenType_GE TokenType = TokenType_GE{}

type TokenType_BRACKET struct {
}

var New_TokenType_BRACKET TokenType = TokenType_BRACKET{}

type TokenType_PIPE struct {
}

var New_TokenType_PIPE TokenType = TokenType_PIPE{}

type TokenType_STRING struct {
}

var New_TokenType_STRING TokenType = TokenType_STRING{}

type TokenType_COLON struct {
}

var New_TokenType_COLON TokenType = TokenType_COLON{}

type TokenType_COMMA struct {
}

var New_TokenType_COMMA TokenType = TokenType_COMMA{}

type TokenType_SEMICOLON struct {
}

var New_TokenType_SEMICOLON TokenType = TokenType_SEMICOLON{}

type TokenType_INT_IMM struct {
}

var New_TokenType_INT_IMM TokenType = TokenType_INT_IMM{}

type TokenType_OF struct {
}

var New_TokenType_OF TokenType = TokenType_OF{}

type TokenType_BAR struct {
}

var New_TokenType_BAR TokenType = TokenType_BAR{}

type TokenType_BARBAR struct {
}

var New_TokenType_BARBAR TokenType = TokenType_BARBAR{}

type TokenType_RARROW struct {
}

var New_TokenType_RARROW TokenType = TokenType_RARROW{}

type TokenType_UNDER_SCORE struct {
}

var New_TokenType_UNDER_SCORE TokenType = TokenType_UNDER_SCORE{}

type TokenType_MATCH struct {
}

var New_TokenType_MATCH TokenType = TokenType_MATCH{}

type TokenType_WITH struct {
}

var New_TokenType_WITH TokenType = TokenType_WITH{}

type TokenType_TRUE struct {
}

var New_TokenType_TRUE TokenType = TokenType_TRUE{}

type TokenType_FALSE struct {
}

var New_TokenType_FALSE TokenType = TokenType_FALSE{}

type TokenType_PACKAGE_INFO struct {
}

var New_TokenType_PACKAGE_INFO TokenType = TokenType_PACKAGE_INFO{}

type TokenType_DOT struct {
}

var New_TokenType_DOT TokenType = TokenType_DOT{}

type TokenType_AND struct {
}

var New_TokenType_AND TokenType = TokenType_AND{}

type TokenType_AMP struct {
}

var New_TokenType_AMP TokenType = TokenType_AMP{}

type TokenType_AMPAMP struct {
}

var New_TokenType_AMPAMP TokenType = TokenType_AMPAMP{}

type TokenType_PLUS struct {
}

var New_TokenType_PLUS TokenType = TokenType_PLUS{}

type TokenType_MINUS struct {
}

var New_TokenType_MINUS TokenType = TokenType_MINUS{}

type TokenType_ASTER struct {
}

var New_TokenType_ASTER TokenType = TokenType_ASTER{}

type TokenType_SLASH struct {
}

var New_TokenType_SLASH TokenType = TokenType_SLASH{}

type TokenType_IF struct {
}

var New_TokenType_IF TokenType = TokenType_IF{}

type TokenType_THEN struct {
}

var New_TokenType_THEN TokenType = TokenType_THEN{}

type TokenType_ELSE struct {
}

var New_TokenType_ELSE TokenType = TokenType_ELSE{}

type TokenType_ELIF struct {
}

var New_TokenType_ELIF TokenType = TokenType_ELIF{}

type TokenType_NOT struct {
}

var New_TokenType_NOT TokenType = TokenType_NOT{}

type Token struct {
	ttype     TokenType
	begin     int
	len       int
	stringVal string
	intVal    int
}

type Tokenizer struct {
	buf     string
	current Token
	col     int
}

func newTkz(buf string) Tokenizer {
	itk := newToken(New_TokenType_ILLEGAL, 0, 0)
	ftk := nextToken(buf, itk)
	return Tokenizer{buf: buf, current: ftk, col: ftk.begin}
}

func tkzNext(tkz Tokenizer) Tokenizer {
	switch (tkz.current.ttype).(type) {
	case TokenType_EOF:
		return tkz
	case TokenType_EOL:
		nt := nextToken(tkz.buf, tkz.current)
		bol := (tkz.current.begin + tkz.current.len)
		return Tokenizer{buf: tkz.buf, current: nt, col: (nt.begin - bol)}
	default:
		nt := nextToken(tkz.buf, tkz.current)
		delta := (nt.begin - tkz.current.begin)
		ncol := (tkz.col + delta)
		return Tokenizer{buf: tkz.buf, current: nt, col: ncol}
	}
}

func tkzNextNOL(tkz Tokenizer) Tokenizer {
	ntkz := tkzNext(tkz)
	return frt.IfElse(frt.OpEqual(ntkz.current.ttype, New_TokenType_EOL), (func() Tokenizer {
		return tkzNextNOL(ntkz)
	}), (func() Tokenizer {
		return ntkz
	}))
}

func tkzIsNeighborLT(tkz Tokenizer) bool {
	return isNeighborLT(tkz.buf, tkz.current)
}
