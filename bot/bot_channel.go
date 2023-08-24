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

	// Unknown stream key
	UnknownStreamKey = "Stream Unknown"
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

type Team struct {
	DisplayName string
}

type Match struct {
	Radiant       Team
	Dire          Team
	ScheduledTime int64
	StreamUrl     string
}

type League struct {
	ID          int
	DisplayName string
}

type LeagueMatchesSet struct {
	League  League
	Matches map[string][]Match
}

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

func (bc *DotaBotChannel) SendMatchesOfTheDayInResponseTo(interaction *discordgo.InteractionCreate) {
	result, leagueMatchesSet := bc.getMatchesToday()

	switch result {
	case ChannelResponseSuccess:
		{
			interactionRespondedTo := false
			for _, leagueMatches := range leagueMatchesSet {
				// Build up the message
				message := ":robot: " + leagueMatches.League.DisplayName + " games today!\n\n"
				matchesMessage := bc.generateDailyMatchMessage(leagueMatches)

				if len(matchesMessage) > 0 {
					fullMessage := message + matchesMessage
					// If we haven't responded to the interaction yet, do that first
					if !interactionRespondedTo {
						err := bc.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fullMessage,
								Flags:   discordgo.MessageFlagsSuppressEmbeds,
							},
						})
						// If there was no error responding to the interaction, it's been responded to
						interactionRespondedTo = err == nil
					} else {
						// Otherwise, just send a regular message
						// Send the message and get the Discord message struct back
						discordMsg, err := bc.session.ChannelMessageSend(bc.config.ChannelID, fullMessage)
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
			}
			break
		}

	case ChannelResponseNoMatches:
		{
			bc.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: ":robot: No matches today!",
					Flags:   discordgo.MessageFlagsSuppressEmbeds,
				},
			})
			break
		}
	}
}

func (bc *DotaBotChannel) SendMatchesOfTheDay() ChannelResponse {
	result, leagueMatchesSet := bc.getMatchesToday()

	if result == ChannelResponseSuccess {
		for _, leagueMatches := range leagueMatchesSet {
			// Build up the message
			message := ":robot: " + leagueMatches.League.DisplayName + " games today!\n\n"
			matchesMessage := bc.generateDailyMatchMessage(leagueMatches)

			if len(matchesMessage) > 0 {
				// Send the message and get the Discord message struct back
				// TODO: Respond to the interaction somehow
				discordMsg, err := bc.session.ChannelMessageSend(bc.config.ChannelID, message+matchesMessage)
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
	}
	return result
}

func (bc *DotaBotChannel) generateDailyMatchMessage(leagueMatches LeagueMatchesSet) string {
	message := ""
	for streamUrl, streamMatches := range leagueMatches.Matches {
		message += "Games on: " + streamUrl + "\n\n"
		for _, streamMatch := range streamMatches {
			convertedTime, err := bc.GetTimeInZone(streamMatch.ScheduledTime)
			if err != nil {
				continue
			}
			message += streamMatch.Radiant.DisplayName + " vs " + streamMatch.Dire.DisplayName + " - " + convertedTime.Format(time.Kitchen) + "\n"
		}
	}
	return message
}

func (bc *DotaBotChannel) getMatchesToday() (ChannelResponse, []LeagueMatchesSet) {
	// If there's no league tiers configured, let the channel know!
	tiers := bc.GetLeagues()
	if len(tiers) == 0 {
		return ChannelResponseNoLeagues, nil
	}

	leagues, err := bc.stratzClient.GetActiveLeagues(tiers)
	if err != nil {
		return ChannelResponseFailedToRetrieveLeagues, nil
	}

	if len(leagues) == 0 {
		return ChannelResponseNoMatches, nil
	}

	var leagueMatches []LeagueMatchesSet

	// TODO: Definitely ways to improve and optimise this code :shrug:
	// Could probably cache these things every X amount of time
	for _, league := range leagues {
		leagueMatch := LeagueMatchesSet{
			League: League{
				ID:          league.Id,
				DisplayName: league.DisplayName,
			},
			Matches: map[string][]Match{},
		}
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
		streamMatchesMap := make(map[string][]Match)
		for _, match := range matches {
			matchToAdd := Match{
				Dire: Team{
					DisplayName: match.TeamTwo.Name,
				},
				Radiant: Team{
					DisplayName: match.TeamOne.Name,
				},
				ScheduledTime: match.ScheduledTime,
			}

			if len(match.Streams) > 0 {
				// Get English stream, if exists, if not :shrug:
				for _, stream := range match.Streams {
					if stream.LanguageId == schema.LanguageEnglish {
						matchToAdd.StreamUrl = stream.StreamUrl
						break
					}
				}
			}

			// If the stream URL for the match is valid, use that as the key in the map and append to that array
			if matchToAdd.StreamUrl != "" {
				streamMatchesMap[matchToAdd.StreamUrl] = append(streamMatchesMap[matchToAdd.StreamUrl], matchToAdd)
			} else { // Otherwise just add it to the UnknownStreamKey array and key
				streamMatchesMap[UnknownStreamKey] = append(streamMatchesMap[UnknownStreamKey], matchToAdd)
			}
		}

		leagueMatch.Matches = streamMatchesMap
		leagueMatches = append(leagueMatches, leagueMatch)
	}

	if len(leagueMatches) == 0 {
		return ChannelResponseNoMatches, nil
	}
	return ChannelResponseSuccess, leagueMatches
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
