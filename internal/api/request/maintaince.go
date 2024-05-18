package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/vebcreatex7/diploma_magister/internal/domain/constant"
	"net/http"
	"strings"
	"time"
)

type AddMaintainceAdmin struct {
	User string
	AddMaintaince
}

func (r *AddMaintainceAdmin) Bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	r.User = strings.TrimSpace(req.Form.Get("maintaince-engineer"))

	return r.AddMaintaince.Bind(req)
}

type AddMaintaince struct {
	Name        string
	Description string
	StartTs     time.Time
	EndTs       time.Time
	Users       []string
	Equipment   []equipmentInMaintaince
}

type equipmentInMaintaince struct {
	Name  string
	Lower time.Time
	Upper time.Time
}

func (r *AddMaintaince) Bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	r.Name = strings.TrimSpace(req.Form.Get("maintaince-name"))
	r.Description = strings.TrimSpace(req.Form.Get("maintaince-description"))

	start, err := time.Parse(constant.Layout, strings.TrimSpace(req.Form.Get("start-ts")))
	if err != nil {
		return fmt.Errorf("parsing start-ts: %w", err)
	}

	end, err := time.Parse(constant.Layout, strings.TrimSpace(req.Form.Get("end-ts")))
	if err != nil {
		return fmt.Errorf("parsing end-ts: %w", err)
	}

	r.StartTs = start
	r.EndTs = end

	users := req.Form["users"]
	equipment := req.Form["equipment-name"]
	lower := req.Form["lower"]
	upper := req.Form["upper"]

	if !(len(equipment) == len(lower) && len(lower) == len(upper)) {
		return fmt.Errorf("unequal number of params for equipment")
	}

	var eq []equipmentInMaintaince
	for i := range equipment {
		l, err := time.Parse(constant.Layout, lower[i])
		if err != nil {
			return fmt.Errorf("parsing %d lower: %w", i, err)
		}
		u, err := time.Parse(constant.Layout, upper[i])
		if err != nil {
			return fmt.Errorf("parsing %d upper: %w", i, err)
		}

		eq = append(eq, equipmentInMaintaince{
			Name:  strings.TrimSpace(equipment[i]),
			Lower: l,
			Upper: u,
		})
	}

	for i := range users {
		users[i] = strings.TrimSpace(users[i])
	}

	r.Users = users
	r.Equipment = eq
	fmt.Println(r)

	return r.validate()
}

func (r *AddMaintaince) validate() error {
	if err := validation.Validate(r.Name, validation.Required); err != nil {
		return fmt.Errorf("validating name: %w", err)
	}
	if err := validation.Validate(r.Description, validation.Required); err != nil {
		return fmt.Errorf("validating description: %w", err)
	}
	//if err := validation.Validate(r.Users, validation.Length(1, 0)); err != nil {
	//	return fmt.Errorf("valdaiting users: %w", err)
	//}

	for i := range r.Equipment {
		if r.Equipment[i].Lower.Before(r.StartTs) {
			return fmt.Errorf("validating equipment %d lower: should be after start-ts", i)
		}

		if r.Equipment[i].Upper.After(r.EndTs) {
			return fmt.Errorf("validating equipment %d upper: should be before end-ts", i)
		}
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

	return nil
}
