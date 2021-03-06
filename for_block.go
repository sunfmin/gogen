package gogen

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type ForBlockBuilder struct {
	cond   Code
	blocks []Code
}

func ForBlock(forTemplate string, vars ...string) (r *ForBlockBuilder) {
	r = &ForBlockBuilder{}
	if strings.Index(forTemplate, "for") != 0 {
		forTemplate = fmt.Sprintf("%s %s", "for", forTemplate)
	}

	r.cond = Snippet(forTemplate, vars...)
	return
}

func (b *ForBlockBuilder) Body(blocks ...Code) (r *ForBlockBuilder) {
	b.blocks = append(b.blocks, blocks...)
	return b
}

func (b *ForBlockBuilder) BodySnippet(template string, vars ...string) (r *ForBlockBuilder) {
	b.Body(Snippet(template, vars...))
	return b
}

func (b *ForBlockBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	buf := bytes.NewBuffer(nil)
	err = Fprint(buf, b.cond, ctx)
	if err != nil {
		return
	}
	buf.WriteString(" {\n")
	err = Fprint(buf, Snippets(b.blocks...), ctx)
	if err != nil {
		panic(err)
	}
	buf.WriteString("}")

	r = buf.Bytes()
	return
}
