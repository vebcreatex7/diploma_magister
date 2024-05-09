package service

import (
	"context"
	"fmt"
	"github.com/vebcreatex7/diploma_magister/internal/api/mapper"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
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
