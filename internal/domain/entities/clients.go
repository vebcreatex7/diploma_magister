package entities

type Client struct {
	UID          string `db:"uid" goqu:"skipinsert,skipupdate"`
	Surname      string `db:"surname"`
	Name         string `db:"name"`
	Patronymic   string `db:"patronymic"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash" goqu:"skipupdate"`
	Email        string `db:"email"`
	Role         string `db:"role"`
	Approved     bool   `db:"approved" goqu:"skipupdate"`
}

type ClientsInExperiment struct {
	ExperimentUID string `db:"experiment_uid"`
	ClientUID     string `db:"client_uid"`
}

type ClientsInMaintaince struct {
	MaintainceUID string `db:"maintaince_uid"`
	ClientUID     string `db:"client_uid"`
}
