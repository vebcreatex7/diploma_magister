package postgres

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
	"github.com/vebcreatex7/diploma_magister/internal/repo/postgres/schema"
)

type inventory struct {
	db *goqu.Database
}

func NewInventory(db *goqu.Database) inventory {
	return inventory{db: db}
}

func (r inventory) GetAllNotCanceled(ctx context.Context) (res []entities.Inventory, err error) {
	return res, r.db.From(schema.Inventory).
		Select(entities.Inventory{}).
		Where(goqu.I("status").Neq("cancel")).
		Order(goqu.C("uid").Desc()).
		Prepared(true).Executor().ScanStructsContext(ctx, &res)
}

func (r inventory) DeleteByUID(ctx context.Context, uid string) error {
	_, err := r.db.Update(schema.Inventory).
		Set(goqu.Record{"status": "cancel"}).
		Where(goqu.I("uid").Eq(uid)).
		Prepared(true).Executor().ExecContext(ctx)

	return err
}

func (r inventory) GetByUID(ctx context.Context, uid string) (entities.Inventory, bool, error) {
	var res entities.Inventory

	q := r.db.From(schema.Inventory).
		Select(entities.Inventory{}).
		Where(goqu.I("uid").Eq(uid)).
		Prepared(true).Executor()

	found, err := q.ScanStructContext(ctx, &res)

	return res, found, err
}

func (r inventory) Edit(ctx context.Context, e entities.Inventory) (entities.Inventory, bool, error) {
	var res entities.Inventory

	q := r.db.Update(schema.Inventory).
		Set(e).
		Where(goqu.I("uid").Eq(e.UID)).
		Returning(entities.Inventory{}).
		Prepared(true).Executor()

	edited, err := q.ScanStructContext(ctx, &res)

	return res, edited, err
}

func (r inventory) Create(ctx context.Context, e entities.Inventory) (entities.Inventory, error) {
	var res entities.Inventory

	q := r.db.Insert(schema.Inventory).
		Rows(e).
		Returning(entities.Inventory{}).
		Prepared(true).Executor()

	_, err := q.
		ScanStructContext(ctx, &res)

	return res, err
}
