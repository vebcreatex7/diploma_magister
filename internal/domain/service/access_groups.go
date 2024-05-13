package service

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
)

type AccessGroup interface {
	GetAllNotCanceled(ctx context.Context) ([]response.AccessGroup, error)
}
