package entities

type Role struct {
	UID         string `db:"uid" goqu:"skipinsert,skipupdate"`
	Name        string `db:"name"`
	Description string `db:"description"`
}
