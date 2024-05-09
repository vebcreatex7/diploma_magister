package mapper

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
)

type Equipment struct {
}

func (m Equipment) MakeResponse(e entities.Equipment) response.Equipment {
	return response.Equipment{
		UID:          e.UID,
		Name:         e.Name,
		Description:  e.Description,
		Type:         e.Type,
		Manufacturer: e.Manufacturer,
		Model:        e.Model,
		Status:       e.Status,
		Room:         e.Room,
	}
}

func (m Equipment) MakeListResponse(e []entities.Equipment) []response.Equipment {
	res := make([]response.Equipment, 0, len(e))

	for i := range e {
		res = append(res, m.MakeResponse(e[i]))
	}

	return res
}
