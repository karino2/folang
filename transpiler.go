package main

import "bytes"

type TypeResolver struct {
	TypeMap map[string]FType
}

func NewResolver() *TypeResolver {
	res := TypeResolver{}
	res.TypeMap = make(map[string]FType)
	return &res
}

func (res *TypeResolver) Register(name string, ftype FType) {
	res.TypeMap[name] = ftype
}

func (res *TypeResolver) Lookup(name string) FType {
	return res.TypeMap[name]
}

func (res *TypeResolver) Resolve(n Node) {
	Walk(n, func(n Node) bool {
		switch n := n.(type) {
		case *Var:
			if n.IsUnresolved() {
				nt := res.Lookup(n.Name)
				if nt != nil {
					n.Type = nt
				}
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
			resolver.Register(n.Name, n.FuncFType())
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
		tp.Resolver.Resolve(stmt)
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
