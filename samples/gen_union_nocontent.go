package main

import "fmt"

type AorB interface {
	AorB_Union()
}

func (AorB_A) AorB_Union() {}
func (AorB_B) AorB_Union() {}

func (v AorB_A) String() string { return "(A)" }
func (v AorB_B) String() string { return "(B)" }

type AorB_A struct {
}

var New_AorB_A AorB = AorB_A{}

type AorB_B struct {
}

var New_AorB_B AorB = AorB_B{}

func ika(ab AorB) string {
	switch (ab).(type) {
	case AorB_A:
		return "a match"
	case AorB_B:
		return "b match"
	default:
		panic("Union pattern fail. Never reached here.")
	}
}
