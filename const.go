package gogen

import (
	"bytes"
	"context"
	"fmt"
)

type ConstBlockBuilder struct {
	name           string
	consts         []Code
	constType      string
	toType         string
	namePrefixType bool
}

func ConstBlock() (r *ConstBlockBuilder) {
	r = &ConstBlockBuilder{}
	return
}

func (b *ConstBlockBuilder) Type(constType string, toType string) (r *ConstBlockBuilder) {
	b.constType = constType
	b.toType = toType
	b.NamePrefixType(true)
	return b
}

func (b *ConstBlockBuilder) NamePrefixType(v bool) (r *ConstBlockBuilder) {
	b.namePrefixType = v
	return b
}

func (b *ConstBlockBuilder) Block(template string, vars ...string) (r *ConstBlockBuilder) {
	b.Consts(Block(template, vars...))
	return b
}

func (b *ConstBlockBuilder) AppendConst(name string, val interface{}) (r *ConstBlockBuilder) {
	if len(b.constType) == 0 {
		panic("To use AppendConst, call Type first")
	}

	if b.namePrefixType {
		name = fmt.Sprintf("%s%s", b.constType, name)
	}

	b.Consts(
		RawCode(fmt.Sprintf("%s %s = %#+v", name, b.constType, val)),
	)
	return b
}

func (b *ConstBlockBuilder) Consts(cs ...Code) (r *ConstBlockBuilder) {
	b.consts = append(b.consts, cs...)
	return b
}

func (b *ConstBlockBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	buf := bytes.NewBuffer(nil)

	if len(b.constType) > 0 {
		buf.WriteString(fmt.Sprintf("type %s %s\n", b.constType, b.toType))
	}
	buf.WriteString("const (\n")
	err = Fprint(buf, Codes(b.consts...), ctx)
	if err != nil {
		return
	}
	buf.WriteString(")\n")

	r = buf.Bytes()
	return
}
