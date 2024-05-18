package service

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
)

type Experiment interface {
	GetSuggestionsForUser(ctx context.Context, userUID string) (response.ExperimentSuggestions, error)
	AddExperiment(ctx context.Context, req request.AddExperiment, userUID string) error
	GetAll(ctx context.Context) ([]response.Experiment, error)
	GetAllForUser(ctx context.Context, userUID string) ([]response.Experiment, error)
	DeleteByUID(ctx context.Context, uid string) error
	DeleteByUIDForUser(ctx context.Context, uid, userUID string) error
}
