package queries

import (
	"github.com/flusaka/dota-tournament-bot/types"
	"github.com/mitchellh/hashstructure/v2"
)

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

func (g GetLeagues) HashCode() (uint64, error) {
	return hashstructure.Hash(g, hashstructure.FormatV2, nil)
}
