package service

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/vebcreatex7/diploma_magister/internal/api/mapper"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/domain/constant"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
	"github.com/vebcreatex7/diploma_magister/internal/domain/repo"
	"github.com/vebcreatex7/diploma_magister/internal/domain/service"
	"github.com/vebcreatex7/diploma_magister/internal/repo/postgres/schema"
	"sort"
	"strings"
	"time"
)

type equipment struct {
	equipmentRepo      repo.Equipment
	clientsRepo        repo.Clients
	accessGroupService service.AccessGroup
	mapper             mapper.Equipment
	db                 *goqu.Database
}

func NewEquipment(
	equipmentRepo repo.Equipment,
	clientsRepo repo.Clients,
	accessGroupService service.AccessGroup,
) equipment {
	return equipment{
		equipmentRepo:      equipmentRepo,
		clientsRepo:        clientsRepo,
		accessGroupService: accessGroupService,
		mapper:             mapper.Equipment{}}
}

func (s equipment) GetAll(ctx context.Context) ([]response.Equipment, error) {
	eq, err := s.equipmentRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting equipment: %w", err)
	}

	return s.mapper.MakeListResponse(eq), nil
}

func (s equipment) GetAllForUser(ctx context.Context, uid string) ([]response.Equipment, error) {
	all, err := s.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all equipment: %w", err)
	}

	ags, err := s.accessGroupService.GetAllForGivenUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("getting acess_groups for user: %w", err)
	}

	var equipment []string

	for i := range ags {
		equipment = append(equipment, strings.Split(ags[i].Equipment, ",")...)
	}

	for i := 0; i < len(all); i++ {
		equipmentFound := false

		for j := range equipment {
			if all[i].Name == equipment[j] {
				equipmentFound = true
				break
			}
		}

		if !equipmentFound {
			all = append(all[:i], all[i+1:]...)
			i--
		}
	}

	return all, nil
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

func (s equipment) GetEquipmentScheduleInRange(ctx context.Context, req request.GetEquipmentSchedule) ([]response.EquipmentSchedule, error) {
	es, err := s.equipmentRepo.SelectScheduleByName(ctx, req.Name, req.Lower, req.Upper)
	if err != nil {
		return nil, fmt.Errorf("getting equipment_schedule: %w", err)
	}

	sort.SliceStable(es, func(i, j int) bool {
		return es[i].TimeInterval.Lower.Time.Before(es[j].TimeInterval.Lower.Time)
	})

	var resp []response.EquipmentSchedule

	for d := req.Lower; d.Before(req.Upper); d = d.Add(time.Hour * 24) {
		var intervals string

		for i := 0; i < len(es); i++ {
			if es[i].TimeInterval.Lower.Time.After(d) &&
				es[i].TimeInterval.Lower.Time.Before(d.Add(time.Hour*24)) {

				intervals += fmt.Sprintf("[%s, %s]",
					es[i].TimeInterval.Lower.Time.Format(time.TimeOnly),
					es[i].TimeInterval.Upper.Time.Format(time.TimeOnly),
				)
			}
		}

		resp = append(resp, response.EquipmentSchedule{
			Date:      d.Format(time.DateOnly),
			Intervals: intervals,
		})
	}

	return resp, nil
}

func (s equipment) GetEquipmentScheduleInRangeForUser(ctx context.Context, req request.GetEquipmentSchedule, userUID string) ([]response.EquipmentSchedule, error) {
	avail, err := s.clientsRepo.IsEquipmentAvailable(ctx, userUID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("checking availble eq: %w", err)
	}

	if !avail {
		return nil, fmt.Errorf("checking availble eq: '%s' not available", req.Name)
	}

	es, err := s.equipmentRepo.SelectScheduleByName(ctx, req.Name, req.Lower, req.Upper)
	if err != nil {
		return nil, fmt.Errorf("getting equipment_schedule: %w", err)
	}

	sort.SliceStable(es, func(i, j int) bool {
		return es[i].TimeInterval.Lower.Time.Before(es[j].TimeInterval.Lower.Time)
	})

	var resp []response.EquipmentSchedule

	for d := req.Lower; d.Before(req.Upper); d = d.Add(time.Hour * 24) {
		var intervals string

		for i := 0; i < len(es); i++ {
			if es[i].TimeInterval.Lower.Time.After(d) &&
				es[i].TimeInterval.Lower.Time.Before(d.Add(time.Hour*24)) {

				intervals += fmt.Sprintf("[%s, %s]",
					es[i].TimeInterval.Lower.Time.Format(time.TimeOnly),
					es[i].TimeInterval.Upper.Time.Format(time.TimeOnly),
				)
			}
		}

		resp = append(resp, response.EquipmentSchedule{
			Date:      d.Format(time.DateOnly),
			Intervals: intervals,
		})
	}

	return resp, nil
}

func (s equipment) getByName(ctx context.Context, tx *goqu.TxDatabase, name string) (entities.Equipment, bool, error) {
	var eq entities.Equipment

	found, err := tx.From(schema.Equipment).
		Select(entities.Equipment{}).
		Where(goqu.I("name").Eq(name)).
		Prepared(true).Executor().
		ScanStructContext(ctx, &eq)

	return eq, found, err
}
