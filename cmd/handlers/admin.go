package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/domain/service"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"html/template"
	"net/http"
)

type admin struct {
	t                *template.Template
	clientsService   service.Clients
	equipmentService service.Equipment
}

func NewAdmin(
	t *template.Template,
	clientsService service.Clients,
	equipmentService service.Equipment,
) admin {
	return admin{
		t:                t,
		clientsService:   clientsService,
		equipmentService: equipmentService,
	}
}

func (h admin) Routes() chi.Router {

	r := chi.NewRouter()

	//r.Use(pkg.ValidateAdminJWTCookies)

	r.Get("/", h.Home)
	r.Get("/home", h.Home)
	r.Get("/users", h.Users)
	r.Get("/equipment", h.Equipment)
	r.Delete("/users/{uid}", h.DeleteUser)
	r.Delete("/equipment/{uid}", h.DeleteEquipment)

	return r
}

func (h admin) BasePrefix() string {
	return "/admin"
}

func (h admin) Home(w http.ResponseWriter, r *http.Request) {
	render.HTMLResponse(w, h.t, "admin.gohtml", r.URL, nil, http.StatusOK)
}

func (h admin) Users(w http.ResponseWriter, r *http.Request) {
	users, err := h.clientsService.GetAllNotCanceled(r.Context())
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
		Data: users,
	}, nil, http.StatusOK)
	return
}

func (h admin) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req request.DeleteClient

	if err := req.Bind(r); err != nil {
		render.HTMLResponse(w, h.t, "admin.gohtml", Message{
			URL:  r.URL,
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	users, err := h.clientsService.DeleteByUID(r.Context(), req.UID)
	if err != nil {
		render.HTMLResponse(w, h.t, "admin.gohtml", Message{
			URL:  r.URL,
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
		return
	}

	render.HTMLResponse(w, h.t, "pages/users.gohtml", Message{
		URL:  r.URL,
		Type: "success",
		Text: "yes",
		Data: users,
	}, nil, http.StatusOK)
	return
}

func (h admin) Equipment(w http.ResponseWriter, r *http.Request) {
	eq, err := h.equipmentService.GetAllNotCanceled(r.Context())
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
		Data: eq,
	}, nil, http.StatusOK)
	return
}

func (h admin) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	var req request.DeleteEquipment

	if err := req.Bind(r); err != nil {
		render.HTMLResponse(w, h.t, "admin.gohtml", Message{
			URL:  r.URL,
			Type: "error",
			Text: err.Error(),
		}, nil, http.StatusBadRequest)
	}

	eq, err := h.equipmentService.DeleteByUID(r.Context(), req.UID)
	if err != nil {
		render.ErrResponse(w, err)
		return
	}

	render.HTMLResponse(w, h.t, "pages/equipment.gohtml", Message{
		URL:  r.URL,
		Type: "success",
		Text: "yes",
		Data: eq,
	}, nil, http.StatusOK)
	return
}
