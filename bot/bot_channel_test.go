package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/flusaka/dota-tournament-bot/models"
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/flusaka/dota-tournament-bot/types"
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

	coordinator.On("GetUpcomingMatches", mock.Anything).Return(matches, nil)
	channelSession.On("InteractionRespond", interactionCreate.Interaction, mock.Anything, mock.Anything).Return(nil)

	botChannel.SendMatchesOfTheDayInResponseTo(interactionCreate)

	coordinator.AssertExpectations(t)
	channelSession.AssertExpectations(t)
}
