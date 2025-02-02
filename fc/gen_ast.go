package main

type GoEval struct {
	goStmt  string
	typeArg FType
}

type Var struct {
	name  string
	ftype FType
}

type FieldAccess struct {
	targetName string
	targetType RecordType
	fieldName  string
}

type ResolvedTypeParam struct {
	name         string
	resolvedType FType
}

type MatchPattern struct {
	caseId  string
	varName string
}

type Expr interface {
	Expr_Union()
}

func (Expr_GoEval) Expr_Union()        {}
func (Expr_StringLiteral) Expr_Union() {}
func (Expr_IntImm) Expr_Union()        {}
func (Expr_Unit) Expr_Union()          {}
func (Expr_BoolLiteral) Expr_Union()   {}
func (Expr_FunCall) Expr_Union()       {}
func (Expr_FieldAccess) Expr_Union()   {}
func (Expr_Var) Expr_Union()           {}
func (Expr_RecordGen) Expr_Union()     {}
func (Expr_Block) Expr_Union()         {}
func (Expr_MatchExpr) Expr_Union()     {}

type Expr_GoEval struct {
	Value GoEval
}

func New_Expr_GoEval(v GoEval) Expr { return Expr_GoEval{v} }

type Expr_StringLiteral struct {
	Value string
}

func New_Expr_StringLiteral(v string) Expr { return Expr_StringLiteral{v} }

type Expr_IntImm struct {
	Value int
}

func New_Expr_IntImm(v int) Expr { return Expr_IntImm{v} }

type Expr_Unit struct {
}

var New_Expr_Unit Expr = Expr_Unit{}

type Expr_BoolLiteral struct {
	Value bool
}

func New_Expr_BoolLiteral(v bool) Expr { return Expr_BoolLiteral{v} }

type Expr_FunCall struct {
	Value FunCall
}

func New_Expr_FunCall(v FunCall) Expr { return Expr_FunCall{v} }

type Expr_FieldAccess struct {
	Value FieldAccess
}

func New_Expr_FieldAccess(v FieldAccess) Expr { return Expr_FieldAccess{v} }

type Expr_Var struct {
	Value Var
}

func New_Expr_Var(v Var) Expr { return Expr_Var{v} }

type Expr_RecordGen struct {
	Value RecordGen
}

func New_Expr_RecordGen(v RecordGen) Expr { return Expr_RecordGen{v} }

type Expr_Block struct {
	Value Block
}

func New_Expr_Block(v Block) Expr { return Expr_Block{v} }

type Expr_MatchExpr struct {
	Value MatchExpr
}

func New_Expr_MatchExpr(v MatchExpr) Expr { return Expr_MatchExpr{v} }

type FunCall struct {
	targetFunc Var
	args       []Expr
	typeParams []ResolvedTypeParam
}
type RecordGen struct {
	fieldNames  []string
	fieldValues []Expr
	recordType  FType
}
type Block struct {
	stms      []Stmt
	finalExpr Expr
	asFunc    bool
}
type MatchRule struct {
	pattern MatchPattern
	body    Block
}
type MatchExpr struct {
	target Expr
	rules  []MatchRule
}
type Stmt interface {
	Stmt_Union()
}

func (Stmt_Import) Stmt_Union() {}

type Stmt_Import struct {
	Value string
}

func New_Stmt_Import(v string) Stmt { return Stmt_Import{v} }
