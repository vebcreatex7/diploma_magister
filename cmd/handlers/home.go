package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/domain/service"
	"github.com/vebcreatex7/diploma_magister/pkg"
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
	log            *logrus.Logger
	clientsService service.Clients
	zxc            *render.Template
}

func NewHome(
	t *template.Template,
	log *logrus.Logger,
	zxc *render.Template,
	clientsService service.Clients,
) home {
	return home{
		t:              t,
		zxc:            zxc,
		log:            log,
		clientsService: clientsService,
	}
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
	p := render.NewPage().
		SetTemplate("home.gohtml").
		SetPath(r.URL.Path).
		SetCode(200)

	h.zxc.Render(w, p)
}

func (h home) Login(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.LoginUser
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())
		h.zxc.Render(w, p)

		return
	}

	u, err := h.clientsService.Login(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("logining")
		p.SetError(err.Error())
		h.zxc.Render(w, p)

		return
	}

	token, err := pkg.GenerateJWT(u)
	if err != nil {
		h.log.WithError(err).Errorf("generating token")
		p.SetError(err.Error())
		h.zxc.Render(w, p)

		return
	}

	render.SetCookie(w, "jwt", token)

	switch u.Role {
	case "admin":
		p.SetHeader("HX-Redirect", "http://localhost:3000/admin/home")
	case "scientist":
		p.SetHeader("HX-Redirect", "http://localhost:3000/scientists/home")
	case "engineer":
		p.SetHeader("HX-Redirect", "http://localhost:3000/engineers/home")
	case "laboratorian":
		p.SetHeader("HX-Redirect", "http://localhost:3000/laboratorians/home")
	}

	p.SetSuccess("user sign in")

	h.zxc.Render(w, p)
}

func (h home) Register(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.CreateUser
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())
		h.zxc.Render(w, p)

		return
	}

	if err := h.clientsService.Create(r.Context(), req); err != nil {
		h.log.WithError(err).Errorf("creating user")
		p.SetError(err.Error())
		h.zxc.Render(w, p)

		return
	}

	p.SetSuccess("user registered")

	h.zxc.RenderEmpty(w, p)
}
