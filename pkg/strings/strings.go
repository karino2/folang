package strings

import (
	"bytes"
	sysstr "strings"
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

func AppendHead(head string, s string) string {
	return head + s
}

func HasSuffix(suffix string, s string) bool {
	return sysstr.HasSuffix(s, suffix)
}

func TrimSuffix(suffix string, s string) string {
	return sysstr.TrimSuffix(s, suffix)
}

func EncloseWith(beg string, end string, center string) string {
	return beg + center + end
}

func Split(sep string, cont string) []string {
	return sysstr.Split(cont, sep)
}

func SplitN(count int, sep string, cont string) []string {
	return sysstr.SplitN(cont, sep, count)
}

func IsEmpty(s string) bool {
	return s == ""
}

func IsNotEmpty(s string) bool {
	return s != ""
}
