package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/flusaka/dota-tournament-bot/datasource"
	"github.com/flusaka/dota-tournament-bot/datasource/queries"
	"github.com/flusaka/dota-tournament-bot/datasource/types"
	"github.com/flusaka/dota-tournament-bot/models"
	"log"
	"sort"
	"strconv"
	"strings"
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
	UnknownStreamKey = "UnknownStream"
	UnknownStreamUrl = "https://twitch.tv (Channel Unknown)"

	NotificationSelectMenuID = "notificationSelectMenu"

	NotificationValueDelimiter = ","
)

var (
	timezoneShortcodeFullMap = map[string]string{
		"GMT": "Europe/London",
		"EET": "Europe/Helsinki",
	}
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

type LeagueMatchesSet struct {
	League  *types.League
	Matches map[string][]*types.Match
}

type DotaBotChannel struct {
	session                    *discordgo.Session
	config                     *models.ChannelConfig
	dataSourceClient           datasource.Client
	notificationTicker         *time.Ticker
	cancelDailyNotifications   chan bool
	cancelMatchEventsListening chan bool
	channel                    *discordgo.Channel
	matchEventNotifier         *MatchEventNotifier
	cachedMatches              map[int]map[int16]*types.Match
}

func NewDotaBotChannel(session *discordgo.Session, channelID string, dataSourceClient datasource.Client) *DotaBotChannel {
	initialConfig := models.NewChannelConfig(channelID)
	cancelMatchEventsListening := make(chan bool, 1)
	dotaBotChannel := &DotaBotChannel{
		session:                    session,
		config:                     initialConfig,
		dataSourceClient:           dataSourceClient,
		matchEventNotifier:         NewMatchEventNotifier(cancelMatchEventsListening),
		cancelMatchEventsListening: cancelMatchEventsListening,
	}
	dotaBotChannel.listenForMatchEvents()
	return dotaBotChannel
}

func NewDotaBotChannelWithConfig(session *discordgo.Session, config *models.ChannelConfig, dataSourceClient datasource.Client) *DotaBotChannel {
	cancelMatchEventsListening := make(chan bool, 1)
	dotaBotChannel := &DotaBotChannel{
		session:                    session,
		config:                     config,
		dataSourceClient:           dataSourceClient,
		matchEventNotifier:         NewMatchEventNotifier(cancelMatchEventsListening),
		cancelMatchEventsListening: cancelMatchEventsListening,
	}

	if config.DailyNotificationsEnabled {
		dotaBotChannel.startDailyNotifications()
	}

	dotaBotChannel.listenForMatchEvents()
	return dotaBotChannel
}

func (bc *DotaBotChannel) Start() {
	bc.config.Create()
}

func (bc *DotaBotChannel) Stop() {
	close(bc.cancelMatchEventsListening)
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
	log.Printf("Restarting daily notifications on %v", bc.getChannelIdentifier())
	timeUntilNotification, err := bc.calculateTimeUntilNextNotification()
	if err != nil {
		return
	}

	bc.cancelDailyNotifications = make(chan bool)
	bc.notificationTicker = time.NewTicker(timeUntilNotification)

	go func() {
		for {
			select {
			case <-bc.notificationTicker.C:
				log.Printf("Ticker alarm activated! Attempting to send matches of the day to channel ID %v", bc.getChannelIdentifier())
				bc.SendMatchesOfTheDay()

				// Reset the daily notifications to 24 hours on from this ticker event
				bc.notificationTicker.Reset(24 * time.Hour)
			case <-bc.cancelDailyNotifications:
				return
			}
		}
	}()
}

func (bc *DotaBotChannel) stopDailyNotifications() {
	log.Printf("Daily notifications stopped for channel ID %v", bc.getChannelIdentifier())
	close(bc.cancelDailyNotifications)
	bc.notificationTicker = nil
}

func (bc *DotaBotChannel) calculateTimeUntilNextNotification() (time.Duration, error) {
	// Calculate time until next daily notifications should be sent
	now := time.Now()
	zone, err := bc.getParsingZone()
	if err != nil {
		return 0, err
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
	return timeUntilNotification, nil
}

func (bc *DotaBotChannel) SendMatchesOfTheDayInResponseTo(interaction *discordgo.InteractionCreate) {
	startingHour := 0
	startingMinute := 0

	if bc.config.DailyNotificationsEnabled {
		startingHour = bc.config.DailyMessageHour
		startingMinute = bc.config.DailyMessageMinute
	}
	result, leagueMatchesSet := bc.getMatchesToday(startingHour, startingMinute, true)

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
					minValues := 1
					matchNotificationOptions := bc.buildNotificationSelectionOptions(leagueMatches)
					if !interactionRespondedTo {
						err := bc.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fullMessage,
								Flags:   discordgo.MessageFlagsSuppressEmbeds,
								Components: []discordgo.MessageComponent{
									discordgo.ActionsRow{
										Components: []discordgo.MessageComponent{
											discordgo.SelectMenu{
												CustomID:    NotificationSelectMenuID,
												Placeholder: "Select matches to be notified of",
												MenuType:    discordgo.StringSelectMenu,
												MinValues:   &minValues,
												MaxValues:   len(matchNotificationOptions),
												Options:     matchNotificationOptions,
											},
										},
									},
								},
							},
						})
						// If there was no error responding to the interaction, it's been responded to
						interactionRespondedTo = err == nil
						if err != nil {
							log.Printf("Error responding to /today %v", err)
						}
					} else {
						// Otherwise, just send a regular message
						// Send the message and get the Discord message struct back
						messageSend := &discordgo.MessageSend{
							Content: fullMessage,
							Components: []discordgo.MessageComponent{
								discordgo.ActionsRow{
									Components: []discordgo.MessageComponent{
										discordgo.SelectMenu{
											CustomID:    NotificationSelectMenuID,
											Placeholder: "Select matches to be notified of",
											MenuType:    discordgo.StringSelectMenu,
											MinValues:   &minValues,
											MaxValues:   len(matchNotificationOptions),
											Options:     matchNotificationOptions,
										},
									},
								},
							},
						}
						discordMsg, err := bc.session.ChannelMessageSendComplex(bc.config.ChannelID, messageSend)
						if err != nil {
							log.Println("Error sending message to", bc.getChannelIdentifier(), err.Error())
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
	case ChannelResponseFailedToRetrieveLeagues:
		{
			log.Printf("Failed to retrieve leagues from DataSource, channel %v", bc.getChannelIdentifier())
			break
		}
	case ChannelResponseNoLeagues:
		{
			bc.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "No leagues have been set up on this channel yet!",
					Flags:   discordgo.MessageFlagsSuppressEmbeds,
				},
			})
			break
		}
	}
}

func (bc *DotaBotChannel) SendMatchesOfTheDay() ChannelResponse {
	result, leagueMatchesSet := bc.getMatchesToday(bc.config.DailyMessageHour, bc.config.DailyMessageMinute, true)

	switch result {
	case ChannelResponseSuccess:
		log.Printf("League matches retrieved for channel %v, number of leagues %v", bc.getChannelIdentifier(), len(leagueMatchesSet))
		for _, leagueMatches := range leagueMatchesSet {
			// Build up the message
			message := ":robot: " + leagueMatches.League.DisplayName + " games today!\n\n"
			matchesMessage := bc.generateDailyMatchMessage(leagueMatches)

			if len(matchesMessage) > 0 {
				// Send the message and get the Discord message struct back
				discordMsg, err := bc.session.ChannelMessageSend(bc.config.ChannelID, message+matchesMessage)
				if err != nil {
					log.Println("Error sending message to", bc.getChannelIdentifier(), err.Error())
				} else {
					// Suppress the embeds on the message from the stream links
					editMessage := discordgo.NewMessageEdit(bc.config.ChannelID, discordMsg.ID)
					editMessage.Flags |= discordgo.MessageFlagsSuppressEmbeds
					bc.session.ChannelMessageEditComplex(editMessage)
				}
			} else {
				log.Printf("Matches message is empty for some reason, channel %v", bc.getChannelIdentifier())
			}
		}
		break
	case ChannelResponseNoMatches:
		log.Printf("Failed to retrieve leagues from DataSource for daily notification for channel %v", bc.getChannelIdentifier())
		break
	case ChannelResponseFailedToRetrieveLeagues:
		log.Printf("Failed to retrieve leagues from DataSource for daily notification for channel %v", bc.getChannelIdentifier())
		break
	case ChannelResponseNoLeagues:
		log.Printf("Channel %v doesn't have any leagues configured for daily notification", bc.getChannelIdentifier())
		break
	}
	return result
}

func (bc *DotaBotChannel) HandleMessageComponentInteraction(interaction *discordgo.InteractionCreate) {
	messageComponentData := interaction.MessageComponentData()
	switch messageComponentData.CustomID {
	case NotificationSelectMenuID:
		{
			selectedValues := messageComponentData.Values
			for _, value := range selectedValues {
				// Split at delimiter to retrieve league and match ID
				split := strings.Split(value, NotificationValueDelimiter)
				leagueIDValue := split[0]
				matchIDValue := split[1]

				leagueIDParsed, _ := strconv.ParseInt(leagueIDValue, 10, 32)
				matchIDParsed, _ := strconv.ParseInt(matchIDValue, 10, 16)

				leagueID := int(leagueIDParsed)
				matchID := int16(matchIDParsed)

				// Find the cached match and add a notification for
				if league, ok := bc.cachedMatches[leagueID]; ok {
					if match, ok := league[matchID]; ok {
						bc.matchEventNotifier.AddUserToNotificationsForMatch(match, interaction.Member.User.ID)
					}
				}
			}

			bc.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
				Data: &discordgo.InteractionResponseData{
					Content: "Notifications have been updated!",
				},
			},
			)
			break
		}
	}
}

func (bc *DotaBotChannel) generateDailyMatchMessage(leagueMatches LeagueMatchesSet) string {
	message := ""
	if len(leagueMatches.Matches) == 0 {
		log.Printf("There are no matches in league %v", leagueMatches.League.DisplayName)
	}
	for streamUrl, streamMatches := range leagueMatches.Matches {
		if streamUrl == UnknownStreamKey {
			streamUrl = "https://twitch.tv (Channel Unknown)"
		}
		message += "Games on: " + streamUrl + "\n"
		for _, streamMatch := range streamMatches {
			convertedTime, err := bc.GetTimeInZone(streamMatch.ScheduledTime)
			if err != nil {
				continue
			}
			message += convertedTime.Format(time.Kitchen) + " - " + streamMatch.Radiant.DisplayName + " vs " + streamMatch.Dire.DisplayName + "\n"
		}
		message += "\n"
	}
	return message
}

func (bc *DotaBotChannel) buildNotificationSelectionOptions(leagueMatches LeagueMatchesSet) []discordgo.SelectMenuOption {
	var options []discordgo.SelectMenuOption
	for _, streamMatches := range leagueMatches.Matches {
		for _, match := range streamMatches {
			valueParts := []string{strconv.Itoa(leagueMatches.League.ID), strconv.Itoa(int(match.ID))}
			options = append(options, discordgo.SelectMenuOption{
				Label:       fmt.Sprintf("%s vs %s", match.Radiant.DisplayName, match.Dire.DisplayName),
				Value:       strings.Join(valueParts, NotificationValueDelimiter),
				Description: "",
				Emoji:       discordgo.ComponentEmoji{},
				Default:     false,
			})
		}
	}
	return options
}

func (bc *DotaBotChannel) getMatchesToday(startingHour int, startingMinute int, cache bool) (ChannelResponse, []LeagueMatchesSet) {
	// If there's no league tiers configured, let the channel know!
	tiers := bc.GetLeagues()
	if len(tiers) == 0 {
		return ChannelResponseNoLeagues, nil
	}

	query := queries.NewGetLeaguesQuery(bc.GetLeagues(), false)
	leagues, err := bc.dataSourceClient.GetLeagues(query)
	if err != nil {
		return ChannelResponseFailedToRetrieveLeagues, nil
	}

	if cache {
		// TODO: the logic to cache matches should be moved elsewhere
		bc.cachedMatches = make(map[int]map[int16]*types.Match, len(leagues))
		for _, league := range leagues {
			bc.cachedMatches[league.ID] = make(map[int16]*types.Match, len(league.Matches))
			for _, match := range league.Matches {
				bc.cachedMatches[league.ID][match.ID] = match
			}
		}
	}

	if len(leagues) == 0 {
		return ChannelResponseNoMatches, nil
	}

	var leagueMatches []LeagueMatchesSet

	// TODO: Definitely ways to improve and optimise this code :shrug:
	// Could probably cache these things every X amount of time
	for _, league := range leagues {
		leagueMatch := LeagueMatchesSet{
			League:  league,
			Matches: map[string][]*types.Match{},
		}

		// If there's no matches today for this league, skip over
		var matches []*types.Match
		for _, match := range league.Matches {
			// TODO: Check actual time if match already completed
			isWithinDay := bc.IsTimeWithinDayFrom(match.ScheduledTime, startingHour, startingMinute)
			if !isWithinDay {
				continue
			}
			matches = append(matches, match)
		}
		if len(matches) == 0 {
			continue
		}

		// Then, let's sort the matches by start time
		sort.Slice(matches, func(i, j int) bool {
			// TODO: Check actual time if match already completed
			return matches[i].ScheduledTime < matches[j].ScheduledTime
		})

		// Finally, make a map of streams to matches
		for _, match := range matches {
			// If the stream URL for the match is valid, use that as the key in the map and append to that array
			if match.StreamUrl != "" {
				leagueMatch.Matches[match.StreamUrl] = append(leagueMatch.Matches[match.StreamUrl], match)
			} else { // Otherwise just add it to the UnknownStreamKey array and key
				leagueMatch.Matches[UnknownStreamKey] = append(leagueMatch.Matches[UnknownStreamKey], match)
			}
		}

		leagueMatches = append(leagueMatches, leagueMatch)
	}

	if len(leagueMatches) == 0 {
		return ChannelResponseNoMatches, nil
	}
	return ChannelResponseSuccess, leagueMatches
}

func (bc *DotaBotChannel) UpdateTimezone(timezone string) error {
	actualTimezone, err := bc.getActualTimezoneLocation(timezone)
	if err != nil {
		return fmt.Errorf("timezone is unsupported")
	}

	_, err = time.LoadLocation(actualTimezone)
	if err != nil {
		return err
	}

	bc.config.Timezone = timezone
	bc.config.Update()

	// Update the existing timer, if it's not nil
	if bc.notificationTicker != nil {
		timeUntilNotification, err := bc.calculateTimeUntilNextNotification()
		if err == nil {
			bc.notificationTicker.Reset(timeUntilNotification)
		}
	}

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

	// Update the existing timer, if it's not nil
	if bc.notificationTicker != nil {
		timeUntilNotification, err := bc.calculateTimeUntilNextNotification()
		if err == nil {
			bc.notificationTicker.Reset(timeUntilNotification)
		}
	}

	return nil
}

func (bc *DotaBotChannel) GetLeagues() []types.Tier {
	return bc.config.Leagues
}

func (bc *DotaBotChannel) AddLeague(league types.Tier) bool {
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

func (bc *DotaBotChannel) RemoveLeague(league types.Tier) {
	var leagues []types.Tier
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

func (bc *DotaBotChannel) IsTimeWithinDayFrom(timestamp int64, startingHour int, startingMinute int) bool {
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
	startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), startingHour, startingMinute, 0, 0, parsingZone)
	endOfDay := startOfDay.Add(time.Hour * 24).Add(-time.Second)
	return convertedTime.After(startOfDay) && convertedTime.Before(endOfDay)
}

func (bc *DotaBotChannel) listenForMatchEvents() {
	go func() {
		for {
			select {
			case matchStarted := <-bc.matchEventNotifier.MatchStarted:
				{
					mentions := ""
					for _, user := range matchStarted.Users {
						mentions += fmt.Sprintf("<@%s> ", user)
					}
					streamUrl := matchStarted.Match.StreamUrl
					if streamUrl == "" {
						streamUrl = UnknownStreamUrl
					}
					bc.sendMessageWithoutEmbeds(
						mentions + fmt.Sprintf("%s vs %s is starting now on %s",
							matchStarted.Match.Radiant.DisplayName,
							matchStarted.Match.Dire.DisplayName,
							streamUrl,
						),
					)
					break
				}
			case <-bc.cancelMatchEventsListening:
				{
					return
				}
			}
		}
	}()
}

func (bc *DotaBotChannel) getParsingZone() (*time.Location, error) {
	actualTimezone, err := bc.getActualTimezoneLocation(bc.config.Timezone)
	activeTimeZone, err := time.LoadLocation(actualTimezone)
	if err != nil {
		return nil, err
	}
	currentTime := time.Now().In(activeTimeZone)
	return time.FixedZone(currentTime.Zone()), nil
}

func (bc *DotaBotChannel) getActualTimezoneLocation(timezone string) (string, error) {
	actualTimezone, timezoneExistsInMap := timezoneShortcodeFullMap[timezone]
	if !timezoneExistsInMap {
		return "", fmt.Errorf("timezone is unsupported")
	}
	return actualTimezone, nil
}

func (bc *DotaBotChannel) getChannel() (*discordgo.Channel, error) {
	if bc.channel == nil {
		channel, err := bc.session.Channel(bc.config.ChannelID)
		if err != nil {
			return nil, err
		}
		bc.channel = channel
	}
	return bc.channel, nil
}

func (bc *DotaBotChannel) getGuild(guildID string) (*discordgo.Guild, error) {
	guild, err := bc.session.Guild(guildID)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

func (bc *DotaBotChannel) getChannelIdentifier() string {
	channel, err := bc.getChannel()
	if err != nil {
		return bc.config.ChannelID
	}
	guild, err := bc.getGuild(channel.GuildID)
	if err != nil {
		return channel.Name
	}
	return fmt.Sprintf("%s:%s", guild.Name, channel.Name)
}

func (bc *DotaBotChannel) sendMessageWithoutEmbeds(messageContent string) {
	discordMsg, err := bc.session.ChannelMessageSend(bc.config.ChannelID, messageContent)
	if err != nil {
		log.Println("Error sending message to", bc.getChannelIdentifier(), err.Error())
	} else {
		// Suppress the embeds on the message from the stream links
		editMessage := discordgo.NewMessageEdit(bc.config.ChannelID, discordMsg.ID)
		editMessage.Flags |= discordgo.MessageFlagsSuppressEmbeds
		bc.session.ChannelMessageEditComplex(editMessage)
	}
}
