package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

var (
	layout      = "2006-01-02T15:04"
	maxDuration = time.Hour * 24 * 2
)

type AddExperiment struct {
	Name        string
	Description string
	Equipment   []equipmentInExperiment
	Inventory   []inventoryInExperiment
}

type equipmentInExperiment struct {
	Name  string
	Lower time.Time
	Upper time.Time
}

type inventoryInExperiment struct {
	Name     string
	Quantity decimal.Decimal
}

func (r *AddExperiment) Bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	r.Name = req.Form.Get("experiment-name")
	r.Description = req.Form.Get("experiment-description")

	equipment := req.Form["equipment-name"]
	lower := req.Form["lower"]
	upper := req.Form["upper"]

	if !(len(equipment) == len(lower) && len(lower) == len(upper)) {
		return fmt.Errorf("unequal number of params for equipment")
	}

	var eq []equipmentInExperiment
	for i := range equipment {
		l, err := time.Parse(layout, lower[i])
		if err != nil {
			return fmt.Errorf("parsing %d lower: %w", i, err)
		}
		u, err := time.Parse(layout, upper[i])
		if err != nil {
			return fmt.Errorf("parsing %d upper: %w", i, err)
		}

		eq = append(eq, equipmentInExperiment{
			Name:  equipment[i],
			Lower: l,
			Upper: u,
		})
	}

	inventory := req.Form["inventory-name"]
	quantity := req.Form["quantity"]

	if len(quantity) != len(inventory) {
		return fmt.Errorf("unequal number of params for inventory")
	}

	var in []inventoryInExperiment
	for i := range inventory {
		q, err := decimal.NewFromString(quantity[i])
		if err != nil {
			return fmt.Errorf("parsing %d quantity: %w", i, err)
		}

		in = append(in, inventoryInExperiment{
			Name:     inventory[i],
			Quantity: q,
		})
	}

	r.Equipment = eq
	r.Inventory = in
	fmt.Println(r)

	return r.validate()
}

func (r *AddExperiment) validate() error {
	if err := validation.Validate(r.Name, validation.Required); err != nil {
		return fmt.Errorf("validating name: %w", err)
	}
	if err := validation.Validate(r.Description, validation.Required); err != nil {
		return fmt.Errorf("validating description: %w", err)
	}

	for i := range r.Equipment {
		if err := validation.Validate(r.Equipment[i].Name, validation.Required); err != nil {
			return fmt.Errorf("validating equipment %d name: %w", i, err)
		}
		if err := validation.Validate(r.Equipment[i].Lower, validation.Required); err != nil {
			return fmt.Errorf("validating equipment %d lower: %w", i, err)
		}
		if err := validation.Validate(r.Equipment[i].Upper, validation.Required); err != nil {
			return fmt.Errorf("validating equipment %d upper: %w", i, err)
		}

		if !r.Equipment[i].Upper.After(r.Equipment[i].Lower) {
			return fmt.Errorf("validating %d interval: upper should be after lower", i)
		}

		if r.Equipment[i].Upper.Sub(r.Equipment[i].Lower) > maxDuration {
			return fmt.Errorf("validating %d interval: shouldn't be bigger than 2 days", i)
		}
	}

	for i := range r.Inventory {
		if err := validation.Validate(r.Inventory[i].Name, validation.Required); err != nil {
			return fmt.Errorf("validating inventory %d name: %w", i, err)
		}

		if err := validation.Validate(r.Inventory[i].Quantity, validation.Required); err != nil {
			return fmt.Errorf("validating inventory %d quantity: %w", i, err)
		}

		if !r.Inventory[i].Quantity.GreaterThan(decimal.NewFromInt(0)) {
			return fmt.Errorf("validating inventory %d quantity: quantity less than 0", i)
		}
	}

	return nil
}
