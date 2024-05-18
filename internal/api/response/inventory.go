package response

import "github.com/shopspring/decimal"

type Inventory struct {
	UID          string
	Name         string
	Description  string
	Type         string
	Manufacturer string
	Quantity     string
	Unit         string
}

type InventoryInExperiment struct {
	Name     string
	Quantity decimal.Decimal
}
