package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/vebcreatex7/diploma_magister/pkg/request"
	"net/http"
	"strings"
	"time"
)

type CreateEquipment struct {
	Name         string
	Description  string
	Type         string
	Manufacturer string
	Model        string
	Room         string
}

func (r *CreateEquipment) Bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("parsing form: %w", err)
	}

	r.Name = strings.TrimSpace(req.Form.Get("name"))
	r.Description = strings.TrimSpace(req.Form.Get("description"))
	r.Type = strings.TrimSpace(req.Form.Get("type"))
	r.Manufacturer = strings.TrimSpace(req.Form.Get("manufacturer"))
	r.Model = strings.TrimSpace(req.Form.Get("model"))
	r.Room = strings.TrimSpace(req.Form.Get("room"))

	return r.validate()
}

func (r *CreateEquipment) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.Type, validation.Required),
		validation.Field(&r.Manufacturer, validation.Required),
		validation.Field(&r.Model, validation.Required),
		validation.Field(&r.Room, validation.Required),
	)
}

type EditEquipment struct {
	UID          string
	Name         string
	Description  string
	Type         string
	Manufacturer string
	Model        string
	Room         string
}

func (r *EditEquipment) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("getting uid from path: %w", err)
	}

	r.UID = uid

	if err = req.ParseForm(); err != nil {
		return fmt.Errorf("parsing form: %w", err)
	}

	r.Name = strings.TrimSpace(req.Form.Get("name"))
	r.Description = strings.TrimSpace(req.Form.Get("description"))
	r.Type = strings.TrimSpace(req.Form.Get("type"))
	r.Manufacturer = strings.TrimSpace(req.Form.Get("manufacturer"))
	r.Model = strings.TrimSpace(req.Form.Get("model"))
	r.Room = strings.TrimSpace(req.Form.Get("room"))

	return nil
}

func (r *EditEquipment) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.Type, validation.Required),
		validation.Field(&r.Manufacturer, validation.Required),
		validation.Field(&r.Model, validation.Required),
		validation.Field(&r.Room, validation.Required),
	)
}

type DeleteEquipment struct {
	UID string
}

func (r *DeleteEquipment) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("getting uid from path: %w", err)
	}

	r.UID = uid

	return nil
}

type GetEquipment struct {
	UID string
}

func (r *GetEquipment) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("getting uid from path: %w", err)
	}

	r.UID = uid

	return nil
}

type GetEquipmentSchedule struct {
	Name  string
	Lower time.Time
	Upper time.Time
}

func (r *GetEquipmentSchedule) Bind(req *http.Request) error {
	name, err := request.ParseStringFromQuery(req, "name", true)
	if err != nil {
		return fmt.Errorf("geting name for query: %w", err)
	}

	lower, err := request.ParseTimeFromQuery(req, "lower", true)
	if err != nil {
		return fmt.Errorf("geting lower for query: %w", err)
	}

	upper, err := request.ParseTimeFromQuery(req, "upper", true)
	if err != nil {
		return fmt.Errorf("geting upper for query: %w", err)
	}

	r.Name = name
	r.Lower = lower
	r.Upper = upper.Add(time.Hour * 24)

	return r.validate()
}

func (r *GetEquipmentSchedule) validate() error {
	if !r.Upper.After(r.Lower) {
		return fmt.Errorf("wrong interval order")
	}

	if r.Upper.Sub(r.Lower) > time.Hour*24*30 {
		return fmt.Errorf("duartion bigger than 30 days")
	}

	return nil
}
