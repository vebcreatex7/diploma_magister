package repo

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
)

type Clients interface {
	Create(ctx context.Context, c entities.Client) error
	GetByLogin(ctx context.Context, login string) (entities.Client, bool, error)
	GetAllNotCanceled(ctx context.Context) (res []entities.Client, err error)
	DeleteByUID(ctx context.Context, uid string) error
}
