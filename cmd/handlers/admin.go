package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vebcreatex7/diploma_magister/internal/pkg/render"
	"html/template"
	"net/http"
)

type admin struct {
	t       *template.Template
	clients clientsService
}

func NewAdmin(t *template.Template, clients clientsService) admin {
	return admin{t: t, clients: clients}
}

func (h admin) Routes() chi.Router {

	r := chi.NewRouter()

	//r.Use(pkg.ValidateAdminJWTCookies)

	r.Get("/", h.Home)
	r.Get("/home", h.Home)
	r.Get("/users", h.Users)

	return r
}

func (h admin) BasePrefix() string {
	return "/admin"
}

func (h admin) Home(w http.ResponseWriter, r *http.Request) {
	render.HTMLResponse(w, h.t, "admin.gohtml", r.URL, nil, http.StatusOK)
}

func (h admin) Users(w http.ResponseWriter, r *http.Request) {
	users, err := h.clients.GetAllNotCanceled(r.Context())
	if err != nil {
		render.HTMLResponse(w, h.t, "admin.gohtml", Message{
			URL:  r.URL,
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	render.HTMLResponse(w, h.t, "admin.gohtml", Message{
		URL:  r.URL,
		Type: "success",
		Text: "yes",
		Data: users,
	}, nil, http.StatusOK)
	return
}
