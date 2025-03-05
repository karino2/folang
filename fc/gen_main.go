package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/sys"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/strings"

import "path/filepath"

func printUsage() {
	frt.Println("Usage: fc file1.fo file2.fo file3.fo ...")
}

func transpileOne(parser ParseState, file string) ParseState {
	frt.Printf1("transpile: %s\n", file)
	src, ok := frt.Destr2(sys.ReadFile(file))
	return frt.IfElse(ok, (func() ParseState {
		defer OnParseError(file)
		ps2, stmts := frt.Destr2(frt.Pipe(psSetNewSrc(src, parser), ParseAll))
		res := RootStmtsToGo(stmts)
		frt.IfOnly(strings.HasSuffix(".fo", file), (func() {
			dir := filepath.Dir(file)
			base := frt.Pipe(filepath.Base(file), (func(_r0 string) string { return strings.TrimSuffix(".fo", _r0) }))
			newFname := (("gen_" + base) + ".go")
			dest := filepath.Join(dir, newFname)
			sys.WriteFile(dest, res)
		}))
		return ps2
	}), (func() ParseState {
		frt.Panicf1("Can't open file: %s", file)
		return parser
	}))
}

func transpileFiles(files []string) {
	parser := initParse("")
	slice.Fold(transpileOne, parser, files)

}

func main() {
	args := frt.Pipe(sys.Args(), slice.Tail)
	frt.IfElseUnit(slice.IsEmpty(args), (func() {
		printUsage()
	}), (func() {
		transpileFiles(args)
	}))
}
