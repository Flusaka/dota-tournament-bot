package queries

import (
	"github.com/flusaka/dota-tournament-bot/types"
	"github.com/mitchellh/hashstructure/v2"
)

type GetMatches struct {
	BeginAt DateRange
	Tiers   []types.Tier
}

type GetUpcomingMatches struct {
	GetMatches
}

func (q GetMatches) HashCode() (uint64, error) {
	return hashstructure.Hash(q, hashstructure.FormatV2, nil)
}
