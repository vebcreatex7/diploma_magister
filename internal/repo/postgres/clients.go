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
		Where(goqu.I("status").Neq("cancel")).
		Order(goqu.C("uid").Desc()).
		Prepared(true).Executor().ScanStructsContext(ctx, &res)
}

func (r clients) DeleteByUID(ctx context.Context, uid string) error {
	_, err := r.db.Update(schema.Client).
		Set(goqu.Record{"status": "cancel"}).
		Where(goqu.I("uid").Eq(uid)).
		Prepared(true).Executor().ExecContext(ctx)

	return err
}

func (r clients) GetByUID(ctx context.Context, uid string) (entities.Client, bool, error) {
	var res entities.Client

	q := r.db.From(schema.Client).
		Select(entities.Client{}).
		Where(goqu.I("uid").Eq(uid)).
		Prepared(true).Executor()

	found, err := q.ScanStructContext(ctx, &res)

	return res, found, err
}

func (r clients) Edit(ctx context.Context, e entities.Client) (entities.Client, bool, error) {
	var res entities.Client

	q := r.db.Update(schema.Client).
		Set(e).
		Where(goqu.I("uid").Eq(e.UID)).
		Returning(entities.Client{}).
		Prepared(true).Executor()

	edited, err := q.ScanStructContext(ctx, &res)

	return res, edited, err
}

func (r clients) GetLoginsByGroupUID(ctx context.Context, uid string) (res []string, err error) {
	return res, r.db.From(schema.Client).
		Select("login").
		Join(goqu.T(schema.ClientsInAccessGroup),
			goqu.On(goqu.I(schema.Client+".uid").Eq(goqu.I(schema.ClientsInAccessGroup+".client_uid"))),
		).
		Where(
			goqu.I(schema.ClientsInAccessGroup+".access_group_uid").Eq(uid),
			goqu.I(schema.Client+".status").Neq("cancel"),
		).
		Prepared(true).Executor().
		ScanValsContext(ctx, &res)
}
