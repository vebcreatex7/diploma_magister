package service

import (
	"context"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
)

type Experiment interface {
	GetSuggestionsForUser(ctx context.Context, userUID string) (response.ExperimentSuggestions, error)
}
