package datasource

import (
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/flusaka/dota-tournament-bot/types"
)

type Client interface {
	GetLeagues(query *queries.GetLeagues) ([]*types.League, error)
}
