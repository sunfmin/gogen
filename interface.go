package gogen

import (
	"bytes"
	"context"
	"fmt"
)

type InterfaceBuilder struct {
	name      string
	funcDecls []Code
}

func Interface(name string) (r *InterfaceBuilder) {
	r = &InterfaceBuilder{}
	r.name = name
	return
}

func (b *InterfaceBuilder) Block(template string, vars ...string) (r *InterfaceBuilder) {
	b.FuncDecls(Block(template, vars...))
	return b
}

func (b *InterfaceBuilder) AppendFuncDecl(template string, vars ...string) (r *InterfaceBuilder) {
	b.FuncDecls(Block(template, vars...))
	return b
}

func (b *InterfaceBuilder) FuncDecls(decls ...Code) (r *InterfaceBuilder) {
	b.funcDecls = append(b.funcDecls, decls...)
	return b
}

func (b *InterfaceBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("\ntype %s interface {\n", b.name))
	err = Fprint(buf, Codes(b.funcDecls...).Separator("\n"), ctx)
	if err != nil {
		return
	}
	buf.WriteString("\n}\n")
	r = buf.Bytes()
	return

}
