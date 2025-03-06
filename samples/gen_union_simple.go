package main

type IntOrString interface {
	IntOrString_Union()
}

func (IntOrString_I) IntOrString_Union() {}
func (IntOrString_S) IntOrString_Union() {}

func (v IntOrString_I) String() string { return frt.Sprintf1("(I: %v)", v.Value) }
func (v IntOrString_S) String() string { return frt.Sprintf1("(S: %v)", v.Value) }

type IntOrString_I struct {
	Value int
}

func New_IntOrString_I(v int) IntOrString { return IntOrString_I{v} }

type IntOrString_S struct {
	Value string
}

func New_IntOrString_S(v string) IntOrString { return IntOrString_S{v} }

func ika() IntOrString {
	return New_IntOrString_I(123)
}
