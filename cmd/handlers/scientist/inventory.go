package scientist

import (
	"github.com/vebcreatex7/diploma_magister/pkg"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h scientist) GetAllMyInventory(w http.ResponseWriter, r *http.Request) {
	var p = render.NewPage()

	uid, err := pkg.GetUIDFromJWT(r)
	if err != nil {
		h.log.WithError(err).Errorf("getting uid from jwt")
		p.SetError(err.Error())
		h.t.Render(w, p)
		return
	}

	res, err := h.inventoryService.GetAllForUser(r.Context(), uid)
	if err != nil {
		h.log.WithError(err).Errorf("getting inventory")
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
