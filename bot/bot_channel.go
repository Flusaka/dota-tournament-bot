package bot

import "github.com/flusaka/dota-tournament-bot/models"

type DotaBotChannel struct {
	config *models.ChannelConfig
}

func NewDotaBotChannel(channelID string) *DotaBotChannel {
	return &DotaBotChannel{
		models.NewChannelConfig(channelID),
	}
}

func (bc *DotaBotChannel) Start() {
	bc.config.Upsert()
}

func (bc *DotaBotChannel) Stop() {
	bc.config.Delete()
}
