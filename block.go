package gogen

import (
	"context"
	"fmt"
	"strings"
)

func Block(template string, vars ...string) (r Code) {
	return CodeFunc(func(ctx context.Context) (r []byte, err error) {
		if len(vars)%2 != 0 {
			vars = append(vars, "")
		}
		val := template

		for i := 0; i < len(vars); i = i + 2 {
			val = strings.ReplaceAll(val, vars[i], vars[i+1])
		}

		if len(strings.Split(val, "\""))%2 == 0 {
			panic(fmt.Sprintf("quote \" not match: %s", val))
		}

		if len(strings.Split(val, "`"))%2 == 0 {
			panic(fmt.Sprintf("quote `` not match: %s", val))
		}

		r = []byte(val)
		return
	})
}
