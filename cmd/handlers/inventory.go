package handlers

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h admin) GetInventory(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage()
	eq, err := h.inventoryService.GetAllNotCanceled(r.Context())
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

func (h admin) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	var req request.DeleteInventory

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		w.WriteHeader(400)

		return
	}

	err := h.inventoryService.DeleteByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("deleting")
		w.WriteHeader(400)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h admin) GetInventoryEditByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetInventory
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		w.WriteHeader(400)

		return
	}

	resp, err := h.inventoryService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		w.WriteHeader(400)

		return
	}

	p.SetTemplate("components/inventory/row_edit.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) GetInventoryByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetInventory
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		w.WriteHeader(400)

		return
	}

	resp, err := h.inventoryService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		w.WriteHeader(400)

		return
	}

	p.SetTemplate("components/inventory/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) EditInventory(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.EditInventory
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		w.WriteHeader(400)

		return
	}

	resp, err := h.inventoryService.Edit(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("editting")
		w.WriteHeader(400)

		return
	}

	p.SetTemplate("components/inventory/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) AddInventoryPage(w http.ResponseWriter, r *http.Request) {
	var (
		p = render.NewPage()
	)

	resp, err := h.inventoryService.GetAllNotCanceled(r.Context())
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		w.WriteHeader(400)

		return
	}

	p.SetTemplate("components/inventory/add.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h admin) AddInventory(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.CreateInventory
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		w.WriteHeader(400)

		return
	}

	resp, err := h.inventoryService.Create(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("creating")
		w.WriteHeader(400)

		return
	}

	p.SetTemplate("components/inventory/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}
