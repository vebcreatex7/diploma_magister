package entities

type AccessGroup struct {
	UID         string `db:"uid" goqu:"skipinsert,skipupdate"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Status      string `db:"status" goqu:"skipinsert,skipupdate"`
}

type AccessGroupExt struct {
	AccessGroup
	Users     []string
	Equipment []string
	Inventory []string
}

type ClientsInAccessGroup struct {
	AccessGroupUID string `db:"access_group_uid"`
	ClientUID      string `db:"client_uid"`
}

type EquipmentInAccessGroup struct {
	AccessGroupUID string `db:"access_group_uid"`
	EquipmentUID   string `db:"equipment_uid"`
}

type InventoryInAccessGroup struct {
	AccessGroupUID string `db:"access_group_uid"`
	InventoryUID   string `db:"inventory_uid"`
}
