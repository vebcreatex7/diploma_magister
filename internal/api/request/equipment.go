package request

import (
	"fmt"
	"github.com/vebcreatex7/diploma_magister/pkg/request"
	"net/http"
)

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

type EditEquipment struct {
	UID          string
	Name         string
	Description  string
	Type         string
	Manufacturer string
	Model        string
	Room         string
	Status       string
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

	r.Name = req.Form.Get("name")
	r.Description = req.Form.Get("description")
	r.Type = req.Form.Get("type")
	r.Manufacturer = req.Form.Get("manufacturer")
	r.Model = req.Form.Get("model")
	r.Room = req.Form.Get("room")
	r.Status = req.Form.Get("status")

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

type CreateEquipment struct {
	Name         string
	Description  string
	Type         string
	Manufacturer string
	Model        string
	Room         string
	Status       string
}

func (r *CreateEquipment) Bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("parsing form: %w", err)
	}

	r.Name = req.Form.Get("name")
	r.Description = req.Form.Get("description")
	r.Type = req.Form.Get("type")
	r.Manufacturer = req.Form.Get("manufacturer")
	r.Model = req.Form.Get("model")
	r.Room = req.Form.Get("room")
	r.Status = req.Form.Get("status")

	return nil
}
