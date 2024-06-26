package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/flusaka/dota-tournament-bot/models"
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/flusaka/dota-tournament-bot/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockChannelSession struct {
	mock.Mock
}

func (s *MockChannelSession) InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error {
	args := s.Called(interaction, resp, options)
	return args.Error(0)
}

func (s *MockChannelSession) ChannelMessageSend(channelID string, message string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	args := s.Called(channelID, message, options)
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

func (s *MockChannelSession) ChannelMessageSendComplex(channelID string, m *discordgo.MessageSend, options ...discordgo.RequestOption) (st *discordgo.Message, err error) {
	args := s.Called(m, options)
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

func (s *MockChannelSession) ChannelMessageEditComplex(m *discordgo.MessageEdit, options ...discordgo.RequestOption) (st *discordgo.Message, err error) {
	args := s.Called(m, options)
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

func (s *MockChannelSession) Channel(channelID string, options ...discordgo.RequestOption) (*discordgo.Channel, error) {
	args := s.Called(channelID, options)
	return args.Get(0).(*discordgo.Channel), args.Error(1)
}

func (s *MockChannelSession) Guild(guildID string, options ...discordgo.RequestOption) (*discordgo.Guild, error) {
	args := s.Called(guildID, options)
	return args.Get(0).(*discordgo.Guild), args.Error(1)
}

type MockCoordinator struct {
	mock.Mock
}

func (m *MockCoordinator) GetMatches(query *queries.GetMatches) ([]types.Match, error) {
	args := m.Called(query)
	return args.Get(0).([]types.Match), args.Error(1)
}

func (m *MockCoordinator) GetTournaments(query *queries.GetTournaments) ([]types.Tournament, error) {
	args := m.Called(query)
	return args.Get(0).([]types.Tournament), args.Error(1)
}

func (m *MockCoordinator) GetUpcomingMatches(query *queries.GetUpcomingMatches) ([]types.Match, error) {
	args := m.Called(query)
	return args.Get(0).([]types.Match), args.Error(1)
}

func TestDotaBotChannel_SendMatchesOfTheDayInResponseTo(t *testing.T) {
	config := models.NewChannelConfig("test_id")
	channelSession := new(MockChannelSession)
	coordinator := new(MockCoordinator)
	botChannel := NewDotaBotChannelWithConfig(channelSession, config, coordinator)

	interactionCreate := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{},
	}

	matches := []types.Match{
		{
			BaseMatch: types.BaseMatch{
				ID:            0,
				Name:          "Tundra Esports vs Team Secret",
				TeamOne:       &types.Team{},
				TeamTwo:       &types.Team{},
				ScheduledTime: time.Now().UTC().Unix(),
				StreamUrl:     "https://www.twitch.tv/pgl_dota2",
			},
			Tournament: types.BaseTournament{
				ID:          0,
				DisplayName: "TI13 Regional Qualifiers WEU",
			},
		},
	}

	coordinator.On("GetMatches", mock.Anything).Return(matches, nil)
	channelSession.On("InteractionRespond", interactionCreate.Interaction, mock.Anything, mock.Anything).Return(nil)

	botChannel.SendMatchesOfTheDayInResponseTo(interactionCreate)

	coordinator.AssertExpectations(t)
	channelSession.AssertExpectations(t)
}

func TestGetTournamentTitle_ReturnsFullCombinedName_IfTournamentSerieAndLeagueHaveNames(t *testing.T) {
	match := types.Match{
		BaseMatch: types.BaseMatch{},
		Tournament: types.BaseTournament{
			DisplayName: "Tournament",
		},
		Serie: types.BaseSerie{
			Name: "Serie",
		},
		League: types.BaseLeague{
			Name: "League",
		},
	}

	title := getTournamentTitle(match)

	assert.Equal(t, "League: Serie - Tournament", title)
}

func TestGetTournamentTitle_ReturnsCombinedLeagueSerieName_IfOnlyLeagueAndSerieHaveNames(t *testing.T) {
	match := types.Match{
		BaseMatch: types.BaseMatch{},
		Tournament: types.BaseTournament{
			DisplayName: "",
		},
		Serie: types.BaseSerie{
			Name: "Serie",
		},
		League: types.BaseLeague{
			Name: "League",
		},
	}

	title := getTournamentTitle(match)

	assert.Equal(t, "League: Serie", title)
}

func TestGetTournamentTitle_ReturnsCombinedLeagueTournamentName_IfOnlyLeagueAndTournamentHaveNames(t *testing.T) {
	match := types.Match{
		BaseMatch: types.BaseMatch{},
		Tournament: types.BaseTournament{
			DisplayName: "Tournament",
		},
		Serie: types.BaseSerie{
			Name: "",
		},
		League: types.BaseLeague{
			Name: "League",
		},
	}

	title := getTournamentTitle(match)

	assert.Equal(t, "League - Tournament", title)
}

func TestGetTournamentTitle_ReturnsCombinedSerieTournamentName_IfOnlySerieAndTournamentHaveNames(t *testing.T) {
	match := types.Match{
		BaseMatch: types.BaseMatch{},
		Tournament: types.BaseTournament{
			DisplayName: "Tournament",
		},
		Serie: types.BaseSerie{
			Name: "Serie",
		},
		League: types.BaseLeague{
			Name: "",
		},
	}

	title := getTournamentTitle(match)

	assert.Equal(t, "Serie - Tournament", title)
}

func TestGetTournamentTitle_ReturnsLeagueName_IfOnlyLeagueHasName(t *testing.T) {
	match := types.Match{
		BaseMatch: types.BaseMatch{},
		Tournament: types.BaseTournament{
			DisplayName: "",
		},
		Serie: types.BaseSerie{
			Name: "",
		},
		League: types.BaseLeague{
			Name: "League",
		},
	}

	title := getTournamentTitle(match)

	assert.Equal(t, "League", title)
}

func TestGetTournamentTitle_ReturnsSerieName_IfOnlySerieHasName(t *testing.T) {
	match := types.Match{
		BaseMatch: types.BaseMatch{},
		Tournament: types.BaseTournament{
			DisplayName: "",
		},
		Serie: types.BaseSerie{
			Name: "",
		},
		League: types.BaseLeague{
			Name: "League",
		},
	}

	title := getTournamentTitle(match)

	assert.Equal(t, "League", title)
}

func TestGetTournamentTitle_ReturnsTournamentName_IfOnlyTournamentHasName(t *testing.T) {
	match := types.Match{
		BaseMatch: types.BaseMatch{},
		Tournament: types.BaseTournament{
			DisplayName: "Tournament",
		},
		Serie: types.BaseSerie{
			Name: "",
		},
		League: types.BaseLeague{
			Name: "",
		},
	}

	title := getTournamentTitle(match)

	assert.Equal(t, "Tournament", title)
}
