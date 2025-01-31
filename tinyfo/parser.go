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
	PIPE
	STRING
	COLON
	COMMA
	SEMICOLON
	INT_IMM
	OF
	BAR
	RARROW
	UNDER_SCORE
	MATCH
	WITH
	TRUE
	FALSE
	PACKAGE_INFO
	DOT
)

var keywordMap = map[string]TokenType{
	"let":          LET,
	"package":      PACKAGE,
	"import":       IMPORT,
	"type":         TYPE,
	"of":           OF,
	"_":            UNDER_SCORE,
	"match":        MATCH,
	"with":         WITH,
	"true":         TRUE,
	"false":        FALSE,
	"package_info": PACKAGE_INFO,
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
	case b == ',':
		cur.setOneChar(COMMA, b)
	case b == '.':
		cur.setOneChar(DOT, b)
	case b == ';':
		cur.setOneChar(SEMICOLON, b)
	case b == '|':
		if tkz.isCharAt(tkz.pos+1, '>') {
			cur.ttype = PIPE
			cur.stringVal = "|>"
			cur.len = 2
			return
		}
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

type Scope struct {
	varMap    map[string]*Var
	recordMap map[string]*FRecord
	typeMap   map[string]FType
	parent    *Scope
}

func NewScope(parent *Scope) *Scope {
	s := &Scope{}
	s.varMap = make(map[string]*Var)
	s.recordMap = make(map[string]*FRecord)
	s.typeMap = make(map[string]FType)
	s.parent = parent
	return s
}

func (s *Scope) DefineVar(name string, v *Var) {
	s.varMap[name] = v
}

func (s *Scope) LookupVar(name string) *Var {
	cur := s
	for cur != nil {
		ret, ok := cur.varMap[name]
		if ok {
			return ret
		}
		cur = cur.parent
	}
	return nil
}

func (s *Scope) LookupType(name string) FType {
	cur := s
	for cur != nil {
		ret, ok := cur.typeMap[name]
		if ok {
			return ret
		}
		cur = cur.parent
	}
	return nil
}

func (s *Scope) lookupRecordCur(fieldNames []string) *FRecord {
	for _, rt := range s.recordMap {
		if rt.Match(fieldNames) {
			return rt
		}
	}
	return nil

}

func (s *Scope) LookupRecord(fieldNames []string) *FRecord {
	cur := s
	for cur != nil {
		ret := cur.lookupRecordCur(fieldNames)
		if ret != nil {
			return ret
		}
		cur = cur.parent
	}
	return nil
}

type Parser struct {
	tokenizer  *Tokenizer
	offsideCol []int
	scope      *Scope
}

func NewParser() *Parser {
	p := &Parser{}
	p.scope = NewScope(nil)
	return p
}

func (p *Parser) pushScope() {
	p.scope = NewScope(p.scope)
}

func (p *Parser) popScope() *Scope {
	ret := p.scope
	p.scope = ret.parent
	return ret
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

func (p *Parser) peekNext() *Token {
	tk := p.Current()
	defer p.tokenizer.RevertTo(tk)

	p.gotoNext()
	return p.Current()
}

func (p *Parser) peekNextNext() *Token {
	tk := p.Current()
	defer p.tokenizer.RevertTo(tk)

	p.gotoNext()
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

func (p *Parser) identName() string {
	p.expect(IDENTIFIER)
	return p.Current().stringVal
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

func (p *Parser) nextIs(ttype TokenType) bool {
	nt := p.peekNext()
	return nt.ttype == ttype
}

func (p *Parser) currentIs(ttype TokenType) bool {
	return p.Current().ttype == ttype
}

func (p *Parser) isEndOfTerm() bool {
	ttype := p.Current().ttype
	return ttype == EOL ||
		ttype == EOF ||
		ttype == SEMICOLON ||
		ttype == RBRACE ||
		ttype == RPAREN ||
		ttype == WITH ||
		p.nextNonEOLIsBinOp()
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

		fnames = append(fnames, p.identName())

		p.gotoNextSL()
		p.consumeSL(EQ)
		one := p.parseExpr()
		fvals = append(fvals, one)
	}

	p.consume(RBRACE)
	rg := NewRecordGen(fnames, fvals)

	rtype := p.scope.LookupRecord(rg.fieldNames)
	rg.recordType = rtype

	return rg
}

/*
Grammar says Application Expression is expr expr.
So create ATOM parser that does not handle application expressions.

ATOM = LITERAL | VARIABLE | REC_GEN | '(' ')' | '(' EXPR ')'
*/
func (p *Parser) parseAtom() Expr {
	tk := p.Current()

	switch {
	case tk.ttype == STRING:
		p.gotoNext()
		return &StringLiteral{tk.stringVal}
	case tk.ttype == INT_IMM:
		p.gotoNext()
		return &IntImm{tk.intVal}
	case tk.ttype == IDENTIFIER:
		fullName := p.parseFullName()
		v := p.scope.LookupVar(fullName)
		if v == nil {
			panic("Undefined var: " + fullName)
		}
		return v
	case tk.ttype == LBRACE:
		return p.parseRecordGen()
	case tk.ttype == LPAREN:
		p.consume(LPAREN)
		if p.currentIs(RPAREN) {
			p.consume(RPAREN)
			return gUnitVal
		} else {
			ret := p.parseExpr()
			p.consume(RPAREN)
			return ret
		}
	case tk.ttype == TRUE:
		p.gotoNext()
		return &BoolLiteral{true}
	case tk.ttype == FALSE:
		p.gotoNext()
		return &BoolLiteral{false}
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

	caseName := p.identName()
	p.gotoNext()
	var varName string
	if p.Current().ttype == RARROW {
		// no content case. use "" for varName.
	} else {
		if p.Current().ttype == UNDER_SCORE {
			varName = "_"
		} else {
			varName = p.identName()
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

/*
TERM = GOEVAL | MATCH_EXPR | ATOM ATOM*
*/
func (p *Parser) parseTerm() Expr {
	tk := p.Current()
	// spcial handling for a while.
	if tk.ttype == IDENTIFIER && tk.stringVal == "GoEval" {
		return p.parseGoEval()
	}

	if tk.ttype == MATCH {
		return p.parseMatchExpr()
	}

	var exprs []Expr
	for !p.isEndOfTerm() {
		expr := p.parseAtom()
		exprs = append(exprs, expr)
	}

	if len(exprs) == 1 {
		return exprs[0]
	}

	v, ok := exprs[0].(*Var)
	if !ok {
		panic("application expr with non variable start, NYI")
	}

	fc := &FunCall{v, exprs[1:], nil}
	fc.ResolveTypeParamByArgs()
	return fc
}

func (p *Parser) peekNextNonEOLToken() *Token {
	if !p.currentIs(EOL) {
		return p.Current()
	}
	tk := p.Current()
	defer p.tokenizer.RevertTo(tk)
	p.skipEOL()
	return p.Current()
}

func (p *Parser) nextNonEOLIsBinOp() bool {
	tk := p.peekNextNonEOLToken()
	return tk.ttype == PIPE
}

/*
EXPR = TERM (BINOP EXPR)*
*/
func (p *Parser) parseExpr() Expr {
	expr := p.parseTerm()
	for p.nextNonEOLIsBinOp() {
		p.skipEOL()
		p.consume(PIPE) // currently, only pipe.
		rhs := p.parseExpr()
		expr = NewPipeCall(expr, rhs)
	}
	return expr
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
	p.pushScope()
	var stmts []Stmt
	var last Expr

	p.skipSpace()

	// limited implementation of offside rule.
	p.pushOffside()

	for {
		if p.Current().ttype == LET {
			stmts = append(stmts, p.parseLet())
			p.skipEOLOne()
		} else {
			last = p.parseExpr()
			p.skipEOLOne()
			if p.isEndOfBlock() {
				p.popOffside()
				return &Block{stmts, last, p.popScope()}
			} else {
				stmts = append(stmts, &ExprStmt{last})
			}
		}
	}
}

/*
FULL_NAME = (IDENTIFIER.)* IDENTIFIER

return concat string like "buf.WriteString"
*/
func (p *Parser) parseFullName() string {
	var buf bytes.Buffer

	last := p.identName()
	p.gotoNext()
	buf.WriteString(last)

	for p.Current().ttype == DOT {
		p.gotoNext()
		buf.WriteString(".")

		last = p.identName()
		p.gotoNext()
		buf.WriteString(last)
	}
	return buf.String()
}

/*
ATOM_TYPE = 'string' | 'int' | '(' ')' | REGISTERED_TYPE | '[' ']' TYPE
*/
func (p *Parser) parseAtomType() FType {
	if p.Current().ttype == LSBRACKET {
		p.consume(LSBRACKET)
		p.consume(RSBRACKET)
		etype := p.parseAtomType()
		return &FSlice{etype}
	}

	if p.Current().ttype == LPAREN {
		p.consume(LPAREN)
		p.consume(RPAREN)
		return FUnit
	}

	tname := p.identName()
	switch tname {
	case "string":
		p.gotoNext()
		return FString
	case "int":
		p.gotoNext()
		return FInt
	default:
		fullName := p.parseFullName()
		resT := p.scope.LookupType(fullName)
		if resT == nil {
			panic(fullName) // unknown type.
		}
		return resT
	}
}

/*
NONE_COMPOUND_FUNC_TYPE = ATOM_TYPE '->' ATOM_TYPE ('->' ATOM_TYPE)*

Pass first ATOM_TYPE as argument.
*/
func (p *Parser) parseNoneCompoundFuncType(first FType) FType {
	types := []FType{first}
	for p.Current().ttype == RARROW {
		p.consume(RARROW)
		types = append(types, p.parseAtomType())
	}
	return &FFunc{types, []string{}}
}

/*
FUNC_ELEM = ATOM_TYPE | '(' NONE_COMPOUND_FUNC_TYPE ')'

Only One nesting is supported like (a->b)->c->d resolved as func (arg1:func(a)b, arg2: c) d
*/
func (p *Parser) parseFuncElemType() FType {
	if p.Current().ttype == LPAREN {
		next := p.peekNext()
		if next.ttype == RPAREN {
			p.consume(LPAREN)
			p.consume(RPAREN)
			return FUnit
		}
		p.consume(LPAREN)
		first := p.parseAtomType()
		ft := p.parseNoneCompoundFuncType(first)
		p.consume(RPAREN)
		return ft
	}
	return p.parseAtomType()
}

/*
FUNC_TYPE = FUNC_ELEM '->' FUNC_ELEM ('->' FUNC_ELEM)*
*/
func (p *Parser) parseFuncType() *FFunc {
	elem := p.parseFuncElemType()
	types := []FType{elem}

	for p.Current().ttype == RARROW {
		p.consume(RARROW)
		types = append(types, p.parseFuncElemType())
	}
	return &FFunc{types, []string{}}
}

/*
TYPE =  ATOM_TYPE | FUNC_TYPE
*/
func (p *Parser) parseType() FType {
	if p.Current().ttype == LPAREN {
		next := p.peekNext()
		if next.ttype == RPAREN {
			p.consume(LPAREN)
			p.consume(RPAREN)
			return FUnit
		}
		return p.parseFuncType()
	}
	one := p.parseAtomType()

	if p.Current().ttype == RARROW {
		// func type
		types := []FType{one}

		for p.Current().ttype == RARROW {
			p.consume(RARROW)
			types = append(types, p.parseFuncElemType())
		}
		return &FFunc{types, []string{}}
	} else {
		return one
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
	varName := p.identName()

	p.gotoNext()
	p.consume(COLON)
	ft := p.parseType()
	p.consume(RPAREN)
	ret := &Var{varName, ft}
	p.scope.DefineVar(varName, ret)
	return ret
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

	// for recursive call, define symbol at first.
	// Use this pointer and replace type after defined.
	v := &Var{fname.stringVal, &FUnresolved{}}
	p.scope.DefineVar(fname.stringVal, v)

	p.gotoNext()
	params := p.parseParams()

	p.consumeSL(EQ)

	block := p.parseBlock()

	ret := &FuncDef{fname.stringVal, params, block}
	// update type.
	v.Type = ret.FuncFType()
	return ret
}

/*
LET_VAR_DEF = 'let' IDENTIFIER '=' expr
*/
func (p *Parser) parseLetDefVar() Stmt {
	p.consume(LET)
	vnameTk := p.Current()
	if vnameTk.ttype != IDENTIFIER {
		panic(vnameTk)
	}
	p.gotoNext()
	p.consume(EQ)

	rhs := p.parseExpr()
	vname := vnameTk.stringVal

	v := &Var{vname, rhs.FType()}
	p.scope.DefineVar(vname, v)

	return &LetVarDef{vname, rhs}
}

/*
LET = LET_VAR_DEF | LET_FUNC_DEF

LET_VAR_DEF = 'let' IDENTIFIER '=' expr
*/
func (p *Parser) parseLet() Stmt {
	nn := p.peekNextNext()
	if nn.ttype == EQ {
		return p.parseLetDefVar()
	} else {
		return p.parseFuncDefLet()
	}
}

/*
UNION_DEF = CASE_DEF (EOL CASE_DEF)* EOL

CASE_DEF = '|' IDENIFIER OF TYPE
*/
func (p *Parser) parseUnionDef(uname string) Stmt {
	var cases []NameTypePair

	for p.Current().ttype == BAR {
		p.gotoNext()
		cname := p.identName()

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
	ret := &UnionDef{uname, cases}
	ret.registerToScope(p.scope)
	return ret
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

		fname := p.identName()

		p.gotoNextSL()
		p.consumeSL(COLON)
		ftype := p.parseType()
		fields = append(fields, NameTypePair{fname, ftype})
	}
	p.consume(RBRACE)

	rd := &RecordDef{rname, fields}

	recType := rd.ToFType()
	p.scope.recordMap[rname] = recType
	p.scope.typeMap[rname] = recType

	return rd
}

// TYPE_DEF = 'type' ID '=' (RECORD_DEF | UNION_DEF)
//
// RECORD_DEF = '{' FIELD_DEFS '}'
//
// UNION_DEF = '|'...
func (p *Parser) parseTypeDef() Stmt {
	p.consume(TYPE)

	tname := p.identName()
	p.gotoNextSL()

	p.consumeSL(EQ)

	if p.Current().ttype == LBRACE {
		return p.parseRecordDef(tname)
	}

	p.expect(BAR)
	return p.parseUnionDef(tname)
}

/*
EXT_TYPE_DEF = 'type' IDENTIFIER
*/
func (p *Parser) parseExtTypeDef(pi *PackageInfo) {
	p.consume(TYPE)
	typeName := p.identName()
	p.gotoNext()
	tp := pi.registerExtType(typeName)
	p.scope.typeMap[typeName] = tp
}

/*
TYPE_PARAM = '<' IDENTIFIER (',' IDENTIFIER)* '>'

Register generic type and return list of type param names.
*/
func (p *Parser) parseTypeParam() []string {
	var res []string
	p.consume(LT)
	ident := p.identName()
	res = append(res, ident)
	p.gotoNext()
	p.scope.typeMap[ident] = &FParametrized{ident}
	for p.Current().ttype == COMMA {
		p.gotoNext()
		ident = p.identName()
		res = append(res, ident)
		p.gotoNext()
		p.scope.typeMap[ident] = &FParametrized{ident}
	}
	p.consume(GT)
	return res
}

/*
EXT_FUNC_DEF = 'let' IDENTIFIER (TYPE_PARAM)? ':' FUNC_TYPE
*/
func (p *Parser) parseExtFuncDef(pi *PackageInfo) {
	p.consume(LET)
	funcName := p.identName()
	p.gotoNext()

	var tps []string

	if p.Current().ttype == LT {
		tps = p.parseTypeParam()
	}

	p.consume(COLON)
	ft := p.parseFuncType()
	ft.TypeParams = tps
	pi.funcInfo[funcName] = ft
	p.scope.varMap[funcName] = &Var{funcName, ft}
}

func (p *Parser) parseExtDef(pi *PackageInfo) {
	switch p.Current().ttype {
	case LET:
		p.parseExtFuncDef(pi)
	case TYPE:
		p.parseExtTypeDef(pi)
	default:
		panic("Unknown pkginfo def")
	}
}

/*
PACKAGE_INFO = 'package_info' IDENTIFIER '=' EOL (OFFSIDE EXT_DEF EOL)*

EXT_DEF = (EXT_TYPE_DEF|EXT_FUNC_DEF)
*/
func (p *Parser) parsePackageInfo() *PackageInfo {
	p.consume(PACKAGE_INFO)
	pkgName := p.identName()
	p.gotoNext()
	p.consume(EQ)
	p.skipEOLOne()
	p.skipSpace()
	p.pushOffside()
	pi := NewPackageInfo(pkgName)

	p.pushScope()
	defer p.popOffside()
	for !p.isEndOfBlock() {
		p.parseExtDef(pi)
		p.skipEOLOne()
	}

	p.popScope()
	pi.registerToScope(p.scope)

	return pi
}

func (p *Parser) parseStmt() Stmt {
	tk := p.Current()
	switch tk.ttype {
	case PACKAGE:
		return p.parsePackage()
	case IMPORT:
		return p.parseImport()
	case PACKAGE_INFO:
		return p.parsePackageInfo()
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
