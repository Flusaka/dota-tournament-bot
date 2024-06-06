package models

import (
	"github.com/flusaka/dota-tournament-bot/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChannelConfig struct {
	ID                        primitive.ObjectID `bson:"_id,omitempty"`
	ChannelID                 string             `bson:"channelID"`
	Timezone                  string             `bson:"tz"`
	DailyMessageHour          int                `bson:"dailyMessageHour"`
	DailyMessageMinute        int                `bson:"dailyMessageMinute"`
	DailyNotificationsEnabled bool               `bson:"dailyNotificationsEnabled"`
	Tiers                     []types.Tier       `bson:"tiers, omitempty"`
}

func NewChannelConfig(channelID string) *ChannelConfig {
	return &ChannelConfig{
		ChannelID: channelID,

		// Default to GMT
		Timezone: "GMT",

		// Default to S and A tiers
		Tiers: []types.Tier{types.TierS, types.TierA},

		DailyMessageHour:          0,
		DailyMessageMinute:        0,
		DailyNotificationsEnabled: false,
	}
}

func (c *ChannelConfig) GetID() string {
	return c.ID.String()
}

func (c *ChannelConfig) GetChannelID() string {
	return c.ChannelID
}

func (c *ChannelConfig) GetTimezone() string {
	return c.Timezone
}

func (c *ChannelConfig) GetDailyMessageHour() int {
	return c.DailyMessageHour
}

func (c *ChannelConfig) GetDailyMessageMinute() int {
	return c.DailyMessageMinute
}

func (c *ChannelConfig) GetDailyNotificationsEnabled() bool {
	return c.DailyNotificationsEnabled
}

func (c *ChannelConfig) GetTiers() []types.Tier {
	return c.Tiers
}
