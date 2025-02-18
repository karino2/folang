package strings

import (
	"bytes"
	syssr "strings"
)

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

func HasSuffix(suffix string, s string) bool {
	return syssr.HasSuffix(s, suffix)
}

func TrimSuffix(suffix string, s string) string {
	return syssr.TrimSuffix(s, suffix)
}
