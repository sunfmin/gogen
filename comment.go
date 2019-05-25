package gogen

import (
	"fmt"
	"strings"
)

func LineComment(comment string) (r Code) {
	if len(comment) == 0 {
		return
	}
	lines := strings.Split(comment, "\n")
	commentLines := []string{}
	for _, l := range lines {
		commentLines = append(commentLines, fmt.Sprintf("// %s", l))
	}
	return RawCode(strings.Join(commentLines, "\n"))
}

func BlockComment(comment string) (r Code) {
	if len(comment) == 0 {
		return
	}

	return RawCode(fmt.Sprintf("/*\n%s\n*/\n", comment))
}
