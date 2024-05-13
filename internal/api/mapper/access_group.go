package mapper

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/domain/constant"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
	"strings"
)

type AccessGroup struct {
}

func (m AccessGroup) MakeResponse(e entities.AccessGroupExt) response.AccessGroup {
	return response.AccessGroup{
		UID:         e.UID,
		Name:        e.Name,
		Description: e.Description,
		Users:       strings.Join(e.Users, ","),
		Equipment:   strings.Join(e.Equipment, ","),
		Inventory:   strings.Join(e.Inventory, ","),
	}
}

func (m AccessGroup) MakeListResponse(e []entities.AccessGroupExt) []response.AccessGroup {
	res := make([]response.AccessGroup, 0, len(e))

	for i := range e {
		res = append(res, m.MakeResponse(e[i]))
	}

	return res
}

func (m AccessGroup) MakeEditEntity(r request.EditAccessGroup) entities.AccessGroupExt {
	return entities.AccessGroupExt{
		AccessGroup: entities.AccessGroup{
			UID:         r.UID,
			Name:        r.Name,
			Description: r.Description,
		},
		Users:     r.Users,
		Equipment: r.Equipment,
		Inventory: r.Inventory,
	}
}

func (m AccessGroup) MakeCreateEntity(r request.CreateAccessGroup) entities.AccessGroupExt {
	return entities.AccessGroupExt{
		AccessGroup: entities.AccessGroup{
			Name:        r.Name,
			Description: r.Description,
			Status:      constant.StatusReady,
		},
		Users:     r.Users,
		Equipment: r.Equipment,
		Inventory: r.Inventory,
	}
}
