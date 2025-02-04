package main

import "github.com/karino2/folang/pkg/frt"

func ika() string {
	return frt.Snd[int, string](frt.Tuple2[int, string]{123, "abc"})
}

func main() {
	s := ika()
	frt.Println(s)
}
