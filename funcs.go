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
	funcSig  *FuncSigBuilder
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

func (b *FuncBuilder) Sig(sig *FuncSigBuilder) (r *FuncBuilder) {
	b.funcSig = sig
	return b
}

func (b *FuncBuilder) Blocks(blocks ...Code) (r *FuncBuilder) {
	b.blocks = append(b.blocks, blocks...)
	return b
}

func (b *FuncBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	buf := bytes.NewBuffer(nil)

	if b.funcSig != nil {
		err = Fprint(buf, b.funcSig, ctx)
		if err != nil {
			return
		}
	} else {
		if strings.Index(b.sig, "func") != 0 {
			if len(b.receiver) > 0 {
				b.sig = fmt.Sprintf("%s %s", b.receiver, b.sig)
			} else {
				b.sig = fmt.Sprintf("func %s", b.sig)
			}
		}

		err = Fprint(buf, Block(b.sig, b.vars...), ctx)
		if err != nil {
			return
		}
	}

	buf.WriteString(" {\n")

	err = Fprint(buf, Codes(b.blocks...), ctx)
	if err != nil {
		return
	}
	buf.WriteString("}\n\n")

	r = buf.Bytes()
	return
}

type FuncSigBuilder struct {
	name        string
	receiverVar string
	receiverTyp string
	parameters  []string
	results     []string
}

func FuncSig(name string) (r *FuncSigBuilder) {
	r = &FuncSigBuilder{}
	r.name = name
	return
}

func (b *FuncSigBuilder) Parameters(vs ...string) (r *FuncSigBuilder) {
	b.parameters = append(b.parameters, vs...)
	return b
}

func (b *FuncSigBuilder) Results(vs ...string) (r *FuncSigBuilder) {
	b.results = append(b.results, vs...)
	return b
}

func (b *FuncSigBuilder) Receiver(varName string, typ string) (r *FuncSigBuilder) {
	b.receiverVar = varName
	b.receiverTyp = typ
	return b
}

func (b *FuncSigBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	buf := bytes.NewBuffer(nil)
	buf.WriteString("func ")
	if len(b.receiverTyp) > 0 {
		if len(b.receiverVar) == 0 {
			b.receiverVar = "this"
		}
		buf.WriteString(fmt.Sprintf("(%s %s) ", b.receiverVar, b.receiverTyp))
	}

	if len(b.parameters)%2 != 0 {
		b.parameters = append(b.parameters, "interface{}")
	}

	params := []string{}
	for i := 0; i < len(b.parameters); i = i + 2 {
		params = append(params, fmt.Sprintf("%s %s", b.parameters[i], b.parameters[i+1]))
	}
	buf.WriteString(b.name + " (")
	buf.WriteString(strings.Join(params, ", "))
	buf.WriteString(")")

	if len(b.results) > 0 {
		buf.WriteString("(")
		rs := []string{}
		for i := 0; i < len(b.results); i = i + 2 {
			rs = append(rs, fmt.Sprintf("%s %s", b.results[i], b.results[i+1]))
		}
		buf.WriteString(strings.Join(rs, ", "))
		buf.WriteString(")")
	}

	r = buf.Bytes()
	return
}
