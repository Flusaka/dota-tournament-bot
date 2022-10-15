package bot

import "github.com/flusaka/dota-tournament-bot/models"

type DotaBotChannel struct {
	config *models.ChannelConfig
}

func NewDotaBotChannel(config *models.ChannelConfig) *DotaBotChannel {
	return &DotaBotChannel{
		config,
	}
}
