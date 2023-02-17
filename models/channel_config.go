package models

import (
	"fmt"
	"github.com/flusaka/dota-tournament-bot/stratz/schema"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type ChannelConfig struct {
	mgm.DefaultModel   `bson:",inline"`
	ChannelID          string              `bson:"channelID"`
	Timezone           string              `bson:"tz"`
	DailyMessageHour   int8                `bson:"dailyMessageHour"`
	DailyMessageMinute int8                `bson:"dailyMessageMinute"`
	Leagues            []schema.LeagueTier `bson:"leagues, omitempty"`
}

func NewChannelConfig(channelID string) *ChannelConfig {
	return &ChannelConfig{
		ChannelID: channelID,

		// Default to GMT
		Timezone: "GMT",

		// Default to DPC League, Majors and The International
		Leagues: []schema.LeagueTier{schema.LeagueTierDpcLeague, schema.LeagueTierMajor, schema.LeagueTierInternational},

		DailyMessageHour:   0,
		DailyMessageMinute: 0,
	}
}

func FetchAllConfigs() ([]*ChannelConfig, error) {
	configs := make([]*ChannelConfig, 0)
	err := mgm.Coll(&ChannelConfig{}).SimpleFind(&configs, bson.D{})
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (cc *ChannelConfig) Create() {
	err := mgm.Coll(cc).Create(cc)
	if err != nil {
		fmt.Println("Error when saving channel config", err)
	}
}

func (cc *ChannelConfig) Update() {
	err := mgm.Coll(cc).Update(cc)
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
