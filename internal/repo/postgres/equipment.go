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
		Order(goqu.C("uid").Desc()).
		Prepared(true).Executor().ScanStructsContext(ctx, &res)
}

func (r equipment) DeleteByUID(ctx context.Context, uid string) error {
	_, err := r.db.Update(schema.Equipment).
		Set(goqu.Record{"status": "cancel"}).
		Where(goqu.I("uid").Eq(uid)).
		Prepared(true).Executor().ExecContext(ctx)

	return err
}

func (r equipment) GetByUID(ctx context.Context, uid string) (entities.Equipment, bool, error) {
	var res entities.Equipment

	q := r.db.From(schema.Equipment).
		Select(entities.Equipment{}).
		Where(goqu.I("uid").Eq(uid)).
		Prepared(true).Executor()

	found, err := q.ScanStructContext(ctx, &res)

	return res, found, err
}

func (r equipment) Edit(ctx context.Context, e entities.Equipment) (entities.Equipment, bool, error) {
	var res entities.Equipment

	q := r.db.Update(schema.Equipment).
		Set(e).
		Where(goqu.I("uid").Eq(e.UID)).
		Returning(entities.Equipment{}).
		Prepared(true).Executor()

	edited, err := q.ScanStructContext(ctx, &res)

	return res, edited, err
}

func (r equipment) Create(ctx context.Context, e entities.Equipment) (entities.Equipment, error) {
	var res entities.Equipment

	q := r.db.Insert(schema.Equipment).
		Rows(e).
		Returning(entities.Equipment{}).
		Prepared(true).Executor()

	_, err := q.
		ScanStructContext(ctx, &res)

	return res, err
}

func (r equipment) GetNamesByGroupUID(ctx context.Context, uid string) (res []string, err error) {
	return res, r.db.From(schema.Equipment).
		Select("name").
		Join(goqu.T(schema.EquipmentInAccessGroup),
			goqu.On(goqu.I(schema.Equipment+".uid").Eq(goqu.I(schema.EquipmentInAccessGroup+".equipment_uid"))),
		).
		Where(
			goqu.I(schema.EquipmentInAccessGroup+".access_group_uid").Eq(uid),
			goqu.I(schema.Equipment+".status").Neq("cancel"),
		).
		Prepared(true).Executor().
		ScanValsContext(ctx, &res)
}
