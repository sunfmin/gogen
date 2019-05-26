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
	return Raw(strings.Join(commentLines, "\n"))
}

func BlockComment(comment string) (r Code) {
	if len(comment) == 0 {
		return
	}

	return Raw(fmt.Sprintf("/*\n%s\n*/\n", comment))
}
