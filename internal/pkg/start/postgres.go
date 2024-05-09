package start

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/vebcreatex7/diploma_magister/internal/config"
)

func Postgres(cfg config.Postgres) (*goqu.Database, error) {
	pgxDB, err := sqlx.Connect("pgx", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	return goqu.New("postgres", pgxDB), nil

}
