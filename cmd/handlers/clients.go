package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/pkg/render"
	"html/template"
	"net/http"
)

type clients struct {
	t       *template.Template
	service clientsService
}

func (h Message) GetMessages() string {
	return h.Type
}

func NewClients(t *template.Template, service clientsService) clients {
	return clients{t: t, service: service}
}

func (h clients) BasePrefix() string {
	return "/clients"
}

func (h clients) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/register", h.RegisterPage)
	r.Post("/register", h.Register)

	r.Get("/login", h.LoginPage)
	r.Post("/login", h.Login)

	return r
}

func (h clients) RegisterPage(w http.ResponseWriter, r *http.Request) {
	render.HTMLResponse(w, h.t, "pages/register.gohtml", nil, nil, http.StatusOK)
}

func (h clients) Register(w http.ResponseWriter, r *http.Request) {
	var req request.CreateClient

	if err := req.Bind(r); err != nil {
		render.HTMLResponse(w, h.t, "pages/register.gohtml", Message{
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), req); err != nil {
		render.HTMLResponse(w, h.t, "pages/register.gohtml", Message{
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	render.HTMLResponse(w, h.t, "pages/register.gohtml", Message{
		Type: "success",
		Text: "Регистрация пройдена, ожидайте подтверждения аккаунта",
	}, nil, http.StatusOK)
	return
}

func (h clients) LoginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	render.HTMLResponse(w, h.t, "pages/login.gohtml", nil, nil, http.StatusOK)
}

func (h clients) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginClient

	if err := req.Bind(r); err != nil {
		render.HTMLResponse(w, h.t, "pages/login.gohtml", Message{
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req)
	if err != nil {
		render.HTMLResponse(w, h.t, "pages/login.gohtml", Message{
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	setCookie(w, "jwt", token)

	render.HTMLResponse(w, h.t, "pages/login.gohtml", Message{
		Type: "success",
		Text: token,
	}, nil, http.StatusOK)

}

func setCookie(w http.ResponseWriter, name, value string) {
	c := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &c)
}
