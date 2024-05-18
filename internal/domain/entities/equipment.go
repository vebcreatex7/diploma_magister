package entities

import "github.com/jackc/pgtype"

type Equipment struct {
	UID          string `db:"uid" goqu:"skipinsert,skipupdate"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	Type         string `db:"type"`
	Manufacturer string `db:"manufacturer"`
	Model        string `db:"model"`
	Room         string `db:"room"`
}

type EquipmentSchedule struct {
	UID            string         `db:"uid" goqu:"skipupdate"`
	EquipmentUID   string         `db:"equipment_uid"`
	TimeInterval   pgtype.Tsrange `db:"time_interval"`
	MaintainceFlag bool           `db:"maintaince_flag"`
}

type EquipmentScheduleInExperiment struct {
	ExperimentUID        string `db:"experiment_uid"`
	EquipmentScheduleUID string `db:"equipment_schedule_uid"`
}

type EquipmentScheduleInMaintaince struct {
	MaintainceUID        string `db:"maintaince_uid"`
	EquipmentScheduleUID string `db:"equipment_schedule_uid"`
}
