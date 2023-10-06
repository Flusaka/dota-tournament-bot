package queries

import "github.com/flusaka/dota-tournament-bot/datasource/types"

type GetLeagues struct {
	Tiers    []types.Tier
	Finished bool
}

func NewGetLeaguesQuery(tiers []types.Tier, finished bool) *GetLeagues {
	return &GetLeagues{
		tiers,
		finished,
	}
}

func NewGetActiveLeaguesQuery(tiers []types.Tier) *GetLeagues {
	return NewGetLeaguesQuery(tiers, false)
}
