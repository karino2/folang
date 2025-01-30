package slice

func Length[T any](s []T) int {
	return len(s)
}

func Take[T any](num int, s []T) []T {
	var res []T

	for i := 0; i < num; i++ {
		res = append(res, s[i])
	}

	return res
}
