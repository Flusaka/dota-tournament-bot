package datasource

import (
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
				streamUrl := getBestStreamUrl(input.StreamsList)

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
				streamUrl := getBestStreamUrl(input.StreamsList)

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

func (ps *PandascoreDataSource) GetMatches(query *queries.GetMatches) ([]types.Match, error) {
	matches, err := ps.pandascoreClient.Dota2.GetMatchesWithParams(clients.MatchParams{
		Range: psquery.MatchRange{
			ScheduledAt: &psquery.DateRange{
				Lower: query.BeginAt.Start,
				Upper: query.BeginAt.End,
			},
		},
		Sort: psquery.NewMatchSort([]psquery.MatchSortField{
			{
				FieldName:  psquery.MatchSortScheduledAt,
				Descending: false,
			},
		}),
	})
	if err != nil {
		return nil, err
	}

	matches = utils.FilterWhere[pstypes.Match](matches, func(element pstypes.Match) bool {
		return isIncludedTier(query.Tiers, element.Tournament.Tier)
	})

	result := make([]types.Match, 0, len(matches))
	for _, match := range matches {
		streamUrl := getBestStreamUrl(match.StreamsList)

		teamOneName := "TBD"
		teamTwoName := "TBD"

		if len(match.Opponents) > 0 {
			teamOneName = match.Opponents[0].Opponent.Name
		}
		if len(match.Opponents) > 1 {
			teamTwoName = match.Opponents[1].Opponent.Name
		}

		result = append(result, types.Match{
			BaseMatch: types.BaseMatch{
				ID: match.Id,
				TeamOne: &types.Team{
					DisplayName: teamOneName,
				},
				TeamTwo: &types.Team{
					DisplayName: teamTwoName,
				},
				ScheduledTime: match.ScheduledAt.Unix(),
				StreamUrl:     streamUrl,
			},
			Tournament: types.BaseTournament{
				ID:          match.TournamentId,
				DisplayName: match.Tournament.Name,
			},
			Serie: types.BaseSerie{
				Name: match.Serie.Name,
			},
			League: types.BaseLeague{
				Name: match.League.Name,
			},
		})
	}
	return result, nil
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

	upcoming = utils.FilterWhere[pstypes.Match](upcoming, func(element pstypes.Match) bool {
		return isIncludedTier(query.Tiers, element.Tournament.Tier)
	})

	matches := make([]types.Match, 0, len(upcoming))
	for _, match := range upcoming {
		streamUrl := getBestStreamUrl(match.StreamsList)

		teamOneName := "TBD"
		teamTwoName := "TBD"

		if len(match.Opponents) > 0 {
			teamOneName = match.Opponents[0].Opponent.Name
		}
		if len(match.Opponents) > 1 {
			teamTwoName = match.Opponents[1].Opponent.Name
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
				ScheduledTime: match.ScheduledAt.Unix(),
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
	return slices.Contains(mappedTiers, actualTier)
}

func getBestStreamUrl(streamList []pstypes.BaseStream) string {
	streamUrl := ""
	for _, stream := range streamList {
		// If we find an official stream in English at any point, return that
		if stream.Language == "en" && stream.Official {
			streamUrl = stream.RawUrl
			break
		} else if stream.Official && streamUrl == "" { // Backup: We at least find the first official stream, even if it's not English, but we continue in case we have an English stream
			streamUrl = stream.RawUrl
		}
	}
	return streamUrl
}
