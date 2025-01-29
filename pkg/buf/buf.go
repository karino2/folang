package buf

import "bytes"

/*
	Byte buffer wrapper.
*/

type Buffer = *bytes.Buffer

func New() Buffer { return &bytes.Buffer{} }

func Write(b Buffer, s string) {
	b.WriteString(s)
}

func String(b Buffer) string {
	return b.String()
}
