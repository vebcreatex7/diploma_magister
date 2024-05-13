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
	equipmentRepo repo.Equipment
	mapper        mapper.Equipment
}

func NewEquipment(equipmentRepo repo.Equipment) equipment {
	return equipment{equipmentRepo: equipmentRepo, mapper: mapper.Equipment{}}
}

func (s equipment) GetAllNotCanceled(ctx context.Context) ([]response.Equipment, error) {
	eq, err := s.equipmentRepo.GetAllNotCanceled(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting equipment: %w", err)
	}

	return s.mapper.MakeListResponse(eq), nil
}

func (s equipment) DeleteByUID(ctx context.Context, uid string) error {
	if err := s.equipmentRepo.DeleteEquipmentInAccessGroupByUID(ctx, uid); err != nil {
		return fmt.Errorf("deleting equipment_in_access_group by uid: %w", err)
	}

	if err := s.equipmentRepo.DeleteByUID(ctx, uid); err != nil {
		return fmt.Errorf("deleting equipment by uid: %w", err)
	}

	return nil
}

func (s equipment) GetByUID(ctx context.Context, uid string) (response.Equipment, error) {
	res, found, err := s.equipmentRepo.GetByUID(ctx, uid)
	if err != nil {
		return response.Equipment{}, fmt.Errorf("getting equipment by uid: %w", err)
	}

	if !found {
		return response.Equipment{}, fmt.Errorf("getting equipment by uid: %w", constant.ErrNotFound)
	}

	return s.mapper.MakeResponse(res), nil
}

func (s equipment) Edit(ctx context.Context, req request.EditEquipment) (response.Equipment, error) {
	res, edited, err := s.equipmentRepo.Edit(ctx, s.mapper.MakeEditEntity(req))
	if err != nil {
		return response.Equipment{}, fmt.Errorf("editing equipement: %w", err)
	}

	if !edited {
		return response.Equipment{}, fmt.Errorf("editing equipement: %w", constant.ErrNotFound)
	}

	return s.mapper.MakeResponse(res), nil
}

func (s equipment) Create(ctx context.Context, req request.CreateEquipment) (response.Equipment, error) {
	res, err := s.equipmentRepo.Create(ctx, s.mapper.MakeCreateEntity(req))
	if err != nil {
		return response.Equipment{}, fmt.Errorf("creating equipment: %w", err)
	}

	return s.mapper.MakeResponse(res), nil
}
