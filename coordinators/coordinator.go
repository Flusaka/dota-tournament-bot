package coordinators

import (
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/flusaka/dota-tournament-bot/types"
)

type QueryCoordinator interface {
	GetLeagues(query *queries.GetLeagues) ([]*types.League, error)
}
