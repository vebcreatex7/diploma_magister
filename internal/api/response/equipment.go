package response

type Equipment struct {
	UID          string
	Name         string
	Description  string
	Type         string
	Manufacturer string
	Model        string
	Room         string
}

type EquipmentSchedule struct {
	Date      string
	Intervals string
}

type EquipmentInExperiment struct {
	Name      string
	Intervals []string
}

type EquipmentInMaintaince struct {
	Name      string
	Intervals []string
}
