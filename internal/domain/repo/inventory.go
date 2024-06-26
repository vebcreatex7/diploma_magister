package repo

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
)

type Inventory interface {
	GetAll(ctx context.Context) (res []entities.Inventory, err error)
	DeleteByUID(ctx context.Context, uid string) error
	GetByUID(ctx context.Context, uid string) (entities.Inventory, bool, error)
	Edit(ctx context.Context, e entities.Inventory) (entities.Inventory, bool, error)
	Create(ctx context.Context, e entities.Inventory) (entities.Inventory, error)
	GetNamesByGroupUID(ctx context.Context, uid string) (res []string, err error)
	GetByName(ctx context.Context, name string) (entities.Inventory, bool, error)
	DeleteInventoryInAccessGroupByUID(ctx context.Context, uid string) error
}
