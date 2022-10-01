package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/flusaka/dota-tournament-bot/bot"
	"github.com/flusaka/dota-tournament-bot/command"
)

func main() {
	cp := command.NewParser("!dotabot")
	b := bot.NewDotaBot(cp)
	err := b.Initialise()
	if err != nil {
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	b.Shutdown()
}
