package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/sys"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/strings"

import "path/filepath"

func printUsage() {
	frt.Println("Usage: fc file1.fo file2.fo file3.fo ...")
}

func transpileOne(parser ParseState, file string) frt.Tuple2[ParseState, bool] {
	src, ok := frt.Destr(sys.ReadFile(file))
	return frt.IfElse(ok, (func() frt.Tuple2[ParseState, bool] {
		ps2, stmts := frt.Destr(frt.Pipe(psSetNewSrc(src, parser), ParseAll))
		res := RootStmtsToGo(stmts)
		frt.IfOnly(strings.HasSuffix(".fo", file), (func() {
			dir := filepath.Dir(file)
			base := frt.Pipe(filepath.Base(file), (func(_r0 string) string { return strings.TrimSuffix(".fo", _r0) }))
			newFname := (("gen_" + base) + ".go")
			dest := filepath.Join(dir, newFname)
			sys.WriteFile(dest, res)
		}))
		return frt.NewTuple2(ps2, true)
	}), (func() frt.Tuple2[ParseState, bool] {
		frt.Panicf1("Can't open file: %s", file)
		return frt.NewTuple2(parser, false)
	}))
}

func transpileRecur(parser ParseState, files []string) frt.Tuple2[ParseState, bool] {
	return frt.IfElse(slice.IsEmpty(files), (func() frt.Tuple2[ParseState, bool] {
		return frt.NewTuple2(parser, true)
	}), (func() frt.Tuple2[ParseState, bool] {
		head := slice.Head(files)
		ps2, ok := frt.Destr(transpileOne(parser, head))
		return frt.IfElse(ok, (func() frt.Tuple2[ParseState, bool] {
			rest := slice.Tail(files)
			return transpileRecur(ps2, rest)
		}), (func() frt.Tuple2[ParseState, bool] {
			return frt.NewTuple2(parser, false)
		}))
	}))
}

func transpileFiles(files []string) {
	parser := initParse("")
	transpileRecur(parser, files)

}

func main() {
	args := frt.Pipe(sys.Args(), slice.Tail)
	frt.IfElseUnit(slice.IsEmpty(args), (func() {
		printUsage()
	}), (func() {
		transpileFiles(args)
	}))
}
