package engineer

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h engineer) GetAllInventory(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage()
	eq, err := h.inventoryService.GetAll(r.Context())
	if err != nil {
		h.log.WithError(err).Errorf("getting all")

		p.SetTemplate("engineer.gohtml").
			SetPath(r.URL.Path).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("engineer.gohtml").
		SetPath(r.URL.Path).
		SetData(eq).
		SetCode(200)
	h.t.Render(w, p)
}

func (h engineer) DeleteInventoryByUID(w http.ResponseWriter, r *http.Request) {
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

	if err := h.inventoryService.DeleteByUID(r.Context(), req.UID); err != nil {
		h.log.WithError(err).Errorf("deleting")
		p.SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	w.WriteHeader(200)
}

func (h engineer) GetInventoryEditByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetInventory
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.inventoryService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/engineer/inventory/row_edit.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h engineer) GetInventoryByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.GetInventory
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.inventoryService.GetByUID(r.Context(), req.UID)
	if err != nil {
		h.log.WithError(err).Errorf("getting")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("components/engineer/inventory/row.gohtml").
		SetData(resp).
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h engineer) EditInventoryByUID(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.EditInventory
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	resp, err := h.inventoryService.Edit(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("editing")
		p.SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	p.SetTemplate("components/engineer/inventory/row.gohtml").
		SetData(resp).
		SetSuccess("inventory edited").
		SetCode(200)

	h.t.RenderData(w, p)
}

func (h engineer) AddInventoryPage(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	p.SetTemplate("components/engineer/inventory/row_add.gohtml").SetCode(200)

	h.t.Render(w, p)
}

func (h engineer) AddInventory(w http.ResponseWriter, r *http.Request) {
	var (
		p   = render.NewPage()
		req request.CreateInventory
	)

	if err := req.Bind(r); err != nil {
		h.log.WithError(err).Errorf("binding request")
		p.SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	resp, err := h.inventoryService.Create(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Errorf("creating")
		p.SetCode(422).
			SetError(err.Error())

		h.t.Render(w, p)

		return
	}

	p.SetTemplate("components/engineer/inventory/row.gohtml").
		SetData(resp).
		SetSuccess("inventory added").
		SetCode(200)

	h.t.RenderData(w, p)
}
