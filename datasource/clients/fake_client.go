package clients

import (
	"encoding/json"
	"github.com/flusaka/dota-tournament-bot/datasource/queries"
	"github.com/flusaka/dota-tournament-bot/datasource/types"
	"math"
	"math/rand"
	"os"
	"time"
)

type FakeDataSourceClient struct {
	leagues []*types.League
}

const (
	leagueStoredFilename = "leagues.json"
	letterBytes          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	idLength             = 6
)

func NewFakeDataSourceClient(reset bool) FakeDataSourceClient {
	if reset {
		deleteStoredLeagues()
	}
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
		now := time.Now().Truncate(time.Minute)
		// If there is no file, generate data and store to it
		leagues = []*types.League{
			{
				ID:          0,
				DisplayName: "The International 2023",
				Matches: []*types.Match{
					{
						ID:            randomID(),
						Radiant:       &types.Team{DisplayName: "Team Liquid"},
						Dire:          &types.Team{DisplayName: "Nigma Galaxy"},
						ScheduledTime: now.Add(time.Minute * 1).UTC().Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						ID:            randomID(),
						Radiant:       &types.Team{DisplayName: "Tundra Esports"},
						Dire:          &types.Team{DisplayName: "Gaimin Gladiators"},
						ScheduledTime: now.Add(time.Minute * 1).UTC().Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						ID:            randomID(),
						Radiant:       &types.Team{DisplayName: "Team Spirit"},
						Dire:          &types.Team{DisplayName: "PSG.LGD"},
						ScheduledTime: now.Add(time.Minute * 2).UTC().Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						ID:            randomID(),
						Radiant:       &types.Team{DisplayName: "Team Aster"},
						Dire:          &types.Team{DisplayName: "Shopify Rebellion"},
						ScheduledTime: now.Add(time.Minute * 2).UTC().Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						ID:            randomID(),
						Radiant:       &types.Team{DisplayName: "Evil Geniuses"},
						Dire:          &types.Team{DisplayName: "OG"},
						ScheduledTime: now.Add(time.Minute * 3).UTC().Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						ID:            randomID(),
						Radiant:       &types.Team{DisplayName: "Team Secret"},
						Dire:          &types.Team{DisplayName: "OG"},
						ScheduledTime: now.Add(time.Minute * 3).UTC().Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						ID:            randomID(),
						Radiant:       &types.Team{DisplayName: "Team Liquid"},
						Dire:          &types.Team{DisplayName: "Newbee"},
						ScheduledTime: now.Add(time.Minute * 4).UTC().Unix(),
						StreamUrl:     "https://twitch.tv/dota2",
					},
					{
						ID:            randomID(),
						Radiant:       &types.Team{DisplayName: "PSG.LGD"},
						Dire:          &types.Team{DisplayName: "OG"},
						ScheduledTime: now.Add(time.Minute * 4).UTC().Unix(),
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

func deleteStoredLeagues() error {
	return os.Remove(leagueStoredFilename)
}

func randomID() int16 {
	return int16(rand.Intn(math.MaxInt16))
}
