package repo

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
)

type Clients interface {
	Create(ctx context.Context, c entities.Client) error
	GetByLogin(ctx context.Context, login string) (entities.Client, bool, error)
	GetAll(ctx context.Context) (res []entities.Client, err error)
	DeleteByUID(ctx context.Context, uid string) error
	GetByUID(ctx context.Context, uid string) (entities.Client, bool, error)
	Edit(ctx context.Context, e entities.Client) (entities.Client, bool, error)
	GetLoginsByGroupUID(ctx context.Context, uid string) (res []string, err error)
	DeleteClientsInAccessGroupByUID(ctx context.Context, uid string) error
	Approve(ctx context.Context, uid string) error
}
