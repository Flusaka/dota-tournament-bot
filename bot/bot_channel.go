package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/flusaka/dota-tournament-bot/models"
	"github.com/flusaka/dota-tournament-bot/stratz"
	"github.com/flusaka/dota-tournament-bot/stratz/schema"
	"sort"
	"time"
)

type ChannelResponse uint8

const (
	// Result constants from any bot method calls

	ChannelResponseSuccess                 ChannelResponse = 0
	ChannelResponseNoLeagues               ChannelResponse = 1
	ChannelResponseFailedToRetrieveLeagues ChannelResponse = 2
	ChannelResponseNoMatches               ChannelResponse = 3
	ChannelResponseFailedToSendToDiscord   ChannelResponse = 4
)

//var (
//	timeLayouts = []string{
//		"0",
//		"0:00",
//		"00:00",
//		"00:00AM",
//		"0:00AM",
//		"0AM",
//		"00AM",
//	}
//)

type DotaBotChannel struct {
	session                  *discordgo.Session
	config                   *models.ChannelConfig
	stratzClient             *stratz.Client
	cancelDailyNotifications chan bool
}

func NewDotaBotChannel(session *discordgo.Session, channelID string, stratzClient *stratz.Client) *DotaBotChannel {
	initialConfig := models.NewChannelConfig(channelID)
	return &DotaBotChannel{
		session,
		initialConfig,
		stratzClient,
		nil,
	}
}

func NewDotaBotChannelWithConfig(session *discordgo.Session, config *models.ChannelConfig, stratzClient *stratz.Client) *DotaBotChannel {
	dotaBotChannel := &DotaBotChannel{
		session,
		config,
		stratzClient,
		nil,
	}

	if config.DailyNotificationsEnabled {
		dotaBotChannel.startDailyNotifications()
	}

	return dotaBotChannel
}

func (bc *DotaBotChannel) Start() {
	bc.config.Create()
}

func (bc *DotaBotChannel) Stop() {
	bc.stopDailyNotifications()
	bc.config.Delete()
}

func (bc *DotaBotChannel) EnableDailyNotifications(enabled bool) ChannelResponse {
	if len(bc.config.Leagues) == 0 {
		return ChannelResponseNoLeagues
	}

	if enabled != bc.config.DailyNotificationsEnabled {
		if enabled {
			bc.startDailyNotifications()
		} else {
			bc.stopDailyNotifications()
		}

		bc.config.DailyNotificationsEnabled = enabled
		bc.config.Update()
	}
	return ChannelResponseSuccess
}

func (bc *DotaBotChannel) startDailyNotifications() {
	bc.cancelDailyNotifications = make(chan bool)

	// Calculate time until next daily notifications should be sent
	now := time.Now()
	zone, err := bc.getParsingZone()
	if err != nil {
		return
	}
	nowInTimezone := now.In(zone)

	// If notification time is still available today, then set notification day to today, otherwise tomorrow
	day := now.Day()
	if nowInTimezone.Hour() > bc.config.DailyMessageHour || (nowInTimezone.Hour() == bc.config.DailyMessageHour && nowInTimezone.Minute() >= bc.config.DailyMessageMinute) {
		day++
	}

	// Create a time object in the current location set to the notification time
	firstNotificationTime := time.Date(nowInTimezone.Year(), nowInTimezone.Month(), day, bc.config.DailyMessageHour, bc.config.DailyMessageMinute, 0, 0, zone)
	timeUntilNotification := firstNotificationTime.Sub(nowInTimezone)

	ticker := time.NewTicker(timeUntilNotification)

	go func() {
		for {
			select {
			case <-ticker.C:
				bc.SendMatchesOfTheDay()

				// Reset the daily notifications to 24 hours on from this ticker event
				ticker.Reset(24 * time.Hour)
			case <-bc.cancelDailyNotifications:
				return
			}
		}
	}()
}

func (bc *DotaBotChannel) stopDailyNotifications() {
	close(bc.cancelDailyNotifications)
}

func (bc *DotaBotChannel) SendMatchesOfTheDay() ChannelResponse {
	// If there's no league tiers configured, let the channel know!
	tiers := bc.GetLeagues()
	if len(tiers) == 0 {
		return ChannelResponseNoLeagues
	}

	leagues, err := bc.stratzClient.GetActiveLeagues(tiers)
	if err != nil {
		return ChannelResponseFailedToRetrieveLeagues
	}

	// TODO: Definitely ways to improve and optimise this code :shrug:
	// Could probably cache these things every X amount of time
	for _, league := range leagues {
		var allMatches []schema.GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType
		for _, nodeGroup := range league.NodeGroups {
			allMatches = append(allMatches, nodeGroup.Nodes...)
		}

		// First, let's remove all matches outside of today
		var matches []schema.GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType
		for _, match := range allMatches {
			// TODO: Check actual time if match already completed
			isWithinDay := bc.IsTimeWithinDay(match.ScheduledTime)
			if !isWithinDay {
				continue
			}
			matches = append(matches, match)
		}

		// If there's no matches today for this league, skip over
		if len(matches) == 0 {
			continue
		}

		// Then, let's sort the matches by start time
		sort.Slice(matches, func(i, j int) bool {
			// TODO: Check actual time if match already completed
			return matches[i].ScheduledTime < matches[j].ScheduledTime
		})

		// Finally, create a map of streams to matches
		streamMatchesMap := make(map[string][]schema.GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType)
		for _, match := range matches {
			// Get English stream, if exists, if not :shrug:
			var foundStream *schema.GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType
			for _, stream := range match.Streams {
				if stream.LanguageId == schema.LanguageEnglish {
					foundStream = &stream
					break
				}
			}

			if foundStream != nil {
				streamMatchesMap[foundStream.StreamUrl] = append(streamMatchesMap[foundStream.StreamUrl], match)
			}
		}

		// If the map is empty, we couldn't find a relevant stream, skip over
		if len(streamMatchesMap) == 0 {
			continue
		}

		// Build up the message
		message := ":robot: " + league.DisplayName + " games today!\n\n"
		for streamUrl, streamMatches := range streamMatchesMap {
			message += "Games on: " + streamUrl + "\n\n"
			for _, streamMatch := range streamMatches {
				convertedTime, err := bc.GetTimeInZone(streamMatch.ScheduledTime)
				if err != nil {
					continue
				}
				message += streamMatch.TeamOne.Name + " vs " + streamMatch.TeamTwo.Name + " - " + convertedTime.Format(time.Kitchen) + "\n"
			}
		}

		if len(message) > 0 {
			// Send the message and get the Discord message struct back
			// TODO: Respond to the interaction somehow
			discordMsg, err := bc.session.ChannelMessageSend(bc.config.ChannelID, message)
			if err != nil {
				fmt.Println("Error sending message to", bc.config.ChannelID, err.Error())
			} else {
				// Suppress the embeds on the message from the stream links
				editMessage := discordgo.NewMessageEdit(bc.config.ChannelID, discordMsg.ID)
				editMessage.Flags |= discordgo.MessageFlagsSuppressEmbeds
				bc.session.ChannelMessageEditComplex(editMessage)
			}
		}
	}
	return ChannelResponseSuccess
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

	bc.config.DailyMessageHour = dailyTime.Hour()
	bc.config.DailyMessageMinute = dailyTime.Minute()
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
