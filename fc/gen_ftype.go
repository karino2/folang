package main

import "github.com/karino2/folang/pkg/frt"

import "github.com/karino2/folang/pkg/buf"

import "github.com/karino2/folang/pkg/slice"

import "github.com/karino2/folang/pkg/strings"

import "github.com/karino2/folang/pkg/dict"

type TypeVar struct {
	Name string
}

type FType interface {
	FType_Union()
}

func (FType_FInt) FType_Union()         {}
func (FType_FString) FType_Union()      {}
func (FType_FBool) FType_Union()        {}
func (FType_FFloat) FType_Union()       {}
func (FType_FUnit) FType_Union()        {}
func (FType_FAny) FType_Union()         {}
func (FType_FFunc) FType_Union()        {}
func (FType_FRecord) FType_Union()      {}
func (FType_FUnion) FType_Union()       {}
func (FType_FSlice) FType_Union()       {}
func (FType_FTuple) FType_Union()       {}
func (FType_FFieldAccess) FType_Union() {}
func (FType_FTypeVar) FType_Union()     {}
func (FType_FParamd) FType_Union()      {}

type FType_FInt struct {
}

var New_FType_FInt FType = FType_FInt{}

type FType_FString struct {
}

var New_FType_FString FType = FType_FString{}

type FType_FBool struct {
}

var New_FType_FBool FType = FType_FBool{}

type FType_FFloat struct {
}

var New_FType_FFloat FType = FType_FFloat{}

type FType_FUnit struct {
}

var New_FType_FUnit FType = FType_FUnit{}

type FType_FAny struct {
}

var New_FType_FAny FType = FType_FAny{}

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

type FType_FSlice struct {
	Value SliceType
}

func New_FType_FSlice(v SliceType) FType { return FType_FSlice{v} }

type FType_FTuple struct {
	Value TupleType
}

func New_FType_FTuple(v TupleType) FType { return FType_FTuple{v} }

type FType_FFieldAccess struct {
	Value FieldAccessType
}

func New_FType_FFieldAccess(v FieldAccessType) FType { return FType_FFieldAccess{v} }

type FType_FTypeVar struct {
	Value TypeVar
}

func New_FType_FTypeVar(v TypeVar) FType { return FType_FTypeVar{v} }

type FType_FParamd struct {
	Value ParamdType
}

func New_FType_FParamd(v ParamdType) FType { return FType_FParamd{v} }

type SliceType struct {
	ElemType FType
}
type FuncType struct {
	Targets []FType
}
type ParamdType struct {
	Name  string
	Targs []FType
}
type FieldAccessType struct {
	RecType   FType
	FieldName string
}
type TupleType struct {
	ElemTypes []FType
}
type NameTypePair struct {
	Name  string
	Ftype FType
}
type RecordType struct {
	Name  string
	Targs []FType
}
type UnionType struct {
	Name  string
	Targs []FType
}

func fargs(ft FuncType) []FType {
	l := slice.Length(ft.Targets)
	return frt.Pipe(ft.Targets, (func(_r0 []FType) []FType { return slice.Take((l - 1), _r0) }))
}

func freturn(ft FuncType) FType {
	return slice.Last(ft.Targets)
}

func funcTypeToGo(ft FuncType, toGo func(FType) string) string {
	last := slice.Last(ft.Targets)
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

func newFFunc(ftypes []FType) FType {
	return frt.Pipe(FuncType{Targets: ftypes}, New_FType_FFunc)
}

func tArgsToGo[T0 any](tGo func(T0) string, targs []T0) string {
	return frt.IfElse(slice.IsEmpty(targs), (func() string {
		return ""
	}), (func() string {
		return frt.Pipe(frt.Pipe(frt.Pipe(targs, (func(_r0 []T0) []string { return slice.Map(tGo, _r0) })), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) string { return strings.EncloseWith("[", "]", _r0) }))
	}))
}

func recordTypeToGo(tGo func(FType) string, frec RecordType) string {
	return (frec.Name + tArgsToGo(tGo, frec.Targs))
}

func utName(ut UnionType) string {
	return ut.Name
}

func fUnionToGo(tGo func(FType) string, ut UnionType) string {
	return (ut.Name + tArgsToGo(tGo, ut.Targs))
}

func fSliceToGo(fs SliceType, toGo func(FType) string) string {
	return ("[]" + toGo(fs.ElemType))
}

func fTupleToGo(toGo func(FType) string, ft TupleType) string {
	args := frt.Pipe(slice.Map(toGo, ft.ElemTypes), (func(_r0 []string) string { return strings.Concat(", ", _r0) }))
	len := slice.Length(ft.ElemTypes)
	return frt.SInterP("frt.Tuple%s[%s]", len, args)
}

func encloseWith(beg string, end string, center string) string {
	return ((beg + center) + end)
}

func fpToGo(tToGo func(FType) string, pt ParamdType) string {
	return frt.IfElse(slice.IsEmpty(pt.Targs), (func() string {
		return pt.Name
	}), (func() string {
		return frt.Pipe(frt.Pipe(slice.Map(tToGo, pt.Targs), (func(_r0 []string) string { return strings.Concat(", ", _r0) })), (func(_r0 string) string { return encloseWith((pt.Name + "["), "]", _r0) }))
	}))
}

func FTypeToGo(ft FType) string {
	switch _v1 := (ft).(type) {
	case FType_FInt:
		return "int"
	case FType_FBool:
		return "bool"
	case FType_FFloat:
		return "float64"
	case FType_FAny:
		return "any"
	case FType_FString:
		return "string"
	case FType_FUnit:
		return ""
	case FType_FFunc:
		ft := _v1.Value
		return funcTypeToGo(ft, FTypeToGo)
	case FType_FRecord:
		fr := _v1.Value
		return recordTypeToGo(FTypeToGo, fr)
	case FType_FUnion:
		fu := _v1.Value
		return fUnionToGo(FTypeToGo, fu)
	case FType_FParamd:
		pt := _v1.Value
		return fpToGo(FTypeToGo, pt)
	case FType_FSlice:
		fs := _v1.Value
		return fSliceToGo(fs, FTypeToGo)
	case FType_FTuple:
		ft := _v1.Value
		return fTupleToGo(FTypeToGo, ft)
	case FType_FFieldAccess:
		return "FieldAccess_Unresoled"
	case FType_FTypeVar:
		fp := _v1.Value
		return fp.Name
	default:
		panic("Union pattern fail. Never reached here.")
	}
}

func frStructName(tGo func(FType) string, frec RecordType) string {
	return (frec.Name + tArgsToGo(tGo, frec.Targs))
}

type RecordTypeInfo struct {
	Fields []NameTypePair
}

var g_recInfoDic = dict.New[string, RecordTypeInfo]()

func encodedKey[T0 any](name T0, targs []FType) string {
	encts := frt.Pipe(slice.Map(FTypeToGo, targs), (func(_r0 []string) string { return strings.Concat("_", _r0) }))
	return frt.SInterP("%s_%s", name, encts)
}

func rtToKey(rt RecordType) string {
	return encodedKey(rt.Name, rt.Targs)
}

func lookupRecInfo(rt RecordType) RecordTypeInfo {
	ri, ok := frt.Destr(dict.TryFind(g_recInfoDic, rtToKey(rt)))
	frt.IfOnly(frt.OpNot(ok), (func() {
		frt.PipeUnit(frt.Sprintf1("Can't find record info: %s.", rt.Name), PanicNow)
	}))
	return ri
}

func updateRecInfo(rt RecordType, rinfo RecordTypeInfo) {
	dict.Add(g_recInfoDic, rtToKey(rt), rinfo)
}

func lookupPairByName(targetName string, pairs []NameTypePair) NameTypePair {
	res := frt.Pipe(pairs, (func(_r0 []NameTypePair) []NameTypePair {
		return slice.Filter(func(x NameTypePair) bool {
			return frt.OpEqual(x.Name, targetName)
		}, _r0)
	}))
	frt.IfOnly(slice.IsEmpty(res), (func() {
		frt.PipeUnit(frt.Sprintf1("Can't find record field of: %s", targetName), PanicNow)
	}))
	return slice.Head(res)
}

func frGetField(frec RecordType, fieldName string) NameTypePair {
	ri := lookupRecInfo(frec)
	return lookupPairByName(fieldName, ri.Fields)
}

func newNTPair(name string, ft FType) NameTypePair {
	return NameTypePair{Name: name, Ftype: ft}
}

func tupToNTPair(tup frt.Tuple2[string, FType]) NameTypePair {
	return newNTPair(frt.Fst(tup), frt.Snd(tup))
}

func frMatch(rt RecordType, fieldNames []string) bool {
	ri := lookupRecInfo(rt)
	return frt.IfElse(frt.OpNotEqual(slice.Length(fieldNames), slice.Length(ri.Fields)), (func() bool {
		return false
	}), (func() bool {
		sortedInput := frt.Pipe(fieldNames, slice.Sort)
		sortedFName := frt.Pipe(slice.Map(func(_v1 NameTypePair) string {
			return _v1.Name
		}, ri.Fields), slice.Sort)
		return frt.OpEqual(sortedInput, sortedFName)
	}))
}

func faResolve(fat FieldAccessType) FType {
	switch _v2 := (fat.RecType).(type) {
	case FType_FRecord:
		rt := _v2.Value
		field := frGetField(rt, fat.FieldName)
		return field.Ftype
	default:
		return frt.Pipe(fat, New_FType_FFieldAccess)
	}
}

type UnionTypeInfo struct {
	Cases []NameTypePair
}

var g_uniInfoDic = dict.New[string, UnionTypeInfo]()

func uniToKey(ut UnionType) string {
	return encodedKey(ut.Name, ut.Targs)
}

func lookupUniInfo(ut UnionType) UnionTypeInfo {
	ui, ok := frt.Destr(dict.TryFind(g_uniInfoDic, uniToKey(ut)))
	frt.IfOnly(frt.OpNot(ok), (func() {
		frt.PipeUnit(frt.Sprintf1("Can't find union info: %s.", ut.Name), PanicNow)
	}))
	return ui
}

func updateUniInfo(ut UnionType, uinfo UnionTypeInfo) {
	dict.Add(g_uniInfoDic, uniToKey(ut), uinfo)
}

func utCases(ut UnionType) []NameTypePair {
	ui := lookupUniInfo(ut)
	return ui.Cases
}

func lookupCase(fu UnionType, caseName string) NameTypePair {
	return frt.Pipe(utCases(fu), (func(_r0 []NameTypePair) NameTypePair { return lookupPairByName(caseName, _r0) }))
}

func unionCSName(unionName string, caseName string) string {
	return ((unionName + "_") + caseName)
}
