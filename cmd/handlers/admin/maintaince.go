package admin

import (
	"errors"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h admin) AddMaintaince(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.AddMaintainceAdmin
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

	if user.Role != "engineer" {
		h.log.WithError(errors.New("wrong user role")).Errorf("getting user")
		p.SetError(errors.New("wrong user role").Error())
		h.t.Render(w, p)
		return
	}

	if err := h.maintainceService.AddMaintaince(r.Context(), req.AddMaintaince, user.UID); err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetSuccess("maintaince added").
		SetCode(200)

	h.t.Render(w, p)
}

func (h admin) AddMaintaincePage(w http.ResponseWriter, r *http.Request) {
	var (
		p = render.NewPage()
	)

	p.SetTemplate("admin.gohtml").
		SetPath(r.URL.Path).
		SetCode(200)
	h.t.Render(w, p)

}

func (h admin) MaintainceEquipmentFormAdd(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	p.SetTemplate("components/maintaince/add_equipment.gohtml").
		SetCode(200)

	h.t.Render(w, p)
}

func (h admin) GetMaintaince(w http.ResponseWriter, r *http.Request) {
	var (
		p = render.NewPage()
	)

	res, err := h.maintainceService.GetAll(r.Context())
	if err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetTemplate("admin.gohtml").
		SetPath(r.URL.Path).
		SetData(res).
		SetCode(200)
	h.t.Render(w, p)
}

func (h admin) DeleteMaintainceByUID(w http.ResponseWriter, r *http.Request) {
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

	if err := h.maintainceService.DeleteByUID(r.Context(), req.UID); err != nil {
		h.log.WithError(err).Errorf("deleting")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	w.WriteHeader(200)
}
