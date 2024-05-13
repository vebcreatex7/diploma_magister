package handlers

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
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

func (h admin) AddAccessGroupPage(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	p.SetTemplate("components/access_groups/row_add.gohtml").SetCode(200)

	h.t.Render(w, p)
}

func (h admin) AddAccessGroup(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.CreateAccessGroup
	)

	if err := req.Bind(r); err != nil {
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.accessGroupsService.Create(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	p.SetTemplate("components/access_groups/row.gohtml").
		SetData(resp).
		SetSuccess("access_group added").
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) GetAccessGroupEditByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetAccessGroup
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.accessGroupsService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/access_groups/row_edit.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) EditAccessGroup(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.EditAccessGroup
	)

	//defer h.t.Render(w, p)

	if err := req.Bind(r); err != nil {
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	resp, err := h.accessGroupsService.Edit(r.Context(), req)
	if err != nil {
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/access_groups/row.gohtml").
		SetData(resp).
		SetSuccess("access_group edited").
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) GetAccessGroupByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetAccessGroup
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.accessGroupsService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/access_groups/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) DeleteAccessGroup(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.DeleteInventory
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	if err := h.accessGroupsService.DeleteByUID(r.Context(), req.UID); err != nil {
		h.log.WithError(err).Errorf("deleting")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	w.WriteHeader(http.StatusOK)
}
