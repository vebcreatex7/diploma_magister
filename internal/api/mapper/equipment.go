package mapper

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
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

func (m Equipment) MakeEditEntity(r request.EditEquipment) entities.Equipment {
	return entities.Equipment{
		UID:          r.UID,
		Name:         r.Name,
		Description:  r.Description,
		Type:         r.Type,
		Manufacturer: r.Manufacturer,
		Model:        r.Model,
		Room:         r.Room,
		Status:       r.Status,
	}
}

func (m Equipment) MakeCreateEntity(r request.CreateEquipment) entities.Equipment {
	return entities.Equipment{
		Name:         r.Name,
		Description:  r.Description,
		Type:         r.Type,
		Manufacturer: r.Manufacturer,
		Model:        r.Model,
		Room:         r.Room,
		Status:       r.Status,
	}
}
