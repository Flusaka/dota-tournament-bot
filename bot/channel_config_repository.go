package bot

import (
	"context"
	"github.com/flusaka/dota-tournament-bot/types"
)

type ChannelConfig interface {
	GetChannelID() string

	GetTimezone() string
	SetTimezone(timezone string)

	GetDailyMessageHour() int
	GetDailyMessageMinute() int
	SetDailyMessageTime(hour int, minute int)

	GetDailyNotificationsEnabled() bool
	SetDailyNotificationsEnabled(enabled bool)

	GetTiers() []types.Tier
	AddTier(tier types.Tier) bool
	RemoveTier(tier types.Tier)
}

type ChannelConfigRepository interface {
	GetAll(ctx context.Context) ([]ChannelConfig, error)
	Create(ctx context.Context, channelID string) (ChannelConfig, error)
	Update(ctx context.Context, config ChannelConfig) error
	Delete(ctx context.Context, config ChannelConfig) error
}
