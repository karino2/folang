package main

import "testing"

func TestParsePackage(t *testing.T) {
	src := `package main
`
	ps := initParse(src)
	gotPair := parsePackage(ps)
	got := gotPair.E1

	switch tgot := got.(type) {
	case Stmt_Package:
		if tgot.Value != "main" {
			t.Errorf("expect main, got %s", tgot.Value)
		}
	default:
		t.Errorf("Unexpected stmt. Expect Package, got %T", got)
	}
}

func TestParseParams(t *testing.T) {
	src := `(a:int) (b:int) =
`
	ps := initParse(src)
	gotPair := parseParams(ps)
	got := gotPair.E1
	if len(got) != 2 {
		t.Errorf("want 2 param, got %d: %v", len(got), got)
	}
	// t.Errorf("%d, %T, %T", len(got), got[0], got[1])

	ps2 := gotPair.E0
	tt := psCurrentTT(ps2)
	if tt != New_TokenType_EQ {
		t.Errorf("want EQ, got %T", tt)
	}
}
