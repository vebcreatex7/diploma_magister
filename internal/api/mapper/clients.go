package mapper

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
)

type Clients struct {
}

func (m Clients) MakeCreateEntity(
	req request.CreateUser,
	passwordHash string,
) entities.Client {
	return entities.Client{
		Surname:      req.Surname,
		Name:         req.Name,
		Patronymic:   req.Password,
		Login:        req.Login,
		PasswordHash: passwordHash,
		Email:        req.Email,
		Role:         req.Role,
	}
}

func (m Clients) MakeResponse(
	e entities.Client,
) response.User {
	return response.User{
		UID:        e.UID,
		Surname:    e.Surname,
		Name:       e.Name,
		Patronymic: e.Patronymic,
		Login:      e.Login,
		Email:      e.Email,
		Role:       e.Role,
		Status:     e.Status,
	}
}

func (m Clients) MakeListResponse(
	e []entities.Client) []response.User {
	res := make([]response.User, 0, len(e))

	for i := range e {
		res = append(res, m.MakeResponse(e[i]))
	}

	return res
}

func (m Clients) MakeEditEntity(r request.EditUser) entities.Client {
	return entities.Client{
		UID:        r.UID,
		Surname:    r.Surname,
		Name:       r.Name,
		Patronymic: r.Patronymic,
		Login:      r.Login,
		Email:      r.Email,
		Role:       r.Role,
		Status:     r.Status,
	}
}
