package queries

import (
	"github.com/flusaka/dota-tournament-bot/types"
	"github.com/mitchellh/hashstructure/v2"
)

type GetTournaments struct {
	Tiers    []types.Tier
	Finished bool
}

func NewGetTournamentsQuery(tiers []types.Tier, finished bool) *GetTournaments {
	return &GetTournaments{
		tiers,
		finished,
	}
}

func NewGetActiveTournamentsQuery(tiers []types.Tier) *GetTournaments {
	return NewGetTournamentsQuery(tiers, false)
}

func (g GetTournaments) HashCode() (uint64, error) {
	return hashstructure.Hash(g, hashstructure.FormatV2, nil)
}
