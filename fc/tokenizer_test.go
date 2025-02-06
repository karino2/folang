package main

import "testing"

func TestTokenizerNormal(t *testing.T) {
	src := `package main`

	tkz := newTkz(src)
	if tkz.col != 0 {
		t.Errorf("col not 0, but %d", tkz.col)
	}
	if tkz.current.ttype != New_TokenType_PACKAGE {
		t.Errorf("first token is not package: %T", tkz.current.ttype)
	}
	tkz = tkzNext(tkz)
	if tkz.col != 8 {
		t.Errorf("next col not 8, but %d", tkz.col)
	}
	if tkz.current.ttype != New_TokenType_IDENTIFIER {
		t.Errorf("second token is not identifier: %T", tkz.current.ttype)
	}
	if tkz.current.stringVal != "main" {
		t.Errorf("second token is not main: %s", tkz.current.stringVal)
	}
}

func TestTokenizerBeginEOL(t *testing.T) {
	src := `  // space
package main`

	tkz := newTkz(src)
	if tkz.current.ttype != New_TokenType_EOL {
		t.Errorf("First empty line token is not EOL: %T", tkz.current.ttype)
	}
	tkz = tkzNext(tkz)
	if tkz.col != 0 {
		t.Errorf("after EOL, col not 0 but %d", tkz.col)
	}
	if tkz.current.ttype != New_TokenType_PACKAGE {
		t.Errorf("after EOL, token is not package: %T", tkz.current.ttype)
	}
}
