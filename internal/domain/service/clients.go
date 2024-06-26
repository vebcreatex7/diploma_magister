package service

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
)

type Clients interface {
	Create(ctx context.Context, req request.CreateUser) error
	Login(ctx context.Context, req request.LoginUser) (response.User, error)
	GetAll(ctx context.Context) ([]response.User, error)
	DeleteByUID(ctx context.Context, uid string) error
	GetByUID(ctx context.Context, uid string) (response.User, error)
	Edit(ctx context.Context, req request.EditUser) (response.User, error)
	Approve(ctx context.Context, uid string) ([]response.User, error)
	GetByLogin(ctx context.Context, login string) (response.User, error)
}
