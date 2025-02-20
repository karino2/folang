package main

import (
	"bytes"
	"fmt"

	"github.com/karino2/folang/pkg/frt"
)

// currently, map is NYI  and generic exxt type is also NYI.
// wrap to standard type for each.
func dictLookup[K comparable, T any](dict map[K]T, key K) frt.Tuple2[T, bool] {
	e, ok := dict[key]
	return frt.NewTuple2(e, ok)
}

func dictPut[T any](dict map[string]T, key string, v T) {
	dict[key] = v
}

func dictKeyValues[K comparable, V any](dict map[K]V) []frt.Tuple2[K, V] {
	var res []frt.Tuple2[K, V]
	for k, v := range dict {
		res = append(res, frt.NewTuple2(k, v))
	}
	return res
}

func toDict[V any](ps []frt.Tuple2[string, V]) map[string]V {
	dic := make(map[string]V)
	for _, tp := range ps {
		pname, v := frt.Destr(tp)
		dic[pname] = v
	}
	return dic
}

type funcFacDict = map[string]FuncFactory

func newFFD() funcFacDict {
	return make(map[string]FuncFactory)
}

func ffdPut(dic funcFacDict, key string, v FuncFactory) {
	dictPut(dic, key, v)
}

func ffdKVs(dic funcFacDict) []frt.Tuple2[string, FuncFactory] {
	return dictKeyValues(dic)
}

type TFDataDict = map[string]TypeFactoryData

func newTFDD() TFDataDict {
	return make(TFDataDict)
}

func tfddPut(dic TFDataDict, key string, v TypeFactoryData) {
	dictPut(dic, key, v)
}

func tfddKVs(dic TFDataDict) []frt.Tuple2[string, TypeFactoryData] {
	return dictKeyValues(dic)
}

/*
  uniqueTmpVarName related.
*/

var uniqueId = 0

func uniqueTmpVarName() string {
	uniqueId++
	return fmt.Sprintf("_v%d", uniqueId)
}

func resetUniqueTmpCounter() {
	uniqueId = 0
}

/*
  Tokenizer related.
  Currently, pattern match has too many NYI and it's easier to implement it in golang world.

	This is mainly porting from tinyfo's tokenizer, but more functional way (create new tokenizer instead of side effect update).
*/

var keywordMap = map[string]TokenType{
	"let":          New_TokenType_LET,
	"package":      New_TokenType_PACKAGE,
	"import":       New_TokenType_IMPORT,
	"type":         New_TokenType_TYPE,
	"of":           New_TokenType_OF,
	"_":            New_TokenType_UNDER_SCORE,
	"match":        New_TokenType_MATCH,
	"with":         New_TokenType_WITH,
	"true":         New_TokenType_TRUE,
	"false":        New_TokenType_FALSE,
	"package_info": New_TokenType_PACKAGE_INFO,
	"and":          New_TokenType_AND,
	"if":           New_TokenType_IF,
	"then":         New_TokenType_THEN,
	"else":         New_TokenType_ELSE,
	"elif":         New_TokenType_ELIF,
	"not":          New_TokenType_NOT,
}

func newToken(ttype TokenType, begin int, len int) Token {
	return Token{ttype, begin, len, "", 0}
}

func newOneCharToken(ttype TokenType, pos int, one byte) Token {
	tk := newToken(ttype, pos, 1)
	tk.stringVal = string(one) // maybe not necessary.
	return tk
}

func newStLikeToken(ttype TokenType, begin int, s string) Token {
	tk := newToken(ttype, begin, len(s))
	tk.stringVal = s
	return tk
}

func isCharAt(buf string, at int, ch byte) bool {
	if at >= len(buf) {
		return false
	}
	return buf[at] == ch
}

func isStringAt(buf string, at int, s string) bool {
	if at+len(s) > len(buf) {
		return false
	}
	for i := range s {
		if s[i] != buf[at+i] {
			return false
		}
	}
	return true
}

/*
	  search s pos from start, then return that position.
		If reach EOF, return -1.
*/
func searchForward(buf string, start int, s string) int {
	pos := start

	for ; pos < len(buf); pos++ {
		if isStringAt(buf, pos, s) {
			return pos
		}
	}
	return -1
}

func scanSpaceToken(buf string, pos int) Token {
	cur := newToken(New_TokenType_SPACE, pos, 0)
	i := 0
	for isCharAt(buf, pos+i, ' ') || isStringAt(buf, pos+i, "/*") || isStringAt(buf, pos+i, "//") || isCharAt(buf, pos+i, '\t') {
		for ; isCharAt(buf, pos+i, ' '); i++ {
		}
		for ; isCharAt(buf, pos+i, '\t'); i++ {
		}
		if isStringAt(buf, pos+i, "/*") {
			end := searchForward(buf, pos+i+2, "*/")
			if end == -1 {
				panic("No comment end found.")
			}
			i = end - pos + 2
		}
		if isStringAt(buf, pos+i, "//") {
			for ; !isCharAt(buf, pos+i, '\n'); i++ {
			}
		}
	}
	cur.len = i
	return cur
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

// exclusive
func (tk Token) end() int { return tk.begin + tk.len }

func scanIdentifierToken(buf string, pos int) Token {
	cur := newToken(New_TokenType_IDENTIFIER, pos, 0)
	i := 1
	for {
		// reach end.
		if pos+i == len(buf) {
			cur.len = i
			break
		}

		c := buf[pos+i]
		if isAlnum(c) || c == '_' {
			i++
			continue
		}

		cur.len = i
		break
	}
	cur.stringVal = string(buf[cur.begin:cur.end()])
	return cur
}

func scanIntImmToken(buf string, pos int) Token {
	cur := newToken(New_TokenType_INT_IMM, pos, 0)
	c := buf[pos]
	n := 0
	i := 0
	for isNumber(c) {
		n = 10*n + int(c-'0')
		i++
		c = buf[pos+i]
	}

	cur.len = i
	cur.intVal = n
	return cur
}

func scanStringLiteralToken(buf string, pos int) Token {
	cur := newToken(New_TokenType_STRING, pos, 0)
	i := 1
	var bb bytes.Buffer
	for {
		// reach end without close. parse error. panic for a while.
		if pos+i == len(buf) {
			panic("unclosed string literal")
		}

		c := buf[pos+i]

		if c == '"' {
			cur.len = i + 1
			cur.stringVal = bb.String()
			return cur
		} else if c == '\\' {
			bb.WriteByte(c)
			i++
			if pos+i == len(buf) {
				panic("escape just before EOF, wrong")
			}
			c2 := buf[pos+i]
			bb.WriteByte(c2)
		} else {
			bb.WriteByte(c)
		}
		i++
	}
}

func scanTokenAt(buf string, pos int) Token {
	if pos == len(buf) {
		return newToken(New_TokenType_EOF, pos, 0)
	}
	b := buf[pos]
	switch {
	case b == ' ' || b == '\t':
		return scanSpaceToken(buf, pos)
	case b == '/':
		if isCharAt(buf, pos+1, '*') || isCharAt(buf, pos+1, '/') {
			return scanSpaceToken(buf, pos)
		} else {
			panic("slash, NYI")
		}
	case 'a' <= b && b <= 'z' ||
		'A' <= b && b <= 'Z' ||
		b == '_':
		cur := scanIdentifierToken(buf, pos)

		// check whether identifier is keyword
		if tt, ok := keywordMap[cur.stringVal]; ok {
			cur.ttype = tt
		}
		return cur
	case isNumber(b):
		return scanIntImmToken(buf, pos)
	case b == '"':
		return scanStringLiteralToken(buf, pos)
	case b == '=':
		return newOneCharToken(New_TokenType_EQ, pos, b)
	case b == '\n':
		return newOneCharToken(New_TokenType_EOL, pos, b)
	case b == '(':
		return newOneCharToken(New_TokenType_LPAREN, pos, b)
	case b == ')':
		return newOneCharToken(New_TokenType_RPAREN, pos, b)
	case b == '{':
		return newOneCharToken(New_TokenType_LBRACE, pos, b)
	case b == '}':
		return newOneCharToken(New_TokenType_RBRACE, pos, b)
	case b == '[':
		return newOneCharToken(New_TokenType_LSBRACKET, pos, b)
	case b == ']':
		return newOneCharToken(New_TokenType_RSBRACKET, pos, b)
	case b == ':':
		return newOneCharToken(New_TokenType_COLON, pos, b)
	case b == ',':
		return newOneCharToken(New_TokenType_COMMA, pos, b)
	case b == '.':
		return newOneCharToken(New_TokenType_DOT, pos, b)
	case b == ';':
		return newOneCharToken(New_TokenType_SEMICOLON, pos, b)
	case b == '|':
		if isCharAt(buf, pos+1, '>') {
			return newStLikeToken(New_TokenType_PIPE, pos, "|>")
		}
		if isCharAt(buf, pos+1, '|') {
			return newStLikeToken(New_TokenType_BARBAR, pos, "||")
		}
		return newOneCharToken(New_TokenType_BAR, pos, b)
	case b == '<':
		if isCharAt(buf, pos+1, '>') {
			return newStLikeToken(New_TokenType_BRACKET, pos, "<>")
		}
		if isCharAt(buf, pos+1, '=') {
			return newStLikeToken(New_TokenType_LE, pos, "<=")
		}
		return newOneCharToken(New_TokenType_LT, pos, b)
	case b == '>':
		if isCharAt(buf, pos+1, '=') {
			return newStLikeToken(New_TokenType_GE, pos, ">=")
		}
		return newOneCharToken(New_TokenType_GT, pos, b)
	case b == '+':
		return newOneCharToken(New_TokenType_PLUS, pos, b)
	case b == '&':
		if isCharAt(buf, pos+1, '&') {
			return newStLikeToken(New_TokenType_AMPAMP, pos, "&&")
		}
		return newOneCharToken(New_TokenType_AMP, pos, b)
	case b == '*':
		return newOneCharToken(New_TokenType_ASTER, pos, b)
	case b == '-':
		if isCharAt(buf, pos+1, '>') {
			return newStLikeToken(New_TokenType_RARROW, pos, "->")
		}
		return newOneCharToken(New_TokenType_MINUS, pos, b)
	default:
		panic(b)
	}
}

// very special.
// No space is important.
// This is the only place where space is matters.
func isNeighborLT(buf string, prev Token) bool {
	if len(buf) <= prev.end() {
		return false
	}
	return buf[prev.end()] == '<'
}

// next non-space token.
func nextToken(buf string, prev Token) Token {
	if len(buf) <= prev.end() {
		return newToken(New_TokenType_EOF, len(buf), 0)
	}
	tk := scanTokenAt(buf, prev.end())
	for tk.ttype == New_TokenType_SPACE {
		tk = scanTokenAt(buf, tk.end())
	}
	return tk
}

/*
	  GoEvalExpr utility.
		It might be possible to implement in folang, but I already have it in tinyfo, so use it.
*/
func reinterpretEscape(buf string) string {
	var b bytes.Buffer
	eof := len(buf)
	i := 0
	for {
		if i == eof {
			break
		}
		c := buf[i]
		if c == '\\' {
			i++
			if i == eof {
				panic("escape just before EOF, wrong")
			}
			c2 := buf[i]
			if c2 == 'n' {
				b.WriteByte('\n')
			} else {
				b.WriteByte(c2)
			}
		} else {
			b.WriteByte(c)
		}
		i++
	}
	return b.String()
}

/*
	Scope implementation.
	Currently, map is not supported, and side effect is hard to write in folang.
	So I write Scope related code in golang, then call it from folang.
*/

type VarFactory = func(stlist []FType, tvgen func() TypeVar) VarRef
type TypeFactory = func(stlist []FType) FType

type scopeImpl struct {
	// var factory Map
	varFacMap  map[string]VarFactory
	recordMap  map[string]RecordType
	typeFacMap map[string]TypeFactory
	parent     *scopeImpl
}

// for folang, show pointer as real type for hidden side effect.
type Scope = *scopeImpl

func newScope(parent *scopeImpl) *scopeImpl {
	s := &scopeImpl{}
	s.varFacMap = make(map[string]VarFactory)
	s.recordMap = make(map[string]RecordType)
	s.typeFacMap = make(map[string]TypeFactory)
	s.parent = parent
	return s
}

func popScope(src Scope) Scope {
	return src.parent
}

func NewScope() Scope {
	return newScope(nil)
}

func scDefVar(s Scope, name string, v Var) {
	s.varFacMap[name] = func(_ []FType, _ func() TypeVar) VarRef { return New_VarRef_VRVar(v) }
}

func scRegisterVarFac(s Scope, name string, fac VarFactory) {
	s.varFacMap[name] = fac
}

// currently, we can't support Result because of absence of generic type.
// We use golang style convention though F# convention is bool is first.
func scLookupVarFac(s Scope, name string) frt.Tuple2[VarFactory, bool] {
	cur := s
	for cur != nil {
		ret, ok := cur.varFacMap[name]
		if ok {
			return frt.NewTuple2(ret, true)
		}
		cur = cur.parent
	}
	return frt.NewTuple2[VarFactory](nil, false)
}

func scRegisterTypeFac(s Scope, name string, tfac TypeFactory) {
	s.typeFacMap[name] = tfac
}

func toTypeFac(ft FType) TypeFactory {
	return func([]FType) FType { return ft }
}

func scRegisterType(s Scope, name string, ftype FType) {
	scRegisterTypeFac(s, name, toTypeFac(ftype))
}

func scRegisterRecType(s Scope, recType RecordType) {
	rname := recType.Name
	s.recordMap[rname] = recType
	scRegisterType(s, rname, New_FType_FRecord(recType))
}

func scLookupRecordCur(s Scope, fieldNames []string) frt.Tuple2[RecordType, bool] {
	for _, rt := range s.recordMap {
		if frMatch(rt, fieldNames) {
			return frt.NewTuple2(rt, true)
		}
	}
	return frt.NewTuple2(RecordType{}, false)

}

func scLookupRecord(s Scope, fieldNames []string) frt.Tuple2[RecordType, bool] {
	cur := s
	for cur != nil {
		ret := scLookupRecordCur(cur, fieldNames)
		if ret.E1 {
			return ret
		}
		cur = cur.parent
	}
	return frt.NewTuple2(RecordType{}, false)
}

func scLookupRecordByName(s Scope, rname string) frt.Tuple2[RecordType, bool] {
	cur := s
	for cur != nil {
		ret, ok := cur.recordMap[rname]
		if ok {
			return frt.NewTuple2(ret, true)
		}
		cur = cur.parent
	}
	return frt.NewTuple2(RecordType{}, false)
}

func scLookupTypeFac(s Scope, name string) frt.Tuple2[TypeFactory, bool] {
	cur := s
	for cur != nil {
		ret, ok := cur.typeFacMap[name]
		if ok {
			return frt.NewTuple2(ret, true)
		}
		cur = cur.parent
	}
	return frt.NewTuple2(toTypeFac(New_FType_FUnit), false)
}

func withPs[T any](ps ParseState, v T) frt.Tuple2[ParseState, T] {
	return frt.NewTuple2(ps, v)
}

/*
inference from funcall to arg side is NYI.

func Cnv1[T any, U any](fn func(T) T, prev frt.Tuple2[T, U]) frt.Tuple2[T, U] {
	t, u := frt.Destr(prev)
	return frt.NewTuple2(fn(t), u)
}

func Cnv2[T any, U any](fn func(U) U, prev frt.Tuple2[T, U]) frt.Tuple2[T, U] {
	t, u := frt.Destr(prev)
	return frt.NewTuple2(t, fn(u))
}
*/

func CnvL[U any](fn func(ParseState) ParseState, prev frt.Tuple2[ParseState, U]) frt.Tuple2[ParseState, U] {
	t, u := frt.Destr(prev)
	return frt.NewTuple2(fn(t), u)
}

func CnvR[T any, U any](fn func(T) U, prev frt.Tuple2[ParseState, T]) frt.Tuple2[ParseState, U] {
	t, u := frt.Destr(prev)
	return frt.NewTuple2(t, fn(u))
}

/*
  TypeVarAllocator
*/

type typeVarAllocator struct {
	seqId     int
	prefix    string
	allocated []TypeVar
}

type TypeVarAllocator = *typeVarAllocator

func NewTypeVarAllocator(prefix string) TypeVarAllocator {
	return &typeVarAllocator{0, prefix, []TypeVar{}}
}

func (tva *typeVarAllocator) Reset() {
	tva.seqId = 0
	tva.allocated = []TypeVar{}
}

func (tva *typeVarAllocator) genVarName() string {
	vname := fmt.Sprintf("%s%d", tva.prefix, tva.seqId)
	tva.seqId++
	if tva.seqId > 100 {
		panic("Too many type var alloc.")
	}
	return vname
}

func (tva *typeVarAllocator) Allocate() TypeVar {
	tvar := TypeVar{tva.genVarName()}
	tva.allocated = append(tva.allocated, tvar)
	return tvar
}

func tvaReset(tva TypeVarAllocator) {
	tva.Reset()
}

func tvaToTypeVarGen(tva TypeVarAllocator) func() TypeVar {
	return func() TypeVar { return tva.Allocate() }
}

func tvaListAlloced(tva TypeVarAllocator) []TypeVar {
	return tva.allocated
}

/*
  TypeVarDict related.
*/

type TypeVarDict = map[string]TypeVar

func toTVDict(ps []frt.Tuple2[string, TypeVar]) TypeVarDict {
	return toDict(ps)
}

// NF: Never fail.
func tvdLookupNF(tvd TypeVarDict, key string) TypeVar {
	return tvd[key]
}

/*
EquivSet
*/

type EquivSet struct {
	dict map[string]bool
}

type EquivInfoDict struct {
	dict map[string]EquivInfo
}

func newEquivSet0() EquivSet {
	s := EquivSet{}
	s.dict = make(map[string]bool)
	return s
}

func NewEquivSet(itype TypeVar) EquivSet {
	s := newEquivSet0()
	s.dict[itype.Name] = true
	return s
}

func NewEquivInfoDict() EquivInfoDict {
	ei := EquivInfoDict{}
	ei.dict = make(map[string]EquivInfo)
	return ei
}

func eidPut(eid EquivInfoDict, key string, v EquivInfo) {
	dictPut(eid.dict, key, v)
}

func eidLookup(eid EquivInfoDict, key string) frt.Tuple2[EquivInfo, bool] {
	return dictLookup(eid.dict, key)
}

func eqsItems(es EquivSet) []string {
	var res []string
	for k := range es.dict {
		res = append(res, k)
	}
	return res
}

func eqsUnion(es1 EquivSet, es2 EquivSet) EquivSet {
	e3 := newEquivSet0()
	for _, k := range eqsItems(es1) {
		e3.dict[k] = true
	}
	for _, k := range eqsItems(es2) {
		e3.dict[k] = true
	}
	return e3
}

/*
BinOp related
*/

// type BinOpInfo = {precedence: int; goFuncName: string}
var binOpMap = map[TokenType]BinOpInfo{
	New_TokenType_PIPE:    {1, "frt.Pipe", false},
	New_TokenType_AMPAMP:  {2, "&&", true},
	New_TokenType_BARBAR:  {2, "||", true},
	New_TokenType_GT:      {2, ">", true},
	New_TokenType_LT:      {2, "<", true},
	New_TokenType_GE:      {2, ">=", true},
	New_TokenType_LE:      {2, "<=", true},
	New_TokenType_EQ:      {3, "frt.OpEqual", true}, // as a comparison operator
	New_TokenType_BRACKET: {3, "frt.OpNotEqual", true},
	New_TokenType_PLUS:    {4, "+", false},
	New_TokenType_MINUS:   {4, "-", false},
}

func lookupBinOp(tk TokenType) frt.Tuple2[BinOpInfo, bool] {
	return dictLookup(binOpMap, tk)
}

/*
TypeDict related
*/
type TypeDict = map[string]FType

func newTD() TypeDict {
	return make(TypeDict)
}
func tdPut(dic TypeDict, key string, v FType) {
	dictPut(dic, key, v)
}

func tdLookup(dic TypeDict, key string) frt.Tuple2[FType, bool] {
	return dictLookup(dic, key)
}

func toTDict(ps []frt.Tuple2[string, FType]) TypeDict {
	return toDict(ps)
}

/*
Facade:
Resolve mutual recursive in golang layer (NYI for and letfunc def).
*/
func parseLetFacade(ps ParseState) frt.Tuple2[ParseState, LLetVarDef] {
	return parseLetVarDef(parseExprFacade, ps)
}

func parseBlockFacade(ps ParseState) frt.Tuple2[ParseState, Block] {
	return parseBlock(parseLetFacade, ps)
}

func parseExprFacade(ps ParseState) frt.Tuple2[ParseState, Expr] {
	return parseExpr(parseBlockFacade, ps)
}

func ParseAll(ps ParseState) frt.Tuple2[ParseState, []RootStmt] {
	var res []RootStmt
	var stmt RootStmt
	ps2 := ps
	for !IsParseEnd(ps2) {
		ps2, stmt = frt.Destr(ParseStmtOne(parseExprFacade, ps2))
		res = append(res, stmt)
	}
	return frt.NewTuple2(ps2, res)
}

// for test backward compat.
func parseAll(ps ParseState) (ParseState, []RootStmt) {
	res := ParseAll(ps)
	return frt.Destr(res)
}
