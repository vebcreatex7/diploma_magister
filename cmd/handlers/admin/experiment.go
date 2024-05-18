package admin

import (
	"errors"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h admin) AddExperimentPage(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	p.SetTemplate("admin.gohtml").
		SetPath(r.URL.Path).
		SetCode(200)

	h.t.Render(w, p)
}

func (h admin) AddExperiment(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.AddExperimentAdmin
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	user, err := h.clientsService.GetByLogin(r.Context(), req.User)
	if err != nil {
		h.log.WithError(err).Errorf("getting user")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	if user.Role != "scientist" {
		h.log.WithError(errors.New("wrong user role")).Errorf("getting user")
		p.SetError(errors.New("wrong user role").Error())
		h.t.Render(w, p)
		return
	}

	if err := h.experimentService.AddExperiment(r.Context(), req.AddExperiment, user.UID); err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetSuccess("experiment added").
		SetCode(200)

	h.t.Render(w, p)
}

func (h admin) ExperimentEquipmentFormAdd(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	p.SetTemplate("components/experiment/add_equipment.gohtml").
		SetCode(200)

	h.t.Render(w, p)
}

func (h admin) ExperimentInventoryFormAdd(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	p.SetTemplate("components/experiment/add_inventory.gohtml").
		SetCode(200)

	h.t.Render(w, p)
}

func (h admin) GetExperiment(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage()
	eq, err := h.experimentService.GetAll(r.Context())
	if err != nil {
		h.log.WithError(err).Errorf("getting all")

		p.SetTemplate("admin.gohtml").
			SetPath(r.URL.Path).
			SetError(err.Error()).
			SetCode(422)

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("admin.gohtml").
		SetPath(r.URL.Path).
		SetData(eq).
		SetCode(200)
	h.t.Render(w, p)
}

func (h admin) DeleteExperimentByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.DeleteExperiment
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	if err := h.experimentService.DeleteByUID(r.Context(), req.UID); err != nil {
		h.log.WithError(err).Errorf("deleting")
		p.SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	w.WriteHeader(200)
}
