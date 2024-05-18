package scientist

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h scientist) AddExperimentPage(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage()
	uid, err := pkg.GetUIDFromJWT(r)
	if err != nil {
		h.log.WithError(err).Errorf("getting uid from jwt")

		p.SetPath(r.URL.Path).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	res, err := h.experimentService.GetSuggestionsForUser(r.Context(), uid)

	if err != nil {
		h.log.WithError(err).Errorf("getting suggestions")

		p.SetTemplate("scientist.gohtml").
			SetPath(r.URL.Path).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("scientist.gohtml").
		SetPath(r.URL.Path).
		SetData(res).
		SetCode(200)

	h.t.Render(w, p)
}

func (h scientist) AddExperiment(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.AddExperiment
	)

	uid, err := pkg.GetUIDFromJWT(r)
	if err != nil {
		h.log.WithError(err).Errorf("getting uid from jwt")

		p.SetPath(r.URL.Path).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	if err := h.experimentService.AddExperiment(r.Context(), req, uid); err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetSuccess("experiment added").
		SetCode(200)

	h.t.Render(w, p)
}

func (h scientist) GetMyExperiments(w http.ResponseWriter, r *http.Request) {
	var (
		p = render.NewPage()
	)

	uid, err := pkg.GetUIDFromJWT(r)
	if err != nil {
		h.log.WithError(err).Errorf("getting uid from jwt")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	res, err := h.experimentService.GetAllForUser(r.Context(), uid)
	if err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetTemplate("scientist.gohtml").
		SetPath(r.URL.Path).
		SetData(res).
		SetCode(200)
	h.t.Render(w, p)
}

func (h scientist) ExperimentEquipmentFormAdd(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	p.SetTemplate("components/scientist/experiment/add_equipment.gohtml").
		SetCode(200)

	h.t.Render(w, p)
}

func (h scientist) ExperimentInventoryFormAdd(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	p.SetTemplate("components/scientist/experiment/add_inventory.gohtml").
		SetCode(200)

	h.t.Render(w, p)
}

func (h scientist) DeleteExperimentByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.DeleteExperiment
	)

	uid, err := pkg.GetUIDFromJWT(r)
	if err != nil {
		h.log.WithError(err).Errorf("getting uid from jwt")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	if err := h.experimentService.DeleteByUIDForUser(r.Context(), req.UID, uid); err != nil {
		h.log.WithError(err).Errorf("deleting")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	w.WriteHeader(200)
}
