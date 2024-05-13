package render

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type Template struct {
	t   *template.Template
	log *logrus.Logger
}

func NewTemplate(
	t *template.Template,
	log *logrus.Logger,
) *Template {
	return &Template{
		t:   t,
		log: log,
	}
}

func (t *Template) Render(w http.ResponseWriter, p *page) {
	if p.toast != nil {
		w.Header().Set("HX-Trigger", p.toast.toJSON())

		if p.toast.Level != "success" {
			w.Header().Set("HX-Reswap", "none")
			if p.code == 200 {
				t.log.Warnf("Unsuccess with code 200")
			}

			w.WriteHeader(p.code)
			return
		}
	}

	for k, v := range p.headers {
		w.Header().Set(k, v)
	}

	if p.code == 200 {
		var buf bytes.Buffer

		if err := t.t.ExecuteTemplate(&buf, p.tmpl, p); err != nil {
			t.log.WithError(err).Errorf("executing template")
			w.WriteHeader(500)

			return
		}

		w.Header().Set("content-type", "text/html")
		w.Write(buf.Bytes())
	}

	w.WriteHeader(p.code)
}

func (t *Template) RenderData(w http.ResponseWriter, p *page) {
	if p.toast != nil {
		w.Header().Set("HX-Trigger", p.toast.toJSON())

		if p.toast.Level != "success" {
			w.Header().Set("HX-Reswap", "none")
			if p.code == 200 {
				t.log.Warnf("Unsuccess with code 200")
			}

			w.WriteHeader(p.code)
			return
		}
	}

	for k, v := range p.headers {
		w.Header().Set(k, v)
	}

	if p.code == 200 {
		var buf bytes.Buffer

		if err := t.t.ExecuteTemplate(&buf, p.tmpl, p.Data); err != nil {
			t.log.WithError(err).Errorf("executing template")
			w.WriteHeader(500)

			return
		}

		w.Header().Set("content-type", "text/html")
		w.Write(buf.Bytes())
	}

	w.WriteHeader(p.code)
}
