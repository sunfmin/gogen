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

func (b *ConstBlockBuilder) Consts(cs ...Code) (r *ConstBlockBuilder) {
	b.consts = append(b.consts, cs...)
	return b
}

func (b *ConstBlockBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	for _, c := range b.consts {
		if ci, ok := c.(*ConstItemBuilder); ok {
			ci.SetConstType(b.constType, b.namePrefixType)
		}
	}

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

type ConstItemBuilder struct {
	name           string
	val            interface{}
	constType      string
	namePrefixType bool
}

func Const(name string, val interface{}) (r *ConstItemBuilder) {
	r = &ConstItemBuilder{
		name: name,
		val:  val,
	}
	return
}

func (b *ConstItemBuilder) SetConstType(constType string, namePrefixType bool) {
	b.constType = constType
	b.namePrefixType = namePrefixType
}

func (b *ConstItemBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	if len(b.constType) == 0 {
		panic("To use Const, Parent must set Type")
	}
	var name = b.name
	if b.namePrefixType {
		name = fmt.Sprintf("%s%s", b.constType, name)
	}

	return RawCode(fmt.Sprintf("%s %s = %#+v", name, b.constType, b.val)).MarshalCode(ctx)
}
