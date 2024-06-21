package datasource

import (
	"fmt"
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/flusaka/dota-tournament-bot/types"
	"github.com/flusaka/dota-tournament-bot/utils"
	"github.com/flusaka/pandascore-go"
	"github.com/flusaka/pandascore-go/clients"
	psquery "github.com/flusaka/pandascore-go/clients/queries"
	pstypes "github.com/flusaka/pandascore-go/types"
	"golang.org/x/exp/slices"
)

type PandascoreDataSource struct {
	pandascoreClient *pandascore.Client
}

func NewPandascoreDataSource(client *pandascore.Client) *PandascoreDataSource {
	return &PandascoreDataSource{
		pandascoreClient: client,
	}
}

func (ps *PandascoreDataSource) GetTournaments(query *queries.GetTournaments) ([]types.Tournament, error) {
	return nil, nil
}

func (ps *PandascoreDataSource) GetRunningTournaments(query *queries.GetTournaments) ([]types.Tournament, error) {
	running, err := ps.pandascoreClient.Dota2.GetRunningTournaments()
	if err != nil {
		return nil, err
	}
	tournaments := utils.MapStructTo[pstypes.Tournament, types.Tournament](running, func(input pstypes.Tournament) types.Tournament {
		return types.Tournament{
			BaseTournament: types.BaseTournament{
				ID:          input.Id,
				DisplayName: input.Name,
			},
			Matches: utils.MapStructTo[pstypes.BaseMatch, types.BaseMatch](input.Matches, func(input pstypes.BaseMatch) types.BaseMatch {
				var streamUrl = ""
				for _, stream := range input.StreamsList {
					if stream.Language == "en" && stream.Official {
						streamUrl = stream.RawUrl
					}
				}

				return types.BaseMatch{
					ID:            input.Id,
					Name:          input.Name,
					ScheduledTime: input.ScheduledAt.Unix(),
					StreamUrl:     streamUrl,
				}
			}),
		}
	})
	return tournaments, nil
}

func (ps *PandascoreDataSource) GetUpcomingTournaments(query *queries.GetTournaments) ([]types.Tournament, error) {
	running, err := ps.pandascoreClient.Dota2.GetUpcomingTournaments()
	if err != nil {
		return nil, err
	}
	tournaments := utils.MapStructTo[pstypes.Tournament, types.Tournament](running, func(input pstypes.Tournament) types.Tournament {
		return types.Tournament{
			BaseTournament: types.BaseTournament{
				ID:          input.Id,
				DisplayName: input.Name,
			},
			Matches: utils.MapStructTo[pstypes.BaseMatch, types.BaseMatch](input.Matches, func(input pstypes.BaseMatch) types.BaseMatch {
				var streamUrl = ""
				for _, stream := range input.StreamsList {
					if stream.Language == "en" && stream.Official {
						streamUrl = stream.RawUrl
					}
				}

				return types.BaseMatch{
					ID:            input.Id,
					Name:          input.Name,
					ScheduledTime: input.ScheduledAt.Unix(),
					StreamUrl:     streamUrl,
				}
			}),
		}
	})
	return tournaments, nil
}

func (ps *PandascoreDataSource) GetUpcomingMatches(query *queries.GetUpcomingMatches) ([]types.Match, error) {
	upcoming, err := ps.pandascoreClient.Dota2.GetUpcomingMatchesWithParams(clients.MatchParams{
		Range: psquery.MatchRange{
			BeginAt: &psquery.DateRange{
				Lower: query.BeginAt.Start,
				Upper: query.BeginAt.End,
			},
		},
		Sort: psquery.NewMatchSort([]psquery.MatchSortField{
			{
				FieldName:  psquery.MatchSortTournamentId,
				Descending: true,
			},
		}),
	})
	if err != nil {
		return nil, err
	}

	for _, match := range upcoming {
		fmt.Printf("Match: %v, Tournament: %v, Tier: %v\n", match.Name, match.Tournament.Name, match.Tournament.Tier)
	}

	upcoming = utils.FilterWhere[pstypes.Match](upcoming, func(element pstypes.Match) bool {
		return isIncludedTier(query.Tiers, pstypes.Tier(element.Tournament.Tier))
	})

	matches := make([]types.Match, 0, len(upcoming))
	for _, match := range upcoming {
		var streamUrl = ""
		for _, stream := range match.StreamsList {
			if stream.Language == "en" && stream.Official {
				streamUrl = stream.RawUrl
			}
		}

		teamOneName := "TBD"
		teamTwoName := "TBD"

		if len(match.Opponents) > 0 {
			teamOneName = match.Opponents[0].Name
		}
		if len(match.Opponents) > 1 {
			teamTwoName = match.Opponents[1].Name
		}

		matches = append(matches, types.Match{
			BaseMatch: types.BaseMatch{
				ID: match.Id,
				TeamOne: &types.Team{
					DisplayName: teamOneName,
				},
				TeamTwo: &types.Team{
					DisplayName: teamTwoName,
				},
				ScheduledTime: match.BeginAt.Unix(),
				StreamUrl:     streamUrl,
			},
			Tournament: types.BaseTournament{
				ID:          match.TournamentId,
				DisplayName: match.Tournament.Name,
			},
		})
	}
	return matches, nil
}

func isIncludedTier(expectedTiers []types.Tier, actualTier pstypes.Tier) bool {
	var mappedTiers = make([]pstypes.Tier, len(expectedTiers))
	for i, tier := range expectedTiers {
		switch tier {
		case types.TierS:
			{
				mappedTiers[i] = pstypes.TierS
			}
		case types.TierA:
			{
				mappedTiers[i] = pstypes.TierA
			}
		case types.TierB:
			{
				mappedTiers[i] = pstypes.TierB
			}
		case types.TierC:
			{
				mappedTiers[i] = pstypes.TierC
			}
		case types.TierD:
			{
				mappedTiers[i] = pstypes.TierD
			}
		}
	}
	fmt.Printf("Expected tiers: %v, Actual Tier: %v\n", mappedTiers, actualTier)
	return slices.Contains(mappedTiers, actualTier)
}
