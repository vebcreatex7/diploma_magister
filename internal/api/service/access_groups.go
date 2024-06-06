package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/vebcreatex7/diploma_magister/internal/api/mapper"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/domain/constant"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
	"github.com/vebcreatex7/diploma_magister/internal/domain/repo"
	"strings"
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

func (s accessGroup) GetAllForGivenUser(ctx context.Context, userUID string) ([]response.AccessGroup, error) {
	all, err := s.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting access_groups: %w", err)
	}

	u, found, err := s.clientsRepo.GetByUID(ctx, userUID)
	if err != nil {
		return nil, fmt.Errorf("getting user by uid: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("getting user by uid: %w", constant.ErrNotFound)
	}

	for i := 0; i < len(all); i++ {
		userFound := false

		for _, login := range strings.Split(all[i].Users, ",") {
			if u.Login == login {
				userFound = true

				break
			}
		}

		if !userFound {
			all = append(all[:i], all[i+1:]...)
			i--
		}
	}

	return all, nil
}

func (s accessGroup) GetAll(ctx context.Context) ([]response.AccessGroup, error) {
	var res []entities.AccessGroupExt
	groups, err := s.accessGroupRepo.GetAll(ctx)
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
			},
			Users:     usersLogins,
			Equipment: equipmentNames,
			Inventory: inventoryNames,
		})
	}

	return s.mapper.MakeListResponse(res), nil
}

func (s accessGroup) GetByUID(ctx context.Context, uid string) (response.AccessGroup, error) {
	g, found, err := s.accessGroupRepo.GetByUID(ctx, uid)
	if err != nil {
		return response.AccessGroup{}, fmt.Errorf("getting access_group by uid: %w", err)
	}

	if !found {
		return response.AccessGroup{}, fmt.Errorf("getting access_group by uid: %w", constant.ErrNotFound)
	}

	usersLogins, err := s.clientsRepo.GetLoginsByGroupUID(ctx, g.UID)
	if err != nil {
		return response.AccessGroup{}, fmt.Errorf("getting users.login for access_group '%s': %w", g.Name, err)
	}

	equipmentNames, err := s.equipmentRepo.GetNamesByGroupUID(ctx, g.UID)
	if err != nil {
		return response.AccessGroup{}, fmt.Errorf("getting equipment.name for access_group '%s': %w", g.Name, err)
	}

	inventoryNames, err := s.inventoryRepo.GetNamesByGroupUID(ctx, g.UID)
	if err != nil {
		return response.AccessGroup{}, fmt.Errorf("getting inventory.name for access_group '%s': %w", g.Name, err)
	}

	return s.mapper.MakeResponse(entities.AccessGroupExt{
		AccessGroup: g,
		Users:       usersLogins,
		Equipment:   equipmentNames,
		Inventory:   inventoryNames,
	}), nil
}

func (s accessGroup) Create(ctx context.Context, r request.CreateAccessGroup) (response.AccessGroup, error) {
	group := s.mapper.MakeCreateEntity(r)
	group.UID = uuid.New().String()

	group.RemoveDuplicates()

	clientsInAG, equipmentInAG, inventoryInAG, err := s.makeRelatedEntities(ctx, group)
	if err != nil {
		return response.AccessGroup{}, err
	}

	resGroup, err := s.accessGroupRepo.Create(ctx, group.AccessGroup)
	if err != nil {
		return response.AccessGroup{}, fmt.Errorf("creating access_group: %w", err)
	}

	err = s.createdRelatedEntities(ctx, clientsInAG, equipmentInAG, inventoryInAG)
	if err != nil {
		return response.AccessGroup{}, err
	}

	return s.GetByUID(ctx, resGroup.UID)
}

func (s accessGroup) Edit(ctx context.Context, r request.EditAccessGroup) (response.AccessGroup, error) {
	group := s.mapper.MakeEditEntity(r)

	group.RemoveDuplicates()

	clientsInAG, equipmentInAG, inventoryInAG, err := s.makeRelatedEntities(ctx, group)
	if err != nil {
		return response.AccessGroup{}, err
	}

	resGroup, edited, err := s.accessGroupRepo.Edit(ctx, group.AccessGroup)
	if err != nil {
		return response.AccessGroup{}, fmt.Errorf("editing access_group '%s': %w", group.Name, err)
	}

	if !edited {
		return response.AccessGroup{}, fmt.Errorf("editing access_group '%s': %w", group.Name, constant.ErrNotFound)
	}

	if err = s.deleteRelatedEntities(ctx, group.UID); err != nil {
		return response.AccessGroup{}, err
	}

	if err = s.createdRelatedEntities(ctx, clientsInAG, equipmentInAG, inventoryInAG); err != nil {
		return response.AccessGroup{}, err
	}

	return s.GetByUID(ctx, resGroup.UID)
}

func (s accessGroup) DeleteByUID(ctx context.Context, uid string) error {
	if err := s.deleteRelatedEntities(ctx, uid); err != nil {
		return err
	}

	if err := s.accessGroupRepo.DeleteByUID(ctx, uid); err != nil {
		return fmt.Errorf("deleting access_group by uid: %w", err)
	}

	return nil
}

func (s accessGroup) makeRelatedEntities(ctx context.Context, group entities.AccessGroupExt) ([]entities.ClientsInAccessGroup, []entities.EquipmentInAccessGroup, []entities.InventoryInAccessGroup, error) {
	var clientsInAG []entities.ClientsInAccessGroup
	var equipmentInAG []entities.EquipmentInAccessGroup
	var inventoryInAG []entities.InventoryInAccessGroup

	for _, login := range group.Users {
		u, f, err := s.clientsRepo.GetByLogin(ctx, login)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("getting user by login '%s': %w", login, err)
		}

		if !f {
			return nil, nil, nil, fmt.Errorf("getting user by login '%s': %w", login, constant.ErrNotFound)
		}

		clientsInAG = append(clientsInAG, entities.ClientsInAccessGroup{
			AccessGroupUID: group.UID,
			ClientUID:      u.UID,
		})
	}

	for _, name := range group.Equipment {
		e, f, err := s.equipmentRepo.GetByName(ctx, name)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("getting equipment by name '%s': %w", name, err)
		}

		if !f {
			return nil, nil, nil, fmt.Errorf("getting equipment by name '%s': %w", name, constant.ErrNotFound)
		}

		equipmentInAG = append(equipmentInAG, entities.EquipmentInAccessGroup{
			AccessGroupUID: group.UID,
			EquipmentUID:   e.UID,
		})
	}

	for _, name := range group.Inventory {
		i, f, err := s.inventoryRepo.GetByName(ctx, name)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("getting inventory by name '%s': %w", name, err)
		}

		if !f {
			return nil, nil, nil, fmt.Errorf("getting inventory by name '%s': %w", name, constant.ErrNotFound)
		}

		inventoryInAG = append(inventoryInAG, entities.InventoryInAccessGroup{
			AccessGroupUID: group.UID,
			InventoryUID:   i.UID,
		})
	}
	return clientsInAG, equipmentInAG, inventoryInAG, nil
}

func (s accessGroup) createdRelatedEntities(
	ctx context.Context,
	clientsInAG []entities.ClientsInAccessGroup,
	equipmentInAG []entities.EquipmentInAccessGroup,
	inventoryInAG []entities.InventoryInAccessGroup,
) error {
	if len(clientsInAG) != 0 {
		if err := s.accessGroupRepo.CreateClientsInAccessGroup(ctx, clientsInAG); err != nil {
			return fmt.Errorf("creating clients_in_access_group: %w", err)
		}
	}

	if len(equipmentInAG) != 0 {
		if err := s.accessGroupRepo.CreateEquipmentInAccessGroup(ctx, equipmentInAG); err != nil {
			return fmt.Errorf("creating equipment_in_access_group: %w", err)
		}
	}

	if len(inventoryInAG) != 0 {
		if err := s.accessGroupRepo.CreateInventoryInAccessGroup(ctx, inventoryInAG); err != nil {
			return fmt.Errorf("creating inventory_in_access_group: %w", err)
		}
	}
	return nil
}

func (s accessGroup) deleteRelatedEntities(ctx context.Context, uid string) error {
	if err := s.accessGroupRepo.DeleteClientsInAccessGroupByUID(ctx, uid); err != nil {
		return fmt.Errorf("deleting releted clients_in_access_groups: %w", err)
	}

	if err := s.accessGroupRepo.DeleteEquipmentInAccessGroupByUID(ctx, uid); err != nil {
		return fmt.Errorf("deleting releted equipment_in_access_groups: %w", err)
	}

	if err := s.accessGroupRepo.DeleteInventoryInAccessGroupByUID(ctx, uid); err != nil {
		return fmt.Errorf("deleting releted inventory_in_access_groups: %w", err)
	}

	return nil
}
