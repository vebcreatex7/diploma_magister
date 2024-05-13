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

type inventory struct {
	inventoryRepo repo.Inventory
	mapper        mapper.Inventory
}

func NewInventory(inventoryRepo repo.Inventory) inventory {
	return inventory{inventoryRepo: inventoryRepo, mapper: mapper.Inventory{}}
}

func (s inventory) GetAllNotCanceled(ctx context.Context) ([]response.Inventory, error) {
	eq, err := s.inventoryRepo.GetAllNotCanceled(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting inventory: %w", err)
	}

	return s.mapper.MakeListResponse(eq), nil
}

func (s inventory) DeleteByUID(ctx context.Context, uid string) error {
	if err := s.inventoryRepo.DeleteInventoryInAccessGroupByUID(ctx, uid); err != nil {
		return fmt.Errorf("deleting inventory_in_access_group by uid: %w", err)
	}

	if err := s.inventoryRepo.DeleteByUID(ctx, uid); err != nil {
		return fmt.Errorf("deleting inventory by uid: %w", err)
	}

	return nil
}

func (s inventory) GetByUID(ctx context.Context, uid string) (response.Inventory, error) {
	res, found, err := s.inventoryRepo.GetByUID(ctx, uid)
	if err != nil {
		return response.Inventory{}, fmt.Errorf("getting inventory by uid: %w", err)
	}

	if !found {
		return response.Inventory{}, fmt.Errorf("getting inventory by uid: %w", constant.ErrNotFound)
	}

	return s.mapper.MakeResponse(res), nil
}

func (s inventory) Edit(ctx context.Context, req request.EditInventory) (response.Inventory, error) {
	res, edited, err := s.inventoryRepo.Edit(ctx, s.mapper.MakeEditEntity(req))
	if err != nil {
		return response.Inventory{}, fmt.Errorf("editing inventory: %w", err)
	}

	if !edited {
		return response.Inventory{}, fmt.Errorf("editing inventory: %w", constant.ErrNotFound)
	}

	return s.mapper.MakeResponse(res), nil
}

func (s inventory) Create(ctx context.Context, req request.CreateInventory) (response.Inventory, error) {
	res, err := s.inventoryRepo.Create(ctx, s.mapper.MakeCreateEntity(req))
	if err != nil {
		return response.Inventory{}, fmt.Errorf("creating inventory: %w", err)
	}

	return s.mapper.MakeResponse(res), nil
}
