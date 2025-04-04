package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/dict"

type GoEvalExpr struct {
	GoStmt  string
	TypeArg FType
}

type Var struct {
	Name  string
	Ftype FType
}

type SpecVar struct {
	Var      Var
	SpecList []FType
}

type VarRef interface {
	VarRef_Union()
}

func (VarRef_VRVar) VarRef_Union()  {}
func (VarRef_VRSVar) VarRef_Union() {}

func (v VarRef_VRVar) String() string  { return frt.Sprintf1("(VRVar: %v)", v.Value) }
func (v VarRef_VRSVar) String() string { return frt.Sprintf1("(VRSVar: %v)", v.Value) }

type VarRef_VRVar struct {
	Value Var
}

func New_VarRef_VRVar(v Var) VarRef { return VarRef_VRVar{v} }

type VarRef_VRSVar struct {
	Value SpecVar
}

func New_VarRef_VRSVar(v SpecVar) VarRef { return VarRef_VRSVar{v} }

func varRefName(vr VarRef) string {
	switch _v1 := (vr).(type) {
	case VarRef_VRVar:
		v := _v1.Value
		return v.Name
	case VarRef_VRSVar:
		sv := _v1.Value
		return sv.Var.Name
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func varRefVar(vr VarRef) Var {
	switch _v2 := (vr).(type) {
	case VarRef_VRVar:
		v := _v2.Value
		return v
	case VarRef_VRSVar:
		sv := _v2.Value
		return sv.Var
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func varRefVarType(vr VarRef) FType {
	v := varRefVar(vr)
	return v.Ftype
}

func transVarVR(transV func(Var) Var, vr VarRef) VarRef {
	switch _v3 := (vr).(type) {
	case VarRef_VRVar:
		v := _v3.Value
		return frt.Pipe(transV(v), New_VarRef_VRVar)
	case VarRef_VRSVar:
		sv := _v3.Value
		nv := transV(sv.Var)
		return frt.Pipe(SpecVar{Var: nv, SpecList: sv.SpecList}, New_VarRef_VRSVar)
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

type UnionMatchPattern struct {
	CaseId  string
	VarName string
}

type FuncFactory struct {
	Tparams []string
	Targets []FType
}

type TypeFactoryData struct {
	Name    string
	Tparams []string
}

type Expr interface {
	Expr_Union()
}

func (Expr_EGoEvalExpr) Expr_Union()     {}
func (Expr_EStringLiteral) Expr_Union()  {}
func (Expr_ESInterP) Expr_Union()        {}
func (Expr_EIntImm) Expr_Union()         {}
func (Expr_EUnit) Expr_Union()           {}
func (Expr_EBoolLiteral) Expr_Union()    {}
func (Expr_EFunCall) Expr_Union()        {}
func (Expr_EFieldAccess) Expr_Union()    {}
func (Expr_EVarRef) Expr_Union()         {}
func (Expr_ESlice) Expr_Union()          {}
func (Expr_ERecordGen) Expr_Union()      {}
func (Expr_ELazyBlock) Expr_Union()      {}
func (Expr_ETupleExpr) Expr_Union()      {}
func (Expr_ELambda) Expr_Union()         {}
func (Expr_EBinOpCall) Expr_Union()      {}
func (Expr_EReturnableExpr) Expr_Union() {}

func (v Expr_EGoEvalExpr) String() string     { return frt.Sprintf1("(EGoEvalExpr: %v)", v.Value) }
func (v Expr_EStringLiteral) String() string  { return frt.Sprintf1("(EStringLiteral: %v)", v.Value) }
func (v Expr_ESInterP) String() string        { return frt.Sprintf1("(ESInterP: %v)", v.Value) }
func (v Expr_EIntImm) String() string         { return frt.Sprintf1("(EIntImm: %v)", v.Value) }
func (v Expr_EUnit) String() string           { return "(EUnit)" }
func (v Expr_EBoolLiteral) String() string    { return frt.Sprintf1("(EBoolLiteral: %v)", v.Value) }
func (v Expr_EFunCall) String() string        { return frt.Sprintf1("(EFunCall: %v)", v.Value) }
func (v Expr_EFieldAccess) String() string    { return frt.Sprintf1("(EFieldAccess: %v)", v.Value) }
func (v Expr_EVarRef) String() string         { return frt.Sprintf1("(EVarRef: %v)", v.Value) }
func (v Expr_ESlice) String() string          { return frt.Sprintf1("(ESlice: %v)", v.Value) }
func (v Expr_ERecordGen) String() string      { return frt.Sprintf1("(ERecordGen: %v)", v.Value) }
func (v Expr_ELazyBlock) String() string      { return frt.Sprintf1("(ELazyBlock: %v)", v.Value) }
func (v Expr_ETupleExpr) String() string      { return frt.Sprintf1("(ETupleExpr: %v)", v.Value) }
func (v Expr_ELambda) String() string         { return frt.Sprintf1("(ELambda: %v)", v.Value) }
func (v Expr_EBinOpCall) String() string      { return frt.Sprintf1("(EBinOpCall: %v)", v.Value) }
func (v Expr_EReturnableExpr) String() string { return frt.Sprintf1("(EReturnableExpr: %v)", v.Value) }

type Expr_EGoEvalExpr struct {
	Value GoEvalExpr
}

func New_Expr_EGoEvalExpr(v GoEvalExpr) Expr { return Expr_EGoEvalExpr{v} }

type Expr_EStringLiteral struct {
	Value string
}

func New_Expr_EStringLiteral(v string) Expr { return Expr_EStringLiteral{v} }

type Expr_ESInterP struct {
	Value string
}

func New_Expr_ESInterP(v string) Expr { return Expr_ESInterP{v} }

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

type Expr_EVarRef struct {
	Value VarRef
}

func New_Expr_EVarRef(v VarRef) Expr { return Expr_EVarRef{v} }

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

type Expr_ETupleExpr struct {
	Value []Expr
}

func New_Expr_ETupleExpr(v []Expr) Expr { return Expr_ETupleExpr{v} }

type Expr_ELambda struct {
	Value LambdaExpr
}

func New_Expr_ELambda(v LambdaExpr) Expr { return Expr_ELambda{v} }

type Expr_EBinOpCall struct {
	Value BinOpCall
}

func New_Expr_EBinOpCall(v BinOpCall) Expr { return Expr_EBinOpCall{v} }

type Expr_EReturnableExpr struct {
	Value ReturnableExpr
}

func New_Expr_EReturnableExpr(v ReturnableExpr) Expr { return Expr_EReturnableExpr{v} }

type FunCall struct {
	TargetFunc VarRef
	Args       []Expr
}
type LambdaExpr struct {
	Params []Var
	Body   Block
}
type FieldAccess struct {
	TargetExpr Expr
	FieldName  string
}
type BinOpCall struct {
	Op    string
	Rtype FType
	Lhs   Expr
	Rhs   Expr
}
type NEPair struct {
	Name string
	Expr Expr
}
type RecordGen struct {
	FieldsNV   []NEPair
	RecordType RecordType
}
type Block struct {
	Stmts     []Stmt
	FinalExpr Expr
}
type LazyBlock struct {
	Block Block
}
type UnionMatchRule struct {
	UnionPattern UnionMatchPattern
	Body         Block
}
type UnionMatchRulesWD struct {
	Unions  []UnionMatchRule
	Default Block
}
type UnionMatchRules interface {
	UnionMatchRules_Union()
}

func (UnionMatchRules_UCaseOnly) UnionMatchRules_Union() {}
func (UnionMatchRules_UCaseWD) UnionMatchRules_Union()   {}

func (v UnionMatchRules_UCaseOnly) String() string { return frt.Sprintf1("(UCaseOnly: %v)", v.Value) }
func (v UnionMatchRules_UCaseWD) String() string   { return frt.Sprintf1("(UCaseWD: %v)", v.Value) }

type UnionMatchRules_UCaseOnly struct {
	Value []UnionMatchRule
}

func New_UnionMatchRules_UCaseOnly(v []UnionMatchRule) UnionMatchRules {
	return UnionMatchRules_UCaseOnly{v}
}

type UnionMatchRules_UCaseWD struct {
	Value UnionMatchRulesWD
}

func New_UnionMatchRules_UCaseWD(v UnionMatchRulesWD) UnionMatchRules {
	return UnionMatchRules_UCaseWD{v}
}

type StringMatchRule struct {
	LiteralPattern string
	Body           Block
}
type StringVarMatchRule struct {
	VarName string
	Body    Block
}
type StringMatchRulesWV struct {
	Literals []StringMatchRule
	VarRule  StringVarMatchRule
}
type StringMatchRulesWD struct {
	Literals []StringMatchRule
	Default  Block
}
type StringMatchRules interface {
	StringMatchRules_Union()
}

func (StringMatchRules_SCaseWV) StringMatchRules_Union() {}
func (StringMatchRules_SCaseWD) StringMatchRules_Union() {}

func (v StringMatchRules_SCaseWV) String() string { return frt.Sprintf1("(SCaseWV: %v)", v.Value) }
func (v StringMatchRules_SCaseWD) String() string { return frt.Sprintf1("(SCaseWD: %v)", v.Value) }

type StringMatchRules_SCaseWV struct {
	Value StringMatchRulesWV
}

func New_StringMatchRules_SCaseWV(v StringMatchRulesWV) StringMatchRules {
	return StringMatchRules_SCaseWV{v}
}

type StringMatchRules_SCaseWD struct {
	Value StringMatchRulesWD
}

func New_StringMatchRules_SCaseWD(v StringMatchRulesWD) StringMatchRules {
	return StringMatchRules_SCaseWD{v}
}

type MatchRules interface {
	MatchRules_Union()
}

func (MatchRules_RUnions) MatchRules_Union()  {}
func (MatchRules_RStrings) MatchRules_Union() {}

func (v MatchRules_RUnions) String() string  { return frt.Sprintf1("(RUnions: %v)", v.Value) }
func (v MatchRules_RStrings) String() string { return frt.Sprintf1("(RStrings: %v)", v.Value) }

type MatchRules_RUnions struct {
	Value UnionMatchRules
}

func New_MatchRules_RUnions(v UnionMatchRules) MatchRules { return MatchRules_RUnions{v} }

type MatchRules_RStrings struct {
	Value StringMatchRules
}

func New_MatchRules_RStrings(v StringMatchRules) MatchRules { return MatchRules_RStrings{v} }

type MatchExpr struct {
	Target Expr
	Rules  MatchRules
}
type ReturnableExpr interface {
	ReturnableExpr_Union()
}

func (ReturnableExpr_RBlock) ReturnableExpr_Union()     {}
func (ReturnableExpr_RMatchExpr) ReturnableExpr_Union() {}

func (v ReturnableExpr_RBlock) String() string     { return frt.Sprintf1("(RBlock: %v)", v.Value) }
func (v ReturnableExpr_RMatchExpr) String() string { return frt.Sprintf1("(RMatchExpr: %v)", v.Value) }

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

func (Stmt_SLetVarDef) Stmt_Union() {}
func (Stmt_SExprStmt) Stmt_Union()  {}

func (v Stmt_SLetVarDef) String() string { return frt.Sprintf1("(SLetVarDef: %v)", v.Value) }
func (v Stmt_SExprStmt) String() string  { return frt.Sprintf1("(SExprStmt: %v)", v.Value) }

type Stmt_SLetVarDef struct {
	Value LLetVarDef
}

func New_Stmt_SLetVarDef(v LLetVarDef) Stmt { return Stmt_SLetVarDef{v} }

type Stmt_SExprStmt struct {
	Value Expr
}

func New_Stmt_SExprStmt(v Expr) Stmt { return Stmt_SExprStmt{v} }

type RootStmt interface {
	RootStmt_Union()
}

func (RootStmt_RSImport) RootStmt_Union()       {}
func (RootStmt_RSPackage) RootStmt_Union()      {}
func (RootStmt_RSPackageInfo) RootStmt_Union()  {}
func (RootStmt_RSDefStmt) RootStmt_Union()      {}
func (RootStmt_RSMultipleDefs) RootStmt_Union() {}
func (RootStmt_RSRootFuncDef) RootStmt_Union()  {}
func (RootStmt_RSRootVarDef) RootStmt_Union()   {}

func (v RootStmt_RSImport) String() string      { return frt.Sprintf1("(RSImport: %v)", v.Value) }
func (v RootStmt_RSPackage) String() string     { return frt.Sprintf1("(RSPackage: %v)", v.Value) }
func (v RootStmt_RSPackageInfo) String() string { return frt.Sprintf1("(RSPackageInfo: %v)", v.Value) }
func (v RootStmt_RSDefStmt) String() string     { return frt.Sprintf1("(RSDefStmt: %v)", v.Value) }
func (v RootStmt_RSMultipleDefs) String() string {
	return frt.Sprintf1("(RSMultipleDefs: %v)", v.Value)
}
func (v RootStmt_RSRootFuncDef) String() string { return frt.Sprintf1("(RSRootFuncDef: %v)", v.Value) }
func (v RootStmt_RSRootVarDef) String() string  { return frt.Sprintf1("(RSRootVarDef: %v)", v.Value) }

type RootStmt_RSImport struct {
	Value string
}

func New_RootStmt_RSImport(v string) RootStmt { return RootStmt_RSImport{v} }

type RootStmt_RSPackage struct {
	Value string
}

func New_RootStmt_RSPackage(v string) RootStmt { return RootStmt_RSPackage{v} }

type RootStmt_RSPackageInfo struct {
	Value PackageInfo
}

func New_RootStmt_RSPackageInfo(v PackageInfo) RootStmt { return RootStmt_RSPackageInfo{v} }

type RootStmt_RSDefStmt struct {
	Value DefStmt
}

func New_RootStmt_RSDefStmt(v DefStmt) RootStmt { return RootStmt_RSDefStmt{v} }

type RootStmt_RSMultipleDefs struct {
	Value MultipleDefs
}

func New_RootStmt_RSMultipleDefs(v MultipleDefs) RootStmt { return RootStmt_RSMultipleDefs{v} }

type RootStmt_RSRootFuncDef struct {
	Value RootFuncDef
}

func New_RootStmt_RSRootFuncDef(v RootFuncDef) RootStmt { return RootStmt_RSRootFuncDef{v} }

type RootStmt_RSRootVarDef struct {
	Value RootVarDef
}

func New_RootStmt_RSRootVarDef(v RootVarDef) RootStmt { return RootStmt_RSRootVarDef{v} }

type LetFuncDef struct {
	Fvar   Var
	Params []Var
	Body   Block
}
type RootFuncDef struct {
	Tparams []string
	Lfd     LetFuncDef
}
type RootVarDef struct {
	Vdef LetVarDef
}
type LetVarDef struct {
	Lvar Var
	Rhs  Expr
}
type LetDestVarDef struct {
	Lvars []Var
	Rhs   Expr
}
type RawLetDef interface {
	RawLetDef_Union()
}

func (RawLetDef_RLetOneVar) RawLetDef_Union()  {}
func (RawLetDef_RLetDestVar) RawLetDef_Union() {}
func (RawLetDef_RLetFunc) RawLetDef_Union()    {}

func (v RawLetDef_RLetOneVar) String() string  { return frt.Sprintf1("(RLetOneVar: %v)", v.Value) }
func (v RawLetDef_RLetDestVar) String() string { return frt.Sprintf1("(RLetDestVar: %v)", v.Value) }
func (v RawLetDef_RLetFunc) String() string    { return frt.Sprintf1("(RLetFunc: %v)", v.Value) }

type RawLetDef_RLetOneVar struct {
	Value LetVarDef
}

func New_RawLetDef_RLetOneVar(v LetVarDef) RawLetDef { return RawLetDef_RLetOneVar{v} }

type RawLetDef_RLetDestVar struct {
	Value LetDestVarDef
}

func New_RawLetDef_RLetDestVar(v LetDestVarDef) RawLetDef { return RawLetDef_RLetDestVar{v} }

type RawLetDef_RLetFunc struct {
	Value LetFuncDef
}

func New_RawLetDef_RLetFunc(v LetFuncDef) RawLetDef { return RawLetDef_RLetFunc{v} }

type LLetVarDef interface {
	LLetVarDef_Union()
}

func (LLetVarDef_LLOneVarDef) LLetVarDef_Union()  {}
func (LLetVarDef_LLDestVarDef) LLetVarDef_Union() {}

func (v LLetVarDef_LLOneVarDef) String() string  { return frt.Sprintf1("(LLOneVarDef: %v)", v.Value) }
func (v LLetVarDef_LLDestVarDef) String() string { return frt.Sprintf1("(LLDestVarDef: %v)", v.Value) }

type LLetVarDef_LLOneVarDef struct {
	Value LetVarDef
}

func New_LLetVarDef_LLOneVarDef(v LetVarDef) LLetVarDef { return LLetVarDef_LLOneVarDef{v} }

type LLetVarDef_LLDestVarDef struct {
	Value LetDestVarDef
}

func New_LLetVarDef_LLDestVarDef(v LetDestVarDef) LLetVarDef { return LLetVarDef_LLDestVarDef{v} }

type PackageInfo struct {
	Name     string
	FuncInfo dict.Dict[string, FuncFactory]
	TypeInfo dict.Dict[string, TypeFactoryData]
}
type RecordDef struct {
	Name    string
	Tparams []string
	Fields  []NameTypePair
}
type UnionDef struct {
	Name    string
	Tparams []string
	Cases   []NameTypePair
}
type DefStmt interface {
	DefStmt_Union()
}

func (DefStmt_DRecordDef) DefStmt_Union() {}
func (DefStmt_DUnionDef) DefStmt_Union()  {}

func (v DefStmt_DRecordDef) String() string { return frt.Sprintf1("(DRecordDef: %v)", v.Value) }
func (v DefStmt_DUnionDef) String() string  { return frt.Sprintf1("(DUnionDef: %v)", v.Value) }

type DefStmt_DRecordDef struct {
	Value RecordDef
}

func New_DefStmt_DRecordDef(v RecordDef) DefStmt { return DefStmt_DRecordDef{v} }

type DefStmt_DUnionDef struct {
	Value UnionDef
}

func New_DefStmt_DUnionDef(v UnionDef) DefStmt { return DefStmt_DUnionDef{v} }

type MultipleDefs struct {
	Defs []DefStmt
}

func NewPackageInfo(name string) PackageInfo {
	ffd := dict.New[string, FuncFactory]()
	tfdd := dict.New[string, TypeFactoryData]()
	return PackageInfo{Name: name, FuncInfo: ffd, TypeInfo: tfdd}
}

func NewUnionDef(name string, tparams []string, cases []NameTypePair) UnionDef {
	return UnionDef{Name: name, Tparams: tparams, Cases: cases}
}

func udUpdate(ud UnionDef, cases []NameTypePair) UnionDef {
	return UnionDef{Name: ud.Name, Tparams: ud.Tparams, Cases: cases}
}

func udCases(ud UnionDef) []NameTypePair {
	return ud.Cases
}

func varToExpr(v Var) Expr {
	return frt.Pipe(New_VarRef_VRVar(v), New_Expr_EVarRef)
}

func newVar(vname string, vtype FType) Var {
	return Var{Name: vname, Ftype: vtype}
}
