package env

import (
	"os"
	"strings"
)

// KeyIsInEnvironment inspects the runtime environment
// it returns the value for the given environment
// key if successful, otherwise it returns false.
// use the ok idiom to check for success.
func KeyIsInEnvironment(name string) (string, bool) {
	for _, str := range os.Environ() {
		before, _, ok := strings.Cut(str, "=")
		if ok && before == name {
			return before, true
		}
	}
	return "", false
}
