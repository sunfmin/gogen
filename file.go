package gogen

import (
	"bytes"
	"context"
)

type FileBuilder struct {
	blocks []Code
	pkg    string
}

func File(name string) (r *FileBuilder) {
	r = &FileBuilder{}
	return
}

func (b *FileBuilder) Blocks(cs ...Code) (r *FileBuilder) {
	b.blocks = append(b.blocks, cs...)
	return b
}

func (b *FileBuilder) Package(pkg string) (r *FileBuilder) {
	b.pkg = pkg
	return b
}

func (b *FileBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("package " + b.pkg)
	buf.WriteString("\n\n")
	err = Fprint(buf, Codes(b.blocks...), ctx)
	if err != nil {
		return
	}
	r = buf.Bytes()
	return
}
