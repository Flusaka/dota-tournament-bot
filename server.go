package main

import (
	"github.com/flusaka/dota-tournament-bot/datasource/clients"
	"github.com/flusaka/dota-tournament-bot/stratz"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
		log.Println("No Discord token specified")
		return
	}
	if mongoUri == "" {
		log.Println("No Mongo URI specified")
		return
	}
	if stratzToken == "" {
		log.Println("No Stratz token specified")
		return
	}

	// Initialise Mongo
	err := mgm.SetDefaultConfig(nil, "bot", options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Println("Error connecting to MongoDB")
		return
	}

	stratzClient := stratz.NewClient(stratzToken)
	stratzClient.Initialise()

	dataSourceClient := clients.NewFakeDataSourceClient(true)
	//dataSourceClient := clients.NewStratzDataSourceClient(stratzClient)

	dotaBot := bot.NewDotaBotWithGuildID(dataSourceClient, guildID)
	err = dotaBot.Initialise(discordToken)
	if err != nil {
		log.Println("Error starting the Discord bot session")
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dotaBot.Shutdown()
}
