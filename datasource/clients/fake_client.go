package clients

import (
	"github.com/flusaka/dota-tournament-bot/datasource/queries"
	"github.com/flusaka/dota-tournament-bot/datasource/types"
	"math"
	"math/rand"
	"time"
)

type FakeDataSourceClient struct {
	leagues []*types.League
}

const (
	leagueStoredFilename = "leagues.json"
)

func NewFakeDataSourceClient() FakeDataSourceClient {
	leagues := loadLeagues()
	return FakeDataSourceClient{
		leagues: leagues,
	}
}

func (receiver FakeDataSourceClient) GetLeagues(query *queries.GetLeagues) ([]*types.League, error) {
	return receiver.leagues, nil
}

func loadLeagues() []*types.League {
	now := time.Now().Round(time.Minute)

	oneTeamProgressedMatch := &types.Match{
		ID:            randomID(),
		TeamOne:       &types.Team{DisplayName: "Gaimin Gladiators"},
		TeamTwo:       nil,
		ScheduledTime: now.Add(time.Minute * 5).UTC().Unix(),
		StreamUrl:     "https://twitch.tv/dota2ti",
	}
	matchProgressing := &types.Match{
		ID:            randomID(),
		TeamOne:       &types.Team{DisplayName: "Team Liquid"},
		TeamTwo:       &types.Team{DisplayName: "Newbee"},
		ScheduledTime: now.Add(time.Minute * 4).UTC().Unix(),
		StreamUrl:     "https://twitch.tv/dota2ti",
	}

	// Set the relevant match references
	oneTeamProgressedMatch.TeamTwoSourceMatch = matchProgressing
	matchProgressing.WinningTeamMatch = oneTeamProgressedMatch

	noTeamsProgressedMatch := &types.Match{
		ID:            randomID(),
		TeamOne:       nil,
		TeamTwo:       nil,
		ScheduledTime: now.Add(time.Minute * 5).UTC().Unix(),
		StreamUrl:     "https://twitch.tv/dota2ti",
	}
	firstMatchProgressing := &types.Match{
		ID:            randomID(),
		TeamOne:       &types.Team{DisplayName: "Team Liquid"},
		TeamTwo:       &types.Team{DisplayName: "Newbee"},
		ScheduledTime: now.Add(time.Minute * 4).UTC().Unix(),
		StreamUrl:     "https://twitch.tv/dota2ti",
	}
	secondMatchProgressing := &types.Match{
		ID:            randomID(),
		TeamOne:       &types.Team{DisplayName: "OG"},
		TeamTwo:       &types.Team{DisplayName: "Shopify Rebellion"},
		ScheduledTime: now.Add(time.Minute * 4).UTC().Unix(),
		StreamUrl:     "https://twitch.tv/dota2ti",
	}

	firstMatchProgressing.WinningTeamMatch = noTeamsProgressedMatch
	secondMatchProgressing.WinningTeamMatch = noTeamsProgressedMatch
	noTeamsProgressedMatch.TeamOneSourceMatch = firstMatchProgressing
	noTeamsProgressedMatch.TeamTwoSourceMatch = secondMatchProgressing

	// If there is no file, generate data and store to it
	leagues := []*types.League{
		{
			ID:          0,
			DisplayName: "The International 2023",
			Matches: []*types.Match{
				{
					ID:            randomID(),
					TeamOne:       &types.Team{DisplayName: "Team Liquid"},
					TeamTwo:       &types.Team{DisplayName: "Nigma Galaxy"},
					ScheduledTime: now.Add(time.Minute * 1).UTC().Unix(),
					StreamUrl:     "https://twitch.tv/dota2ti",
				},
				{
					ID:            randomID(),
					TeamOne:       &types.Team{DisplayName: "Tundra Esports"},
					TeamTwo:       &types.Team{DisplayName: "Gaimin Gladiators"},
					ScheduledTime: now.Add(time.Minute * 1).UTC().Unix(),
					StreamUrl:     "https://twitch.tv/dota2ti",
				},
				{
					ID:            randomID(),
					TeamOne:       &types.Team{DisplayName: "Team Spirit"},
					TeamTwo:       &types.Team{DisplayName: "PSG.LGD"},
					ScheduledTime: now.Add(time.Minute * 2).UTC().Unix(),
					StreamUrl:     "https://twitch.tv/dota2ti",
				},
				{
					ID:            randomID(),
					TeamOne:       &types.Team{DisplayName: "Team Aster"},
					TeamTwo:       &types.Team{DisplayName: "Shopify Rebellion"},
					ScheduledTime: now.Add(time.Minute * 2).UTC().Unix(),
					StreamUrl:     "https://twitch.tv/dota2ti",
				},
				{
					ID:            randomID(),
					TeamOne:       &types.Team{DisplayName: "Evil Geniuses"},
					TeamTwo:       &types.Team{DisplayName: "OG"},
					ScheduledTime: now.Add(time.Minute * 3).UTC().Unix(),
					StreamUrl:     "https://twitch.tv/dota2ti",
				},
				{
					ID:            randomID(),
					TeamOne:       &types.Team{DisplayName: "Team Secret"},
					TeamTwo:       &types.Team{DisplayName: "OG"},
					ScheduledTime: now.Add(time.Minute * 3).UTC().Unix(),
					StreamUrl:     "https://twitch.tv/dota2ti",
				},
				{
					ID:            randomID(),
					TeamOne:       &types.Team{DisplayName: "PSG.LGD"},
					TeamTwo:       &types.Team{DisplayName: "OG"},
					ScheduledTime: now.Add(time.Minute * 4).UTC().Unix(),
					StreamUrl:     "https://twitch.tv/dota2ti",
				},
				matchProgressing,
				oneTeamProgressedMatch,
				firstMatchProgressing,
				secondMatchProgressing,
				noTeamsProgressedMatch,
			},
		},
	}
	return leagues
}

func randomID() int16 {
	return int16(rand.Intn(math.MaxInt16))
}
