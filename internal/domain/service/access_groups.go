package service

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
)

type AccessGroup interface {
	GetAllNotCanceled(ctx context.Context) ([]response.AccessGroup, error)
	Create(ctx context.Context, r request.CreateAccessGroup) (response.AccessGroup, error)
	Edit(ctx context.Context, r request.EditAccessGroup) (response.AccessGroup, error)
	GetByUID(ctx context.Context, uid string) (response.AccessGroup, error)
	DeleteByUID(ctx context.Context, uid string) error
}
