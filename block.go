package gogen

import (
	"context"
	"strings"
)

type BlockBuilder struct {
	template  string
	vars      []string
	varBlocks map[string]Code
}

func Block(template string, vars ...string) (r *BlockBuilder) {
	r = &BlockBuilder{
		template:  template,
		vars:      vars,
		varBlocks: make(map[string]Code),
	}
	return
}

func (b *BlockBuilder) Var(varName string, val string) (r *BlockBuilder) {
	b.vars = append(b.vars, varName, val)
	return b
}

func (b *BlockBuilder) VarBlock(varName string, c Code) (r *BlockBuilder) {
	b.varBlocks[varName] = c
	return b
}

func (b *BlockBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {
	vars := b.vars
	if len(vars)%2 != 0 {
		vars = append(vars, "")
	}
	val := b.template

	for i := 0; i < len(vars); i = i + 2 {
		val = strings.ReplaceAll(val, vars[i], vars[i+1])
	}

	for varName, c := range b.varBlocks {
		val = strings.ReplaceAll(val, varName, MustString(c, ctx))
	}

	r = []byte(val)
	return
}
