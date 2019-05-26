package gogen

import (
	"context"
	"strings"
)

type SnippetBuilder struct {
	template    string
	vars        []string
	varSnippets map[string]Code
}

func Snippet(template string, vars ...string) (r *SnippetBuilder) {
	r = &SnippetBuilder{
		template:    template,
		vars:        vars,
		varSnippets: make(map[string]Code),
	}
	return
}

func (b *SnippetBuilder) Var(varName string, val string) (r *SnippetBuilder) {
	b.vars = append(b.vars, varName, val)
	return b
}

func (b *SnippetBuilder) VarCode(varName string, c Code) (r *SnippetBuilder) {
	b.varSnippets[varName] = c
	return b
}

func (b *SnippetBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {
	vars := b.vars
	if len(vars)%2 != 0 {
		vars = append(vars, "")
	}
	val := b.template

	for i := 0; i < len(vars); i = i + 2 {
		val = strings.ReplaceAll(val, vars[i], vars[i+1])
	}

	for varName, c := range b.varSnippets {
		val = strings.ReplaceAll(val, varName, MustString(c, ctx))
	}

	r = []byte(val)
	return
}
