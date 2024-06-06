package repositories

import (
	"context"
	"fmt"
	"github.com/flusaka/dota-tournament-bot/bot"
	"github.com/flusaka/dota-tournament-bot/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	return nil, nil
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
	hexId := config.GetChannelID()
	id, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		return err
	}

	_, err = r.Collection.UpdateByID(ctx, id, config)
	if err != nil {
		return err
	}
	fmt.Println("Config updated!")
	return nil
}

func (r *MongoChannelConfigRepository) Delete(ctx context.Context, config bot.ChannelConfig) error {
	hexId := config.GetChannelID()
	id, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		return err
	}

	_, err = r.Collection.DeleteOne(ctx, bson.D{{
		"_id", id,
	}})
	if err != nil {
		return err
	}
	fmt.Println("Config deleted!")
	return nil
}
