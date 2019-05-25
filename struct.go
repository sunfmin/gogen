package gogen

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type StructBuilder struct {
	name             string
	fields           []Code
	funcs            []Code
	funcsReceiverVar string
	pointer          bool
}

func Struct(name string) (r *StructBuilder) {
	r = &StructBuilder{}
	r.name = name
	r.pointer = true
	return
}

func (b *StructBuilder) Block(template string, vars ...string) (r *StructBuilder) {
	b.Fields(Block(template, vars...))
	return b
}

func (b *StructBuilder) Fields(imps ...Code) (r *StructBuilder) {
	b.fields = append(b.fields, imps...)
	return b
}

func (b *StructBuilder) Funcs(funcs ...Code) (r *StructBuilder) {
	b.funcs = append(b.funcs, funcs...)
	return b
}

func (b *StructBuilder) ReceiverVar(varName string) (r *StructBuilder) {
	b.funcsReceiverVar = varName
	return b
}

func (b *StructBuilder) Pointer(v bool) (r *StructBuilder) {
	b.pointer = v
	return b
}

func (b *StructBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("type %s struct {\n", b.name))
	err = Fprint(buf, Codes(b.fields...), ctx)
	if err != nil {
		return
	}
	buf.WriteString("}\n")

	for _, f := range b.funcs {
		if fb, ok := f.(*FuncBuilder); ok {
			varName := b.funcsReceiverVar
			if len(varName) == 0 {
				varName = "this"
			}
			tp := b.name
			if b.pointer {
				tp = "*" + b.name
			}
			fb.Receiver(varName, tp)
		}
	}

	err = Fprint(buf, Codes(b.funcs...), ctx)
	if err != nil {
		return
	}
	r = buf.Bytes()
	return

}

func Tag(key string, val string) (r string) {
	r = fmt.Sprintf(`%s "%s"`, key, val)
	return
}

func Field(name, typ string, tags ...string) (r Code) {
	theTags := ""
	if len(tags) > 0 {
		theTags = " `" + strings.Join(tags, " ") + "`"
	}

	r = RawCode(fmt.Sprintf("%s %s%s",
		name,
		typ,
		theTags,
	))
	return
}
