package datasource

import (
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/flusaka/dota-tournament-bot/types"
	"github.com/flusaka/pandascore-go"
	"github.com/flusaka/pandascore-go/clients"
	psquery "github.com/flusaka/pandascore-go/clients/queries"
)

type PandascoreDataSource struct {
	pandascoreClient *pandascore.Client
}

func NewPandascoreDataSource(client *pandascore.Client) *PandascoreDataSource {
	return &PandascoreDataSource{
		pandascoreClient: client,
	}
}

func (ps *PandascoreDataSource) GetTournaments(query *queries.GetTournaments) ([]*types.Tournament, error) {
	return nil, nil
}

func (ps *PandascoreDataSource) GetUpcomingMatches(query *queries.GetUpcomingMatches) ([]*types.Match, error) {
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
	matches := make([]*types.Match, len(upcoming))
	for _, match := range upcoming {
		var streamUrl = ""
		for _, stream := range match.StreamsList {
			if stream.Language == "en" && stream.Official {
				streamUrl = stream.RawUrl
			}
		}

		matches = append(matches, &types.Match{
			ID: match.Id,
			TeamOne: &types.Team{
				DisplayName: match.Opponents[0].Name,
			},
			TeamTwo: &types.Team{
				DisplayName: match.Opponents[1].Name,
			},
			ScheduledTime: match.BeginAt.Unix(),
			StreamUrl:     streamUrl,
		})
	}
	return matches, nil
}
