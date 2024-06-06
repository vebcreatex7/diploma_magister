package render

import "encoding/json"

type toast struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func (r toast) ToJSON() string {
	var tmp = make(map[string]toast)

	tmp["makeToast"] = r

	res, _ := json.Marshal(tmp)

	return string(res)
}

type page struct {
	tmpl    string
	code    int
	Path    string
	tp      string
	message string
	Data    any
	headers map[string]string
	toast   *toast
}

func NewPage() *page {
	return &page{headers: make(map[string]string)}
}

func (r *page) SetTemplate(t string) *page {
	r.tmpl = t

	return r
}

func (r *page) SetCode(c int) *page {
	r.code = c

	return r
}

func (r *page) SetPath(p string) *page {
	r.Path = p

	return r
}

func (r *page) SetType(t string) *page {
	r.tp = t

	return r
}

func (r *page) SetMessage(m string) *page {
	r.message = m

	return r
}

func (r *page) SetData(d any) *page {
	r.Data = d

	return r
}

func (r *page) SetError(m string) *page {
	r.tp = "error"
	r.message = m

	r.toast = &toast{
		Level:   "error",
		Message: m,
	}

	r.code = 422

	return r
}

func (r *page) SetWarning(m string) *page {
	r.tp = "warning"
	r.message = m

	r.toast = &toast{
		Level:   "warning",
		Message: m,
	}

	return r
}

func (r *page) SetSuccess(m string) *page {
	r.tp = "success"
	r.message = m

	r.toast = &toast{
		Level:   "success",
		Message: m,
	}

	r.code = 200

	return r
}

func (r *page) SetHeader(k, v string) *page {
	r.headers[k] = v

	return r
}

func (r *page) Toast() *toast {
	return r.toast
}
