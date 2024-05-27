package laboratorian

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h laboratorian) GetExperiments(w http.ResponseWriter, r *http.Request) {
	var (
		p = render.NewPage()
	)
	res, err := h.experimentService.GetAllFinishedNotMarked(r.Context())
	if err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetTemplate("laboratorian.gohtml").
		SetPath(r.URL.Path).
		SetData(res).
		SetCode(200)
	h.t.Render(w, p)
}

func (h laboratorian) GetExperimentFinishByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetExperimentByUID
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	exp, err := h.experimentService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting exp")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/laboratorian/experiment/finished_experiment_inventory.gohtml").
		SetData(exp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h laboratorian) GetExperimentByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetExperimentByUID
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	exp, err := h.experimentService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting exp")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/laboratorian/experiment/row.gohtml").
		SetData(exp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h laboratorian) FinishExperimentByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.FinishExperiment
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	if err := h.experimentService.Finish(r.Context(), req); err != nil {
		h.log.WithError(err).Errorf("finishing")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	w.WriteHeader(200)
}
