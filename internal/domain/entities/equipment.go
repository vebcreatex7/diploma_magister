package entities

type Equipment struct {
	UID          string `db:"uid" goqu:"skipinsert,skipupdate"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	Type         string `db:"type"`
	Manufacturer string `db:"manufacturer"`
	Model        string `db:"model"`
	Room         string `db:"room"`
}
