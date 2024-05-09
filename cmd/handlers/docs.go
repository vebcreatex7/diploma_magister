package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vebcreatex7/diploma_magister/internal/pkg/render"
	"html/template"
	"net/http"
)

type docs struct {
	t *template.Template
}

func (h docs) BasePrefix() string {
	return "/docs"
}

func (h docs) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.Docs)

	return r
}

func NewDocs(t *template.Template) docs {
	return docs{t: t}
}

func (h docs) Docs(w http.ResponseWriter, r *http.Request) {
	render.HTMLResponse(w, h.t, "docs.gohtml", nil, nil, http.StatusOK)
}
