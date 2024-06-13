package types

type BaseTournament struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

type Tournament struct {
	BaseTournament
	Matches []BaseMatch
}
