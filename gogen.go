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

type RawCode string

func (rc RawCode) MarshalCode(ctx context.Context) (r []byte, err error) {
	r = []byte(rc)
	return
}

func (cf CodeFunc) MarshalCode(ctx context.Context) (r []byte, err error) {
	return cf(ctx)
}

type CodesBuilder struct {
	cs         []Code
	sep        string
	appendLast bool
}

func Codes(cs ...Code) (r *CodesBuilder) {
	r = &CodesBuilder{cs: cs}
	r.sep = "\n"
	r.appendLast = true
	return
}

func (b *CodesBuilder) Separator(sep string, appendLast bool) (r *CodesBuilder) {
	b.sep = sep
	b.appendLast = appendLast
	return b
}

func (b *CodesBuilder) Clone() (r *CodesBuilder) {
	r = &CodesBuilder{
		cs:         b.cs,
		sep:        b.sep,
		appendLast: b.appendLast,
	}
	return
}

func (b *CodesBuilder) Append(codes ...Code) (r *CodesBuilder) {
	b.cs = append(b.cs, codes...)
	return b
}

func (b *CodesBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {
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
