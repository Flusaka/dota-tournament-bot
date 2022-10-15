package models

import (
	"fmt"

	"github.com/kamva/mgm/v3"
)

type ChannelConfig struct {
	mgm.DefaultModel `bson:",inline"`
	ChannelID        string `bson:"channelID"`
	Timezone         string `bson:"tz"`
}

func NewChannelConfig(channelID string) *ChannelConfig {
	return &ChannelConfig{
		ChannelID: channelID,
		Timezone:  "Europe/London",
	}
}

func (cc *ChannelConfig) Upsert() {
	err := mgm.Coll(cc).Update(cc, mgm.UpsertTrueOption())
	if err != nil {
		fmt.Println("Error when saving channel config", err)
	}
}

func (cc *ChannelConfig) Delete() {
	err := mgm.Coll(cc).Delete(cc)
	if err != nil {
		fmt.Println("Error when deleting channel config", err)
	}
}
