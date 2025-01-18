package main

import "bytes"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	SPACE
	IDENTIFIER
	EQ
	LET
	EOL
	PACKAGE
	IMPORT
	LPAREN
	RPAREN
	STRING
	/*
		INDENT
		UNDENT

	*/
)

type Token struct {
	ttype     TokenType
	begin     int
	len       int
	stringVal string // also used in Identifier.
	/*
		intVal    int
	*/
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
	/*
		line         int // 0 origin
		offsideCol   int
		dents        int
	*/
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

var keywordMap = map[string]TokenType{
	"let":     LET,
	"package": PACKAGE,
	"import":  IMPORT,
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
	case b == '"':
		tkz.analyzeCurAsStringLiteral()
	case b == '=':
		cur.ttype = EQ
		cur.stringVal = "="
		cur.len = 1
	case b == '\n':
		cur.ttype = EOL
		cur.stringVal = "\n"
		cur.len = 1
	case b == '(':
		cur.ttype = LPAREN
		cur.stringVal = "("
		cur.len = 1
	case b == ')':
		cur.ttype = RPAREN
		cur.stringVal = ")"
		cur.len = 1
	default:
		panic(b)
	}

}

func (tkz *Tokenizer) Current() *Token { return tkz.currentToken }
func (tkz *Tokenizer) Setup()          { tkz.analyzeCur() }
func (tkz *Tokenizer) GotoNext() {
	if tkz.currentToken.ttype == EOF {
		return
	}
	tkz.pos = tkz.currentToken.end()
	tkz.analyzeCur()
}

type Parser struct {
	tokenizer *Tokenizer
}

func (p *Parser) Current() *Token {
	return p.tokenizer.Current()
}

func (p *Parser) skipSpace() {
	tk := p.Current()
	for tk.ttype == SPACE {
		p.tokenizer.GotoNext()
		tk = p.Current()
	}
}

func (p *Parser) gotoNext() {
	p.tokenizer.GotoNext()
	p.skipSpace()
}

func (p *Parser) skipEOL() {
	tk := p.Current()
	for tk.ttype == EOL {
		p.gotoNext()
		tk = p.Current()
	}
}

func (p *Parser) consume(ttype TokenType) {
	tk := p.Current()
	if tk.ttype != ttype {
		panic(tk)
	}
	p.gotoNext()
}

func (p *Parser) parsePackage() *Package {
	p.consume(PACKAGE)
	ident := p.Current()
	if ident.ttype != IDENTIFIER {
		panic(ident)
	}
	p.gotoNext()
	p.skipEOL()
	return &Package{ident.stringVal}
}

func (p *Parser) parseImport() *Import {
	p.consume(IMPORT)
	tk := p.Current()
	if tk.ttype != STRING {
		panic(tk)
	}
	p.gotoNext()
	p.skipEOL()
	return &Import{tk.stringVal}
}

func (p *Parser) parseExpr() Expr {
	// only support
	// GoEval "xxx"
	tk := p.Current()
	if tk.ttype != IDENTIFIER || tk.stringVal != "GoEval" {
		panic(*tk)
	}

	p.gotoNext()
	arg := p.Current()
	if arg.ttype != STRING {
		panic(arg)
	}
	p.gotoNext()
	return &GoEval{arg.stringVal}
}

func (p *Parser) parseFuncDefLet() Stmt {
	/*
		expext

		let hoge () =
			expr

		for a while.
	*/
	p.consume(LET)
	fname := p.Current()
	if fname.ttype != IDENTIFIER {
		panic(fname)
	}
	p.gotoNext()
	p.consume(LPAREN)
	p.consume(RPAREN)
	p.consume(EQ)
	p.skipEOL()
	// parse func body.
	p.skipSpace()
	expr := p.parseExpr()
	return &FuncDef{fname.stringVal, []Var{}, expr}
}

func (p *Parser) parseLet() Stmt {
	return p.parseFuncDefLet()
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

func (p *Parser) Parse(fileName string, buf []byte) []Stmt {
	p.tokenizer = NewTokenizer(fileName, buf)
	p.tokenizer.Setup()
	return p.parseStmts()
}
