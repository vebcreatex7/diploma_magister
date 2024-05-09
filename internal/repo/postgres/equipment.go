package postgres

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
	"github.com/vebcreatex7/diploma_magister/internal/repo/postgres/schema"
)

type equipment struct {
	db *goqu.Database
}

func NewEquipment(db *goqu.Database) equipment {
	return equipment{db: db}
}

func (r equipment) GetAllNotCanceled(ctx context.Context) (res []entities.Equipment, err error) {
	return res, r.db.From(schema.Equipment).
		Select(entities.Equipment{}).
		Where(goqu.I("status").Neq("cancel")).
		Prepared(true).Executor().ScanStructsContext(ctx, &res)
}

func (r equipment) DeleteByUID(ctx context.Context, uid string) error {
	_, err := r.db.Update(schema.Equipment).
		Set(goqu.Record{"status": "cancel"}).
		Where(goqu.I("uid").Eq(uid)).
		Prepared(true).Executor().ExecContext(ctx)

	return err
}
