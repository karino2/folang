package main

import (
	"fmt"
)

func main() {
	tp := NewTranspiler()
	f := NewFile(
		[]Stmt{
			&Package{"main"},
			&Import{"fmt"},
			&FuncDef{"main", nil, &GoEval{"fmt.Println(\"Hello World\")"}},
		},
	)
	fmt.Println(tp.Transpile(f))
}
