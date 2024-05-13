package service

import (
	"context"
	"fmt"
	"github.com/vebcreatex7/diploma_magister/internal/api/mapper"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
	"github.com/vebcreatex7/diploma_magister/internal/domain/repo"
)

type accessGroup struct {
	accessGroupRepo repo.AccessGroups
	clientsRepo     repo.Clients
	equipmentRepo   repo.Equipment
	inventoryRepo   repo.Inventory
	mapper          mapper.AccessGroup
}

func NewAccessGroup(
	accessGroupRepo repo.AccessGroups,
	clientsRepo repo.Clients,
	equipmentRepo repo.Equipment,
	inventoryRepo repo.Inventory,
) accessGroup {
	return accessGroup{
		accessGroupRepo: accessGroupRepo,
		clientsRepo:     clientsRepo,
		equipmentRepo:   equipmentRepo,
		inventoryRepo:   inventoryRepo,
		mapper:          mapper.AccessGroup{},
	}
}

func (s accessGroup) GetAllNotCanceled(ctx context.Context) ([]response.AccessGroup, error) {
	var res []entities.AccessGroupExt
	groups, err := s.accessGroupRepo.GetAllNotCanceled(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting access_groups: %w", err)
	}

	for _, g := range groups {
		usersLogins, err := s.clientsRepo.GetLoginsByGroupUID(ctx, g.UID)
		if err != nil {
			return nil, fmt.Errorf("getting users.login for access_group '%s': %w", g.Name, err)
		}

		equipmentNames, err := s.equipmentRepo.GetNamesByGroupUID(ctx, g.UID)
		if err != nil {
			return nil, fmt.Errorf("getting equipment.name for access_group '%s': %w", g.Name, err)
		}

		inventoryNames, err := s.inventoryRepo.GetNamesByGroupUID(ctx, g.UID)
		if err != nil {
			return nil, fmt.Errorf("getting inventory.name for access_group '%s': %w", g.Name, err)
		}

		res = append(res, entities.AccessGroupExt{
			AccessGroup: entities.AccessGroup{
				UID:         g.UID,
				Name:        g.Name,
				Description: g.Description,
				Status:      g.Status,
			},
			Users:     usersLogins,
			Equipment: equipmentNames,
			Inventory: inventoryNames,
		})
	}

	return s.mapper.MakeListResponse(res), nil
}

func (s accessGroup) Create(ctx context.Context, r request.CreateAccessGroup) (response.AccessGroup, error) {
	e := s.mapper.MakeCreateEntity(r)

	for _, login := range e.Users {
		s.clientsRepo.GetByLogin(ctx, login)

	}

	resGroup, err := s.accessGroupRepo.Create(ctx, e.AccessGroup)
	if err != nil {
		return response.AccessGroup{}, fmt.Errorf("creating access_group: %w", err)
	}

	e.UID = resGroup.UID

}
