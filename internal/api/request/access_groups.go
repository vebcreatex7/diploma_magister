package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/vebcreatex7/diploma_magister/pkg/request"
	"net/http"
	"strings"
)

type CreateAccessGroup struct {
	Name        string
	Description string
	Users       []string
	Equipment   []string
	Inventory   []string
}

func (r *CreateAccessGroup) Bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		{
			return fmt.Errorf("parsing form: %w", err)
		}
	}

	r.Name = strings.TrimSpace(req.Form.Get("name"))
	r.Description = strings.TrimSpace(req.Form.Get("description"))

	users := strings.ReplaceAll(req.Form.Get("users"), " ", "")
	equipment := strings.ReplaceAll(req.Form.Get("equipment"), " ", "")
	inventory := strings.ReplaceAll(req.Form.Get("inventory"), " ", "")

	if len(users) != 0 {
		r.Users = strings.Split(users, ",")
	}

	if len(equipment) != 0 {
		r.Equipment = strings.Split(equipment, ",")
	}

	if len(inventory) != 0 {
		r.Inventory = strings.Split(inventory, ",")
	}

	return r.validate()
}

func (r *CreateAccessGroup) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description, validation.Required))
}

type EditAccessGroup struct {
	UID         string
	Name        string
	Description string
	Users       []string
	Equipment   []string
	Inventory   []string
}

func (r *EditAccessGroup) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("parsing uid from path: %w", err)
	}

	r.UID = uid

	if err := req.ParseForm(); err != nil {
		{
			return fmt.Errorf("parsing form: %w", err)
		}
	}

	r.Name = strings.TrimSpace(req.Form.Get("name"))
	r.Description = strings.TrimSpace(req.Form.Get("description"))
	users := strings.ReplaceAll(req.Form.Get("users"), " ", "")
	equipment := strings.ReplaceAll(req.Form.Get("equipment"), " ", "")
	inventory := strings.ReplaceAll(req.Form.Get("inventory"), " ", "")

	if len(users) != 0 {
		r.Users = strings.Split(users, ",")
	}

	if len(equipment) != 0 {
		r.Equipment = strings.Split(equipment, ",")
	}

	if len(inventory) != 0 {
		r.Inventory = strings.Split(inventory, ",")
	}

	return r.validate()
}

func (r *EditAccessGroup) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description, validation.Required))
}

type GetAccessGroup struct {
	UID string
}

func (r *GetAccessGroup) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("getting uid from path: %w", err)
	}

	r.UID = uid

	return nil
}

type DeleteAccessGroup struct {
	UID string
}

func (r *DeleteAccessGroup) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("getting uid from path: %w", err)
	}

	r.UID = uid

	return nil
}
