package entities

import (
	"time"
)

type Experiment struct {
	UID         string    `db:"uid" goqu:"skipupdate"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	StartTs     time.Time `db:"start_ts"`
	EndTs       time.Time `db:"end_ts"`
	Finished    bool      `db:"finished"`
}
