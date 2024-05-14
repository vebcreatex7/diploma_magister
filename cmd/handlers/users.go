package handlers

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h admin) GetUsers(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	users, err := h.clientsService.GetAll(r.Context())
	if err != nil {
		h.log.WithError(err).Errorf("getting all")

		p.SetTemplate("admin.gohtml").
			SetPath(r.URL.Path).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("admin.gohtml").
		SetPath(r.URL.Path).
		SetData(users).
		SetCode(200)

	h.t.Render(w, p)
	return
}

func (h admin) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.DeleteInventory
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	if err := h.clientsService.DeleteByUID(r.Context(), req.UID); err != nil {
		h.log.WithError(err).Errorf("deleting")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	w.WriteHeader(200)
}

func (h admin) GetUserEditByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetUser
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.clientsService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/user/row_edit.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) GetUserByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetUser
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.clientsService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/user/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) EditUser(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.EditUser
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.clientsService.Edit(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("editting")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/user/row.gohtml").
		SetData(resp).
		SetSuccess("user edited")

	h.t.RenderData(w, p)
}

func (h admin) ApproveUser(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.ApproveUser
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.clientsService.Approve(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("approving")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("admin.gohtml").
		SetPath(r.URL.Path).
		SetData(resp).
		SetCode(200)

	h.t.Render(w, p)
}
