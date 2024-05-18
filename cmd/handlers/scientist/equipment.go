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

	resp, err := h.equipmentService.GetEquipmentScheduleInRangeForUser(r.Context(), req, uid)
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

func (h scientist) GetAllMyEquipment(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	uid, err := pkg.GetUIDFromJWT(r)
	if err != nil {
		h.log.WithError(err).Errorf("getting uid from jwt")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	res, err := h.equipmentService.GetAllForUser(r.Context(), uid)
	if err != nil {
		h.log.WithError(err).Errorf("getting equipment")
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
