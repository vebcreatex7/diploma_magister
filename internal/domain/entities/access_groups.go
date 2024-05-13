package entities

import "slices"

type AccessGroup struct {
	UID         string `db:"uid" goqu:"skipupdate"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type AccessGroupExt struct {
	AccessGroup
	Users     []string
	Equipment []string
	Inventory []string
}

func (e *AccessGroupExt) RemoveDuplicates() {
	slices.Sort(e.Users)
	e.Users = slices.Compact(e.Users)
	slices.Sort(e.Equipment)
	e.Equipment = slices.Compact(e.Equipment)
	slices.Sort(e.Inventory)
	e.Inventory = slices.Compact(e.Inventory)
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
