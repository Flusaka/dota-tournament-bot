package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flusaka/dota-tournament-bot/command"
	"github.com/flusaka/dota-tournament-bot/models"
)

type DotaBot struct {
	commandParser  *command.Parser
	discordSession *discordgo.Session
	channels       map[string]*DotaBotChannel
}

func NewDotaBot(commandParser *command.Parser) *DotaBot {
	b := new(DotaBot)
	b.commandParser = commandParser
	b.channels = make(map[string]*DotaBotChannel)
	return b
}

func (b *DotaBot) Initialise(token string) error {
	configs, err := models.FetchAllConfigs()
	if err != nil {
		fmt.Println("Could not retrieve configs", err)
	}

	// Setup existing bot channels
	for _, config := range configs {
		fmt.Println("Restarting channel on ID", config.ChannelID)
		b.channels[config.ChannelID] = NewDotaBotChannelWithConfig(config)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session", err)
		return err
	}

	b.commandParser.Register("start", func(params *command.ParseParameters) {
		fmt.Println("Start called with channel ID", params.ChannelID)
		if _, ok := b.channels[params.ChannelID]; ok {
			fmt.Println("Bot already started in this channel")
		} else {
			fmt.Println("Starting bot on channel", params.ChannelID, "with", len(params.Parameters), "parameters")
			channel := NewDotaBotChannel(params.ChannelID)
			channel.Start()
			b.channels[params.ChannelID] = channel
		}

	})

	b.commandParser.Register("stop", func(params *command.ParseParameters) {
		fmt.Println("Stopping bot on channel", params.ChannelID)
		if channel, ok := b.channels[params.ChannelID]; ok {
			channel.Stop()
		}
		delete(b.channels, params.ChannelID)
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return
		}

		if !b.commandParser.Parse(m.Message) {
			fmt.Println("Ignoring unparsed message")
		}
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()

	if err == nil {
		b.discordSession = dg
	} else {
		fmt.Println("Error when opening session", err)
	}
	return err
}

func (b *DotaBot) Shutdown() {
	err := b.discordSession.Close()
	if err != nil {
		fmt.Println("Error when closing Discord session", err)
	}
}
