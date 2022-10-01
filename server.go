package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/flusaka/dota-tournament-bot/bot"
)

func main() {
	b := bot.DotaBot{}
	err := b.Initialise()
	if err != nil {
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	b.Shutdown()
}
