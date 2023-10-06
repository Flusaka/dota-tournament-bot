package main

import (
	"fmt"
	"github.com/flusaka/dota-tournament-bot/datasource/clients"
	"github.com/flusaka/dota-tournament-bot/stratz"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"os/signal"
	"syscall"

	"github.com/flusaka/dota-tournament-bot/bot"
)

func main() {
	discordToken := os.Getenv("DISCORD_TOKEN")
	mongoUri := os.Getenv("MONGO_URI")
	stratzToken := os.Getenv("STRATZ_TOKEN")
	guildID := os.Getenv("GUILD_ID")

	if discordToken == "" {
		fmt.Println("No Discord token specified")
		return
	}
	if mongoUri == "" {
		fmt.Println("No Mongo URI specified")
		return
	}
	if stratzToken == "" {
		fmt.Println("No Stratz token specified")
		return
	}

	// Initialise Mongo
	err := mgm.SetDefaultConfig(nil, "bot", options.Client().ApplyURI(mongoUri))
	if err != nil {
		fmt.Println("Error connecting to MongoDB")
		return
	}

	stratzClient := stratz.NewClient(stratzToken)
	stratzClient.Initialise()

	dataSourceClient := clients.NewStratzDataSourceClient(stratzClient)

	dotaBot := bot.NewDotaBotWithGuildID(dataSourceClient, guildID)
	err = dotaBot.Initialise(discordToken)
	if err != nil {
		fmt.Println("Error starting the Discord bot session")
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dotaBot.Shutdown()
}
