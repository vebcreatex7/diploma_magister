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

func (s clients) Create(ctx context.Context, req request.CreateClient) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generating password hash: %w", err)
	}

	if err = s.clientsRepo.Create(ctx, s.mapper.MakeCreateEntity(req, string(hash))); err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	return nil
}

func (s clients) Login(ctx context.Context, req request.LoginClient) (string, error) {
	c, found, err := s.clientsRepo.GetByLogin(ctx, req.Login)
	if err != nil {
		return "", fmt.Errorf("getting client by login: %w", err)
	}

	if !found {
		return "", fmt.Errorf("getting client by login: %w", constant.ErrNotFound)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("comparing password: %w", err)
	}

	jwt, err := pkg.GenerateJWT(c)
	if err != nil {
		return "", fmt.Errorf("generating gwt token: %w", err)
	}

	return jwt, nil
}

func (s clients) GetAllNotCanceled(ctx context.Context) ([]response.Client, error) {
	e, err := s.clientsRepo.GetAllNotCanceled(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting clients: %w", err)
	}

	return s.mapper.MakeListResponse(e), nil
}

func (s clients) DeleteByUID(ctx context.Context, uid string) ([]response.Client, error) {
	if err := s.clientsRepo.DeleteByUID(ctx, uid); err != nil {
		return nil, fmt.Errorf("deleting client by uid: %w", err)
	}

	e, err := s.clientsRepo.GetAllNotCanceled(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting clients: %w", err)
	}

	return s.mapper.MakeListResponse(e), nil
}
