package bot

import (
	"github.com/flusaka/dota-tournament-bot/models"
	"time"
)

type DotaBotChannel struct {
	config *models.ChannelConfig
}

func NewDotaBotChannel(channelID string) *DotaBotChannel {
	return &DotaBotChannel{
		models.NewChannelConfig(channelID),
	}
}

func NewDotaBotChannelWithConfig(config *models.ChannelConfig) *DotaBotChannel {
	return &DotaBotChannel{
		config,
	}
}

func (bc *DotaBotChannel) Start() {
	bc.config.Upsert()
}

func (bc *DotaBotChannel) Stop() {
	bc.config.Delete()
}

func (bc *DotaBotChannel) UpdateTimezone(timezone string) error {
	_, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}
	
	bc.config.Timezone = timezone
	bc.config.Upsert()

	return nil
}
