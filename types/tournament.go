package types

type Tournament struct {
	ID          int      `json:"id"`
	DisplayName string   `json:"displayName"`
	Matches     []*Match `json:"matches"`
}

func NewTournament(id int, displayName string) *Tournament {
	return &Tournament{
		id, displayName, make([]*Match, 1),
	}
}

func NewTournamentWithMatches(id int, displayName string, matches []*Match) *Tournament {
	return &Tournament{
		id, displayName, matches,
	}
}
