package strings

import "bytes"

func Concat(sep string, strs []string) string {
	var buf bytes.Buffer

	for i, s := range strs {
		if i != 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(s)
	}
	return buf.String()
}

func Length(str string) int {
	return len(str)
}

func AppendTail(tail string, s string) string {
	return s + tail
}
