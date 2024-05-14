package entities

import "github.com/shopspring/decimal"

type Inventory struct {
	UID          string          `db:"uid" goqu:"skipinsert,skipupdate"`
	Name         string          `db:"name"`
	Description  string          `db:"description"`
	Type         string          `db:"type"`
	Manufacturer string          `db:"manufacturer"`
	Quantity     decimal.Decimal `db:"quantity"`
	Unit         string          `db:"unit"`
}
