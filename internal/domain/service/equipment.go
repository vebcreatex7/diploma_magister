package service

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
)

type Equipment interface {
	GetAllNotCanceled(ctx context.Context) ([]response.Equipment, error)
	DeleteByUID(ctx context.Context, uid string) ([]response.Equipment, error)
}
