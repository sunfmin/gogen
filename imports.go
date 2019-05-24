package gogen

import (
	"bytes"
	"context"
	"fmt"
)

type ImportBuilder struct {
	blocks []Code
}

func Imports(imps ...string) (r *ImportBuilder) {
	r = &ImportBuilder{}
	for _, im := range imps {
		imr := []rune(im)
		if imr[len(imr)-1] != '"' {
			im = fmt.Sprintf("%q", im)
		}
		r.blocks = append(r.blocks, RawCode(im))
	}
	return
}

func (b *ImportBuilder) Blocks(imps ...Code) (r *ImportBuilder) {
	b.blocks = append(b.blocks, imps...)
	return b
}

func (b *ImportBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {
	if len(b.blocks) == 0 {
		return
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteString("import (\n")
	err = Fprint(buf, Codes(b.blocks...).Separator("\n"), ctx)
	if err != nil {
		panic(err)
	}
	buf.WriteString("\n)\n")
	r = buf.Bytes()
	return

}

func ImportAs(as string, imp string) (r Code) {
	if len(as) == 0 {
		return RawCode(fmt.Sprintf(`"%s"`, imp))
	}
	return RawCode(fmt.Sprintf(`%s "%s"`, as, imp))
}
