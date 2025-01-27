package main

type hoge struct {
	X string
	Y string
}

func ika() hoge {
	return hoge{X: "abc", Y: "def"}
}
