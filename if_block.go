package gogen

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type condBlocks struct {
	cond   Code
	blocks []Code
}
type IfBlockBuilder struct {
	ifBlock        *condBlocks
	elseIfBlocks   []*condBlocks
	elseBlocks     []Code
	lastCondBlocks *condBlocks
}

func IfBlock(ifTemplate string, vars ...string) (r *IfBlockBuilder) {
	r = &IfBlockBuilder{}
	if strings.Index(ifTemplate, "if") != 0 {
		ifTemplate = fmt.Sprintf("%s %s", "if", ifTemplate)
	}

	r.ifBlock = &condBlocks{cond: Snippet(ifTemplate, vars...)}
	r.lastCondBlocks = r.ifBlock
	return
}

func (b *IfBlockBuilder) Then(blocks ...Code) (r *IfBlockBuilder) {
	b.lastCondBlocks.blocks = append(b.lastCondBlocks.blocks, blocks...)
	return b
}

func (b *IfBlockBuilder) ThenSnippet(template string, vars ...string) (r *IfBlockBuilder) {
	b.Then(Snippet(template, vars...))
	return b
}

func (b *IfBlockBuilder) ElseIf(elseIfTemplate string, vars ...string) (r *IfBlockBuilder) {
	if strings.Index(elseIfTemplate, "else if") != 0 {
		elseIfTemplate = fmt.Sprintf("%s %s", "else if", elseIfTemplate)
	}

	b.lastCondBlocks = &condBlocks{cond: Snippet(elseIfTemplate, vars...)}
	b.elseIfBlocks = append(b.elseIfBlocks, b.lastCondBlocks)
	return b
}

func (b *IfBlockBuilder) Else(blocks ...Code) (r *IfBlockBuilder) {
	b.elseBlocks = append(b.elseBlocks, blocks...)
	return b
}

func (b *IfBlockBuilder) ElseSnippet(template string, vars ...string) (r *IfBlockBuilder) {
	b.Else(Snippet(template, vars...))
	return b
}

func (b *IfBlockBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	buf := bytes.NewBuffer(nil)
	err = Fprint(buf, b.ifBlock.cond, ctx)
	if err != nil {
		return
	}
	buf.WriteString(" {\n")
	err = Fprint(buf, Snippets(b.ifBlock.blocks...), ctx)
	if err != nil {
		panic(err)
	}
	buf.WriteString("} ")

	for _, elsIf := range b.elseIfBlocks {
		err = Fprint(buf, elsIf.cond, ctx)
		if err != nil {
			return
		}

		buf.WriteString(" {\n")
		err = Fprint(buf, Snippets(elsIf.blocks...), ctx)
		if err != nil {
			panic(err)
		}
		buf.WriteString("} ")

	}

	if len(b.elseBlocks) > 0 {
		buf.WriteString(" else {\n")
		err = Fprint(buf, Snippets(b.elseBlocks...), ctx)
		if err != nil {
			panic(err)
		}

		buf.WriteString("}\n")
	}

	r = buf.Bytes()
	return
}
