package postgres

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
	"github.com/vebcreatex7/diploma_magister/internal/repo/postgres/schema"
)

type clients struct {
	db *goqu.Database
}

func NewClients(db *goqu.Database) clients {
	return clients{db: db}
}

func (r clients) Create(ctx context.Context, c entities.Client) error {
	q := r.db.Insert(schema.Client).
		Rows(c).
		Prepared(true).Executor()

	if _, err := q.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (r clients) GetByLogin(ctx context.Context, login string) (entities.Client, bool, error) {
	var res entities.Client

	q := r.db.From(schema.Client).
		Select(entities.Client{}).
		Where(goqu.I("login").Eq(login)).
		Prepared(true).Executor()

	found, err := q.ScanStructContext(ctx, &res)

	return res, found, err
}

func (r clients) GetAllNotCanceled(ctx context.Context) (res []entities.Client, err error) {
	return res, r.db.From(schema.Client).
		Select(entities.Client{}).
		Prepared(true).Executor().ScanStructsContext(ctx, &res)
}
