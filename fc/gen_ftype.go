package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/buf"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/strings"

type FType interface {
	FType_Union()
}

func (FType_FInt) FType_Union()          {}
func (FType_FString) FType_Union()       {}
func (FType_FBool) FType_Union()         {}
func (FType_FUnit) FType_Union()         {}
func (FType_FUnresolved) FType_Union()   {}
func (FType_FFunc) FType_Union()         {}
func (FType_FRecord) FType_Union()       {}
func (FType_FUnion) FType_Union()        {}
func (FType_FExtType) FType_Union()      {}
func (FType_FSlice) FType_Union()        {}
func (FType_FPreUsed) FType_Union()      {}
func (FType_FParametrized) FType_Union() {}

type FType_FInt struct {
}

var New_FType_FInt FType = FType_FInt{}

type FType_FString struct {
}

var New_FType_FString FType = FType_FString{}

type FType_FBool struct {
}

var New_FType_FBool FType = FType_FBool{}

type FType_FUnit struct {
}

var New_FType_FUnit FType = FType_FUnit{}

type FType_FUnresolved struct {
}

var New_FType_FUnresolved FType = FType_FUnresolved{}

type FType_FFunc struct {
	Value FuncType
}

func New_FType_FFunc(v FuncType) FType { return FType_FFunc{v} }

type FType_FRecord struct {
	Value RecordType
}

func New_FType_FRecord(v RecordType) FType { return FType_FRecord{v} }

type FType_FUnion struct {
	Value UnionType
}

func New_FType_FUnion(v UnionType) FType { return FType_FUnion{v} }

type FType_FExtType struct {
	Value string
}

func New_FType_FExtType(v string) FType { return FType_FExtType{v} }

type FType_FSlice struct {
	Value SliceType
}

func New_FType_FSlice(v SliceType) FType { return FType_FSlice{v} }

type FType_FPreUsed struct {
	Value string
}

func New_FType_FPreUsed(v string) FType { return FType_FPreUsed{v} }

type FType_FParametrized struct {
	Value string
}

func New_FType_FParametrized(v string) FType { return FType_FParametrized{v} }

type SliceType struct {
	elemType FType
}
type FuncType struct {
	targets []FType
}
type NameTypePair struct {
	name  string
	ftype FType
}
type RecordType struct {
	name   string
	fields []NameTypePair
}
type UnionType struct {
	name  string
	cases []NameTypePair
}

func IsUnresolved(ft FType) bool {
	switch (ft).(type) {
	case FType_FUnresolved:
		return true
	default:
		return false
	}
}

func fargs(ft FuncType) []FType {
	l := slice.Length[FType](ft.targets)
	return frt.Pipe[[]FType, []FType](ft.targets, (func(_r0 []FType) []FType { return slice.Take(frt.OpMinus[int](l, 1), _r0) }))
}

func freturn(ft FuncType) FType {
	return slice.Last[FType](ft.targets)
}

func FFuncToGo(ft FuncType, toGo func(FType) string) string {
	last := slice.Last[FType](ft.targets)
	args := fargs(ft)
	bw := buf.New()
	buf.Write(bw, "func (")
	frt.PipeUnit[string](frt.Pipe[[]string, string](frt.Pipe[[]FType, []string](args, (func(_r0 []FType) []string { return slice.Map(toGo, _r0) })), (func(_r0 []string) string { return strings.Concat(",", _r0) })), (func(_r0 string) { buf.Write(bw, _r0) }))
	buf.Write(bw, ")")
	ret := (func() string {
		switch (last).(type) {
		case FType_FUnit:
			return ""
		default:
			return frt.OpPlus[string](" ", toGo(last))
		}
	})()
	buf.Write(bw, ret)
	return buf.String(bw)
}

func FRecordToGo(frec RecordType) string {
	return frec.name
}

func frStructName(frec RecordType) string {
	return frec.name
}

func namePairMatch(targetName string, pair NameTypePair) bool {
	return frt.OpEqual[string](targetName, pair.name)
}

func frGetField(frec RecordType, fieldName string) NameTypePair {
	res := frt.Pipe[[]NameTypePair, []NameTypePair](frec.fields, (func(_r0 []NameTypePair) []NameTypePair {
		return slice.Filter((func(_r0 NameTypePair) bool { return namePairMatch(fieldName, _r0) }), _r0)
	}))
	return slice.Head[NameTypePair](res)
}

func FUnionToGo(fu UnionType) string {
	return fu.name
}

func FSliceToGo(fs SliceType, toGo func(FType) string) string {
	return frt.OpPlus[string]("[]", toGo(fs.elemType))
}

func FTypeToGo(ft FType) string {
	switch _v17 := (ft).(type) {
	case FType_FInt:
		return "int"
	case FType_FBool:
		return "bool"
	case FType_FString:
		return "string"
	case FType_FUnit:
		return ""
	case FType_FUnresolved:
		return ""
	case FType_FFunc:
		ft := _v17.Value
		return FFuncToGo(ft, FTypeToGo)
	case FType_FRecord:
		fr := _v17.Value
		return FRecordToGo(fr)
	case FType_FUnion:
		fu := _v17.Value
		return FUnionToGo(fu)
	case FType_FExtType:
		fe := _v17.Value
		return fe
	case FType_FSlice:
		fs := _v17.Value
		return FSliceToGo(fs, FTypeToGo)
	case FType_FPreUsed:
		fp := _v17.Value
		return fp
	case FType_FParametrized:
		fp := _v17.Value
		return fp
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
