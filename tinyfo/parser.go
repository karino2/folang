package main

import (
	"bytes"
)

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	SPACE
	IDENTIFIER
	EQ
	LET
	TYPE
	EOL
	PACKAGE
	IMPORT
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LSBRACKET
	RSBRACKET
	LT
	GT
	STRING
	COLON
	SEMICOLON
	INT_IMM
	OF
	BAR
	RARROW
	UNDER_SCORE
	MATCH
	WITH
	/*
		INDENT
		UNDENT

	*/
)

var keywordMap = map[string]TokenType{
	"let":     LET,
	"package": PACKAGE,
	"import":  IMPORT,
	"type":    TYPE,
	"of":      OF,
	"_":       UNDER_SCORE,
	"match":   MATCH,
	"with":    WITH,
}

type Token struct {
	ttype     TokenType
	begin     int
	len       int
	stringVal string // also used in Identifier.
	intVal    int
}

func NewToken(ttype TokenType, begin int, len int) *Token {
	return &Token{
		ttype: ttype,
		begin: begin,
		len:   len,
	}
}

// exclusive
func (tk *Token) end() int { return tk.begin + tk.len }

type Tokenizer struct {
	buf          []byte
	fileName     string
	pos          int
	currentToken *Token
	col          int
}

func NewTokenizer(fileName string, buf []byte) *Tokenizer {
	return &Tokenizer{
		buf:      buf,
		fileName: fileName,
	}
}

func isAlpha(b byte) bool {
	return 'a' <= b && b <= 'z' ||
		'A' <= b && b <= 'Z'
}

func isNumber(b byte) bool {
	return '0' <= b && b <= '9'
}

func isAlnum(b byte) bool {
	return isAlpha(b) || isNumber(b)
}

func (tkz *Tokenizer) analyzeCurAsIdentifier() {
	cur := tkz.currentToken
	cur.ttype = IDENTIFIER
	i := 1
	for {
		// reach end.
		if tkz.pos+i == len(tkz.buf) {
			cur.len = i
			break
		}

		c := tkz.buf[tkz.pos+i]
		if isAlnum(c) || c == '_' {
			i++
			continue
		}

		cur.len = i
		break
	}
	cur.stringVal = string(tkz.buf[cur.begin:cur.end()])
}

func (tkz *Tokenizer) isCharAt(at int, ch byte) bool {
	if at >= len(tkz.buf) {
		return false
	}
	return tkz.buf[at] == ch
}

func (tkz *Tokenizer) analyzeCurAsStringLiteral() {
	cur := tkz.currentToken
	cur.ttype = STRING
	i := 1
	var buf bytes.Buffer
	for {
		// reach end without close. parse error. panic for a while.
		if tkz.pos+i == len(tkz.buf) {
			panic("unclosed string literal")
		}

		c := tkz.buf[tkz.pos+i]

		if c == '"' {
			cur.len = i + 1
			cur.stringVal = buf.String()
			return
		} else if c == '\\' {
			i++
			if tkz.pos+i == len(tkz.buf) {
				panic("escape just before EOF, wrong")
			}
			c2 := tkz.buf[tkz.pos+i]
			switch c2 {
			case 'n':
				buf.WriteByte('\n')
			default:
				buf.WriteByte(c2)
			}
		} else {
			buf.WriteByte(c)
		}
		i++
	}
}

func (tkz *Tokenizer) analyzeCurAsIntImm() {
	cur := tkz.currentToken
	cur.ttype = INT_IMM
	c := tkz.buf[tkz.pos]
	n := 0
	i := 0
	for isNumber(c) {
		n = 10*n + int(c-'0')
		i++
		c = tkz.buf[tkz.pos+i]
	}

	cur.len = i
	cur.intVal = n
}

func (tk *Token) setOneChar(tid TokenType, one byte) {
	tk.ttype = tid
	tk.stringVal = string(one)
	tk.len = 1
}

/*
Identify token of current pos.
*/
func (tkz *Tokenizer) analyzeCur() {
	if tkz.pos == len(tkz.buf) {
		tkz.currentToken = NewToken(EOF, tkz.pos, 0)
		return
	}

	b := tkz.buf[tkz.pos]
	tkz.currentToken = NewToken(ILLEGAL, tkz.pos, 0)
	cur := tkz.currentToken

	switch {
	case b == ' ':
		cur.ttype = SPACE
		i := 1
		for ; tkz.pos+i < len(tkz.buf) && tkz.buf[tkz.pos+i] == ' '; i++ {
		}
		cur.len = i
	case 'a' <= b && b <= 'z' ||
		'A' <= b && b <= 'Z' ||
		b == '_':
		tkz.analyzeCurAsIdentifier()

		// check whether identifier is keyword
		if tt, ok := keywordMap[cur.stringVal]; ok {
			cur.ttype = tt
		}
	case isNumber(b):
		tkz.analyzeCurAsIntImm()
	case b == '"':
		tkz.analyzeCurAsStringLiteral()
	case b == '=':
		cur.setOneChar(EQ, b)
	case b == '\n':
		cur.setOneChar(EOL, b)
	case b == '(':
		cur.setOneChar(LPAREN, b)
	case b == ')':
		cur.setOneChar(RPAREN, b)
	case b == '{':
		cur.setOneChar(LBRACE, b)
	case b == '}':
		cur.setOneChar(RBRACE, b)
	case b == '[':
		cur.setOneChar(LSBRACKET, b)
	case b == ']':
		cur.setOneChar(RSBRACKET, b)
	case b == ':':
		cur.setOneChar(COLON, b)
	case b == ';':
		cur.setOneChar(SEMICOLON, b)
	case b == '|':
		cur.setOneChar(BAR, b)
	case b == '<':
		cur.setOneChar(LT, b)
	case b == '>':
		cur.setOneChar(GT, b)
	case b == '-':
		if tkz.isCharAt(tkz.pos+1, '>') {
			cur.ttype = RARROW
			cur.stringVal = "->"
			cur.len = 2
			return
		}
		panic("NYI")
	default:
		panic(b)
	}

}

func (tkz *Tokenizer) Col() int        { return tkz.col }
func (tkz *Tokenizer) Current() *Token { return tkz.currentToken }
func (tkz *Tokenizer) Setup()          { tkz.analyzeCur() }
func (tkz *Tokenizer) GotoNext() {
	if tkz.currentToken.ttype == EOF {
		return
	}
	if tkz.currentToken.ttype == EOL {
		tkz.col = 0
	} else {
		tkz.col += tkz.currentToken.len
	}
	tkz.pos = tkz.currentToken.end()
	tkz.analyzeCur()
}

func (tkz *Tokenizer) RevertTo(token *Token) {
	tkz.pos = token.begin
	tkz.currentToken = token
}

type Parser struct {
	tokenizer  *Tokenizer
	offsideCol []int
}

func (p *Parser) Current() *Token {
	return p.tokenizer.Current()
}

func (p *Parser) currentCol() int {
	return p.tokenizer.Col()
}

func (p *Parser) pushOffside() {
	if p.currentOffside() >= p.currentCol() {
		panic("Overrun offside line.")
	}
	p.offsideCol = append(p.offsideCol, p.currentCol())
}

func (p *Parser) popOffside() {
	p.offsideCol = p.offsideCol[0 : len(p.offsideCol)-1]
}

func (p *Parser) currentOffside() int {
	if len(p.offsideCol) == 0 {
		return 0
	}
	return p.offsideCol[len(p.offsideCol)-1]
}

func (p *Parser) skipSpace() {
	tk := p.Current()
	for tk.ttype == SPACE {
		p.tokenizer.GotoNext()
		tk = p.Current()
	}
}

func (p *Parser) PeekNext() *Token {
	tk := p.Current()
	defer p.tokenizer.RevertTo(tk)

	p.gotoNext()
	return p.Current()
}

func (p *Parser) gotoNext() {
	p.tokenizer.GotoNext()
	p.skipSpace()
}

// goto next with Skip EOL
func (p *Parser) gotoNextSL() {
	p.gotoNext()
	p.skipEOL()
}

func (p *Parser) skipEOL() {
	tk := p.Current()
	for tk.ttype == EOL {
		p.gotoNext()
		tk = p.Current()
	}
}

func (p *Parser) skipEOLOne() {
	tk := p.Current()
	if tk.ttype == EOL {
		p.gotoNext()
	}
}

func (p *Parser) expect(ttype TokenType) {
	tk := p.Current()
	if tk.ttype != ttype {
		panic(tk)
	}
}

func (p *Parser) consume(ttype TokenType) {
	tk := p.Current()
	if tk.ttype != ttype {
		panic(tk)
	}
	p.gotoNext()
}

// consume with skip EOL
func (p *Parser) consumeSL(ttype TokenType) {
	p.consume(ttype)
	p.skipEOL()
}

func (p *Parser) parsePackage() *Package {
	p.consume(PACKAGE)
	ident := p.Current()
	if ident.ttype != IDENTIFIER {
		panic(ident)
	}
	p.gotoNextSL()
	return &Package{ident.stringVal}
}

func (p *Parser) parseImport() *Import {
	p.consume(IMPORT)
	tk := p.Current()
	if tk.ttype != STRING {
		panic(tk)
	}
	p.gotoNextSL()
	return &Import{tk.stringVal}
}

/*
GO_EVAL = 'GoEval' string | 'GoEval' '<' TYPE '>' string

It should not have space between 'GoEval' and '<' though currently we just don't care those differences.
*/
func (p *Parser) parseGoEval() Expr {
	p.gotoNext()
	arg := p.Current()
	switch arg.ttype {
	case STRING:
		p.gotoNext()
		return NewGoEval(arg.stringVal)
	case LT:
		p.gotoNext()
		ft := p.parseType()
		p.consume(GT)
		arg = p.Current()
		if arg.ttype != STRING {
			panic(arg)
		}
		p.gotoNext()
		return &GoEval{arg.stringVal, ft}
	default:
		panic(arg)
	}
}

func (p *Parser) isEndOfExpr() bool {
	ttype := p.Current().ttype
	return ttype == EOL ||
		ttype == EOF ||
		ttype == SEMICOLON ||
		ttype == RBRACE ||
		ttype == RPAREN ||
		ttype == WITH
}

/*
RECORD_EXPRESISONN = '{' FIELD_INITIALIZERS '}'

FIELD_INITIALIZERS = FIELD_INITIALIZER (';' FIELD_INITIALIZER)*

FIELD_INITIALIZER = IDENTIFIER '=' expr
*/
func (p *Parser) parseRecordGen() Expr {
	p.consumeSL(LBRACE)

	var fnames []string
	var fvals []Expr
	for i := 0; p.Current().ttype != RBRACE; i++ {
		if i != 0 {
			p.consumeSL(SEMICOLON)
		}

		p.expect(IDENTIFIER)
		fnames = append(fnames, p.Current().stringVal)

		p.gotoNextSL()
		p.consumeSL(EQ)
		one := p.parseExpr()
		fvals = append(fvals, one)
	}

	p.consume(RBRACE)
	return NewRecordGen(fnames, fvals)
}

/*
Grammar says Application Expression is expr expr.
So create singleExpr parser that does not handle application expressions.
*/
func (p *Parser) parseSingleExpr() Expr {
	tk := p.Current()

	switch {
	case tk.ttype == STRING:
		p.gotoNext()
		return &StringLiteral{tk.stringVal}
	case tk.ttype == INT_IMM:
		p.gotoNext()
		return &IntImm{tk.intVal}
	case tk.ttype == IDENTIFIER:
		v := &Var{tk.stringVal, &FUnresolved{}}
		p.gotoNext()
		return v
	case tk.ttype == LBRACE:
		return p.parseRecordGen()
	default:
		panic("NYI")
	}
}

/*
MATCH_RULE = '|' IDENTIFIER IDENTIFIER '->' BLOCK

	| '|' '_' '->' BLOCK

ex:

	| Record r -> hoge r
*/
func (p *Parser) parseMatchRule() *MatchRule {
	p.consume(BAR)

	if p.Current().ttype == UNDER_SCORE {
		// default case.
		p.gotoNext()
		p.consume(RARROW)
		p.skipEOLOne()
		p.skipSpace()
		block := p.parseBlock()
		return &MatchRule{&MatchPattern{"_", ""}, block}
	}

	p.expect(IDENTIFIER)
	caseName := p.Current().stringVal
	p.gotoNext()
	var varName string
	if p.Current().ttype == RARROW {
		// no content case. use "" for varName.
	} else {
		if p.Current().ttype == UNDER_SCORE {
			varName = "_"
		} else {
			p.expect(IDENTIFIER)
			varName = p.Current().stringVal
		}
		p.gotoNext()
	}
	p.consume(RARROW)
	p.skipEOLOne()
	p.skipSpace()
	block := p.parseBlock()
	return &MatchRule{&MatchPattern{caseName, varName}, block}
}

/*
MATCH_EXPR = 'match' EXPR 'with' EOL MATCH_RULE+
*/
func (p *Parser) parseMatchExpr() Expr {
	p.consume(MATCH)
	target := p.parseExpr()
	p.consume(WITH)
	p.skipEOLOne()
	p.skipSpace()

	var rules []*MatchRule
	for p.Current().ttype == BAR {
		one := p.parseMatchRule()
		rules = append(rules, one)
	}
	return &MatchExpr{target, rules}
}

func (p *Parser) parseExpr() Expr {
	tk := p.Current()
	// spcial handling for a while.
	if tk.ttype == IDENTIFIER && tk.stringVal == "GoEval" {
		return p.parseGoEval()
	}

	if tk.ttype == MATCH {
		return p.parseMatchExpr()
	}

	var exprs []Expr
	for !p.isEndOfExpr() {
		expr := p.parseSingleExpr()
		exprs = append(exprs, expr)
	}

	if len(exprs) == 1 {
		return exprs[0]
	}

	v, ok := exprs[0].(*Var)
	if !ok {
		panic("application expr with non variable start, NYI")
	}

	return &FunCall{v, exprs[1:]}
}

func (p *Parser) isEndOfBlock() bool {
	p.skipSpace()
	return p.currentCol() < p.currentOffside()
}

/*
BLOCK = EXPR | (STMT_LIKE EOL)* EXPR

STMT_LIKE = LET_STMT | EXPR
*/
func (p *Parser) parseBlock() *Block {
	var stmts []Stmt

	p.skipSpace()

	// limited implementation of offside rule.
	p.pushOffside()

	expr := p.parseExpr()
	p.skipEOLOne()

	if p.isEndOfBlock() {
		p.popOffside()
		return &Block{stmts, expr}
	}

	for !p.isEndOfBlock() {
		stmts = append(stmts, &ExprStmt{expr})
		expr = p.parseExpr()
		p.skipEOLOne()
	}
	p.popOffside()

	return &Block{stmts, expr}
}

/*
TYPE = 'string' | 'int' | IDENTIFIER | '[”]' TYPE
*/
func (p *Parser) parseType() FType {

	if p.Current().ttype == LSBRACKET {
		p.consume(LSBRACKET)
		p.consume(RSBRACKET)
		etype := p.parseType()
		return &FSlice{etype}
	}

	p.expect(IDENTIFIER)
	tname := p.Current().stringVal
	switch tname {
	case "string":
		p.gotoNext()
		return FString
	case "int":
		p.gotoNext()
		return FInt
	default:
		p.gotoNext()
		return &FCustom{tname}
	}
}

/*
	PARAM = '(' ')'
				|	'(' IDENTIFIER : TYPE ')'

for '(' ')' case, return nil
*/
func (p *Parser) parseParam() *Var {
	p.consume(LPAREN)
	if p.Current().ttype == RPAREN {
		p.consume(RPAREN)
		return nil
	}
	p.expect(IDENTIFIER)
	varName := p.Current().stringVal

	p.gotoNext()
	p.consume(COLON)
	ft := p.parseType()
	p.consume(RPAREN)
	return &Var{varName, ft}
}

/*
PARAMS = PARAM*
*/
func (p *Parser) parseParams() []*Var {
	var params []*Var

	one := p.parseParam()
	if one == nil {
		return params
	}

	params = append(params, one)

	for p.Current().ttype == LPAREN {
		one = p.parseParam()
		params = append(params, one)
	}
	return params
}

/*
FUNC_DEF_LET = 'let' IDENTIFIER PARAMS '=' EOL BLOCK
*/
func (p *Parser) parseFuncDefLet() Stmt {
	p.consume(LET)
	fname := p.Current()
	if fname.ttype != IDENTIFIER {
		panic(fname)
	}
	p.gotoNext()
	params := p.parseParams()

	p.consumeSL(EQ)

	block := p.parseBlock()

	return &FuncDef{fname.stringVal, params, block}
}

func (p *Parser) parseLet() Stmt {
	return p.parseFuncDefLet()
}

/*
UNION_DEF = CASE_DEF (EOL CASE_DEF)* EOL

CASE_DEF = '|' IDENIFIER OF TYPE
*/
func (p *Parser) parseUnionDef(uname string) Stmt {
	var cases []NameTypePair

	for p.Current().ttype == BAR {
		p.gotoNext()
		p.expect(IDENTIFIER)
		cname := p.Current().stringVal

		p.gotoNext()
		if p.Current().ttype == OF {
			p.consume(OF)

			tp := p.parseType()
			cases = append(cases, NameTypePair{cname, tp})
			p.consume(EOL)
		} else {
			// no "of", unit case.
			cases = append(cases, NameTypePair{cname, FUnit})
			p.consume(EOL)
		}
	}
	return &UnionDef{uname, cases}
}

// RECORD_DEF = '{' FIELD_DEFS '}'
//
// FIELD_DEFS = FIELD_DEF  (';' FIELD_DEF)*
//
// FIELD_DEF = ID ':' TYPE
func (p *Parser) parseRecordDef(rname string) Stmt {
	p.consumeSL(LBRACE)

	var fields []NameTypePair
	for i := 0; p.Current().ttype != RBRACE; i++ {
		if i != 0 {
			p.consumeSL(SEMICOLON)
		}

		p.expect(IDENTIFIER)
		fname := p.Current().stringVal

		p.gotoNextSL()
		p.consumeSL(COLON)
		ftype := p.parseType()
		fields = append(fields, NameTypePair{fname, ftype})
	}
	p.consume(RBRACE)

	return &RecordDef{rname, fields}
}

// TYPE_DEF = 'type' ID '=' (RECORD_DEF | UNION_DEF)
//
// RECORD_DEF = '{' FIELD_DEFS '}'
//
// UNION_DEF = '|'...
func (p *Parser) parseTypeDef() Stmt {
	p.consume(TYPE)
	p.expect(IDENTIFIER)

	tname := p.Current().stringVal
	p.gotoNextSL()

	p.consumeSL(EQ)

	if p.Current().ttype == LBRACE {
		return p.parseRecordDef(tname)
	}

	p.expect(BAR)
	return p.parseUnionDef(tname)
}

func (p *Parser) parseStmt() Stmt {
	tk := p.Current()
	switch tk.ttype {
	case PACKAGE:
		return p.parsePackage()
	case IMPORT:
		return p.parseImport()
	case LET:
		return p.parseLet()
	case TYPE:
		return p.parseTypeDef()
	default:
		panic(tk)
	}
}

func (p *Parser) parseStmts() []Stmt {
	var stmts []Stmt
	p.skipEOL()
	for p.Current().ttype != EOF {
		one := p.parseStmt()
		stmts = append(stmts, one)
		p.skipEOL()
	}

	return stmts
}

func (p *Parser) Setup(fileName string, buf []byte) {
	p.tokenizer = NewTokenizer(fileName, buf)
	p.tokenizer.Setup()
}

func (p *Parser) Parse(fileName string, buf []byte) []Stmt {
	p.Setup(fileName, buf)
	return p.parseStmts()
}
