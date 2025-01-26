package main

import "bytes"

type TypeResolver struct {
	// variable type map.
	VarTypeMap map[string]FType

	RecordMap map[string]*FRecord
	TypeMap   map[string]FType
	/*
		Union case constructor use internal name (like "I") and go lang name ("New_IntOrString_I") differently.
		So register "I": &Var{"New_IntOrbool_I", int->IntOrString} for such case.
	*/
	AliasMap map[string]*Var
}

func NewResolver() *TypeResolver {
	res := TypeResolver{}
	res.VarTypeMap = make(map[string]FType)
	res.RecordMap = make(map[string]*FRecord)
	res.TypeMap = make(map[string]FType)
	res.AliasMap = make(map[string]*Var)
	return &res
}

func (res *TypeResolver) RegisterVarType(name string, ftype FType) {
	res.VarTypeMap[name] = ftype
}

func (res *TypeResolver) RegisterRecord(name string, frtype *FRecord) {
	res.RecordMap[name] = frtype
	res.TypeMap[name] = frtype
}

func (res *TypeResolver) RegisterType(name string, ftype FType) {
	res.TypeMap[name] = ftype
}

// LookupVarType type of variable by variable name.
func (res *TypeResolver) LookupVarType(name string) FType {
	v := res.VarTypeMap[name]
	if v != nil {
		if vt, ok := v.(*FCustom); ok {
			nt := res.LookupByTypeName(vt.name)
			if nt != nil {
				return nt
			}
		}
	}
	return v
}

// Lookup custom defined type by typename. Not by variable name.
func (res *TypeResolver) LookupByTypeName(tname string) FType {
	return res.TypeMap[tname]
}

func (res *TypeResolver) LookupRecord(fieldNames []string) *FRecord {
	for _, rt := range res.RecordMap {
		if rt.Match(fieldNames) {
			return rt
		}
	}
	return nil
}

func (res *TypeResolver) Resolve(n Node) {
	Walk(n, func(n Node) bool {
		switch n := n.(type) {
		case *Var:
			switch vt := n.Type.(type) {
			case *FUnresolved:
				nt := res.LookupVarType(n.Name)
				if nt != nil {
					n.Type = nt
				} else if alt, ok := res.AliasMap[n.Name]; ok {
					*n = *alt
				}
			case *FCustom:
				nt := res.LookupByTypeName(vt.name)
				if nt != nil {
					n.Type = nt
				}
			}
			return true
		case *FunCall:
			if n.Func.IsUnresolved() {
				if alt, ok := res.AliasMap[n.Func.Name]; ok {
					*(n.Func) = *alt
				}
			}
			return true
		case *RecordGen:
			rt := res.LookupRecord(n.fieldNames)
			if rt != nil {
				n.recordType = rt
			}
			return true
		default:
			return true
		}
	})
}

func registerType(resolver *TypeResolver, root Stmt) {
	Walk(root, func(n Node) bool {
		switch n := n.(type) {
		case *FuncDef:
			resolver.RegisterVarType(n.Name, n.FuncFType())
			// TODO* support scope
			for _, p := range n.Params {
				resolver.RegisterVarType(p.Name, p.Type)
			}
			return false
		case *RecordDef:
			resolver.RegisterRecord(n.Name, n.ToFType())
			return false
		case *UnionDef:
			n.ResiterUnionTypeInfo(resolver)
			return false
		default:
			return true
		}
	})
}

type File struct {
	FileName string
	Stmts    []Stmt
}

func (f *File) addStmt(stmt Stmt) {
	f.Stmts = append(f.Stmts, stmt)
}

func NewFile(stmts []Stmt) *File {
	var f File
	for _, s := range stmts {
		f.addStmt(s)
	}
	return &f
}

type Transpiler struct {
	Resolver *TypeResolver
}

func (tp *Transpiler) resolveAndRegisterType(stmts []Stmt) {
	for _, stmt := range stmts {
		// for param register
		registerType(tp.Resolver, stmt)
		tp.Resolver.Resolve(stmt)
		// for result func register.
		registerType(tp.Resolver, stmt)
	}
}

func (tp *Transpiler) Transpile(file *File) string {
	var buf bytes.Buffer
	tp.resolveAndRegisterType(file.Stmts)
	for _, stmt := range file.Stmts {
		buf.WriteString(stmt.ToGo())
		buf.WriteString("\n\n")
	}
	return buf.String()
}

func NewTranspiler() *Transpiler {
	var tp Transpiler
	tp.Resolver = NewResolver()
	return &tp
}
