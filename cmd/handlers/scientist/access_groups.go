package scientist

import (
	"fmt"
	"github.com/vebcreatex7/diploma_magister/pkg"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

func (h scientist) GetMyAccessGroups(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage()
	uid, err := pkg.GetUIDFromJWT(r)
	if err != nil {
		h.log.WithError(err).Errorf("getting uid from jwt")

		p.SetTemplate("scientist.gohtml").
			SetPath(r.URL.Path).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	res, err := h.accessGroupsService.GetAllForGivenUser(r.Context(), uid)
	if err != nil {
		h.log.WithError(err).Errorf("getting access groups")

		p.SetTemplate("scientist.gohtml").
			SetPath(r.URL.Path).
			SetError(err.Error())

		h.t.Render(w, p)
		return
	}

	p.SetTemplate("scientist.gohtml").
		SetPath(r.URL.Path).
		SetData(res).
		SetCode(200)

	h.t.Render(w, p)

	fmt.Println(uid)
}
