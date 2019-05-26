package gogen

import (
	"bytes"
	"context"
	"fmt"
)

type InterfaceBuilder struct {
	name string
	body []Code
}

func Interface(name string) (r *InterfaceBuilder) {
	r = &InterfaceBuilder{}
	r.name = name
	return
}

func (b *InterfaceBuilder) BodySnippet(template string, vars ...string) (r *InterfaceBuilder) {
	b.Body(Snippet(template, vars...))
	return b
}

func (b *InterfaceBuilder) Body(decls ...Code) (r *InterfaceBuilder) {
	b.body = append(b.body, decls...)
	return b
}

func (b *InterfaceBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("\ntype %s interface {\n", b.name))
	err = Fprint(buf, Snippets(b.body...), ctx)
	if err != nil {
		return
	}
	buf.WriteString("}\n")
	r = buf.Bytes()
	return

}
