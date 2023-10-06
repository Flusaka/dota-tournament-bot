package datasource

import (
	"github.com/flusaka/dota-tournament-bot/datasource/queries"
	"github.com/flusaka/dota-tournament-bot/datasource/types"
)

type Client interface {
	GetLeagues(query *queries.GetLeagues) ([]*types.League, error)
}
