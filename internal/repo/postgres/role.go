package postgres

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
	"github.com/vebcreatex7/diploma_magister/internal/repo/postgres/schema"
)

type roles struct {
	db *goqu.Database
}

func NewRoles(db *goqu.Database) roles {
	return roles{db: db}
}

func (r roles) GetByName(ctx context.Context, name string) (entities.Role, bool, error) {
	var res entities.Role

	q := r.db.From(schema.Role).
		Select(entities.Role{}).
		Where(goqu.C("name").Eq(name)).
		Prepared(true).Executor()

	found, err := q.ScanStructContext(ctx, &res)

	return res, found, err
}
