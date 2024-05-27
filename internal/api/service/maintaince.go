package service

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/vebcreatex7/diploma_magister/internal/api/request"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"github.com/vebcreatex7/diploma_magister/internal/domain/constant"
	"github.com/vebcreatex7/diploma_magister/internal/domain/entities"
	"github.com/vebcreatex7/diploma_magister/internal/domain/repo"
	"github.com/vebcreatex7/diploma_magister/internal/repo/postgres/schema"
	"github.com/vebcreatex7/diploma_magister/pkg/mailer"
	"sort"
)

type maintaince struct {
	db               *goqu.Database
	clientsRepo      repo.Clients
	equipmentService equipment
	m                *mailer.Mailer
}

func NewMaintaince(
	db *goqu.Database,
	clientsRepo repo.Clients,
	equipmentService equipment,
	m *mailer.Mailer,
) maintaince {
	return maintaince{
		db:               db,
		clientsRepo:      clientsRepo,
		equipmentService: equipmentService,
		m:                m,
	}
}

func (s maintaince) GetSuggestions(ctx context.Context) (response.MaintainceSuggestions, error) {
	var eq []string

	if err := s.db.ScanValsContext(
		ctx,
		&eq,
		`select name from equipment`,
	); err != nil {
		return response.MaintainceSuggestions{}, fmt.Errorf("getting equipment suggestions: %w", err)
	}

	return response.MaintainceSuggestions{
		Equipment: eq,
	}, nil
}

func (s maintaince) GetAll(ctx context.Context) ([]response.Maintaince, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	var res []response.Maintaince

	var mts []entities.Maintaince

	if err := tx.From(schema.Maintaince).
		Select(entities.Maintaince{}).
		Order(goqu.C("start_ts").Asc()).
		Prepared(true).
		Executor().ScanStructsContext(ctx, &mts); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("getting maintainces: %w", err)
	}

	for _, mt := range mts {
		var users []string

		if err := tx.ScanValsContext(
			ctx,
			&users,
			`select login from client c 
join clients_in_maintaince cinmt on c.uid = cinmt.client_uid
where cinmt.maintaince_uid = $1`, mt.UID); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("getting users in mt: %w", err)
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
join equipment_schedule_in_maintaince esinmt on es.uid = esinmt.equipment_schedule_uid
where esinmt.maintaince_uid = $1`, mt.UID); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("getting equipment in mt: %w", err)
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

		var eqRes []response.EquipmentInMaintaince

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
			eqRes = append(eqRes, response.EquipmentInMaintaince{
				Name:      n,
				Intervals: intervalsStr,
			})
		}

		sort.Slice(eqRes, func(i, j int) bool {
			return eqRes[i].Name < eqRes[j].Name
		})

		res = append(res, response.Maintaince{
			UID:         mt.UID,
			Name:        mt.Name,
			Description: mt.Description,
			StartTs:     mt.StartTs.Format(constant.Layout),
			EndTs:       mt.EndTs.Format(constant.Layout),
			Users:       users,
			Equipment:   eqRes,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s maintaince) GetAllForUser(ctx context.Context, userUID string) ([]response.Maintaince, error) {
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

func (s maintaince) AddMaintaince(ctx context.Context, req request.AddMaintaince, userUID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	mt := entities.Maintaince{
		UID:         uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		StartTs:     req.StartTs,
		EndTs:       req.EndTs,
	}

	var ess []entities.EquipmentSchedule
	var esinmt []entities.EquipmentScheduleInMaintaince

	for i := range req.Equipment {
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

		ess = append(ess, tmp)
		esinmt = append(esinmt, entities.EquipmentScheduleInMaintaince{
			MaintainceUID:        mt.UID,
			EquipmentScheduleUID: tmp.UID,
		})
	}

	for _, es := range ess {
		var essIntersected []entities.EquipmentSchedule

		if err := s.db.ScanStructsContext(
			ctx,
			&essIntersected,
			`select * from equipment_schedule
where equipment_uid = $1 and time_interval && $2`,
			es.EquipmentUID,
			es.TimeInterval,
		); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("checking equipment_schedule intersections: %w", err)
		}

		if err := s.alertEquipmentScheduleChanges(ctx, tx, essIntersected); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("alerting: %w", err)
		}

		if err = s.deleteEquipmentSchedule(ctx, tx, essIntersected); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("deleting equipment_schedule: %w", err)
		}
	}

	if _, err := tx.Insert(schema.Maintaince).
		Rows(mt).
		Prepared(true).Executor().
		ExecContext(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("inserting maintaince: %w", err)
	}

	if _, err := tx.Insert(schema.EquipmentSchedule).
		Rows(ess).
		Prepared(true).Executor().
		ExecContext(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("inserting equipment_schedule: %w", err)
	}

	if _, err := tx.Insert(schema.EquipmentScheduleInMaintaince).
		Rows(esinmt).
		Prepared(true).Executor().
		ExecContext(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("inserting equipment_schedule_in_maintaince: %w", err)
	}

	if _, err := tx.Insert(schema.ClientsInMaintaince).
		Rows(entities.ClientsInMaintaince{
			MaintainceUID: mt.UID,
			ClientUID:     userUID,
		}).
		Prepared(true).Executor().
		ExecContext(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("inserting clients_in_experiment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s maintaince) DeleteByUID(ctx context.Context, uid string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	var curEsinmt []entities.EquipmentScheduleInMaintaince
	if err = tx.ScanStructsContext(
		ctx,
		&curEsinmt,
		`select * from equipment_schedule_in_maintaince
where maintaince_uid = $1`,
		uid,
	); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("getting equipment_schedule_in_maintaince: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`delete from equipment_schedule_in_maintaince
where maintaince_uid = $1`,
		uid,
	); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("deleting equipment_schedule_in_maintaince: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`delete from clients_in_maintaince
where maintaince_uid = $1`,
		uid,
	); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("deleting clients_in_maintaince: %w", err)
	}

	for i := range curEsinmt {
		if _, err := tx.ExecContext(
			ctx,
			`delete from equipment_schedule
where uid = $1`,
			curEsinmt[i].EquipmentScheduleUID,
		); err != nil {
			if err = tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("deleting equipment_schedule: %w", err)
		}
	}

	if _, err = tx.ExecContext(
		ctx,
		`delete from maintaince
where uid = $1`,
		uid,
	); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("deleting maintaince: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s maintaince) DeleteByUIDForUser(ctx context.Context, uid, userUID string) error {
	var tmp string

	av, err := s.db.ScanValContext(
		ctx,
		&tmp,
		`select client_uid from clients_in_maintaince
where maintaince_uid = $1 and client_uid = $2`,
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

func (s maintaince) deleteEquipmentSchedule(ctx context.Context, tx *goqu.TxDatabase, essIntersected []entities.EquipmentSchedule) error {
	for _, esIntersected := range essIntersected {
		if _, err := s.db.ExecContext(
			ctx,
			`delete from equipment_schedule_in_experiment
where equipment_schedule_uid = $1`,
			esIntersected.UID,
		); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("deleting equipment_schedule_in_experiment intersections: %w", err)
		}

		if _, err := s.db.ExecContext(
			ctx,
			`delete from equipment_schedule_in_maintaince
where equipment_schedule_uid = $1`,
			esIntersected.UID,
		); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("deleting equipment_schedule_in_maintaince intersections: %w", err)
		}

		if _, err := s.db.ExecContext(
			ctx,
			`delete from equipment_schedule
where uid = $1`,
			esIntersected.UID,
		); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("deleting equipment_schedule intersections: %w", err)
		}
	}
	return nil
}

func (s maintaince) alertEquipmentScheduleChanges(ctx context.Context, tx *goqu.TxDatabase, essIntersected []entities.EquipmentSchedule) error {
	for _, esIntersected := range essIntersected {
		var eqName string

		if _, err := tx.ScanValContext(
			ctx,
			&eqName,
			`select name from equipment where uid = $1`,
			esIntersected.EquipmentUID,
		); err != nil {
			return fmt.Errorf("getting eq name: %w", err)
		}

		type emailWithProcessName struct {
			Email string `db:"email"`
			Name  string `db:"name"`
		}

		var en emailWithProcessName

		enFound, err := tx.ScanStructContext(
			ctx,
			&en,
			`select email, m.name from client
join clients_in_maintaince cinmt on client.uid = cinmt.client_uid
join equipment_schedule_in_maintaince esinmt on cinmt.maintaince_uid = esinmt.maintaince_uid
join equipment_schedule es on esinmt.equipment_schedule_uid = es.uid
join maintaince m on cinmt.maintaince_uid = m.uid
where es.uid = $1`,
			esIntersected.UID,
		)
		if err != nil {
			return fmt.Errorf("getting engineers emails: %w", err)
		}

		var sc emailWithProcessName
		scFound, err := tx.ScanStructContext(
			ctx,
			&sc,
			`select email, ex.name from client
join clients_in_experiment cinex on client.uid = cinex.client_uid
join equipment_schedule_in_experiment esinex on cinex.experiment_uid = esinex.experiment_uid
join equipment_schedule es on esinex.equipment_schedule_uid = es.uid
join experiment ex on cinex.experiment_uid = ex.uid
where es.uid = $1`,
			esIntersected.UID,
		)
		if err != nil {
			return fmt.Errorf("getting engineers emails: %w", err)
		}

		if enFound {
			if err := s.sendMessageToEngineer(en.Email, en.Name, eqName, esIntersected); err != nil {
				return fmt.Errorf("sending cancellation message to engineer: %w", err)
			}
		}

		if scFound {
			if err := s.sendMessageToScientist(sc.Email, sc.Name, eqName, esIntersected); err != nil {
				return fmt.Errorf("sending cancellation message to scientist: %w", err)
			}
		}

	}
	return nil
}

func (s maintaince) sendMessageToScientist(email, expName, eqName string, es entities.EquipmentSchedule) error {
	return s.m.Send(
		email,
		fmt.Sprintf("Отмена брони в эксперименте %s", expName),
		fmt.Sprintf(`В эксперименте %s была отменена бронь оборудования %s в интервале (%s. %s) в связи с запланированными работами`,
			expName, eqName, es.TimeInterval.Lower.Time.Format(constant.Layout), es.TimeInterval.Upper.Time.Format(constant.Layout),
		),
	)
}

func (s maintaince) sendMessageToEngineer(email, mtName, eqName string, es entities.EquipmentSchedule) error {
	return s.m.Send(
		email,
		fmt.Sprintf("Отмена брони в обслуживание %s", mtName),
		fmt.Sprintf(`В обслуживание %s была отменена бронь оборудования %s в интервале (%s. %s) в связи с запланированными работами`,
			mtName, eqName, es.TimeInterval.Lower.Time.Format(constant.Layout), es.TimeInterval.Upper.Time.Format(constant.Layout),
		),
	)
}
