package service

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
)

type experiment struct {
	db *goqu.Database
}

func NewExperiment(db *goqu.Database) experiment {
	return experiment{db: db}
}

func (s experiment) GetSuggestionsForUser(ctx context.Context, userUID string) (response.ExperimentSuggestions, error) {
	var eq []string

	if err := s.db.ScanValsContext(
		ctx,
		&eq,
		`select eq.name from equipment eq
join equipment_in_access_group eqag on eq.uid = eqag.equipment_uid
join clients_in_access_group cag on eqag.access_group_uid = cag.access_group_uid
join client c on cag.client_uid = c.uid
where c.uid = $1`,
		userUID,
	); err != nil {
		return response.ExperimentSuggestions{}, fmt.Errorf("getting equipment suggestions: %w", err)
	}

	var in []string
	if err := s.db.ScanValsContext(
		ctx,
		&in,
		`select i.name from inventory i
                        join inventory_in_access_group iag on i.uid = iag.inventory_uid
                        join clients_in_access_group cag on iag.access_group_uid = cag.access_group_uid
                        join client c on cag.client_uid = c.uid
where c.uid = $1`,
		userUID,
	); err != nil {
		return response.ExperimentSuggestions{}, fmt.Errorf("getting inventory suggestions: %w", err)
	}

	return response.ExperimentSuggestions{
		Equipment: eq,
		Inventory: in,
	}, nil
}
