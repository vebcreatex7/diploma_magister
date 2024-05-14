package service

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
)

type Equipment interface {
	GetAll(ctx context.Context) ([]response.Equipment, error)
	DeleteByUID(ctx context.Context, uid string) error
	Edit(ctx context.Context, req request.EditEquipment) (response.Equipment, error)
	GetByUID(ctx context.Context, uid string) (response.Equipment, error)
	Create(ctx context.Context, req request.CreateEquipment) (response.Equipment, error)
}
