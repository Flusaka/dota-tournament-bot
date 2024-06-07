package bot

import (
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/flusaka/dota-tournament-bot/types"
)

type DataSource interface {
	GetTournaments(query *queries.GetTournaments) ([]*types.Tournament, error)
	GetUpcomingMatches(query *queries.GetUpcomingMatches) ([]*types.Match, error)
}