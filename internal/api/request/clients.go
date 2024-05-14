package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/vebcreatex7/diploma_magister/internal/domain/constant"
	"github.com/vebcreatex7/diploma_magister/pkg/request"
	"log"
	"net/http"
	"strings"
)

type CreateUser struct {
	Surname    string
	Name       string
	Patronymic string
	Email      string
	Login      string
	Password   string
	Role       string
}

func (r *CreateUser) Bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("parsing request: %w", err)
	}

	log.Println(req.Form)

	r.Surname = strings.TrimSpace(req.Form.Get("surname"))
	r.Name = strings.TrimSpace(req.Form.Get("name"))
	r.Patronymic = strings.TrimSpace(req.Form.Get("patronymic"))
	r.Login = strings.TrimSpace(req.Form.Get("login"))
	r.Password = strings.TrimSpace(req.Form.Get("password"))
	r.Email = strings.TrimSpace(req.Form.Get("email"))
	r.Role = strings.TrimSpace(req.Form.Get("role"))

	if err := r.validate(); err != nil {
		return fmt.Errorf("validating request: %w", err)
	}

	return nil
}

func (r *CreateUser) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Surname, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Patronymic, validation.Required),
		validation.Field(&r.Login, validation.Required),
		validation.Field(&r.Password, validation.Required),
		validation.Field(&r.Email, validation.Required),
		validation.Field(&r.Role, validation.Required, validation.In(
			constant.EngineerRole,
			constant.ScientistRole,
			constant.LaboratorianRole,
		)),
	)
}

type LoginUser struct {
	Login    string
	Password string
}

func (r *LoginUser) Bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("parsing request: %w", err)
	}

	r.Login = strings.TrimSpace(req.Form.Get("login"))
	r.Password = strings.TrimSpace(req.Form.Get("password"))

	if err := r.validate(); err != nil {
		return fmt.Errorf("validating request: %w", err)
	}

	return nil
}

func (r *LoginUser) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Login, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}

type DeleteUser struct {
	UID string
}

func (r *DeleteUser) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("parsing uid from path: %w", err)
	}

	r.UID = uid

	return nil
}

type EditUser struct {
	UID        string
	Surname    string
	Name       string
	Patronymic string
	Login      string
	Email      string
	Role       string
}

func (r *EditUser) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("parsing uid from path: %w", err)
	}

	r.UID = uid

	if err = req.ParseForm(); err != nil {
		return fmt.Errorf("parsing form: %w", err)
	}

	log.Println(req.Form)

	r.Surname = strings.TrimSpace(req.Form.Get("surname"))
	r.Name = strings.TrimSpace(req.Form.Get("name"))
	r.Patronymic = strings.TrimSpace(req.Form.Get("patronymic"))
	r.Login = strings.TrimSpace(req.Form.Get("login"))
	r.Email = strings.TrimSpace(req.Form.Get("email"))
	r.Role = strings.TrimSpace(req.Form.Get("role"))

	return r.validate()
}

func (r *EditUser) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Surname, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Patronymic, validation.Required),
		validation.Field(&r.Login, validation.Required),
		validation.Field(&r.Email, validation.Required),
		validation.Field(&r.Role, validation.Required, validation.In(
			constant.EngineerRole,
			constant.ScientistRole,
			constant.LaboratorianRole,
		)),
	)
}

type GetUser struct {
	UID string
}

func (r *GetUser) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("parsing uid from path: %w", err)
	}

	r.UID = uid

	return nil
}

type ApproveUser struct {
	UID string
}

func (r *ApproveUser) Bind(req *http.Request) error {
	uid, err := request.ParseUIDFromPath(req, true)
	if err != nil {
		return fmt.Errorf("parsing uid from path: %w", err)
	}

	r.UID = uid

	return nil
}
