package main

import (
	"context"
	"github.com/flusaka/dota-tournament-bot/cache"
	"github.com/flusaka/dota-tournament-bot/coordinators"
	"github.com/flusaka/dota-tournament-bot/datasource"
	"github.com/flusaka/dota-tournament-bot/repositories"
	"github.com/flusaka/pandascore-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/flusaka/dota-tournament-bot/bot"
)

const (
	defaultQueryCacheTimeInMinutes = 5
)

func main() {
	discordToken := os.Getenv("DISCORD_TOKEN")
	mongoUri := os.Getenv("MONGO_URI")
	pandascoreToken := os.Getenv("PANDASCORE_TOKEN")
	guildID := os.Getenv("GUILD_ID")
	queryCacheTimeInMinutesEnv := os.Getenv("QUERY_CACHE_TIME_IN_MINUTES")

	if discordToken == "" {
		log.Println("No Discord token specified")
		return
	}
	if mongoUri == "" {
		log.Println("No Mongo URI specified")
		return
	}
	if pandascoreToken == "" {
		log.Println("No Pandascore token specified")
		return
	}
	queryCacheTimeInMinutes := defaultQueryCacheTimeInMinutes
	if queryCacheTimeInMinutesEnv != "" {
		var err error = nil
		queryCacheTimeInMinutes, err = strconv.Atoi(queryCacheTimeInMinutesEnv)
		if err != nil || queryCacheTimeInMinutes < 0 {
			queryCacheTimeInMinutes = defaultQueryCacheTimeInMinutes
		}
	}

	// Initialise Mongo
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverApi)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Println("Error connecting to MongoDB")
		return
	}

	pandascoreClient := pandascore.NewClient(pandascoreToken)
	dataSourceClient := datasource.NewPandascoreDataSource(pandascoreClient)

	defaultQueryCache := cache.NewDefaultQueryResultCache(time.Minute * time.Duration(queryCacheTimeInMinutes))
	queryCoordinator := coordinators.NewDefaultQueryCoordinator(dataSourceClient, defaultQueryCache)

	mongoChannelConfigRepository := repositories.NewMongoChannelConfigRepository(client.Database("bot"))

	dotaBot := bot.NewDotaBotWithGuildID(queryCoordinator, mongoChannelConfigRepository, guildID)
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
