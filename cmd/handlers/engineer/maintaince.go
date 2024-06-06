package engineer

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h engineer) AddMaintaincePage(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage()

	res, err := h.maintainceService.GetSuggestions(r.Context())
	if err != nil {
		h.log.WithError(err).Errorf("getting suggestions")

		p.SetTemplate("engineer.gohtml").
			SetPath(r.URL.Path).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("engineer.gohtml").
		SetPath(r.URL.Path).
		SetData(res).
		SetCode(200)

	h.t.Render(w, p)

}

func (h engineer) AddMaintaince(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.AddMaintaince
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

	if err := h.maintainceService.AddMaintaince(r.Context(), req, uid); err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetSuccess("maintaince added").
		SetCode(200)

	h.t.Render(w, p)
}

func (h engineer) MaintainceEquipmentFormAdd(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	p.SetTemplate("components/engineer/maintaince/add_equipment.gohtml").
		SetCode(200)

	h.t.Render(w, p)
}

func (h engineer) GetMyMaintaince(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.maintainceService.GetAllForUser(r.Context(), uid)
	if err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetTemplate("engineer.gohtml").
		SetPath(r.URL.Path).
		SetData(res).
		SetCode(200)
	h.t.Render(w, p)
}

func (h engineer) DeleteMaintainceByUID(w http.ResponseWriter, r *http.Request) {
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

	if err := h.maintainceService.DeleteByUIDForUser(r.Context(), req.UID, uid); err != nil {
		h.log.WithError(err).Errorf("deleting")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	w.WriteHeader(200)
}
