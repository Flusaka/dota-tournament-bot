package models

import "github.com/kamva/mgm/v3"

type ChannelConfig struct {
	mgm.DefaultModel `bson:",inline"`
	Timezone         string `bson:"tz"`
}

func NewChannelConfig() *ChannelConfig {
	return &ChannelConfig{
		Timezone: "Europe/London",
	}
}
