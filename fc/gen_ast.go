package main

type GoEvalExpr struct {
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

type MatchPattern struct {
	caseId  string
	varName string
}

type Expr interface {
	Expr_Union()
}

func (Expr_GoEvalExpr) Expr_Union()     {}
func (Expr_StringLiteral) Expr_Union()  {}
func (Expr_IntImm) Expr_Union()         {}
func (Expr_Unit) Expr_Union()           {}
func (Expr_BoolLiteral) Expr_Union()    {}
func (Expr_FunCall) Expr_Union()        {}
func (Expr_FieldAccess) Expr_Union()    {}
func (Expr_Var) Expr_Union()            {}
func (Expr_RecordGen) Expr_Union()      {}
func (Expr_LazyBlock) Expr_Union()      {}
func (Expr_ReturnableExpr) Expr_Union() {}

type Expr_GoEvalExpr struct {
	Value GoEvalExpr
}

func New_Expr_GoEvalExpr(v GoEvalExpr) Expr { return Expr_GoEvalExpr{v} }

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

type Expr_LazyBlock struct {
	Value LazyBlock
}

func New_Expr_LazyBlock(v LazyBlock) Expr { return Expr_LazyBlock{v} }

type Expr_ReturnableExpr struct {
	Value ReturnableExpr
}

func New_Expr_ReturnableExpr(v ReturnableExpr) Expr { return Expr_ReturnableExpr{v} }

type FunCall struct {
	targetFunc Var
	args       []Expr
}
type RecordGen struct {
	fieldNames  []string
	fieldValues []Expr
	recordType  RecordType
}
type Block struct {
	stmts     []Stmt
	finalExpr Expr
}
type LazyBlock struct {
	stmts     []Stmt
	finalExpr Expr
}
type MatchRule struct {
	pattern MatchPattern
	body    Block
}
type MatchExpr struct {
	target Expr
	rules  []MatchRule
}
type ReturnableExpr interface {
	ReturnableExpr_Union()
}

func (ReturnableExpr_Block) ReturnableExpr_Union()     {}
func (ReturnableExpr_MatchExpr) ReturnableExpr_Union() {}

type ReturnableExpr_Block struct {
	Value Block
}

func New_ReturnableExpr_Block(v Block) ReturnableExpr { return ReturnableExpr_Block{v} }

type ReturnableExpr_MatchExpr struct {
	Value MatchExpr
}

func New_ReturnableExpr_MatchExpr(v MatchExpr) ReturnableExpr { return ReturnableExpr_MatchExpr{v} }

type Stmt interface {
	Stmt_Union()
}

func (Stmt_Import) Stmt_Union()       {}
func (Stmt_Package) Stmt_Union()      {}
func (Stmt_PackageInfo) Stmt_Union()  {}
func (Stmt_LetFuncDef) Stmt_Union()   {}
func (Stmt_LetVarDef) Stmt_Union()    {}
func (Stmt_ExprStmt) Stmt_Union()     {}
func (Stmt_DefStmt) Stmt_Union()      {}
func (Stmt_MultipleDefs) Stmt_Union() {}

type Stmt_Import struct {
	Value string
}

func New_Stmt_Import(v string) Stmt { return Stmt_Import{v} }

type Stmt_Package struct {
	Value string
}

func New_Stmt_Package(v string) Stmt { return Stmt_Package{v} }

type Stmt_PackageInfo struct {
	Value PackageInfo
}

func New_Stmt_PackageInfo(v PackageInfo) Stmt { return Stmt_PackageInfo{v} }

type Stmt_LetFuncDef struct {
	Value LetFuncDef
}

func New_Stmt_LetFuncDef(v LetFuncDef) Stmt { return Stmt_LetFuncDef{v} }

type Stmt_LetVarDef struct {
	Value LetVarDef
}

func New_Stmt_LetVarDef(v LetVarDef) Stmt { return Stmt_LetVarDef{v} }

type Stmt_ExprStmt struct {
	Value Expr
}

func New_Stmt_ExprStmt(v Expr) Stmt { return Stmt_ExprStmt{v} }

type Stmt_DefStmt struct {
	Value DefStmt
}

func New_Stmt_DefStmt(v DefStmt) Stmt { return Stmt_DefStmt{v} }

type Stmt_MultipleDefs struct {
	Value MultipleDefs
}

func New_Stmt_MultipleDefs(v MultipleDefs) Stmt { return Stmt_MultipleDefs{v} }

type LetFuncDef struct {
	name   string
	params []Var
	body   Block
}
type LetVarDef struct {
	name string
	rhs  Expr
}
type PackageInfo struct {
	name     string
	funcInfo funcTypeDict
	typeInfo extTypeDict
}
type RecordDef struct {
	name   string
	fields []NameTypePair
}
type UnionDef struct {
	name  string
	cases []NameTypePair
}
type DefStmt interface {
	DefStmt_Union()
}

func (DefStmt_RecordDef) DefStmt_Union() {}
func (DefStmt_UnionDef) DefStmt_Union()  {}

type DefStmt_RecordDef struct {
	Value RecordDef
}

func New_DefStmt_RecordDef(v RecordDef) DefStmt { return DefStmt_RecordDef{v} }

type DefStmt_UnionDef struct {
	Value UnionDef
}

func New_DefStmt_UnionDef(v UnionDef) DefStmt { return DefStmt_UnionDef{v} }

type MultipleDefs struct {
	defs []DefStmt
}
