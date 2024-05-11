package service

import (
	"context"
	"fmt"
	"github.com/vebcreatex7/diploma_magister/internal/api/mapper"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/domain/constant"
	"github.com/vebcreatex7/diploma_magister/internal/domain/repo"
)

type equipment struct {
	repo   repo.Equipment
	mapper mapper.Equipment
}

func NewEquipment(repo repo.Equipment) equipment {
	return equipment{repo: repo, mapper: mapper.Equipment{}}
}

func (s equipment) GetAllNotCanceled(ctx context.Context) ([]response.Equipment, error) {
	eq, err := s.repo.GetAllNotCanceled(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting equipment: %w", err)
	}

	return s.mapper.MakeListResponse(eq), nil
}

func (s equipment) DeleteByUID(ctx context.Context, uid string) ([]response.Equipment, error) {
	if err := s.repo.DeleteByUID(ctx, uid); err != nil {
		return nil, fmt.Errorf("deleting equipment by uid: %w", err)
	}

	eq, err := s.repo.GetAllNotCanceled(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting")
	}

	return s.mapper.MakeListResponse(eq), nil
}

func (s equipment) GetByUID(ctx context.Context, uid string) (response.Equipment, error) {
	res, found, err := s.repo.GetByUID(ctx, uid)
	if err != nil {
		return response.Equipment{}, fmt.Errorf("getting equipment by uid: %w", err)
	}

	if !found {
		return response.Equipment{}, fmt.Errorf("getting equipment by uid: %w", constant.ErrNotFound)
	}

	return s.mapper.MakeResponse(res), nil
}

func (s equipment) Edit(ctx context.Context, req request.EditEquipment) (response.Equipment, error) {
	res, edited, err := s.repo.Edit(ctx, s.mapper.MakeEditEntity(req))
	if err != nil {
		return response.Equipment{}, fmt.Errorf("editing equipement: %w", err)
	}

	if !edited {
		return response.Equipment{}, fmt.Errorf("editing equipement: %w", constant.ErrNotFound)
	}

	return s.mapper.MakeResponse(res), nil
}

func (s equipment) Create(ctx context.Context, req request.CreateEquipment) (response.Equipment, error) {
	res, err := s.repo.Create(ctx, s.mapper.MakeCreateEntity(req))
	if err != nil {
		return response.Equipment{}, fmt.Errorf("creating equipment: %w", err)
	}

	return s.mapper.MakeResponse(res), nil
}
