package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/domain/service"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"html/template"
	"net/http"
	"net/url"
)

type Message struct {
	*url.URL
	Type string
	Text string
	Data any
}
type home struct {
	t              *template.Template
	clientsService service.Clients
}

func NewHome(t *template.Template, clientsService service.Clients) home {
	return home{t: t, clientsService: clientsService}
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

	token, err := h.clientsService.Login(r.Context(), req)
	if err != nil {
		render.HTMLResponse(w, h.t, "home.gohtml", Message{
			URL:  r.URL,
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	render.SetCookie(w, "jwt", token)

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

	if err := h.clientsService.Create(r.Context(), req); err != nil {
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
