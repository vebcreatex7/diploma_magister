package service

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
)

type Clients interface {
	Create(ctx context.Context, req request.CreateClient) error
	Login(ctx context.Context, req request.LoginClient) (string, error)
	GetAllNotCanceled(ctx context.Context) ([]response.Client, error)
	DeleteByUID(ctx context.Context, uid string) ([]response.Client, error)
}
