package bot

import (
	"github.com/flusaka/dota-tournament-bot/models"
	"time"
)

var (
	timeLayouts = []string{
		"0",
		"0:00",
		"00:00",
		"00:00AM",
		"0:00AM",
		"0AM",
		"00AM",
	}
)

type DotaBotChannel struct {
	config *models.ChannelConfig
}

func NewDotaBotChannel(channelID string) *DotaBotChannel {
	initialConfig := models.NewChannelConfig(channelID)
	return &DotaBotChannel{
		initialConfig,
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

func (bc *DotaBotChannel) UpdateDailyMessageTime(timeString string) error {
	/** Parse from string, possible formats:
		- 9:30
		- 15:45
		- 12PM
		- 1AM
	**/
	activeTimeZone, err := time.LoadLocation(bc.config.Timezone)
	parsingZone := time.FixedZone(time.Now().In(activeTimeZone).Zone())
	dailyTime, err := time.ParseInLocation(time.Kitchen, timeString, parsingZone)

	if err != nil {
		return err
	}

	dailyTimeUtc := dailyTime.UTC()

	bc.config.DailyMessageTime = dailyTimeUtc
	bc.config.Upsert()

	return nil
}

func (bc *DotaBotChannel) IsTimeWithinDay(timestamp int64) (bool, time.Time) {
	activeTimeZone, err := time.LoadLocation(bc.config.Timezone)
	convertedTime := time.Now()
	if err != nil {
		return false, convertedTime
	}
	currentTime := time.Now().In(activeTimeZone)
	parsingZone := time.FixedZone(currentTime.Zone())
	startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, parsingZone)
	endOfDay := startOfDay.Add(time.Hour * 24).Add(-time.Second)
	convertedTime = time.Unix(timestamp, 0).In(parsingZone)
	return convertedTime.After(startOfDay) && convertedTime.Before(endOfDay), convertedTime
}
