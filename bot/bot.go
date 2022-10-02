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

func (b *DotaBot) Initialise() error {
	dg, err := discordgo.New("Bot ODk4Njc0ODgxMTIwNTY3MzI3.Gfyw1i.R3cVZteQih5BWhVoO7lyRF4S9UXCGvqP4jfK-M")
	if err != nil {
		fmt.Println("Error creating Discord session", err)
		return err
	}

	b.commandParser.Register("start", func(params ...string) {
		fmt.Println("Start command parsed! Num params:", len(params))
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return
		}

		if !b.commandParser.Parse(m.Content) {
			fmt.Println("Ignoring unparsed message")
		}
	})
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err == nil {
		b.discordSession = dg
	}
	return err
}

func (b *DotaBot) Shutdown() {
	err := b.discordSession.Close()
	if err != nil {
		fmt.Println("Error when closing Discord session", err)
	}
}
