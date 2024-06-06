package response

type Maintaince struct {
	UID         string
	Name        string
	Description string
	StartTs     string
	EndTs       string
	Users       []string
	Equipment   []EquipmentInMaintaince
}

type MaintainceSuggestions struct {
	Equipment []string
}
