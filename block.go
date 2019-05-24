package gogen

import (
	"context"
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
		r = []byte(val)
		return
	})
}
