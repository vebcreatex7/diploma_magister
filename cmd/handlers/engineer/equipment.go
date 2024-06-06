package engineer

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h engineer) GetAllEquipment(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage()
	eq, err := h.equipmentService.GetAll(r.Context())
	if err != nil {
		h.log.WithError(err).Errorf("getting all")

		p.SetTemplate("engineer.gohtml").
			SetPath(r.URL.Path).
			SetError(err.Error()).
			SetCode(422)

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("engineer.gohtml").
		SetPath(r.URL.Path).
		SetData(eq).
		SetCode(200)
	h.t.Render(w, p)
}

func (h engineer) GetEquipmentSchedule(w http.ResponseWriter, r *http.Request) {
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

	p.SetTemplate("components/engineer/equipment/schedule_response.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h engineer) GetEquipmentScheduleEmpty(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<div id="equipment-schedule-response"><div>`))
	return
}
