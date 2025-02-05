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

func fargs(ft FuncType) []FType {
	l := slice.Length(ft.targets)
	return frt.Pipe(ft.targets, (func(_r0 []FType) []FType { return slice.Take((l - 1), _r0) }))
}

func freturn(ft FuncType) FType {
	return slice.Last(ft.targets)
}

func funcTypeToGo(ft FuncType, toGo func(FType) string) string {
	last := slice.Last(ft.targets)
	args := fargs(ft)
	bw := buf.New()
	buf.Write(bw, "func (")
	frt.PipeUnit(frt.Pipe(frt.Pipe(args, (func(_r0 []FType) []string { return slice.Map(toGo, _r0) })), (func(_r0 []string) string { return strings.Concat(",", _r0) })), (func(_r0 string) { buf.Write(bw, _r0) }))
	buf.Write(bw, ")")
	ret := (func() string {
		switch (last).(type) {
		case FType_FUnit:
			return ""
		default:
			return (" " + toGo(last))
		}
	})()
	buf.Write(bw, ret)
	return buf.String(bw)
}

func recordTypeToGo(frec RecordType) string {
	return frec.name
}

func frStructName(frec RecordType) string {
	return frec.name
}

func namePairMatch(targetName string, pair NameTypePair) bool {
	return frt.OpEqual(targetName, pair.name)
}

func lookupPairByName(targetName string, pairs []NameTypePair) NameTypePair {
	res := frt.Pipe(pairs, (func(_r0 []NameTypePair) []NameTypePair {
		return slice.Filter((func(_r0 NameTypePair) bool { return namePairMatch(targetName, _r0) }), _r0)
	}))
	return slice.Head(res)
}

func frGetField(frec RecordType, fieldName string) NameTypePair {
	return lookupPairByName(fieldName, frec.fields)
}

func npName(pair NameTypePair) string {
	return pair.name
}

func frMatch(frec RecordType, fieldNames []string) bool {
	return frt.IfElse(frt.OpNotEqual(slice.Length(fieldNames), slice.Length(frec.fields)), (func() bool {
		return false
	}), (func() bool {
		sortedInput := frt.Pipe(fieldNames, slice.Sort)
		sortedFName := frt.Pipe(slice.Map(npName, frec.fields), slice.Sort)
		return frt.OpEqual(sortedInput, sortedFName)
	}))
}

func funionToGo(fu UnionType) string {
	return fu.name
}

func lookupCase(fu UnionType, caseName string) NameTypePair {
	return lookupPairByName(caseName, fu.cases)
}

func unionCaseStructName(unionName string, caseName string) string {
	return ((unionName + "_") + caseName)
}

func fSliceToGo(fs SliceType, toGo func(FType) string) string {
	return ("[]" + toGo(fs.elemType))
}

func FTypeToGo(ft FType) string {
	switch _v19 := (ft).(type) {
	case FType_FInt:
		return "int"
	case FType_FBool:
		return "bool"
	case FType_FString:
		return "string"
	case FType_FUnit:
		return ""
	case FType_FFunc:
		ft := _v19.Value
		return funcTypeToGo(ft, FTypeToGo)
	case FType_FRecord:
		fr := _v19.Value
		return recordTypeToGo(fr)
	case FType_FUnion:
		fu := _v19.Value
		return funionToGo(fu)
	case FType_FExtType:
		fe := _v19.Value
		return fe
	case FType_FSlice:
		fs := _v19.Value
		return fSliceToGo(fs, FTypeToGo)
	case FType_FPreUsed:
		fp := _v19.Value
		return fp
	case FType_FParametrized:
		fp := _v19.Value
		return fp
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
