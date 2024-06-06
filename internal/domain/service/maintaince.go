package service

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
)

type Maintaince interface {
	GetSuggestions(ctx context.Context) (response.MaintainceSuggestions, error)
	GetAll(ctx context.Context) ([]response.Maintaince, error)
	GetAllForUser(ctx context.Context, userUID string) ([]response.Maintaince, error)
	AddMaintaince(ctx context.Context, req request.AddMaintaince, userUID string) error
	DeleteByUIDForUser(ctx context.Context, uid, userUID string) error
	DeleteByUID(ctx context.Context, uid string) error
}
