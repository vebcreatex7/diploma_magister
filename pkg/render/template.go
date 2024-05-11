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
	if err := t.t.ExecuteTemplate(w, p.Tmpl, *p); err != nil {
		t.log.WithError(err).Errorf("executing template")
		w.WriteHeader(500)

		return
	}

	w.Header().Set("content-type", "text/html")
	w.WriteHeader(p.Code)
}

func (t *Template) RenderData(w http.ResponseWriter, p *page) {
	var buf bytes.Buffer

	if err := t.t.ExecuteTemplate(&buf, p.Tmpl, p.Data); err != nil {
		t.log.WithError(err).Errorf("executing template")
		w.WriteHeader(500)

		return
	}

	w.Header().Set("content-type", "text/html")
	w.WriteHeader(p.Code)
	w.Write(buf.Bytes())
}
