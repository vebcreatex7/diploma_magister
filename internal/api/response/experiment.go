package response

type ExperimentSuggestions struct {
	Equipment []string
	Inventory []string
}

type Experiment struct {
	UID         string
	Name        string
	Description string
	StartTs     string
	EndTs       string
	Users       []string
	Equipment   []EquipmentInExperiment
	Inventory   []InventoryInExperiment
}
