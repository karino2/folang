package slice

func Length[T any](s []T) int {
	return len(s)
}

func Last[T any](s []T) T {
	return s[len(s)-1]
}

func Head[T any](s []T) T {
	if len(s) == 0 {
		panic("call Head to empty list")
	}
	return s[0]
}

func Take[T any](num int, s []T) []T {
	var res []T

	for i := 0; i < num; i++ {
		res = append(res, s[i])
	}

	return res
}

func Map[T any, U any](f func(T) U, s []T) []U {
	var res []U
	for _, e := range s {
		res = append(res, f(e))
	}
	return res
}

func Filter[T any](pred func(T) bool, s []T) []T {
	var res []T
	for _, e := range s {
		if pred(e) {
			res = append(res, e)
		}
	}
	return res

}
