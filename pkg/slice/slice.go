package slice

func Length[T any](s []T) int {
	return len(s)
}

func Last[T any](s []T) T {
	return s[len(s)-1]
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
