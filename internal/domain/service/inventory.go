package service

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
)

type Inventory interface {
	GetAllNotCanceled(ctx context.Context) ([]response.Inventory, error)
	DeleteByUID(ctx context.Context, uid string) error
	Edit(ctx context.Context, req request.EditInventory) (response.Inventory, error)
	GetByUID(ctx context.Context, uid string) (response.Inventory, error)
	Create(ctx context.Context, req request.CreateInventory) (response.Inventory, error)
}
