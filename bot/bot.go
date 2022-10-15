package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/flusaka/dota-tournament-bot/command"
)

type DotaBot struct {
	commandParser  *command.Parser
	discordSession *discordgo.Session
}

func NewDotaBot(commandParser *command.Parser) *DotaBot {
	b := new(DotaBot)
	b.commandParser = commandParser
	return b
}

func (b *DotaBot) Initialise(token string) error {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session", err)
		return err
	}

	b.commandParser.Register("start", func(params *command.ParseParameters) {
		fmt.Println("Starting bot on channel", params.ChannelID, "with", len(params.Parameters), "parameters")

		// Open channel???
	})

	b.commandParser.Register("stop", func(params *command.ParseParameters) {

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
