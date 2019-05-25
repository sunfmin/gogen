package gogen

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"strconv"
	"strings"
)

type FileBuilder struct {
	name   string
	blocks []Code
	pkg    string
}

func File(name string) (r *FileBuilder) {
	r = &FileBuilder{}
	r.name = name
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

func (b *FileBuilder) Fprint(w io.Writer, ctx context.Context) (err error) {
	src := MustString(b, ctx)
	fset := token.NewFileSet()
	var f *ast.File
	f, err = parser.ParseFile(fset, b.name, src, 0)
	if err != nil {
		return
	}
	err = printer.Fprint(w, fset, f)
	return
}

func (b *FileBuilder) MustString(ctx context.Context) (r string) {
	buf := bytes.NewBuffer(nil)
	err := b.Fprint(buf, ctx)
	if err != nil {
		hl, _ := strconv.ParseInt(strings.Split(err.Error(), ":")[0], 10, 64)

		panic(fmt.Sprintf("%s\n%s", err, codeWithLineNumber(b, hl, ctx)))
	}
	return buf.String()
}

func codeWithLineNumber(c Code, highlightLine int64, ctx context.Context) (r string) {
	src := MustString(c, ctx)
	lines := strings.Split(src, "\n")
	linesWithNumber := []string{}
	for i, l := range lines {
		hl := "   "
		if int64(i+1) == highlightLine {
			hl = ">> "
		}
		linesWithNumber = append(linesWithNumber, fmt.Sprintf("%s%d: %s", hl, i+1, l))
	}
	r = strings.Join(linesWithNumber, "\n")
	return
}
