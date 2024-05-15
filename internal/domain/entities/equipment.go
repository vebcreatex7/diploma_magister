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
	UID            string         `db:"uid" goqu:"skipinsert,skipupdate"`
	EquipmentUID   string         `db:"equipment_uid"`
	TimeInterval   pgtype.Tsrange `db:"time_interval"`
	Status         string         // remove me
	MaintainceFlag bool           `db:"maintaince_flag"`
}
