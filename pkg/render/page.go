package render

import (
	"github.com/vebcreatex7/diploma_magister/pkg/ptr"
)

type page struct {
	Tmpl    string
	Code    int
	Path    string
	Type    *string
	Message *string
	Data    any
}

func NewPage() *page {
	return &page{}
}

func (r *page) SetTemplate(t string) *page {
	r.Tmpl = t

	return r
}

func (r *page) SetCode(c int) *page {
	r.Code = c

	return r
}

func (r *page) SetPath(p string) *page {
	r.Path = p

	return r
}

func (r *page) SetType(t string) *page {
	r.Type = ptr.Ref(t)

	return r
}

func (r *page) SetMessage(m string) *page {
	r.Message = ptr.Ref(m)

	return r
}

func (r *page) SetData(d any) *page {
	r.Data = d

	return r
}

func (r *page) SetError(m string) *page {
	r.Type = ptr.Ref("error")
	r.Message = ptr.Ref(m)

	return r
}

func (r *page) SetWarning(m string) *page {
	r.Type = ptr.Ref("warning")
	r.Message = ptr.Ref(m)

	return r
}

func (r *page) SetSuccess(m string) *page {
	r.Type = ptr.Ref("success")
	r.Message = ptr.Ref(m)

	return r
}
