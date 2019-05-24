package gogen

import "context"

type IfBuilder struct {
	cs  []Code
	set bool
}

func If(v bool, cs ...Code) (r *IfBuilder) {
	r = &IfBuilder{}
	if v {
		r.cs = cs
		r.set = true
	}
	return
}

func (b *IfBuilder) ElseIf(v bool, cs ...Code) (r *IfBuilder) {
	if b.set {
		return b
	}
	if v {
		b.cs = cs
		b.set = true
	}
	return b
}

func (b *IfBuilder) Else(cs ...Code) (r *IfBuilder) {
	if b.set {
		return b
	}
	b.set = true
	b.cs = cs
	return b
}

func (b *IfBuilder) MarshalCode(ctx context.Context) (r []byte, err error) {
	if len(b.cs) == 0 {
		return
	}
	return Codes(b.cs...).MarshalCode(ctx)
}
