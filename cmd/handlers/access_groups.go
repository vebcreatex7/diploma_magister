package handlers

import (
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h admin) GetAccessGroups(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage()

	grs, err := h.accessGroupsService.GetAllNotCanceled(r.Context())
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
		SetData(grs).
		SetCode(200)
	h.t.Render(w, p)
}
