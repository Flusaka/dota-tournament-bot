package queries

import (
	"github.com/flusaka/dota-tournament-bot/types"
	"github.com/mitchellh/hashstructure/v2"
)

type GetTournaments struct {
	Tiers []types.Tier
}

func NewGetTournamentsQuery(tiers []types.Tier) *GetTournaments {
	return &GetTournaments{
		tiers,
	}
}

func (g GetTournaments) HashCode() (uint64, error) {
	return hashstructure.Hash(g, hashstructure.FormatV2, nil)
}
