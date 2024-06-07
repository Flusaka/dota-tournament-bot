package repositories

import (
	"context"
	"fmt"
	"github.com/flusaka/dota-tournament-bot/bot"
	"github.com/flusaka/dota-tournament-bot/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoChannelConfigRepository struct {
	Collection *mongo.Collection
}

func NewMongoChannelConfigRepository(database *mongo.Database) *MongoChannelConfigRepository {
	return &MongoChannelConfigRepository{
		Collection: database.Collection("channel_configs"),
	}
}

func (r *MongoChannelConfigRepository) GetAll(ctx context.Context) ([]bot.ChannelConfig, error) {
	filter := bson.D{{}}
	opts := options.Find().SetSort(bson.D{{
		"createdAt", -1,
	}})
	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []bot.ChannelConfig
	for cursor.Next(ctx) {
		var result models.ChannelConfig
		if decodeErr := cursor.Decode(&result); decodeErr != nil {
			continue
		}
		configs = append(configs, &result)
	}
	return configs, nil
}

func (r *MongoChannelConfigRepository) Create(ctx context.Context, channelID string) (bot.ChannelConfig, error) {
	config := models.NewChannelConfig(channelID)
	result, err := r.Collection.InsertOne(ctx, config)
	if err != nil {
		fmt.Printf("Error when creating config %v\n", err)
		return nil, err
	}
	fmt.Printf("Config created %v\n", result.InsertedID)
	return config, nil
}

func (r *MongoChannelConfigRepository) Update(ctx context.Context, config bot.ChannelConfig) error {
	channelId := config.GetChannelID()

	_, err := r.Collection.UpdateOne(ctx, bson.D{{
		"channelID", channelId,
	}}, config)
	if err != nil {
		return err
	}
	fmt.Println("Config updated!")
	return nil
}

func (r *MongoChannelConfigRepository) Delete(ctx context.Context, config bot.ChannelConfig) error {
	channelID := config.GetChannelID()

	_, err := r.Collection.DeleteOne(ctx, bson.D{{
		"channelID", channelID,
	}})
	if err != nil {
		return err
	}
	fmt.Println("Config deleted!")
	return nil
}
