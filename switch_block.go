package gogen

import (
	"bytes"
	"context"
)

type SwitchBlockBuilder struct {
	switchBlock  Code
	cases        []Code
	defaultBlock Code
}

func SwitchBlock(switchTemplate string, vars ...string) (r *SwitchBlockBuilder) {
	r = &SwitchBlockBuilder{}
	r.switchBlock = Block(switchTemplate, vars...)
	return
}

func (b *SwitchBlockBuilder) Cases(caseBlocks ...Code) (r *SwitchBlockBuilder) {
	b.cases = append(b.cases, caseBlocks...)
	return b
}

func (b *SwitchBlockBuilder) Default(defaultTemplate string, vars ...string) (r *SwitchBlockBuilder) {
	b.defaultBlock = Block(defaultTemplate, vars...)
	return b
}

func (b *SwitchBlockBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {

	buf := bytes.NewBuffer(nil)
	err = Fprint(buf, b.switchBlock, ctx)
	if err != nil {
		return
	}
	buf.WriteString(" {\n")
	err = Fprint(buf, Codes(b.cases...).Separator("\n"), ctx)
	if err != nil {
		panic(err)
	}
	buf.WriteString("\n")

	err = Fprint(buf, b.defaultBlock, ctx)
	if err != nil {
		panic(err)
	}

	buf.WriteString("\n}\n")
	r = buf.Bytes()
	return

}
