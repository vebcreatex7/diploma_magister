package admin

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h admin) GetEquipment(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage()
	eq, err := h.equipmentService.GetAll(r.Context())
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

func (h admin) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.DeleteEquipment
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	if err := h.equipmentService.DeleteByUID(r.Context(), req.UID); err != nil {
		h.log.WithError(err).Errorf("deleting")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	w.WriteHeader(200)
}

func (h admin) GetEquipmentEditByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetEquipment
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.equipmentService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		p.SetError(err.Error())

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
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.equipmentService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		p.SetError(err.Error())

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
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.equipmentService.Edit(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("editing")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/equipment/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) AddEquipmentPage(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	p.SetTemplate("components/equipment/row_add.gohtml").SetCode(200)

	h.t.Render(w, p)
}

func (h admin) GetEquipmentSchedule(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetEquipmentSchedule
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	resp, err := h.equipmentService.GetEquipmentScheduleInRange(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("getting equipment_schedule")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/equipment/schedule_response.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) GetEquipmentScheduleEmpty(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<div id="equipment-schedule-response"><div>`))
	return
}

func (h admin) AddEquipment(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.CreateEquipment
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.equipmentService.Create(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/equipment/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}
