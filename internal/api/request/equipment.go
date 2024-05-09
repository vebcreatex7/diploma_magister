package request

import (
	"fmt"
	"github.com/vebcreatex7/diploma_magister/pkg/request"
	"net/http"
)

type DeleteEquipment struct {
	UID string
}

func (r *DeleteEquipment) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("getting uid from path: %w", err)
	}

	r.UID = uid

	return nil
}
