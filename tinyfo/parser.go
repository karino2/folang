package main

import (
	"bytes"
	"fmt"
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
	LE
	GE
	BRACKET // <>
	PIPE
	STRING
	COLON
	COMMA
	SEMICOLON
	INT_IMM
	OF
	BAR
	BARBAR
	RARROW
	UNDER_SCORE
	MATCH
	WITH
	TRUE
	FALSE
	PACKAGE_INFO
	DOT
	AND    // "and"
	AMP    // "&"
	AMPAMP // "&&"
	PLUS
	MINUS
	ASTER
	IF
	THEN
	ELSE
	ELIF
	NOT
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
	"and":          AND,
	"if":           IF,
	"then":         THEN,
	"else":         ELSE,
	"elif":         ELIF,
	"not":          NOT,
}

type binOpInfo struct {
	precedence int
	goFuncName string
}

// int is preedance
var binOpMap = map[TokenType]binOpInfo{
	PIPE:    {1, "frt.Pipe"},
	AMPAMP:  {2, "&&"},
	BARBAR:  {2, "||"},
	GT:      {2, ">"},
	LT:      {2, "<"},
	GE:      {2, ">="},
	LE:      {2, "<="},
	EQ:      {3, "frt.OpEqual"}, // as a comparison operator
	BRACKET: {3, "frt.OpNotEqual"},
	PLUS:    {4, "+"},
	MINUS:   {4, "-"},
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

func (tkz *Tokenizer) isStringAt(at int, s string) bool {
	if at+len(s) > len(tkz.buf) {
		return false
	}
	for i := range s {
		if s[i] != tkz.buf[at+i] {
			return false
		}
	}
	return true
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
			buf.WriteByte(c)
			i++
			if tkz.pos+i == len(tkz.buf) {
				panic("escape just before EOF, wrong")
			}
			c2 := tkz.buf[tkz.pos+i]
			buf.WriteByte(c2)
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
	  search s pos from start, then return that position.
		If reach EOF, return -1.
*/
func (tkz *Tokenizer) searchForward(start int, s string) int {
	pos := start

	for ; pos < len(tkz.buf); pos++ {
		if tkz.isStringAt(pos, s) {
			return pos
		}
	}
	return -1
}

func (tkz *Tokenizer) analyzeCurAsSpace() {
	cur := tkz.currentToken
	cur.ttype = SPACE
	i := 0
	for tkz.isCharAt(tkz.pos+i, ' ') || tkz.isStringAt(tkz.pos+i, "/*") || tkz.isStringAt(tkz.pos+i, "//") || tkz.isCharAt(tkz.pos+i, '\t') {
		for ; tkz.isCharAt(tkz.pos+i, ' '); i++ {
		}
		for ; tkz.isCharAt(tkz.pos+i, '\t'); i++ {
		}
		if tkz.isStringAt(tkz.pos+i, "/*") {
			end := tkz.searchForward(tkz.pos+i+2, "*/")
			if end == -1 {
				panic("No comment end found.")
			}
			i = end - tkz.pos + 2
		}
		if tkz.isStringAt(tkz.pos+i, "//") {
			for ; !tkz.isCharAt(tkz.pos+i, '\n'); i++ {
			}
		}
	}
	cur.len = i
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
	case b == ' ' || b == '\t':
		tkz.analyzeCurAsSpace()
	case b == '/':
		if tkz.isCharAt(tkz.pos+1, '*') || tkz.isCharAt(tkz.pos+1, '/') {
			tkz.analyzeCurAsSpace()
		} else {
			panic("slash, NYI")
		}
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
		if tkz.isCharAt(tkz.pos+1, '|') {
			cur.ttype = BARBAR
			cur.stringVal = "||"
			cur.len = 2
			return
		}
		cur.setOneChar(BAR, b)
	case b == '<':
		if tkz.isCharAt(tkz.pos+1, '>') {
			cur.ttype = BRACKET
			cur.stringVal = "<>"
			cur.len = 2
			return
		}
		if tkz.isCharAt(tkz.pos+1, '=') {
			cur.ttype = LE
			cur.stringVal = "<="
			cur.len = 2
			return
		}
		cur.setOneChar(LT, b)
	case b == '>':
		if tkz.isCharAt(tkz.pos+1, '=') {
			cur.ttype = GE
			cur.stringVal = ">="
			cur.len = 2
			return
		}
		cur.setOneChar(GT, b)
	case b == '+':
		cur.setOneChar(PLUS, b)
	case b == '&':
		if tkz.isCharAt(tkz.pos+1, '&') {
			cur.ttype = AMPAMP
			cur.stringVal = "&&"
			cur.len = 2
			return
		}
		cur.setOneChar(AMP, b)
	case b == '*':
		cur.setOneChar(ASTER, b)
	case b == '-':
		if tkz.isCharAt(tkz.pos+1, '>') {
			cur.ttype = RARROW
			cur.stringVal = "->"
			cur.len = 2
			return
		}
		cur.setOneChar(MINUS, b)
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

func (s *Scope) LookupRecordByName(rname string) *FRecord {
	cur := s
	for cur != nil {
		ret, ok := cur.recordMap[rname]
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

// parser context inside type definition.
type typeDefCtx struct {
	defined   map[string]FType
	inTypeDef bool
}

func newTypeDefCtx() *typeDefCtx {
	tdc := &typeDefCtx{}
	tdc.defined = make(map[string]FType)
	return tdc
}

func (tdc *typeDefCtx) clear() {
	clear(tdc.defined)
	tdc.inTypeDef = false
}

func (tdc *typeDefCtx) resolvePreUsedType(md *MultipleDefs) {
	Walk(md, func(n Node) bool {
		switch dt := n.(type) {
		case *MultipleDefs:
			return true
		case *RecordDef:
			for i, np := range dt.Fields {
				if pt, ok := np.Type.(*FPreUsed); ok {
					if rt, ok := tdc.defined[pt.name]; ok {
						dt.Fields[i].Type = rt
					} else {
						panic("type used but not found: " + pt.name)
					}
				}
			}
			return false
		case *UnionDef:
			for i, cs := range dt.Cases {
				if pt, ok := cs.Type.(*FPreUsed); ok {
					if rt, ok := tdc.defined[pt.name]; ok {
						dt.Cases[i].Type = rt
					} else {
						panic("type used but not found2: " + pt.name)
					}
				}
			}
			return false
		default:
			return false
		}
	})
}

type Parser struct {
	tokenizer  *Tokenizer
	offsideCol []int
	scope      *Scope
	typeDefCtx *typeDefCtx
}

func NewParser() *Parser {
	p := &Parser{}
	p.scope = NewScope(nil)
	p.typeDefCtx = newTypeDefCtx()
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

func (p *Parser) isInTypeDef() bool {
	return p.typeDefCtx.inTypeDef
}

func (p *Parser) enterTypeDef() {
	p.typeDefCtx.inTypeDef = true
}

func (p *Parser) leaveTypeDef() {
	p.typeDefCtx.clear()
	p.typeDefCtx.inTypeDef = false
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

func (p *Parser) revertTo(tk *Token) {
	p.tokenizer.RevertTo(tk)
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

func (p *Parser) identOrUSName() string {
	if p.currentIs(UNDER_SCORE) {
		return "_"
	}
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

/*
IMPORT_STMT = 'import' (STRING_LITERAL | IDENTIFIER)

For identifier case, add prefix of folang pkg path.
ex:
import frt
=>
import "github.com/karino2/folang/pkg/frt"
*/
func (p *Parser) parseImport() *Import {
	p.consume(IMPORT)
	tk := p.Current()
	var impath string
	if tk.ttype == STRING {
		impath = tk.stringVal
	} else if tk.ttype == IDENTIFIER {
		impath = fmt.Sprintf("github.com/karino2/folang/pkg/%s", tk.stringVal)
	} else {
		panic(tk)
	}
	p.gotoNextSL()
	return &Import{impath}
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
		ttype == RSBRACKET ||
		ttype == WITH ||
		ttype == THEN ||
		ttype == ELSE ||
		ttype == COMMA ||
		p.nextNonEOLIsBinOp()
}

/*
RECORD_EXPRESISONN = '{' FIELD_INITIALIZERS '}'

FIELD_INITIALIZERS = FIELD_INITIALIZER (';' FIELD_INITIALIZER)*

FIELD_INITIALIZER = IDENTIFIER '=' expr | SPECIFIED_INITIALIZER

Specified initializer specify record name like: {myRec.X = 3; Y=4; Z=5}

SPECIFIED_INITIALIZER = IDENTIFIER '.' IDENTIFIER '=' expr
*/
func (p *Parser) parseRecordGen() Expr {
	p.consumeSL(LBRACE)

	var fnames []string
	var fvals []Expr
	recName := ""
	for i := 0; p.Current().ttype != RBRACE; i++ {
		if i != 0 {
			p.consumeSL(SEMICOLON)
		}

		fnameCand := p.identName()
		p.gotoNextSL()

		if p.currentIs(DOT) {
			p.consume(DOT)
			recName = fnameCand
			fnameCand = p.identName()
			p.gotoNextSL()
		}

		fnames = append(fnames, fnameCand)

		p.consumeSL(EQ)
		one := p.parseExpr()
		fvals = append(fvals, one)
	}

	p.consume(RBRACE)
	rg := NewRecordGen(fnames, fvals)

	if recName != "" {
		// specified initializer
		rg.recordType = p.scope.LookupRecordByName(recName)
	} else {
		rg.recordType = p.scope.LookupRecord(rg.fieldNames)
	}

	return rg
}

func (p *Parser) referenceVar(vname string) *Var {
	v := p.scope.LookupVar(vname)
	if v == nil {
		panic("Undefined var: " + vname)
	}
	return v
}

/*
VARIABLE_REF = IDENTIFIER | IDENTIFER '.' IDENTIFIER

For first case, IDENTIFIER must be in scope.
For later case, there are 2 types:

  - external pkg access: slice.Take
  - record field access: rec.Person
*/
func (p *Parser) parseVariableReference() Expr {
	firstId := p.identName()
	if !p.nextIs(DOT) {
		p.gotoNext()
		return p.referenceVar(firstId)
	}

	// Next is dot. Check whether rec field access or pkg access.
	v := p.scope.LookupVar(firstId)
	if v != nil {
		// symbol found, record field access.
		var expr Expr = v
		p.gotoNext()

		// Special handling for a.b.c.d
		// Target expr should be more versatile.
		// But I just want to use a.b.c.d.
		// So just support that case.
		for p.currentIs(DOT) {
			p.consume(DOT)

			fname := p.identName()
			p.gotoNext()

			// v must be record.
			rtype := expr.FType().(*FRecord)
			expr = &FieldAccess{expr, rtype, fname}
		}
		return expr
	} else {
		// external pkg access.
		fullName := p.parseFullName()
		return p.referenceVar(fullName)
	}
}

/*
Grammar says Application Expression is expr expr.
So create ATOM parser that does not handle application expressions.

ATOM = LITERAL | VARIABLE_REF | REC_GEN | '(' ')' | '(' PAREN_EXPR ')'

PAREN_EXPR = EXPR (',' EXPR)*
that is, EXPR or TUPLE_EXPR.
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
		return p.parseVariableReference()
	case tk.ttype == LBRACE:
		return p.parseRecordGen()
	case tk.ttype == LPAREN:
		p.consume(LPAREN)
		if p.currentIs(RPAREN) {
			p.consume(RPAREN)
			return gUnitVal
		} else {
			ret := p.parseExpr()
			first := true
			for p.currentIs(COMMA) {
				if !first {
					panic("only pair is suppoted for tuple")
				}
				first = false
				p.consume(COMMA)
				ne := p.parseExpr()
				ret = &TupleExpr{[]Expr{ret, ne}}
			}
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
func (p *Parser) parseMatchRule(target Expr) *MatchRule {
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

	// add varName to scope
	p.pushScope()
	if varName != "" && varName != "_" {
		// currently, match expect Union type.
		fu := target.FType().(*FUnion)
		uc := fu.lookupCase(caseName)
		p.scope.DefineVar(varName, &Var{varName, uc.Type})
	}

	block := p.parseBlockAfterPushScope()
	return &MatchRule{&MatchPattern{caseName, varName}, block}
}

func (p *Parser) nextIsBarInsideOffside() bool {
	tk := p.Current()
	defer p.tokenizer.RevertTo(tk)

	p.skipEOL()
	if !p.currentIs(BAR) {
		return false
	}
	return p.currentCol() >= p.currentOffside()
}

/*
MATCH_EXPR = 'match' EXPR 'with' EOL MATCH_RULE+
*/
func (p *Parser) parseMatchExpr() Expr {
	p.consume(MATCH)
	target := p.parseExpr()
	p.consume(WITH)

	var rules []*MatchRule
	for p.currentIs(EOL) && p.nextIsBarInsideOffside() {
		p.skipEOL()
		one := p.parseMatchRule(target)
		rules = append(rules, one)
	}
	return &MatchExpr{target, rules}
}

/*
parse single expr and return as block.
This is occur like 'if c then XX else YY' onliner.
*/
func (p *Parser) parseInlineBlock() *Block {
	p.pushScope()
	expr := p.parseExpr()
	return &Block{[]Stmt{}, expr, p.popScope(), true}
}

/*
IF_EXPR = 'if' expr 'then' block (('elif' expr 'then' block)* 'else' block)?

But for elif, this function parse just after if symbol.
*/
func (p *Parser) parseIfAfterIfExpr() Expr {
	cond := p.parseExpr()
	if cond.FType() != FBool {
		panic("cond is not bool")
	}

	p.consume(THEN)
	if p.currentIs(EOL) {
		// offside
		p.skipEOLOne()
		tbody := p.parseBlock()
		savePos := p.Current()

		p.skipEOLOne()
		if !p.currentIs(ELSE) && !p.currentIs(ELIF) {
			// no else block.
			p.revertTo(savePos)
			return NewIfOnlyCall(cond, tbody)
		}
		if p.currentIs(ELIF) {
			p.consume(ELIF)
			// 'elif' means 'else' block which start with 'if' expression.
			// but offside rule is a little tricky.
			// So regard this else block as ifblock start from here.
			ebody := p.parseIfAfterIfExpr()
			return NewIfElseCall(cond, tbody, NewBlock(ebody))
		}

		p.consume(ELSE)
		p.skipEOLOne()
		fbody := p.parseBlock()
		return NewIfElseCall(cond, tbody, fbody)
	} else {
		// one line case: if COND then TBODY else FBODY
		tbody := p.parseInlineBlock()
		if !p.currentIs(ELSE) {
			return NewIfOnlyCall(cond, tbody)
		}
		p.consume(ELSE)
		fbody := p.parseInlineBlock()
		return NewIfElseCall(cond, tbody, fbody)
	}
}

/*
IF_EXPR = 'if' expr 'then' block (('elif' expr 'then' block)* 'else' block)?

- expr must bo bool type.
- both block must be the same return type.
- if there is EOL after 'then', it assumue offside block
- if there is no else block, block return type must be unit.
*/
func (p *Parser) parseIfExpr() Expr {
	p.consume(IF)
	return p.parseIfAfterIfExpr()
}

/*
SLICE_EXPR = '[' expr (; expr)* ']'
*/
func (p *Parser) parseSliceExpr() Expr {
	var exprs []Expr
	p.consume(LSBRACKET)
	expr := p.parseExpr()
	exprs = append(exprs, expr)

	for p.currentIs(SEMICOLON) {
		p.consume(SEMICOLON)
		expr = p.parseExpr()
		exprs = append(exprs, expr)
	}

	p.consume(RSBRACKET)
	return &SliceExpr{exprs}
}

/*
TERM = GOEVAL | MATCH_EXPR | IF_EXPR | SLICE_EXPR | ATOM ATOM* | 'not' TERM

'not' in F# has much lower precedence, but I handle here because it's easier and OK for most of the case.
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

	if tk.ttype == LSBRACKET {
		return p.parseSliceExpr()
	}

	if tk.ttype == IF {
		return p.parseIfExpr()
	}

	if tk.ttype == NOT {
		p.consume(NOT)
		target := p.parseTerm()

		// treat as frt.OpNot call.
		pvar := &Var{"frt.OpNot", NewFFunc(FBool, FBool)}
		return &FunCall{pvar, []Expr{target}, []ResolvedTypeParam{}}
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
	_, ok := binOpMap[tk.ttype]
	return ok
}

/*
EXPR = TERM (BINOP EXPR)*
*/
func (p *Parser) parseExprWithPrecedence(minPrec int) Expr {
	expr := p.parseTerm()
	for p.nextNonEOLIsBinOp() {
		p.skipEOL()
		binOpType := p.Current().ttype
		binInfo := binOpMap[binOpType]

		if binInfo.precedence < minPrec {
			return expr
		}

		p.consume(binOpType)
		rhs := p.parseExprWithPrecedence(binInfo.precedence + 1)

		expr = NewBinOpCall(binOpType, binInfo, expr, rhs)
	}
	return expr
}

func (p *Parser) parseExpr() Expr {
	return p.parseExprWithPrecedence(1)
}

func (p *Parser) isEndOfBlock() bool {
	p.skipSpace()
	return p.currentCol() < p.currentOffside() || p.currentIs(EOF)
}

/*
BLOCK = EXPR | (STMT_LIKE EOL)* EXPR

STMT_LIKE = LET_STMT | EXPR

In some case, we wanto add local variables in block scope,
in that case, it's handy to pushScope before calling this function.
*/
func (p *Parser) parseBlockAfterPushScope() *Block {
	var stmts []Stmt
	var last Expr

	p.skipSpace()
	p.skipEOL() // skip comment only line by this skipEOL

	// limited implementation of offside rule.
	p.pushOffside()

	for {
		if p.Current().ttype == LET {
			stmts = append(stmts, p.parseLet())
			p.skipEOL()
		} else {
			last = p.parseExpr()

			// Check next line is end of offside line.
			// If so, go back to current pos and end block.
			// If not, go to next stmt parse.
			savePos := p.Current()

			p.skipEOL()
			if p.isEndOfBlock() {
				p.popOffside()
				p.revertTo(savePos)
				return &Block{stmts, last, p.popScope(), false}
			} else {
				stmts = append(stmts, &ExprStmt{last})
			}
		}
	}
}

func (p *Parser) parseBlock() *Block {
	p.pushScope()
	return p.parseBlockAfterPushScope()
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
ATOM_TYPE = 'string' | 'int' | '(' ')' | REGISTERED_TYPE | '[' ']' TERM_TYPE | '(' TYPE ')'
*/
func (p *Parser) parseAtomType() FType {
	if p.Current().ttype == LSBRACKET {
		p.consume(LSBRACKET)
		p.consume(RSBRACKET)
		etype := p.parseTermType()
		return &FSlice{etype}
	}

	if p.Current().ttype == LPAREN {
		p.consume(LPAREN)
		if p.currentIs(RPAREN) {
			p.consume(RPAREN)
			return FUnit
		}
		res := p.parseType()
		p.consume(RPAREN)
		return res
	}

	tname := p.identName()
	switch tname {
	case "string":
		p.gotoNext()
		return FString
	case "int":
		p.gotoNext()
		return FInt
	case "bool":
		p.gotoNext()
		return FBool
	default:
		fullName := p.parseFullName()
		resT := p.scope.LookupType(fullName)
		if resT == nil {
			// unknown type.

			// if inside typedef,
			// this type might be defined in later 'and' definition.
			// So return FUndefined and resolve later.
			if p.isInTypeDef() {
				return &FPreUsed{fullName}
			}

			// not in typedef and unknown, unknown type.
			panic(fullName)
		}
		return resT
	}
}

/*
TERM_TYPE = ATOM_TYPE | TUPLE_TYPE

TUPLE_TYPE = ATOME_TYPE '*' ATOM_TYPE ('*' ATOME_TYPE)*
*/
func (p *Parser) parseTermType() FType {
	one := p.parseAtomType()

	if !p.currentIs(ASTER) {
		return one
	}
	// only support pair for a while.
	p.consume(ASTER)
	two := p.parseAtomType()
	return NewFTuple(one, two)
}

/*
For backward compat.
*/
func (p *Parser) parseFuncType() *FFunc {
	ft := p.parseType()
	return ft.(*FFunc)
}

/*
TYPE =  TERM_TYPE | FUNC_TYPE

TERM_TYPE = ATOM_TYPE | TUPLE_TYPE | '(' TYPE ')'
FUNC_TYPE = TERM_TYPE '->' TERM_TYPE ('->' TERM_TYPE)*
*/
func (p *Parser) parseType() FType {
	one := p.parseTermType()

	if p.Current().ttype == RARROW {
		// func type
		types := []FType{one}

		for p.Current().ttype == RARROW {
			p.consume(RARROW)
			types = append(types, p.parseTermType())
		}
		return NewFFunc(types...)
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
FUNC_DEF_LET = 'let' IDENTIFIER PARAMS (':' TYPE)? '=' EOL BLOCK
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

	if p.currentIs(COLON) {
		// return type annotation.
		// In this case, we might register current func type before starting parse func body.
		// Good for recursive call.
		p.consume(COLON)

		rt := p.parseType()
		fts := varsToFTypes(params)
		fts = append(fts, rt)

		// NYI for type parameters
		v.Type = NewFFunc(fts...)
	}

	p.consumeSL(EQ)

	block := p.parseBlock()

	ret := &FuncDef{fname.stringVal, params, block}
	// update type.
	v.Type = ret.FuncFType()
	return ret
}

/*
Destructuring, only support pair for a while like:
let (a, b) = ...

LET_DEST_VAR_DEF = 'let' '(' IDENTIFIER ',' IDENTIFIER ')' '=' expr
*/
func (p *Parser) parseDestLetDefVar() Stmt {
	p.consume(LET)
	p.consume(LPAREN)
	vname1 := p.identOrUSName()
	p.gotoNext()
	p.consume(COMMA)
	vname2 := p.identOrUSName()
	p.gotoNext()
	p.consume(RPAREN)
	p.consume(EQ)

	rhs := p.parseExpr()
	if tp, ok := rhs.FType().(*FTuple); ok {
		if len(tp.Elems) != 2 {
			panic("Destructuring let, rhs tuple not 2D.")
		}
		if vname1 != "_" {
			v1 := &Var{vname1, tp.Elems[0]}
			p.scope.DefineVar(vname1, v1)
		}
		if vname2 != "_" {
			v2 := &Var{vname2, tp.Elems[1]}
			p.scope.DefineVar(vname2, v2)
		}

		return &LetDestVarDef{[]string{vname1, vname2}, rhs}
	} else {
		panic("destructuring with righ expr not tuple.")
	}
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
LET = LET_DEST_VAR_DEF | LET_VAR_DEF | LET_FUNC_DEF

LET_VAR_DEF = 'let' IDENTIFIER '=' expr
LET_DEST_VAR_DEF = 'let' '(' IDENTIFIER, IDENTIFIER ')' '=' expr
*/
func (p *Parser) parseLet() Stmt {
	nt := p.peekNext()
	if nt.ttype == LPAREN {
		return p.parseDestLetDefVar()
	} else {
		nn := p.peekNextNext()
		if nn.ttype == EQ {
			return p.parseLetDefVar()
		} else {
			return p.parseFuncDefLet()
		}

	}
}

func (p *Parser) registerUnionType(ud *UnionDef) {
	ud.registerToScope(p.scope)
	p.typeDefCtx.defined[ud.Name] = ud.UnionFType()
}

/*
UNION_DEF = CASE_DEF (EOL CASE_DEF)* EOL

CASE_DEF = '|' IDENIFIER OF TYPE
*/
func (p *Parser) parseUnionDef(uname string) Stmt {
	var cases []NameTypePair

	lastSavedPos := p.Current()
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
		lastSavedPos = p.Current()
		p.skipEOL() // try skip empty line, eg. comment only line.
	}
	// revert last skip EOL.
	p.revertTo(lastSavedPos)

	ret := &UnionDef{uname, cases}
	p.registerUnionType(ret)
	return ret
}

func (p *Parser) registerRecordType(rname string, recType *FRecord) {
	p.scope.recordMap[rname] = recType
	p.scope.typeMap[rname] = recType
	p.typeDefCtx.defined[rname] = recType
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
			if p.Current().ttype == RBRACE {
				break
			}
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
	p.registerRecordType(rname, recType)

	return rd
}

// TYPE_DEF_BODY = ID '=' (RECORD_DEF | UNION_DEF)
//
// RECORD_DEF = '{' FIELD_DEFS '}'
//
// UNION_DEF = '|'...
func (p *Parser) parseTypeDefBody() Stmt {
	tname := p.identName()
	p.gotoNextSL()

	p.consumeSL(EQ)

	if p.Current().ttype == LBRACE {
		return p.parseRecordDef(tname)
	}

	p.expect(BAR)
	return p.parseUnionDef(tname)
}

// TYPE_DEF = 'type' TYPE_DEF_BODY ('and' TYPE_DEF_BODY)*
func (p *Parser) parseTypeDef() Stmt {
	p.enterTypeDef()
	defer p.leaveTypeDef()

	p.consume(TYPE)
	stmt := p.parseTypeDefBody()

	if !p.currentIs(AND) {
		return stmt
	}

	var stmts []Stmt
	stmts = append(stmts, stmt)
	for p.currentIs(AND) {
		p.consume(AND)
		stmt = p.parseTypeDefBody()
		stmts = append(stmts, stmt)
		p.skipEOLOne()
	}

	md := &MultipleDefs{stmts}
	p.typeDefCtx.resolvePreUsedType(md)
	return md
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
	var pkgName string

	if p.currentIs(IDENTIFIER) {
		pkgName = p.identName()
		p.gotoNext()
	} else {
		if !p.currentIs(UNDER_SCORE) {
			panic("Wrong package info name type.")
		}
		pkgName = "_"
		p.gotoNext()
	}
	p.consume(EQ)
	p.skipEOL()
	p.pushOffside()
	pi := NewPackageInfo(pkgName)

	p.pushScope()
	defer p.popOffside()
	for !p.isEndOfBlock() {
		p.parseExtDef(pi)
		p.skipEOL()
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
