package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		fmt.Println("Usage: fc file1.fo file2.fo file3.fo ...")
		return
	}

	tp := NewTranspiler()
	parser := &Parser{}
	for _, file := range files {
		buf, err := os.ReadFile(file)
		check(err)
		stmts := parser.Parse(file, buf)
		f := NewFile(stmts)
		res := tp.Transpile(f)

		dir := filepath.Dir(file)
		base := strings.TrimSuffix(filepath.Base(file), ".fo")
		dest := filepath.Join(dir, base+"_gen.go")
		os.WriteFile(dest, []byte(res), 0644)
	}
}
