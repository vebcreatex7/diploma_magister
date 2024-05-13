package handlers

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h admin) GetUsers(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	users, err := h.clientsService.GetAllNotCanceled(r.Context())
	if err != nil {
		h.log.WithError(err).Errorf("getting all")

		p.SetTemplate("admin.gohtml").
			SetPath(r.URL.Path).
			SetError(err.Error()).
			SetCode(400)

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
	var req request.DeleteUser

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		w.WriteHeader(400)

		return
	}

	if err := h.clientsService.DeleteByUID(r.Context(), req.UID); err != nil {
		h.log.WithError(err).Errorf("deleting")
		w.WriteHeader(400)

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
		w.WriteHeader(400)

		return
	}

	resp, err := h.clientsService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		w.WriteHeader(400)

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
		w.WriteHeader(400)

		return
	}

	resp, err := h.clientsService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		w.WriteHeader(400)

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
		w.WriteHeader(400)

		return
	}

	resp, err := h.clientsService.Edit(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("editting")
		w.WriteHeader(400)

		return
	}

	p.SetTemplate("components/user/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}
