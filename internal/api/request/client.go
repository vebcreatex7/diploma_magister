package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/vebcreatex7/diploma_magister/internal/domain/constant"
	"log"
	"net/http"
)

type CreateClient struct {
	Surname    string
	Name       string
	Patronymic string
	Email      string
	Login      string
	Password   string
	Role       string
}

func (r *CreateClient) Bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("parsing request: %w", err)
	}

	log.Println(req.Form)

	r.Surname = req.Form.Get("surname")
	r.Name = req.Form.Get("name")
	r.Patronymic = req.Form.Get("patronymic")
	r.Login = req.Form.Get("login")
	r.Password = req.Form.Get("password")
	r.Email = req.Form.Get("email")
	r.Role = req.Form.Get("role")

	if err := r.validate(); err != nil {
		return fmt.Errorf("validating request: %w", err)
	}

	return nil
}

func (r *CreateClient) validate() error {
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

type LoginClient struct {
	Login    string
	Password string
}

func (r *LoginClient) Bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("parsing request: %w", err)
	}

	r.Login = req.Form.Get("login")
	r.Password = req.Form.Get("password")

	if err := r.validate(); err != nil {
		return fmt.Errorf("validating request: %w", err)
	}

	return nil
}

func (r *LoginClient) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Login, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}
