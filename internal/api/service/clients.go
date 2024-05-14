package service

import (
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/domain/constant"
	"github.com/vebcreatex7/diploma_magister/pkg"
	"golang.org/x/crypto/bcrypt"

	"context"
	"fmt"
	"github.com/vebcreatex7/diploma_magister/internal/api/mapper"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/domain/repo"
)

type clients struct {
	clientsRepo repo.Clients
	mapper      mapper.Clients
}

func NewClients(clientsRepo repo.Clients) clients {
	return clients{
		clientsRepo: clientsRepo,
		mapper:      mapper.Clients{},
	}
}

func (s clients) Create(ctx context.Context, req request.CreateUser) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generating password hash: %w", err)
	}

	if err = s.clientsRepo.Create(ctx, s.mapper.MakeCreateEntity(req, string(hash))); err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	return nil
}

func (s clients) Login(ctx context.Context, req request.LoginUser) (string, error) {
	c, found, err := s.clientsRepo.GetByLogin(ctx, req.Login)
	if err != nil {
		return "", fmt.Errorf("getting client by user: %w", err)
	}

	if !found {
		return "", fmt.Errorf("getting client by user: %w", constant.ErrNotFound)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("comparing password: %w", err)
	}

	if !c.Approved {
		return "", fmt.Errorf("user is waiting for approve")
	}

	jwt, err := pkg.GenerateJWT(c)
	if err != nil {
		return "", fmt.Errorf("generating gwt token: %w", err)
	}

	return jwt, nil
}

func (s clients) GetAll(ctx context.Context) ([]response.User, error) {
	e, err := s.clientsRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting users: %w", err)
	}

	return s.mapper.MakeListResponse(e), nil
}

func (s clients) DeleteByUID(ctx context.Context, uid string) error {
	if err := s.clientsRepo.DeleteClientsInAccessGroupByUID(ctx, uid); err != nil {
		return fmt.Errorf("deleting clients_in_access_group by uid: %w", err)
	}

	if err := s.clientsRepo.DeleteByUID(ctx, uid); err != nil {
		return fmt.Errorf("deleting client by uid: %w", err)
	}

	return nil
}

func (s clients) GetByUID(ctx context.Context, uid string) (response.User, error) {
	res, found, err := s.clientsRepo.GetByUID(ctx, uid)
	if err != nil {
		return response.User{}, fmt.Errorf("getting user by uid: %w", err)
	}

	if !found {
		return response.User{}, fmt.Errorf("getting user by uid: %w", constant.ErrNotFound)
	}

	return s.mapper.MakeResponse(res), nil
}

func (s clients) Edit(ctx context.Context, req request.EditUser) (response.User, error) {
	res, edited, err := s.clientsRepo.Edit(ctx, s.mapper.MakeEditEntity(req))
	if err != nil {
		return response.User{}, fmt.Errorf("editing user: %w", err)
	}

	if !edited {
		return response.User{}, fmt.Errorf("editing user: %w", constant.ErrNotFound)
	}

	return s.mapper.MakeResponse(res), nil
}

func (s clients) Approve(ctx context.Context, uid string) ([]response.User, error) {
	if err := s.clientsRepo.Approve(ctx, uid); err != nil {
		return nil, fmt.Errorf("approving user by uid: %w", err)
	}

	res, err := s.clientsRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting users: %w", err)
	}

	return s.mapper.MakeListResponse(res), nil
}
