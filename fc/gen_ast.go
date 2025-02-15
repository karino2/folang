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

func (Expr_EGoEvalExpr) Expr_Union()     {}
func (Expr_EStringLiteral) Expr_Union()  {}
func (Expr_EIntImm) Expr_Union()         {}
func (Expr_EUnit) Expr_Union()           {}
func (Expr_EBoolLiteral) Expr_Union()    {}
func (Expr_EFunCall) Expr_Union()        {}
func (Expr_EFieldAccess) Expr_Union()    {}
func (Expr_EVar) Expr_Union()            {}
func (Expr_ESlice) Expr_Union()          {}
func (Expr_ERecordGen) Expr_Union()      {}
func (Expr_ELazyBlock) Expr_Union()      {}
func (Expr_EReturnableExpr) Expr_Union() {}

type Expr_EGoEvalExpr struct {
	Value GoEvalExpr
}

func New_Expr_EGoEvalExpr(v GoEvalExpr) Expr { return Expr_EGoEvalExpr{v} }

type Expr_EStringLiteral struct {
	Value string
}

func New_Expr_EStringLiteral(v string) Expr { return Expr_EStringLiteral{v} }

type Expr_EIntImm struct {
	Value int
}

func New_Expr_EIntImm(v int) Expr { return Expr_EIntImm{v} }

type Expr_EUnit struct {
}

var New_Expr_EUnit Expr = Expr_EUnit{}

type Expr_EBoolLiteral struct {
	Value bool
}

func New_Expr_EBoolLiteral(v bool) Expr { return Expr_EBoolLiteral{v} }

type Expr_EFunCall struct {
	Value FunCall
}

func New_Expr_EFunCall(v FunCall) Expr { return Expr_EFunCall{v} }

type Expr_EFieldAccess struct {
	Value FieldAccess
}

func New_Expr_EFieldAccess(v FieldAccess) Expr { return Expr_EFieldAccess{v} }

type Expr_EVar struct {
	Value Var
}

func New_Expr_EVar(v Var) Expr { return Expr_EVar{v} }

type Expr_ESlice struct {
	Value []Expr
}

func New_Expr_ESlice(v []Expr) Expr { return Expr_ESlice{v} }

type Expr_ERecordGen struct {
	Value RecordGen
}

func New_Expr_ERecordGen(v RecordGen) Expr { return Expr_ERecordGen{v} }

type Expr_ELazyBlock struct {
	Value LazyBlock
}

func New_Expr_ELazyBlock(v LazyBlock) Expr { return Expr_ELazyBlock{v} }

type Expr_EReturnableExpr struct {
	Value ReturnableExpr
}

func New_Expr_EReturnableExpr(v ReturnableExpr) Expr { return Expr_EReturnableExpr{v} }

type FunCall struct {
	targetFunc Var
	args       []Expr
}
type NEPair struct {
	name string
	expr Expr
}
type RecordGen struct {
	fieldsNV   []NEPair
	recordType RecordType
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

func (ReturnableExpr_RBlock) ReturnableExpr_Union()     {}
func (ReturnableExpr_RMatchExpr) ReturnableExpr_Union() {}

type ReturnableExpr_RBlock struct {
	Value Block
}

func New_ReturnableExpr_RBlock(v Block) ReturnableExpr { return ReturnableExpr_RBlock{v} }

type ReturnableExpr_RMatchExpr struct {
	Value MatchExpr
}

func New_ReturnableExpr_RMatchExpr(v MatchExpr) ReturnableExpr { return ReturnableExpr_RMatchExpr{v} }

type Stmt interface {
	Stmt_Union()
}

func (Stmt_SImport) Stmt_Union()       {}
func (Stmt_SPackage) Stmt_Union()      {}
func (Stmt_SPackageInfo) Stmt_Union()  {}
func (Stmt_SLetFuncDef) Stmt_Union()   {}
func (Stmt_SLetVarDef) Stmt_Union()    {}
func (Stmt_SExprStmt) Stmt_Union()     {}
func (Stmt_SDefStmt) Stmt_Union()      {}
func (Stmt_SMultipleDefs) Stmt_Union() {}

type Stmt_SImport struct {
	Value string
}

func New_Stmt_SImport(v string) Stmt { return Stmt_SImport{v} }

type Stmt_SPackage struct {
	Value string
}

func New_Stmt_SPackage(v string) Stmt { return Stmt_SPackage{v} }

type Stmt_SPackageInfo struct {
	Value PackageInfo
}

func New_Stmt_SPackageInfo(v PackageInfo) Stmt { return Stmt_SPackageInfo{v} }

type Stmt_SLetFuncDef struct {
	Value LetFuncDef
}

func New_Stmt_SLetFuncDef(v LetFuncDef) Stmt { return Stmt_SLetFuncDef{v} }

type Stmt_SLetVarDef struct {
	Value LetVarDef
}

func New_Stmt_SLetVarDef(v LetVarDef) Stmt { return Stmt_SLetVarDef{v} }

type Stmt_SExprStmt struct {
	Value Expr
}

func New_Stmt_SExprStmt(v Expr) Stmt { return Stmt_SExprStmt{v} }

type Stmt_SDefStmt struct {
	Value DefStmt
}

func New_Stmt_SDefStmt(v DefStmt) Stmt { return Stmt_SDefStmt{v} }

type Stmt_SMultipleDefs struct {
	Value MultipleDefs
}

func New_Stmt_SMultipleDefs(v MultipleDefs) Stmt { return Stmt_SMultipleDefs{v} }

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

func (DefStmt_DRecordDef) DefStmt_Union() {}
func (DefStmt_DUnionDef) DefStmt_Union()  {}

type DefStmt_DRecordDef struct {
	Value RecordDef
}

func New_DefStmt_DRecordDef(v RecordDef) DefStmt { return DefStmt_DRecordDef{v} }

type DefStmt_DUnionDef struct {
	Value UnionDef
}

func New_DefStmt_DUnionDef(v UnionDef) DefStmt { return DefStmt_DUnionDef{v} }

type MultipleDefs struct {
	defs []DefStmt
}
