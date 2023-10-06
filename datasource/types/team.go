package types

type Team struct {
	DisplayName string
}

func NewTeam(displayName string) *Team {
	return &Team{
		displayName,
	}
}
