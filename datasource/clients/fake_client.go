package clients

import (
	"encoding/json"
	"github.com/flusaka/dota-tournament-bot/datasource/queries"
	"github.com/flusaka/dota-tournament-bot/datasource/types"
	"os"
	"time"
)

type FakeDataSourceClient struct {
	leagues []*types.League
}

const (
	leagueStoredFilename = "leagues.json"
)

func NewFakeDataSourceClient() FakeDataSourceClient {
	leagues, _ := loadLeagues()
	return FakeDataSourceClient{
		leagues: leagues,
	}
}

func (receiver FakeDataSourceClient) GetLeagues(query *queries.GetLeagues) ([]*types.League, error) {
	return receiver.leagues, nil
}

func loadLeagues() ([]*types.League, error) {
	var leagues []*types.League
	if data, err := os.ReadFile(leagueStoredFilename); err != nil {
		// If there is no file, generate data and store to it
		leagues = []*types.League{
			{
				ID:          0,
				DisplayName: "The International 2023",
				Matches: []*types.Match{
					{
						Radiant:       &types.Team{DisplayName: "Team Liquid"},
						Dire:          &types.Team{DisplayName: "Nigma Galaxy"},
						ScheduledTime: time.Now().Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						Radiant:       &types.Team{DisplayName: "Tundra Esports"},
						Dire:          &types.Team{DisplayName: "Gaimin Gladiators"},
						ScheduledTime: time.Now().Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						Radiant:       &types.Team{DisplayName: "Team Spirit"},
						Dire:          &types.Team{DisplayName: "PSG.LGD"},
						ScheduledTime: time.Now().Add(time.Hour * 24).Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						Radiant:       &types.Team{DisplayName: "Team Aster"},
						Dire:          &types.Team{DisplayName: "Shopify Rebellion"},
						ScheduledTime: time.Now().Add(time.Hour * 24).Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						Radiant:       &types.Team{DisplayName: "Evil Geniuses"},
						Dire:          &types.Team{DisplayName: "OG"},
						ScheduledTime: time.Now().Add(time.Hour * 24 * 2).Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						Radiant:       &types.Team{DisplayName: "Team Secret"},
						Dire:          &types.Team{DisplayName: "OG"},
						ScheduledTime: time.Now().Add(time.Hour * 24 * 2).Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						Radiant:       &types.Team{DisplayName: "Team Liquid"},
						Dire:          &types.Team{DisplayName: "Newbee"},
						ScheduledTime: time.Now().Add(time.Hour * 24 * 3).Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						Radiant:       &types.Team{DisplayName: "PSG.LGD"},
						Dire:          &types.Team{DisplayName: "OG"},
						ScheduledTime: time.Now().Add(time.Hour * 24 * 3).Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
				},
			},
		}
		leaguesJson, _ := json.MarshalIndent(leagues, "", "    ")
		os.WriteFile(leagueStoredFilename, leaguesJson, 0644)
	} else {
		json.Unmarshal(data, &leagues)
	}
	return leagues, nil
}
