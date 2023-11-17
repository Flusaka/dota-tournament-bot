package types

type Match struct {
	ID            int16  `json:"id"`
	TeamOne       *Team  `json:"teamOne"`
	TeamTwo       *Team  `json:"teamTwo"`
	ScheduledTime int64  `json:"scheduledTime"`
	StreamUrl     string `json:"streamUrl"`
	// Reverse lookup from this match to the one TeamOne came from (if any)
	TeamOneSourceMatch *Match `json:"teamOneSourceMatch"`
	// Reverse lookup from this match to the one TeamTwo came from (if any)
	TeamTwoSourceMatch *Match `json:"teamTwoSourceMatch"`
	// Forward lookup to the match the winner of this match will progress to (if any)
	WinningTeamMatch *Match `json:"winningTeamMatch"`
	// Forward lookup to the match the loser of this match will progress to (if any)
	LosingTeamMatch *Match `json:"losingTeamMatch"`
}

func NewMatch(id int16, teamOne *Team, teamTwo *Team, scheduledTime int64, streamUrl string) *Match {
	return &Match{
		ID:            id,
		TeamOne:       teamOne,
		TeamTwo:       teamTwo,
		ScheduledTime: scheduledTime,
		StreamUrl:     streamUrl,
	}
}
