package repo

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
)

type AccessGroups interface {
	Create(ctx context.Context, e entities.AccessGroup) (entities.AccessGroup, error)
	GetAllNotCanceled(ctx context.Context) (res []entities.AccessGroup, err error)
	DeleteByUID(ctx context.Context, uid string) error
	GetByUID(ctx context.Context, uid string) (entities.AccessGroup, bool, error)
	Edit(ctx context.Context, e entities.AccessGroup) (entities.AccessGroup, bool, error)
	CreateClientsInAccessGroup(ctx context.Context, e []entities.ClientsInAccessGroup) error
	CreateEquipmentInAccessGroup(ctx context.Context, e []entities.EquipmentInAccessGroup) error
	CreateInventoryInAccessGroup(ctx context.Context, e []entities.InventoryInAccessGroup) error
	DeleteClientsInAccessGroupByUID(ctx context.Context, uid string) error
	DeleteEquipmentInAccessGroupByUID(ctx context.Context, uid string) error
	DeleteInventoryInAccessGroupByUID(ctx context.Context, uid string) error
}
