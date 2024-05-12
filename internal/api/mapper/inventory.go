package mapper

import (
	"github.com/shopspring/decimal"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/domain/constant"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
)

type Inventory struct {
}

func (m Inventory) MakeResponse(e entities.Inventory) response.Inventory {
	return response.Inventory{
		UID:          e.UID,
		Name:         e.Name,
		Description:  e.Description,
		Type:         e.Type,
		Manufacturer: e.Manufacturer,
		Quantity:     e.Quantity.String(),
		Unit:         e.Unit,
		Status:       e.Status,
	}
}

func (m Inventory) MakeListResponse(e []entities.Inventory) []response.Inventory {
	res := make([]response.Inventory, 0, len(e))

	for i := range e {
		res = append(res, m.MakeResponse(e[i]))
	}

	return res
}

func (m Inventory) MakeEditEntity(r request.EditInventory) entities.Inventory {
	return entities.Inventory{
		UID:          r.UID,
		Name:         r.Name,
		Description:  r.Description,
		Type:         r.Type,
		Manufacturer: r.Manufacturer,
		Quantity:     decimal.RequireFromString(r.Quantity),
		Unit:         r.Unit,
		Status:       r.Status,
	}
}

func (m Inventory) MakeCreateEntity(r request.CreateInventory) entities.Inventory {
	return entities.Inventory{
		Name:         r.Name,
		Description:  r.Description,
		Type:         r.Type,
		Manufacturer: r.Manufacturer,
		Quantity:     decimal.RequireFromString(r.Quantity),
		Unit:         r.Unit,
		Status:       constant.StatusWaitApprove,
	}
}
