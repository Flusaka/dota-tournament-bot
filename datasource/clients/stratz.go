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

	// Make a cache of matches for quicker look up later
	matchMap := make(map[int16]*types.Match)
	for _, league := range leagues {
		for _, nodeGroup := range league.NodeGroups {
			for _, node := range nodeGroup.Nodes {
				var teamOne *types.Team
				if node.TeamOne != nil {
					teamOne = types.NewTeam(*node.TeamOne.Name)
				}
				var teamTwo *types.Team
				if node.TeamTwo != nil {
					teamTwo = types.NewTeam(*node.TeamTwo.Name)
				}
				streamUrl := ""
				if len(node.Streams) > 0 {
					// Get English stream, if exists, if not :shrug:
					for _, stream := range node.Streams {
						if *stream.LanguageId == schema.LanguageEnglish {
							streamUrl = *stream.StreamUrl
							break
						}
					}
				}
				matchMap[*node.Id] = types.NewMatch(*node.Id, teamOne, teamTwo, *node.ScheduledTime, streamUrl)
			}
		}
	}

	var mappedLeagues []*types.League
	for _, league := range leagues {
		var matches []*types.Match
		for _, nodeGroup := range league.NodeGroups {
			for _, node := range nodeGroup.Nodes {
				match := matchMap[*node.Id]

				// Now work out which matches are connected to each other, if any
				if node.WinningNodeId != nil {
					winningMatch := matchMap[*node.WinningNodeId]
					if winningMatch.TeamOne == nil && winningMatch.TeamOneSourceMatch == nil {
						winningMatch.TeamOneSourceMatch = match
					} else if winningMatch.TeamTwo == nil && winningMatch.TeamTwoSourceMatch == nil {
						winningMatch.TeamTwoSourceMatch = match
					}
					match.WinningTeamMatch = winningMatch
				}

				if node.LosingNodeId != nil {
					losingMatch := matchMap[*node.LosingNodeId]
					if losingMatch.TeamOne == nil && losingMatch.TeamOneSourceMatch == nil {
						losingMatch.TeamOneSourceMatch = match
					} else if losingMatch.TeamTwo == nil && losingMatch.TeamTwoSourceMatch == nil {
						losingMatch.TeamTwoSourceMatch = match
					}
					match.LosingTeamMatch = losingMatch
				}

				matches = append(matches, match)
			}
		}
		mappedLeague := types.NewLeagueWithMatches(*league.Id, *league.DisplayName, matches)
		mappedLeagues = append(mappedLeagues, mappedLeague)
	}
	return mappedLeagues, nil
}

func convertTiers(tiers []types.Tier) []*schema.LeagueTier {
	var leagueTiers []*schema.LeagueTier
	for _, t := range tiers {
		leagueTiers = append(leagueTiers, convertTier(t))
	}
	return leagueTiers
}

func convertTier(tier types.Tier) *schema.LeagueTier {
	convertedTier := schema.LeagueTierUnset
	switch tier {
	case types.TierAmateur:
		convertedTier = schema.LeagueTierAmateur
	case types.TierProfessional:
		convertedTier = schema.LeagueTierProfessional
	case types.TierMinor:
		convertedTier = schema.LeagueTierMinor
	case types.TierMajor:
		convertedTier = schema.LeagueTierMajor
	case types.TierInternational:
		convertedTier = schema.LeagueTierInternational
	case types.TierDpcQualifier:
		convertedTier = schema.LeagueTierDpcQualifier
	case types.TierDpcLeagueQualifier:
		convertedTier = schema.LeagueTierDpcLeagueQualifier
	case types.TierDpcLeague:
		convertedTier = schema.LeagueTierDpcLeague
	case types.TierDpcLeagueFinals:
		convertedTier = schema.LeagueTierDpcLeagueFinals
	}
	return &convertedTier
}
