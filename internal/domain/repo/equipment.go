package repo

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
)

type Equipment interface {
	GetAllNotCanceled(ctx context.Context) (res []entities.Equipment, err error)
	DeleteByUID(ctx context.Context, uid string) error
	GetByUID(ctx context.Context, uid string) (entities.Equipment, bool, error)
	Edit(ctx context.Context, e entities.Equipment) (entities.Equipment, bool, error)
	Create(ctx context.Context, e entities.Equipment) (entities.Equipment, error)
	GetNamesByGroupUID(ctx context.Context, uid string) (res []string, err error)
	GetByName(ctx context.Context, name string) (entities.Equipment, bool, error)
	DeleteEquipmentInAccessGroupByUID(ctx context.Context, uid string) error
}
