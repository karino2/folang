package buf

import "testing"

func TestBufBasic(t *testing.T) {
	b := New()
	Write(b, "hoge")
	Write(b, "ika")
	got := String(b)
	want := "hogeika"
	if got != want {
		t.Errorf("want %s, got %s", want, got)
	}

}
