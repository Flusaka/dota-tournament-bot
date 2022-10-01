package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	dg *discordgo.Session
)

type DotaBot struct {
}

func (b *DotaBot) Initialise() error {
	dg, err := discordgo.New("Bot ODk4Njc0ODgxMTIwNTY3MzI3.Gfyw1i.R3cVZteQih5BWhVoO7lyRF4S9UXCGvqP4jfK-M")
	if err != nil {
		fmt.Println("Error creating Discord session", err)
		return err
	}

	dg.AddHandler(onMessage)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	return err
}

func (b *DotaBot) Shutdown() {
	dg.Close()
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	fmt.Println("Message received", m.Author.Username, m.Content)
}
