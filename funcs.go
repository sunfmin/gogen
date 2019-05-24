package gogen

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type FuncBuilder struct {
	receiver string
	sig      string
	vars     []string
	blocks   []Code
}

func Func(sig string, vars ...string) (r *FuncBuilder) {
	r = &FuncBuilder{}
	r.sig = sig
	r.vars = vars
	return
}

func (b *FuncBuilder) Block(template string, vars ...string) (r *FuncBuilder) {
	b.Blocks(Block(template, vars...))
	return b
}

func (b *FuncBuilder) Receiver(varName string, typ string) (r *FuncBuilder) {
	b.receiver = fmt.Sprintf("func (%s %s)", varName, typ)
	return b
}

func (b *FuncBuilder) Blocks(blocks ...Code) (r *FuncBuilder) {
	b.blocks = append(b.blocks, blocks...)
	return b
}

func (b *FuncBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	buf := bytes.NewBuffer(nil)

	if strings.Index(b.sig, "func") != 0 && len(b.receiver) > 0 {
		b.sig = fmt.Sprintf("%s %s", b.receiver, b.sig)
	}

	err = Fprint(buf, Block(b.sig, b.vars...), ctx)
	if err != nil {
		return
	}

	buf.WriteString(" {\n")

	err = Fprint(buf, Codes(b.blocks...).Separator("\n"), ctx)
	if err != nil {
		return
	}
	buf.WriteString("\n}\n\n")

	r = buf.Bytes()
	return
}
