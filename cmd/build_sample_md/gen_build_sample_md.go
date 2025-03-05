package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/sys"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/buf"

import "github.com/karino2/folang/pkg/strings"

import "path/filepath"

func printUsage() {
	frt.Println("Usage: build_sample_md filelist.txt")
}

func convOne(dir string, oneline string) string {
	cols := strings.SplitN(2, " ", oneline)
	foFname, title := frt.Destr2(frt.NewTuple2(slice.Head(cols), slice.Last(cols)))
	b := buf.New()
	frt.Printf1("process: %s\n", foFname)
	content, ok := frt.Destr2(frt.Pipe(filepath.Join(dir, foFname), sys.ReadFile))
	frt.IfOnly(frt.OpNot(ok), (func() {
		frt.Panicf1("Can't open file %s", foFname)
	}))
	frt.PipeUnit(frt.Sprintf1("### %s\n\n", title), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "```\n")
	buf.Write(b, content)
	buf.Write(b, "\n```\n\n")
	base := frt.Pipe(foFname, (func(_r0 string) string { return strings.TrimSuffix(".fo", _r0) }))
	genName := (("gen_" + base) + ".go")
	frt.PipeUnit(frt.Sprintf1("generated go: [%s]", genName), (func(_r0 string) { buf.Write(b, _r0) }))
	frt.PipeUnit(frt.Sprintf1("(./%s)", genName), (func(_r0 string) { buf.Write(b, _r0) }))
	buf.Write(b, "\n\n")
	return buf.String(b)
}

func processListFile(destName string, listPath string) {
	dir := filepath.Dir(listPath)
	content, ok := frt.Destr2(sys.ReadFile(listPath))
	frt.IfOnly(frt.OpNot(ok), (func() {
		frt.Panicf1("Can't open list file: %s", listPath)
	}))
	frt.Pipe(frt.Pipe(frt.Pipe(frt.Pipe(frt.Pipe(frt.Pipe(content, (func(_r0 string) []string { return strings.Split("\n", _r0) })), (func(_r0 []string) []string { return slice.Filter(strings.IsNotEmpty, _r0) })), (func(_r0 []string) []string {
		return slice.Map((func(_r0 string) string { return convOne(dir, _r0) }), _r0)
	})), (func(_r0 []string) string { return strings.Concat("\n", _r0) })), (func(_r0 string) string { return strings.AppendHead("## Folang Sample \n\n\n", _r0) })), (func(_r0 string) bool { return sys.WriteFile(filepath.Join(dir, destName), _r0) }))

}

func main() {
	args := frt.Pipe(sys.Args(), slice.Tail)
	frt.IfElseUnit(frt.OpNotEqual(slice.Len(args), 1), (func() {
		printUsage()
	}), (func() {
		frt.PipeUnit(slice.Head(args), (func(_r0 string) { processListFile("README.md", _r0) }))
	}))
}
