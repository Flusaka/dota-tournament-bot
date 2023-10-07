package clients

import (
	"github.com/flusaka/dota-tournament-bot/datasource/queries"
	"github.com/flusaka/dota-tournament-bot/datasource/types"
	"github.com/flusaka/dota-tournament-bot/stratz"
	"github.com/flusaka/dota-tournament-bot/stratz/schema"
)

type StratzDataSourceClient struct {
	stratzClient *stratz.Client
}

func NewStratzDataSourceClient(stratzClient *stratz.Client) StratzDataSourceClient {
	return StratzDataSourceClient{
		stratzClient,
	}
}

func (receiver StratzDataSourceClient) GetLeagues(query *queries.GetLeagues) ([]*types.League, error) {
	leagueTiers := convertTiers(query.Tiers)
	leagues, err := receiver.stratzClient.GetLeagues(leagueTiers, query.Finished)
	if err != nil {
		return nil, err
	}

	var mappedLeagues []*types.League
	for _, league := range leagues {
		var matches []*types.Match
		for _, nodeGroup := range league.NodeGroups {
			for _, node := range nodeGroup.Nodes {
				radiantTeam := types.NewTeam(node.TeamOne.Name)
				direTeam := types.NewTeam(node.TeamTwo.Name)
				streamUrl := ""
				if len(node.Streams) > 0 {
					// Get English stream, if exists, if not :shrug:
					for _, stream := range node.Streams {
						if stream.LanguageId == schema.LanguageEnglish {
							streamUrl = stream.StreamUrl
							break
						}
					}
				}
				match := types.NewMatch(radiantTeam, direTeam, node.ScheduledTime, streamUrl)
				matches = append(matches, match)
			}
		}
		mappedLeague := types.NewLeagueWithMatches(league.Id, league.DisplayName, matches)
		mappedLeagues = append(mappedLeagues, mappedLeague)
	}
	return mappedLeagues, nil
}

func convertTiers(tiers []types.Tier) []schema.LeagueTier {
	var leagueTiers []schema.LeagueTier
	for _, t := range tiers {
		leagueTiers = append(leagueTiers, convertTier(t))
	}
	return leagueTiers
}

func convertTier(tier types.Tier) schema.LeagueTier {
	switch tier {
	case types.TierAmateur:
		return schema.LeagueTierAmateur
	case types.TierProfessional:
		return schema.LeagueTierProfessional
	case types.TierMinor:
		return schema.LeagueTierMinor
	case types.TierMajor:
		return schema.LeagueTierMajor
	case types.TierInternational:
		return schema.LeagueTierInternational
	case types.TierDpcQualifier:
		return schema.LeagueTierDpcQualifier
	case types.TierDpcLeagueQualifier:
		return schema.LeagueTierDpcLeagueQualifier
	case types.TierDpcLeague:
		return schema.LeagueTierDpcLeague
	case types.TierDpcLeagueFinals:
		return schema.LeagueTierDpcLeagueFinals
	}
	return schema.LeagueTierUnset
}
