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
		r.blocks = append(r.blocks, Raw(im))
	}
	return
}

func (b *ImportBuilder) Body(imps ...Code) (r *ImportBuilder) {
	b.blocks = append(b.blocks, imps...)
	return b
}

func (b *ImportBuilder) BodySnippet(template string, vars ...string) (r *ImportBuilder) {
	b.Body(Snippet(template, vars...))
	return b
}

func (b *ImportBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {
	if len(b.blocks) == 0 {
		return
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteString("import (\n")
	err = Fprint(buf, Snippets(b.blocks...), ctx)
	if err != nil {
		panic(err)
	}
	buf.WriteString(")\n")
	r = buf.Bytes()
	return

}

func ImportAs(as string, imp string) (r Code) {
	if len(as) == 0 {
		return Raw(fmt.Sprintf(`"%s"`, imp))
	}
	return Raw(fmt.Sprintf(`%s "%s"`, as, imp))
}
