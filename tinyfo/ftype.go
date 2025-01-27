package main

import (
	"bytes"
	"fmt"
)

type FType interface {
	ftype()
	ToGo() string // Goでの型を表す文字列、表せないものをどうするかはあとで考える
}

func (*FPrimitive) ftype()  {}
func (*FUnitType) ftype()   {}
func (*FFunc) ftype()       {}
func (*FUnresolved) ftype() {}
func (*FRecord) ftype()     {}
func (*FUnion) ftype()      {}
func (*FCustom) ftype()     {}
func (*FSlice) ftype()      {}

func IsUnresolved(ft FType) bool {
	_, ok := ft.(*FUnresolved)
	return ok
}

func IsCustom(ft FType) bool {
	_, ok := ft.(*FCustom)
	return ok
}

type FPrimitive struct {
	Name string // "int", "string" etc.
}

func (p *FPrimitive) ToGo() string {
	return p.Name
}

type FUnitType struct {
}

func (p *FUnitType) ToGo() string {
	// unitはなにも無しとしておく。関数のreturnなどはそれで動くので。
	return ""
}

var (
	FInt    = &FPrimitive{"int"}
	FString = &FPrimitive{"string"}
	FBool   = &FPrimitive{"bool"}
	FUnit   = &FUnitType{}
)

// type inferenceでまだ未解決の状態。
// 最終的にはこれは無くならないとinference error。
type FUnresolved struct {
}

// ToGoの段階では全てresolveされていないと型エラーなので、これはコンパイルエラーのケースとなる。
// コンパイルエラーをどう扱うかはあとで考える。
// panicの方がいいかもしれないが、とりあえず""を返しておく。
func (p *FUnresolved) ToGo() string {
	return ""
}

/*
Custom type. Mostly user defined type.
Some type is defined after using.
So first keep identifier as is, then resolve.
*/
type FCustom struct {
	name string
}

func (p *FCustom) ToGo() string {
	return p.name
}

type FSlice struct {
	elemType FType
}

func (s *FSlice) ToGo() string {
	return "[]" + s.elemType.ToGo()
}

type FFunc struct {
	// カリー化されている前提で、引数も戻りも区別せず持つ
	Targets []FType
}

func (p *FFunc) Args() []FType {
	return p.Targets[0 : len(p.Targets)-1]
}

func (p *FFunc) ReturnType() FType {
	return p.Targets[len(p.Targets)-1]
}

func (p *FFunc) ToGo() string {
	var buf bytes.Buffer
	buf.WriteString("func (")
	for i, tp := range p.Args() {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(tp.ToGo())
	}
	buf.WriteString(")")

	last := p.ReturnType()
	if last != FUnit {
		buf.WriteString(" ")
		buf.WriteString(last.ToGo())
	}
	return buf.String()
}

func (f *FFunc) String() string {
	var buf bytes.Buffer
	for i, ft := range f.Targets {
		if i != 0 {
			buf.WriteString(" -> ")
		}
		// ToGo is not good.
		// It should be folang type name. But currently there are no such method.
		buf.WriteString(ft.ToGo())
	}
	return buf.String()
}

type NameTypePair struct {
	Name string
	Type FType
}

type FRecord struct {
	name   string
	fields []NameTypePair
}

func (fr *FRecord) ToGo() string {
	return fr.name
}

func (fr *FRecord) StructName() string {
	return fr.name
}

func (fr *FRecord) Match(fieldNames []string) bool {
	if len(fieldNames) != len(fr.fields) {
		return false
	}
	m := make(map[string]bool)
	for _, f := range fr.fields {
		m[f.Name] = true
	}

	for _, fn := range fieldNames {
		_, ok := m[fn]
		if !ok {
			return false
		}
	}
	return true
}

type FUnion struct {
	name  string
	cases []NameTypePair
}

// FUnion is interface in go. No need to add *.
func (fu *FUnion) ToGo() string {
	return fu.name
}

// Return IntOrString_I
func UnionCaseStructName(unionName string, caseName string) string {
	return fmt.Sprintf("%s_%s", unionName, caseName)
}
