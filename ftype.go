package main

import "bytes"

type FType interface {
	ftype()
	ToGo() string // Goでの型を表す文字列、表せないものをどうするかはあとで考える
}

func (*FPrimitive) ftype()  {}
func (*FUnitType) ftype()   {}
func (*FFunc) ftype()       {}
func (*FUnresolved) ftype() {}
func (*FRecord) ftype()     {}

func IsUnresolved(ft FType) bool {
	_, ok := ft.(*FUnresolved)
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

var (
	FInt    = &FPrimitive{"int"}
	FString = &FPrimitive{"string"}
	FUnit   = &FUnitType{}
)

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

type RecordField struct {
	Name string
	Type FType
}

type FRecord struct {
	name   string
	fields []RecordField
}

// Use for cast and variable declaration.
// Use pointer type as record type.
func (fr *FRecord) ToGo() string {
	return "*" + fr.name
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
