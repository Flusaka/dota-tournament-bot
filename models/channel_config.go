package models

import (
	"github.com/flusaka/dota-tournament-bot/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ChannelConfig struct {
	ID                        primitive.ObjectID `bson:"_id,omitempty"`
	ChannelID                 string             `bson:"channelID"`
	Timezone                  string             `bson:"tz"`
	DailyMessageHour          int                `bson:"dailyMessageHour"`
	DailyMessageMinute        int                `bson:"dailyMessageMinute"`
	DailyNotificationsEnabled bool               `bson:"dailyNotificationsEnabled"`
	Tiers                     []types.Tier       `bson:"tiers, omitempty"`
	CreatedAt                 time.Time          `bson:"createdAt"`
	UpdatedAt                 time.Time          `json:"updatedAt"`
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

		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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

func (c *ChannelConfig) SetTimezone(timezone string) {
	// TODO: Validate time zone
	c.Timezone = timezone
	c.UpdatedAt = time.Now().UTC()
}

func (c *ChannelConfig) GetDailyMessageHour() int {
	return c.DailyMessageHour
}

func (c *ChannelConfig) GetDailyMessageMinute() int {
	return c.DailyMessageMinute
}

func (c *ChannelConfig) SetDailyMessageTime(hour int, minute int) {
	c.DailyMessageHour = hour
	c.DailyMessageMinute = minute
	c.UpdatedAt = time.Now().UTC()
}

func (c *ChannelConfig) GetDailyNotificationsEnabled() bool {
	return c.DailyNotificationsEnabled
}

func (c *ChannelConfig) SetDailyNotificationsEnabled(enabled bool) {
	c.DailyNotificationsEnabled = enabled
	c.UpdatedAt = time.Now().UTC()
}

func (c *ChannelConfig) GetTiers() []types.Tier {
	return c.Tiers
}

func (c *ChannelConfig) AddTier(tier types.Tier) bool {
	// Check it doesn't already exist, and if it does, return false
	for _, existingTier := range c.Tiers {
		if existingTier == tier {
			return false
		}
	}
	// Otherwise append it and return true
	c.Tiers = append(c.Tiers, tier)
	c.UpdatedAt = time.Now().UTC()
	return true
}

func (c *ChannelConfig) RemoveTier(tier types.Tier) bool {
	var tiers []types.Tier
	for _, existingTier := range c.Tiers {
		if existingTier != tier {
			tiers = append(tiers, existingTier)
		}
	}
	wasRemoved := len(tiers) < len(c.Tiers)
	c.Tiers = tiers
	if wasRemoved {
		c.UpdatedAt = time.Now().UTC()
	}
	return wasRemoved
}
