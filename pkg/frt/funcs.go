package frt

import "fmt"

func Pipe[T any, U any](elem T, f func(T) U) U {
	return f(elem)
}

func Println(str string) {
	fmt.Println(str)
}
