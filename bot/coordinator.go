package bot

import (
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/flusaka/dota-tournament-bot/types"
)

type QueryCoordinator interface {
	GetTournaments(query *queries.GetTournaments) ([]types.Tournament, error)
	GetMatches(query *queries.GetMatches) ([]types.Match, error)
	GetUpcomingMatches(query *queries.GetUpcomingMatches) ([]types.Match, error)
}
