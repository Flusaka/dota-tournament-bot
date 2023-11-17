package types

type Team struct {
	DisplayName string `json:"displayName"`
}

func NewTeam(displayName string) *Team {
	return &Team{
		displayName,
	}
}
