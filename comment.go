package gogen

import (
	"fmt"
	"strings"
)

type InterfaceBuilder struct {
	name      string
	funcDecls []Code
}

func LineComment(comment string) (r Code) {
	lines := strings.Split(comment, "\n")
	commentLines := []string{}
	for _, l := range lines {
		commentLines = append(commentLines, fmt.Sprintf("// %s", l))
	}
	return RawCode(strings.Join(commentLines, "\n"))
}

func BlockComment(comment string) (r Code) {
	return RawCode(fmt.Sprintf("/*\n%s\n*/\n", comment))
}
