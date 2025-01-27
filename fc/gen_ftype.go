package main

import "bytes"

import "fmt"

type FType interface {
	FType_Union()
}

func (FType_FInt) FType_Union()        {}
func (FType_FString) FType_Union()     {}
func (FType_FUnit) FType_Union()       {}
func (FType_FUnresolved) FType_Union() {}
func (FType_FFunc) FType_Union()       {}
func (FType_FRecord) FType_Union()     {}
func (FType_FUnion) FType_Union()      {}

type FType_FInt struct {
}

var New_FType_FInt FType = FType_FInt{}

type FType_FString struct {
}

var New_FType_FString FType = FType_FString{}

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

type NameTypePair struct {
	Name string
	Type FType
}

type RecordType struct {
	name   string
	fiedls []NameTypePair
}

type UnionType struct {
	name  string
	cases []NameTypePair
}

type FuncType struct {
	targets []FType
}

func IsUnresolved(ft FType) {
	switch (ft).(type) {
	case FType_FUnresolved:
		return true
	default:
		return false
	}
}
