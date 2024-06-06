package bot

import (
	"context"
	"github.com/flusaka/dota-tournament-bot/types"
)

type ChannelConfig interface {
	GetChannelID() string
	GetTimezone() string
	GetDailyMessageHour() int
	GetDailyMessageMinute() int
	GetDailyNotificationsEnabled() bool
	GetTiers() []types.Tier
}

type ChannelConfigRepository interface {
	GetAll(ctx context.Context) ([]ChannelConfig, error)
	Create(ctx context.Context, channelID string) (ChannelConfig, error)
	Update(ctx context.Context, config ChannelConfig) error
	Delete(ctx context.Context, config ChannelConfig) error
}
