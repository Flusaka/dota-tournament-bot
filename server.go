package main

import (
	"flag"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"os/signal"
	"syscall"

	"github.com/flusaka/dota-tournament-bot/bot"
	"github.com/flusaka/dota-tournament-bot/command"
)

func main() {
	discordToken := flag.String("t", "", "The token for the Discord Bot")
	mongoUri := flag.String("db", "", "The URI for the MongoDB database instance")
	flag.Parse()
	if *discordToken == "" {
		fmt.Println("No Discord token specified")
		return
	}

	// Initialise Mongo
	err := mgm.SetDefaultConfig(nil, "bot", options.Client().ApplyURI(*mongoUri))
	if err != nil {
		fmt.Println("Error connecting to MongoDB")
		return
	}

	cp := command.NewParser("!dotabot")
	b := bot.NewDotaBot(cp)
	err = b.Initialise(*discordToken)
	if err != nil {
		fmt.Println("Error starting the Discord bot session")
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	b.Shutdown()
}
