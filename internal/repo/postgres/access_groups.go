package postgres

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
	"github.com/vebcreatex7/diploma_magister/internal/repo/postgres/schema"
)

type accessGroups struct {
	db *goqu.Database
}

func NewAccessGroups(db *goqu.Database) accessGroups {
	return accessGroups{db: db}
}

func (r accessGroups) Create(ctx context.Context, e entities.AccessGroup) (entities.AccessGroup, error) {
	var res entities.AccessGroup

	_, err := r.db.Insert(schema.AccessGroup).
		Rows(e).
		Returning(entities.AccessGroup{}).
		Prepared(true).Executor().ScanStructContext(ctx, &res)

	return res, err
}

func (r accessGroups) GetAllNotCanceled(ctx context.Context) (res []entities.AccessGroup, err error) {
	return res, r.db.From(schema.AccessGroup).
		Select(entities.AccessGroup{}).
		Order(goqu.C("uid").Desc()).
		Prepared(true).Executor().
		ScanStructsContext(ctx, &res)
}

func (r accessGroups) DeleteByUID(ctx context.Context, uid string) error {
	_, err := r.db.Delete(schema.AccessGroup).
		Where(goqu.I("uid").Eq(uid)).
		Prepared(true).Executor().ExecContext(ctx)

	return err
}

func (r accessGroups) GetByUID(ctx context.Context, uid string) (entities.AccessGroup, bool, error) {
	var res entities.AccessGroup

	found, err := r.db.From(schema.AccessGroup).
		Select(entities.AccessGroup{}).
		Where(goqu.I("uid").Eq(uid)).
		Prepared(true).Executor().
		ScanStructContext(ctx, &res)

	return res, found, err
}

func (r accessGroups) Edit(ctx context.Context, e entities.AccessGroup) (entities.AccessGroup, bool, error) {
	var res entities.AccessGroup

	q := r.db.Update(schema.AccessGroup).
		Set(e).
		Where(goqu.I("uid").Eq(e.UID)).
		Returning(entities.AccessGroup{}).
		Prepared(true).Executor()

	edited, err := q.ScanStructContext(ctx, &res)

	return res, edited, err
}

func (r accessGroups) CreateClientsInAccessGroup(ctx context.Context, e []entities.ClientsInAccessGroup) error {
	_, err := r.db.Insert(schema.ClientsInAccessGroup).
		Rows(e).
		Prepared(true).Executor().ExecContext(ctx)

	return err
}

func (r accessGroups) CreateEquipmentInAccessGroup(ctx context.Context, e []entities.EquipmentInAccessGroup) error {
	q := r.db.Insert(schema.EquipmentInAccessGroup).
		Rows(e).
		Prepared(true).Executor()

	_, err := q.ExecContext(ctx)

	return err
}

func (r accessGroups) CreateInventoryInAccessGroup(ctx context.Context, e []entities.InventoryInAccessGroup) error {

	_, err := r.db.Insert(schema.InventoryInAccessGroup).
		Rows(e).
		Prepared(true).Executor().ExecContext(ctx)

	return err
}

func (r accessGroups) DeleteClientsInAccessGroupByUID(ctx context.Context, uid string) error {
	_, err := r.db.Delete(schema.ClientsInAccessGroup).
		Where(goqu.I("access_group_uid").Eq(uid)).
		Prepared(true).Executor().ExecContext(ctx)

	return err
}

func (r accessGroups) DeleteEquipmentInAccessGroupByUID(ctx context.Context, uid string) error {
	_, err := r.db.Delete(schema.EquipmentInAccessGroup).
		Where(goqu.I("access_group_uid").Eq(uid)).
		Prepared(true).Executor().ExecContext(ctx)

	return err
}

func (r accessGroups) DeleteInventoryInAccessGroupByUID(ctx context.Context, uid string) error {
	_, err := r.db.Delete(schema.InventoryInAccessGroup).
		Where(goqu.I("access_group_uid").Eq(uid)).
		Prepared(true).Executor().ExecContext(ctx)

	return err
}
