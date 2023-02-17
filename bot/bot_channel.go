package bot

import (
	"github.com/flusaka/dota-tournament-bot/models"
	"github.com/flusaka/dota-tournament-bot/stratz/schema"
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
	bc.config.Create()
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
	bc.config.Update()

	return nil
}

func (bc *DotaBotChannel) UpdateDailyMessageTime(timeString string) error {
	activeTimeZone, err := time.LoadLocation(bc.config.Timezone)
	parsingZone := time.FixedZone(time.Now().In(activeTimeZone).Zone())
	dailyTime, err := time.ParseInLocation(time.Kitchen, timeString, parsingZone)

	if err != nil {
		return err
	}

	dailyTimeUTC := dailyTime.UTC()
	bc.config.DailyMessageHour = int8(dailyTimeUTC.Hour())
	bc.config.DailyMessageMinute = int8(dailyTimeUTC.Minute())
	bc.config.Update()

	return nil
}

func (bc *DotaBotChannel) GetLeagues() []schema.LeagueTier {
	return bc.config.Leagues
}

func (bc *DotaBotChannel) AddLeague(league schema.LeagueTier) bool {
	// TODO: Maybe change to error response
	for _, existingLeague := range bc.config.Leagues {
		if existingLeague == league {
			return false
		}
	}
	bc.config.Leagues = append(bc.config.Leagues, league)
	bc.config.Update()
	return true
}

func (bc *DotaBotChannel) RemoveLeague(league schema.LeagueTier) {
	var leagues []schema.LeagueTier
	for _, existingLeague := range bc.config.Leagues {
		if existingLeague != league {
			leagues = append(leagues, existingLeague)
		}
	}
	bc.config.Leagues = leagues
	bc.config.Update()
}

func (bc *DotaBotChannel) GetTimeInZone(timestamp int64) (time.Time, error) {
	parsingZone, err := bc.getParsingZone()
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(timestamp, 0).In(parsingZone), nil
}

func (bc *DotaBotChannel) IsTimeWithinDay(timestamp int64) bool {
	// This is a bit awkward, need to think of a better way to break down this logic
	parsingZone, err := bc.getParsingZone()
	if err != nil {
		return false
	}
	convertedTime, err := bc.GetTimeInZone(timestamp)
	if err != nil {
		return false
	}

	currentTime := time.Now()
	startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, parsingZone)
	endOfDay := startOfDay.Add(time.Hour * 24).Add(-time.Second)
	return convertedTime.After(startOfDay) && convertedTime.Before(endOfDay)
}

func (bc *DotaBotChannel) getParsingZone() (*time.Location, error) {
	activeTimeZone, err := time.LoadLocation(bc.config.Timezone)
	if err != nil {
		return nil, err
	}
	currentTime := time.Now().In(activeTimeZone)
	return time.FixedZone(currentTime.Zone()), nil
}
