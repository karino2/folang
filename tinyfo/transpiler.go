package main

import "bytes"

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
}

func (tp *Transpiler) Transpile(file *File) string {
	var buf bytes.Buffer
	for _, stmt := range file.Stmts {
		buf.WriteString(stmt.ToGo())
		buf.WriteString("\n\n")
	}
	return buf.String()
}

func NewTranspiler() *Transpiler {
	var tp Transpiler
	return &tp
}
