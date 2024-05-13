package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/vebcreatex7/diploma_magister/pkg/request"
	"net/http"
)

type CreateInventory struct {
	Name         string
	Description  string
	Type         string
	Manufacturer string
	Quantity     string
	Unit         string
}

func (r *CreateInventory) Bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		{
			return fmt.Errorf("parsing form: %w", err)
		}
	}

	r.Name = req.Form.Get("name")
	r.Description = req.Form.Get("description")
	r.Type = req.Form.Get("type")
	r.Manufacturer = req.Form.Get("manufacturer")
	r.Quantity = req.Form.Get("quantity")
	r.Unit = req.Form.Get("unit")

	return r.validate()
}

func (r *CreateInventory) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.Type, validation.Required),
		validation.Field(&r.Manufacturer, validation.Required),
		validation.Field(&r.Quantity, validation.Required),
		validation.Field(&r.Unit, validation.Required),
	)
}

type EditInventory struct {
	UID          string
	Name         string
	Description  string
	Type         string
	Manufacturer string
	Quantity     string
	Unit         string
	Status       string
}

func (r *EditInventory) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("getting uid from path: %w", err)
	}

	r.UID = uid

	if err = req.ParseForm(); err != nil {
		return fmt.Errorf("parsing form: %w", err)
	}

	r.Name = req.Form.Get("name")
	r.Description = req.Form.Get("description")
	r.Type = req.Form.Get("type")
	r.Manufacturer = req.Form.Get("manufacturer")
	r.Quantity = req.Form.Get("quantity")
	r.Unit = req.Form.Get("unit")
	r.Status = req.Form.Get("status")

	return nil
}

type GetInventory struct {
	UID string
}

func (r *GetInventory) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("getting uid from path: %w", err)
	}

	r.UID = uid

	return nil
}

type DeleteInventory struct {
	UID string
}

func (r *DeleteInventory) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("getting uid from path: %w", err)
	}

	r.UID = uid

	return nil
}
