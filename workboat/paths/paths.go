package paths

import (
	"fmt"
	"strings"
)

const (
	Index = "/"
)

func Make(path string, replacements ...any) string {
	x := strings.Split(path, "/")
	for i, p := range x {
		if strings.HasPrefix(p, ":") {
			x[i] = "%s"
		}
	}
	return fmt.Sprintf(strings.Join(x, "/"), replacements...)
}
