package service

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/domain/constant"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
	"github.com/vebcreatex7/diploma_magister/internal/domain/repo"
	"github.com/vebcreatex7/diploma_magister/internal/repo/postgres/schema"
	"sort"
	"time"
)

type experiment struct {
	db               *goqu.Database
	clientsRepo      repo.Clients
	equipmentService equipment
}

func NewExperiment(
	db *goqu.Database,
	clientsRepo repo.Clients,
	equipmentService equipment,
) experiment {
	return experiment{
		db:               db,
		clientsRepo:      clientsRepo,
		equipmentService: equipmentService,
	}
}

func (s experiment) GetSuggestionsForUser(ctx context.Context, userUID string) (response.ExperimentSuggestions, error) {
	var eq []string

	if err := s.db.ScanValsContext(
		ctx,
		&eq,
		`select eq.name from equipment eq
join equipment_in_access_group eqag on eq.uid = eqag.equipment_uid
join clients_in_access_group cag on eqag.access_group_uid = cag.access_group_uid
join client c on cag.client_uid = c.uid
where c.uid = $1`,
		userUID,
	); err != nil {
		return response.ExperimentSuggestions{}, fmt.Errorf("getting equipment suggestions: %w", err)
	}

	var in []string
	if err := s.db.ScanValsContext(
		ctx,
		&in,
		`select i.name from inventory i
                        join inventory_in_access_group iag on i.uid = iag.inventory_uid
                        join clients_in_access_group cag on iag.access_group_uid = cag.access_group_uid
                        join client c on cag.client_uid = c.uid
where c.uid = $1`,
		userUID,
	); err != nil {
		return response.ExperimentSuggestions{}, fmt.Errorf("getting inventory suggestions: %w", err)
	}

	return response.ExperimentSuggestions{
		Equipment: eq,
		Inventory: in,
	}, nil
}

func (s experiment) AddExperiment(ctx context.Context, req request.AddExperiment, userUID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	exp := entities.Experiment{
		UID:         uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		StartTs:     req.StartTs,
		EndTs:       req.EndTs,
	}

	var es []entities.EquipmentSchedule
	var esinex []entities.EquipmentScheduleInExperiment
	var ininex []entities.InventoryInExperiment

	for i := range req.Equipment {
		av, err := s.isEquipmentAvailableForUser(ctx, tx, req.Equipment[i].Name, userUID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("checking available for %d equipment: %w", i, err)

		}

		if !av {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("checking available for %d equipment: not available", i)
		}

		eq, f, err := s.equipmentService.getByName(ctx, tx, req.Equipment[i].Name)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("getting uid for %d equipment: %w", i, err)
		}
		if !f {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("getting uid for %d equipment: %w", i, constant.ErrNotFound)
		}

		tmp := entities.EquipmentSchedule{
			UID:          uuid.New().String(),
			EquipmentUID: eq.UID,
			TimeInterval: pgtype.Tsrange{
				Lower:     pgtype.Timestamp{},
				Upper:     pgtype.Timestamp{},
				LowerType: pgtype.Exclusive,
				UpperType: pgtype.Exclusive,
				Status:    pgtype.Present,
			},
			MaintainceFlag: false,
		}

		if err := tmp.TimeInterval.Upper.Set(req.Equipment[i].Upper); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("setting upper for %d equipment: %w", i, err)
		}

		if err := tmp.TimeInterval.Lower.Set(req.Equipment[i].Lower); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("setting lower for %d equipment: %w", i, err)
		}

		es = append(es, tmp)
		esinex = append(esinex, entities.EquipmentScheduleInExperiment{
			ExperimentUID:        exp.UID,
			EquipmentScheduleUID: tmp.UID,
		})

	}

	for i := range req.Inventory {
		av, err := s.isInventoryAvailableForUser(ctx, tx, req.Inventory[i].Name, userUID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("checking available for %d inventory: %w", i, err)

		}

		if !av {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("checking available for %d inventory: not available", i)
		}

		invUID, f, err := s.getInventoryUIDByName(ctx, tx, req.Inventory[i].Name)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("getting uid for %d inventory: %w", i, err)
		}
		if !f {
			return fmt.Errorf("getting uid for %d inventory: %w", i, constant.ErrNotFound)
		}

		ininex = append(ininex, entities.InventoryInExperiment{
			ExperimentUID: exp.UID,
			InventoryUID:  invUID,
			Quantity:      req.Inventory[i].Quantity,
		})
	}

	if _, err := tx.Insert(schema.Experiment).
		Rows(exp).
		Prepared(true).Executor().
		ExecContext(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("inserting experiment: %w", err)
	}

	if _, err := tx.Insert(schema.EquipmentSchedule).
		Rows(es).
		Prepared(true).Executor().
		ExecContext(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("inserting equipment_schedule: %w", err)
	}

	if _, err := tx.Insert(schema.EquipmentScheduleInExperiment).
		Rows(esinex).
		Prepared(true).Executor().
		ExecContext(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("inserting equipment_schedule_in_experiment: %w", err)
	}

	if _, err := tx.Insert(schema.ClientsInExperiment).
		Rows(entities.ClientsInExperiment{
			ExperimentUID: exp.UID,
			ClientUID:     userUID,
		}).
		Prepared(true).Executor().
		ExecContext(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("inserting clients_in_experiment: %w", err)
	}

	if _, err := tx.Insert(schema.InventoryInExperiment).
		Rows(ininex).
		Prepared(true).Executor().
		ExecContext(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("inserting inventory_in_experiment: %w", err)
	}

	for i := range ininex {
		if _, err := tx.ExecContext(ctx, `update inventory set quantity = quantity - $1 where uid = $2`, ininex[i].Quantity, ininex[i].InventoryUID); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("updating inventory quantity: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s experiment) GetAll(ctx context.Context) ([]response.Experiment, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	var res []response.Experiment

	var exps []entities.Experiment

	if err := tx.From(schema.Experiment).
		Select(entities.Experiment{}).
		Where(goqu.I("finished").IsNotTrue()).
		Order(goqu.C("start_ts").Asc()).
		Prepared(true).
		Executor().ScanStructsContext(ctx, &exps); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("getting experiments: %w", err)
	}

	for _, exp := range exps {
		var users []string

		if err := tx.ScanValsContext(
			ctx,
			&users,
			`select login from client c 
join clients_in_experiment cinex on c.uid = cinex.client_uid
where cinex.experiment_uid = $1`, exp.UID); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("getting users in exp: %w", err)
		}

		type tmpInv struct {
			Name     string          `db:"name"`
			Quantity decimal.Decimal `db:"quantity"`
		}

		var ins []tmpInv
		if err := tx.ScanStructsContext(
			ctx,
			&ins,
			`select i.name, iinex.quantity from inventory i
join inventory_in_experiment iinex on i.uid = iinex.inventory_uid
where iinex.experiment_uid = $1`, exp.UID); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("getting inventory in exp: %w", err)
		}

		type tmpES struct {
			Name     string         `db:"name"`
			Interval pgtype.Tsrange `db:"time_interval"`
		}
		var ess []tmpES
		if err := tx.ScanStructsContext(
			ctx,
			&ess,
			`select eq.name, es.time_interval from equipment eq
join equipment_schedule es on eq.uid = es.equipment_uid
join equipment_schedule_in_experiment esinex on es.uid = esinex.equipment_schedule_uid
where esinex.experiment_uid = $1`, exp.UID); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("getting equipment in exp: %w", err)
		}

		eqNameToIntervals := make(map[string][]pgtype.Tsrange)

		for _, es := range ess {
			eqNameToIntervals[es.Name] = append(eqNameToIntervals[es.Name], es.Interval)
		}

		for _, intervals := range eqNameToIntervals {
			sort.Slice(intervals, func(i, j int) bool {
				return intervals[i].Lower.Time.Before(intervals[j].Lower.Time)
			})
		}

		var eqRes []response.EquipmentInExperiment

		for n, intervals := range eqNameToIntervals {
			var intervalsStr []string

			for i := range intervals {
				intervalsStr = append(intervalsStr,
					fmt.Sprintf(
						"(%s, %s)",
						intervals[i].Lower.Time.Format(constant.Layout),
						intervals[i].Upper.Time.Format(constant.Layout),
					),
				)
			}
			eqRes = append(eqRes, response.EquipmentInExperiment{
				Name:      n,
				Intervals: intervalsStr,
			})
		}

		sort.Slice(eqRes, func(i, j int) bool {
			return eqRes[i].Name < eqRes[j].Name
		})

		var inRes []response.InventoryInExperiment

		for i := range ins {
			inRes = append(inRes, response.InventoryInExperiment{
				Name:     ins[i].Name,
				Quantity: ins[i].Quantity,
			})
		}

		sort.Slice(inRes, func(i, j int) bool {
			return inRes[i].Name < inRes[j].Name
		})

		res = append(res, response.Experiment{
			UID:         exp.UID,
			Name:        exp.Name,
			Description: exp.Description,
			StartTs:     exp.StartTs.Format(constant.Layout),
			EndTs:       exp.EndTs.Format(constant.Layout),
			Users:       users,
			Equipment:   eqRes,
			Inventory:   inRes,
		})

	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s experiment) GetAllForUser(ctx context.Context, userUID string) ([]response.Experiment, error) {
	all, err := s.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting experiments: %w", err)
	}

	u, f, err := s.clientsRepo.GetByUID(ctx, userUID)
	if err != nil {
		return nil, fmt.Errorf("get user by uid: %w", err)
	}
	if !f {
		return nil, fmt.Errorf("get user by uid: %w", constant.ErrNotFound)
	}

	for i := 0; i < len(all); i++ {
		var userFound = false

		for _, l := range all[i].Users {
			if l == u.Login {
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

func (s experiment) DeleteByUID(ctx context.Context, uid string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	var curEsinex []entities.EquipmentScheduleInExperiment
	if err = tx.ScanStructsContext(
		ctx,
		&curEsinex,
		`select * from equipment_schedule_in_experiment
where experiment_uid = $1`,
		uid,
	); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("getting equipment_schedule_in_experiment: %w", err)
	}

	var curIninex []entities.InventoryInExperiment
	if err = tx.ScanStructsContext(
		ctx,
		&curIninex,
		`select * from inventory_in_experiment
where experiment_uid = $1`,
		uid,
	); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("getting inventory_in_experiment: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`delete from equipment_schedule_in_experiment
where experiment_uid = $1`,
		uid,
	); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("deleting equipment_schedule_in_experiment: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`delete from inventory_in_experiment
where experiment_uid = $1`,
		uid,
	); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("deleting inventory_in_experiment: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`delete from clients_in_experiment
where experiment_uid = $1`,
		uid,
	); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("deleting clients_in_experiment: %w", err)
	}

	for i := range curEsinex {
		if _, err := tx.ExecContext(
			ctx,
			`delete from equipment_schedule
where uid = $1`,
			curEsinex[i].EquipmentScheduleUID,
		); err != nil {
			if err = tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("deleting equipment_schedule: %w", err)
		}
	}

	for i := range curIninex {
		if _, err = tx.ExecContext(
			ctx,
			`update inventory set quantity = quantity + $1
where uid = $2`,
			curIninex[i].Quantity,
			curIninex[i].InventoryUID,
		); err != nil {
			if err = tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("updating inventory quantity: %w", err)
		}
	}

	if _, err = tx.ExecContext(
		ctx,
		`delete from experiment
where uid = $1`,
		uid,
	); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("deleting experiment: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s experiment) DeleteByUIDForUser(ctx context.Context, uid, userUID string) error {
	var tmp string

	av, err := s.db.ScanValContext(
		ctx,
		&tmp,
		`select client_uid from clients_in_experiment
where experiment_uid = $1 and client_uid = $2`,
		uid, userUID,
	)
	if err != nil {
		return fmt.Errorf("checking user: %w", err)
	}

	if !av {
		return fmt.Errorf("checking user: not available")
	}

	return s.DeleteByUID(ctx, uid)
}

func (s experiment) GetAllFinishedNotMarked(ctx context.Context) ([]response.Experiment, error) {
	all, err := s.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all exp: %w", err)
	}

	now := time.Now()

	for i := 0; i < len(all); i++ {
		endTs, err := time.Parse(constant.Layout, all[i].EndTs)
		if err != nil {
			return nil, fmt.Errorf("parsing end_ts: %w", err)
		}

		if endTs.After(now) {
			all = append(all[:i], all[i+1:]...)
			i--
		}
	}

	return all, nil
}

func (s experiment) GetByUID(ctx context.Context, uid string) (response.Experiment, error) {
	all, err := s.GetAll(ctx)
	if err != nil {
		return response.Experiment{}, fmt.Errorf("getting exp: %w", err)
	}

	for i := range all {
		if all[i].UID == uid {
			return all[i], nil
		}
	}

	return response.Experiment{}, constant.ErrNotFound
}

func (s experiment) Finish(ctx context.Context, req request.FinishExperiment) error {
	exp, err := s.GetByUID(ctx, req.UID)
	if err != nil {
		return fmt.Errorf("getting exp: %w", err)
	}

	for i := range req.InventoryName {
		if _, err := s.db.ExecContext(
			ctx,
			`update inventory set quantity = quantity + $1 where
                                                  name = $2`,
			req.InventoryQuantity[i], req.InventoryName[i],
		); err != nil {
			return fmt.Errorf("updating inventory.quantity: %w", err)
		}
	}

	if _, err := s.db.ExecContext(
		ctx,
		`update experiment set finished = true where
                                          uid = $1`,
		exp.UID,
	); err != nil {
		return fmt.Errorf("setting exp finished: %w", err)
	}

	return err

}

func (s experiment) isEquipmentAvailableForUser(ctx context.Context, tx *goqu.TxDatabase, eqName, userUID string) (bool, error) {
	var tmp string

	return tx.ScanValContext(ctx, &tmp,
		`select eq.uid from equipment eq
join equipment_in_access_group eqag on eqag.equipment_uid = eq.uid
join clients_in_access_group cag on cag.access_group_uid = eqag.access_group_uid
join client c on c.uid = cag.client_uid
where c.uid = $1 and eq.name = $2`, userUID, eqName)
}

func (s experiment) isInventoryAvailableForUser(ctx context.Context, tx *goqu.TxDatabase, inName, userUID string) (bool, error) {
	var tmp string

	return tx.ScanValContext(ctx, &tmp,
		`select i.uid from inventory i 
    join inventory_in_access_group iag on iag.inventory_uid = i.uid
    join clients_in_access_group cag on cag.access_group_uid = iag.access_group_uid
    join client c on c.uid = cag.client_uid
where c.uid = $1 and i.name = $2`,
		userUID, inName,
	)
}

func (s experiment) getEquipmentUIDByName(ctx context.Context, tx *goqu.TxDatabase, name string) (string, bool, error) {
	var uid string

	found, err := tx.From(schema.Equipment).
		Select("uid").
		Where(goqu.I("name").Eq(name)).
		Prepared(true).Executor().
		ScanValContext(ctx, &uid)

	return uid, found, err
}

func (s experiment) getInventoryUIDByName(ctx context.Context, tx *goqu.TxDatabase, name string) (string, bool, error) {
	var uid string

	found, err := tx.From(schema.Inventory).
		Select("uid").
		Where(goqu.I("name").Eq(name)).
		Prepared(true).Executor().
		ScanValContext(ctx, &uid)

	return uid, found, err
}
