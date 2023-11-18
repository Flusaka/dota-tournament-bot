package types

type League struct {
	ID          int      `json:"id"`
	DisplayName string   `json:"displayName"`
	Matches     []*Match `json:"matches"`
}

func NewLeague(id int, displayName string) *League {
	return &League{
		id, displayName, make([]*Match, 1),
	}
}

func NewLeagueWithMatches(id int, displayName string, matches []*Match) *League {
	return &League{
		id, displayName, matches,
	}
}
