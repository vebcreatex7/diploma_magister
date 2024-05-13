package handlers

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h admin) GetEquipment(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage()
	eq, err := h.equipmentService.GetAllNotCanceled(r.Context())
	if err != nil {
		h.log.WithError(err).Errorf("getting all")

		p.SetTemplate("admin.gohtml").
			SetPath(r.URL.Path).
			SetError(err.Error()).
			SetCode(200)

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("admin.gohtml").
		SetPath(r.URL.Path).
		SetData(eq).
		SetCode(200)
	h.t.Render(w, p)
}

func (h admin) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.DeleteEquipment
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	_, err := h.equipmentService.DeleteByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("deleting")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h admin) GetEquipmentEditByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetEquipment
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	resp, err := h.equipmentService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	p.SetTemplate("components/equipment/row_edit.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) GetEquipmentByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetEquipment
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	resp, err := h.equipmentService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	p.SetTemplate("components/equipment/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) EditEquipment(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.EditEquipment
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	resp, err := h.equipmentService.Edit(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("editting")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	p.SetTemplate("components/equipment/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) AddEquipmentPage(w http.ResponseWriter, r *http.Request) {
	var (
		p = render.NewPage()
	)

	resp, err := h.equipmentService.GetAllNotCanceled(r.Context())
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	p.SetTemplate("components/equipment/add.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) AddEquipment(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.CreateEquipment
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	resp, err := h.equipmentService.Create(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	p.SetTemplate("components/equipment/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}
