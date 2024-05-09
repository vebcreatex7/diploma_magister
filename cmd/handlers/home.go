package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/pkg/render"
	"html/template"
	"net/http"
	"net/url"
)

type clientsService interface {
	Create(ctx context.Context, req request.CreateClient) error
	Login(ctx context.Context, req request.LoginClient) (string, error)
	GetAllNotCanceled(ctx context.Context) ([]response.Client, error)
}
type Message struct {
	*url.URL
	Type string
	Text string
	Data any
}
type home struct {
	t       *template.Template
	service clientsService
}

func NewHome(t *template.Template, service clientsService) home {
	return home{t: t, service: service}
}

func (h home) BasePrefix() string {
	return "/"
}

func (h home) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.Home)
	r.Get("/home", h.Home)
	r.Get("/docs", h.Home)
	r.Get("/register", h.Home)
	r.Get("/login", h.Home)
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	return r
}

func (h home) Home(w http.ResponseWriter, r *http.Request) {
	render.HTMLResponse(w, h.t, "home.gohtml", Message{URL: r.URL}, nil, http.StatusOK)
}

func (h home) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginClient

	if err := req.Bind(r); err != nil {
		render.HTMLResponse(w, h.t, "home.gohtml", Message{
			URL:  r.URL,
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req)
	if err != nil {
		render.HTMLResponse(w, h.t, "home.gohtml", Message{
			URL:  r.URL,
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	setCookie(w, "jwt", token)

	render.HTMLResponse(w, h.t, "home.gohtml", Message{
		URL:  r.URL,
		Type: "success",
		Text: token,
	}, nil, http.StatusOK)
}

func (h home) Register(w http.ResponseWriter, r *http.Request) {
	var req request.CreateClient

	if err := req.Bind(r); err != nil {
		render.HTMLResponse(w, h.t, "home.gohtml", Message{
			URL:  r.URL,
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), req); err != nil {
		render.HTMLResponse(w, h.t, "home.gohtml", Message{
			URL:  r.URL,
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	render.HTMLResponse(w, h.t, "home.gohtml", Message{
		URL:  r.URL,
		Type: "success",
		Text: "Регистрация пройдена, ожидайте подтверждения аккаунта",
	}, nil, http.StatusOK)
	return
}
