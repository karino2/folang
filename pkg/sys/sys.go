package sys

import (
	"os"

	"github.com/karino2/folang/pkg/frt"
)

func Args() []string {
	return os.Args
}

func ReadFile(file string) frt.Tuple2[string, bool] {
	cont, err := os.ReadFile(file)
	ok := err == nil
	return frt.NewTuple2(string(cont), ok)
}

func WriteFile(file string, content string) bool {
	err := os.WriteFile(file, []byte(content), 0644)
	ok := err == nil
	return ok
}
