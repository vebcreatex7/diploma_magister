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

type InventoryInExperiment struct {
	ExperimentUID string          `db:"experiment_uid"`
	InventoryUID  string          `db:"inventory_uid"`
	Quantity      decimal.Decimal `db:"quantity"`
}
