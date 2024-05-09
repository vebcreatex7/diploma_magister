package repo

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
)

type Equipment interface {
	GetAllNotCanceled(ctx context.Context) (res []entities.Equipment, err error)
	DeleteByUID(ctx context.Context, uid string) error
}
