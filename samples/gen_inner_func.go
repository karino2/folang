package main

import "github.com/karino2/folang/pkg/frt"

func main() {
	inner := func(s string) string {
		switch v := (s); v {
		case "abc":
			return "hit"
		default:
			return frt.SInterP("%s does not hit", v)
		}
	}
	frt.Println(inner("abc"))
}
