/*
# gogen: Generate Go code with structured blocks and composition

	import . "github.com/sunfmin/gogen"

*/
package gogen

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

type Code interface {
	MarshalCode(ctx context.Context) (r []byte, err error)
}

type CodeFunc func(ctx context.Context) (r []byte, err error)

type Raw string

func (rc Raw) MarshalCode(ctx context.Context) (r []byte, err error) {
	r = []byte(rc)
	return
}

func (cf CodeFunc) MarshalCode(ctx context.Context) (r []byte, err error) {
	return cf(ctx)
}

type SnippetsBuilder struct {
	cs         []Code
	sep        string
	appendLast bool
}

func Snippets(cs ...Code) (r *SnippetsBuilder) {
	r = &SnippetsBuilder{cs: cs}
	r.sep = "\n"
	r.appendLast = true
	return
}

func (b *SnippetsBuilder) Separator(sep string, appendLast bool) (r *SnippetsBuilder) {
	b.sep = sep
	b.appendLast = appendLast
	return b
}

func (b *SnippetsBuilder) Clone() (r *SnippetsBuilder) {
	r = &SnippetsBuilder{
		cs:         b.cs,
		sep:        b.sep,
		appendLast: b.appendLast,
	}
	return
}

func (b *SnippetsBuilder) Append(Snippets ...Code) (r *SnippetsBuilder) {
	b.cs = append(b.cs, Snippets...)
	return b
}

func (b *SnippetsBuilder) AppendSnippet(template string, vars ...string) (r *SnippetsBuilder) {
	b.Append(Snippet(template, vars...))
	return b
}

func (b *SnippetsBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {
	buf := bytes.NewBuffer(nil)
	l := len(b.cs)
	for i, c := range b.cs {
		if c == nil {
			continue
		}
		err = Fprint(buf, c, ctx)
		if err != nil {
			return
		}
		if b.appendLast || i+1 < l {
			buf.WriteString(b.sep)
		}
	}
	r = buf.Bytes()
	return

}

func Fprint(w io.Writer, c Code, ctx context.Context) (err error) {
	if c == nil {
		return
	}
	var b []byte
	b, err = c.MarshalCode(ctx)
	if err != nil {
		return
	}
	_, err = fmt.Fprint(w, string(b))
	return
}

func MustString(c Code, ctx context.Context) (r string) {
	buf := bytes.NewBuffer(nil)
	err := Fprint(buf, c, ctx)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func Quote(q string) (r string) {
	return fmt.Sprintf("%q", q)
}
