package scientist

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h scientist) GetEquipmentSchedule(w http.ResponseWriter, r *http.Request) {
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

	uid, err := pkg.GetUIDFromJWT(r)
	if err != nil {
		h.log.WithError(err).Errorf("getting uid from jwt")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	resp, err := h.equipmentService.GetEquipmentScheduleInRange(r.Context(), req, uid)
	if err != nil {
		h.log.WithError(err).Errorf("getting equipment_schedule")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/scientist/equipment/schedule_response.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h scientist) GetEquipmentScheduleEmpty(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<div id="equipment-schedule-response"><div>`))
	return
}
