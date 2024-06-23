package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/flusaka/dota-tournament-bot/types"
	"github.com/flusaka/dota-tournament-bot/utils"
	"log"
	"time"
)

type ChannelResponse uint8

const (
	// Result constants from any bot method calls
	ChannelResponseSuccess                 ChannelResponse = 0
	ChannelResponseNoTiers                 ChannelResponse = 1
	ChannelResponseFailedToRetrieveLeagues ChannelResponse = 2
	ChannelResponseNoMatches               ChannelResponse = 3
	ChannelResponseFailedToSendToDiscord   ChannelResponse = 4

	// Unknown stream key
	UnknownStreamKey = "UnknownStream"
	UnknownStreamUrl = "https://twitch.tv (Channel Unknown)"

	NotificationSelectMenuID = "notificationSelectMenu"

	NotificationValueDelimiter = ","
)

type StreamMatchMap = map[string][]types.Match
type TournamentMatchDetails struct {
	StreamMatches StreamMatchMap
	Title         string
}

type TournamentDetailsSet struct {
	Tournament types.BaseTournament
	Details    TournamentMatchDetails
}

type ChannelSession interface {
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error
	ChannelMessageSend(channelID string, message string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelMessageSendComplex(channelID string, message *discordgo.MessageSend, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelMessageEditComplex(m *discordgo.MessageEdit, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
	Channel(channelID string, options ...discordgo.RequestOption) (*discordgo.Channel, error)
	Guild(guildID string, options ...discordgo.RequestOption) (*discordgo.Guild, error)
}

type DotaBotChannel struct {
	session                    ChannelSession
	config                     ChannelConfig
	queryCoordinator           QueryCoordinator
	notificationTicker         *time.Ticker
	cancelDailyNotifications   chan bool
	cancelMatchEventsListening chan bool
	channel                    *discordgo.Channel
	matchEventNotifier         *MatchEventNotifier
	cachedMatches              map[int]map[int]*types.Match
}

func NewDotaBotChannelWithConfig(session ChannelSession, config ChannelConfig, queryCoordinator QueryCoordinator) *DotaBotChannel {
	cancelMatchEventsListening := make(chan bool, 1)
	dotaBotChannel := &DotaBotChannel{
		session:                    session,
		config:                     config,
		queryCoordinator:           queryCoordinator,
		matchEventNotifier:         NewMatchEventNotifier(cancelMatchEventsListening),
		cancelMatchEventsListening: cancelMatchEventsListening,
	}

	if config.GetDailyNotificationsEnabled() {
		dotaBotChannel.startDailyNotifications()
	}

	dotaBotChannel.listenForMatchEvents()
	return dotaBotChannel
}

func (bc *DotaBotChannel) Close() {
	close(bc.cancelMatchEventsListening)
	bc.stopDailyNotifications()
}

func (bc *DotaBotChannel) EnableDailyNotifications(enabled bool) {
	// TODO: Move this check to Bot
	//if len(bc.config.Tiers) == 0 {
	//	return ChannelResponseNoTiers
	//}
	//
	//if enabled != bc.config.DailyNotificationsEnabled {
	//	if enabled {
	//		bc.startDailyNotifications()
	//	} else {
	//		bc.stopDailyNotifications()
	//	}
	//
	//	bc.config.DailyNotificationsEnabled = enabled
	//	bc.config.Update()
	//}

	if enabled {
		bc.startDailyNotifications()
	} else {
		bc.stopDailyNotifications()
	}
}

func (bc *DotaBotChannel) startDailyNotifications() {
	log.Printf("Starting daily notifications on %v\n", bc.getChannelIdentifier())
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
				log.Printf("Ticker alarm activated! Attempting to send matches of the day to channel ID %v\n", bc.getChannelIdentifier())
				bc.sendMatchesOfTheDay()

				// Reset the daily notifications to 24 hours on from this ticker event
				bc.notificationTicker.Reset(24 * time.Hour)
			case <-bc.cancelDailyNotifications:
				return
			}
		}
	}()
}

func (bc *DotaBotChannel) stopDailyNotifications() {
	log.Printf("Daily notifications stopped for channel ID %v\n", bc.getChannelIdentifier())
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
	if nowInTimezone.Hour() > bc.config.GetDailyMessageHour() || (nowInTimezone.Hour() == bc.config.GetDailyMessageMinute() && nowInTimezone.Minute() >= bc.config.GetDailyMessageMinute()) {
		day++
	}

	// Create a time object in the current location set to the notification time
	firstNotificationTime := time.Date(nowInTimezone.Year(), nowInTimezone.Month(), day, bc.config.GetDailyMessageHour(), bc.config.GetDailyMessageMinute(), 0, 0, zone)
	timeUntilNotification := firstNotificationTime.Sub(nowInTimezone)
	return timeUntilNotification, nil
}

func (bc *DotaBotChannel) SendMatchesOfTheDayInResponseTo(interaction *discordgo.InteractionCreate) {
	startingHour := 0
	startingMinute := 0

	if bc.config.GetDailyNotificationsEnabled() {
		startingHour = bc.config.GetDailyMessageHour()
		startingMinute = bc.config.GetDailyMessageMinute()
	}
	result, tournamentMatchesSet := bc.getMatchesToday(startingHour, startingMinute, false)

	switch result {
	case ChannelResponseSuccess:
		{
			interactionRespondedTo := false
			for _, tournamentMatches := range tournamentMatchesSet {
				if len(tournamentMatches.Details.StreamMatches) == 0 {
					continue
				}

				// Build up the message
				message := ":robot: " + tournamentMatches.Details.Title + " games today!\n\n"
				matchesMessage := bc.generateDailyMatchMessage(tournamentMatches)

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
					if err != nil {
						log.Printf("Error responding to /today %v\n", err)
					}
				} else {
					// Otherwise, just send a regular message
					// Send the message and get the Discord message struct back
					discordMsg, err := bc.session.ChannelMessageSend(bc.config.GetChannelID(), fullMessage)
					if err != nil {
						log.Println("Error sending message to", bc.getChannelIdentifier(), err.Error())
					} else {
						// Suppress the embeds on the message from the stream links
						editMessage := discordgo.NewMessageEdit(bc.config.GetChannelID(), discordMsg.ID)
						editMessage.Flags |= discordgo.MessageFlagsSuppressEmbeds
						bc.session.ChannelMessageEditComplex(editMessage)
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
			log.Printf("Failed to retrieve leagues from DataSource, channel %v\n", bc.getChannelIdentifier())
			break
		}
	case ChannelResponseNoTiers:
		{
			bc.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "No tiers have been set up on this channel yet!",
					Flags:   discordgo.MessageFlagsSuppressEmbeds,
				},
			})
			break
		}
	}
}

func (bc *DotaBotChannel) sendMatchesOfTheDay() {
	result, tournamentDetailsSet := bc.getMatchesToday(bc.config.GetDailyMessageHour(), bc.config.GetDailyMessageMinute(), true)

	switch result {
	case ChannelResponseSuccess:
		log.Printf("Matches retrieved for channel %v, number of sets %v\n", bc.getChannelIdentifier(), len(tournamentDetailsSet))
		for _, tournamentMatches := range tournamentDetailsSet {
			// Build up the message
			message := ":robot: " + tournamentMatches.Details.Title + " games today!\n\n"
			matchesMessage := bc.generateDailyMatchMessage(tournamentMatches)

			if len(matchesMessage) > 0 {
				//minValues := 0
				//matchNotificationOptions := bc.buildNotificationSelectionOptions(tournamentMatches)

				// Send the message and get the Discord message struct back
				messageSend := &discordgo.MessageSend{
					Content: message + matchesMessage,
					//Components: []discordgo.MessageComponent{
					//	discordgo.ActionsRow{
					//		Components: []discordgo.MessageComponent{
					//			discordgo.SelectMenu{
					//				CustomID:    NotificationSelectMenuID,
					//				Placeholder: "Select matches to be notified of",
					//				MenuType:    discordgo.StringSelectMenu,
					//				MinValues:   &minValues,
					//				MaxValues:   len(matchNotificationOptions),
					//				Options:     matchNotificationOptions,
					//			},
					//		},
					//	},
					//},
				}
				discordMsg, err := bc.session.ChannelMessageSendComplex(bc.config.GetChannelID(), messageSend)
				if err != nil {
					log.Println("Error sending message to", bc.getChannelIdentifier(), err.Error())
				} else {
					// Suppress the embeds on the message from the stream links
					editMessage := discordgo.NewMessageEdit(bc.config.GetChannelID(), discordMsg.ID)
					editMessage.Flags |= discordgo.MessageFlagsSuppressEmbeds
					bc.session.ChannelMessageEditComplex(editMessage)
				}
			} else {
				log.Printf("Matches message is empty for some reason, channel %v\n", bc.getChannelIdentifier())
			}
		}
		break
	case ChannelResponseNoMatches:
		log.Printf("Failed to retrieve leagues from DataSource for daily notification for channel %v\n", bc.getChannelIdentifier())
		break
	case ChannelResponseFailedToRetrieveLeagues:
		log.Printf("Failed to retrieve leagues from DataSource for daily notification for channel %v\n", bc.getChannelIdentifier())
		break
	case ChannelResponseNoTiers:
		log.Printf("Channel %v doesn't have any leagues configured for daily notification\n", bc.getChannelIdentifier())
		break
	}
}

func (bc *DotaBotChannel) HandleMessageComponentInteraction(interaction *discordgo.InteractionCreate) {
	//messageComponentData := interaction.MessageComponentData()
	//switch messageComponentData.CustomID {
	//case NotificationSelectMenuID:
	//	{
	//		subscribedMatches := bc.matchEventNotifier.GetSubscribedMatchesForUser(interaction.Member.User.ID)
	//
	//		selectedValues := messageComponentData.Values
	//		for _, value := range selectedValues {
	//			// Split at delimiter to retrieve league and match ID
	//			split := strings.Split(value, NotificationValueDelimiter)
	//			leagueIDValue := split[0]
	//			matchIDValue := split[1]
	//
	//			leagueIDParsed, _ := strconv.ParseInt(leagueIDValue, 10, 32)
	//			matchIDParsed, _ := strconv.ParseInt(matchIDValue, 10, 16)
	//
	//			leagueID := int(leagueIDParsed)
	//			matchID := int16(matchIDParsed)
	//
	//			// Find the cached match and add a notification for
	//			if league, ok := bc.cachedMatches[leagueID]; ok {
	//				if match, ok := league[matchID]; ok {
	//					bc.matchEventNotifier.AddUserToNotificationsForMatch(match, interaction.Member.User.ID)
	//				}
	//			}
	//
	//			subscribedMatches = slices.DeleteFunc(subscribedMatches, func(el *types.Match) bool {
	//				return el.ID == matchID
	//			})
	//		}
	//
	//		// Any leftover subscribed matches are ones that were removed, so we should delete our subscription for notifications
	//		for _, subscription := range subscribedMatches {
	//			log.Printf("Removing notification subscription for match %d for user %s", subscription.ID, interaction.Member.User.ID)
	//			bc.matchEventNotifier.RemoveUserFromNotificationsForMatch(subscription, interaction.Member.User.ID)
	//		}
	//
	//		bc.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseDeferredMessageUpdate,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: "Notifications have been updated!",
	//			},
	//		},
	//		)
	//		break
	//	}
	//}
}

func (bc *DotaBotChannel) generateDailyMatchMessage(tournamentMatches TournamentDetailsSet) string {
	message := ""
	if len(tournamentMatches.Details.StreamMatches) == 0 {
		log.Printf("There are no matches in tournament %v\n", tournamentMatches.Details.Title)
	}
	for streamUrl, streamMatches := range tournamentMatches.Details.StreamMatches {
		if streamUrl == UnknownStreamKey {
			streamUrl = "https://twitch.tv (Channel Unknown)"
		}
		message += "Games on: " + streamUrl + "\n"
		for _, streamMatch := range streamMatches {
			startTime := fmt.Sprintf("<t:%d:t>", streamMatch.ScheduledTime)

			// If TeamOne is undetermined, use the TeamOneSourceMatch field to determine the teams to display
			teamOneComponent := bc.generateTeamMessageComponent(streamMatch.TeamOne, nil)
			teamTwoComponent := bc.generateTeamMessageComponent(streamMatch.TeamTwo, nil)
			message += startTime + " - " + teamOneComponent + " vs " + teamTwoComponent + "\n"
		}
		message += "\n"
	}
	return message
}

func (bc *DotaBotChannel) generateTeamMessageComponent(team *types.Team, teamSourceMatch *types.Match) string {
	if team != nil {
		return team.DisplayName
	} else if teamSourceMatch != nil && teamSourceMatch.TeamOne != nil && teamSourceMatch.TeamTwo != nil {
		return fmt.Sprintf("%s/%s", teamSourceMatch.TeamOne.DisplayName, teamSourceMatch.TeamTwo.DisplayName)
	}

	return "TBD"
}

func (bc *DotaBotChannel) buildNotificationSelectionOptions(tournamentMatches TournamentDetailsSet) []discordgo.SelectMenuOption {
	var options []discordgo.SelectMenuOption
	//for _, streamMatches := range leagueMatches.Matches {
	//	for _, match := range streamMatches {
	//		valueParts := []string{strconv.Itoa(leagueMatches.League.ID), strconv.Itoa(int(match.ID))}
	//
	//		teamOneComponent := bc.generateTeamMessageComponent(match.TeamOne, match.TeamOneSourceMatch)
	//		teamTwoComponent := bc.generateTeamMessageComponent(match.TeamTwo, match.TeamTwoSourceMatch)
	//
	//		options = append(options, discordgo.SelectMenuOption{
	//			Label:       fmt.Sprintf("%s vs %s", teamOneComponent, teamTwoComponent),
	//			Value:       strings.Join(valueParts, NotificationValueDelimiter),
	//			Description: "",
	//			Emoji:       discordgo.ComponentEmoji{},
	//			Default:     false,
	//		})
	//	}
	//}
	return options
}

func (bc *DotaBotChannel) getMatchesToday(startingHour int, startingMinute int, cache bool) (ChannelResponse, []TournamentDetailsSet) {
	// If there's no tournament tiers configured, let the channel know!
	tiers := bc.getTiers()
	if len(tiers) == 0 {
		return ChannelResponseNoTiers, nil
	}

	parsingZone, err := bc.getParsingZone()
	if err != nil {
		return ChannelResponseNoMatches, nil
	}

	currentTime := time.Now()
	startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), startingHour, startingMinute, 0, 0, parsingZone)
	endOfDay := startOfDay.Add(time.Hour * 24).Add(-time.Second)

	query := &queries.GetMatches{
		BeginAt: queries.DateRange{
			Start: startOfDay,
			End:   endOfDay,
		},
		Tiers: bc.getTiers(),
	}
	upcomingMatches, err := bc.queryCoordinator.GetMatches(query)
	if err != nil {
		return ChannelResponseFailedToRetrieveLeagues, nil
	}

	//if cache {
	//	// TODO: the logic to cache matches should be moved elsewhere
	//	bc.cachedMatches = make(map[int]map[int]*types.Match, len(leagues))
	//	for _, tournament := range leagues {
	//		bc.cachedMatches[tournament.ID] = make(map[int16]*types.Match, len(tournament.Matches))
	//		for _, match := range tournament.Matches {
	//			bc.cachedMatches[tournament.ID][match.ID] = match
	//		}
	//	}
	//}

	if len(upcomingMatches) == 0 {
		return ChannelResponseNoMatches, nil
	}

	// Create a map of tournaments to matches first, for the sake of ease...
	tournamentMatchesMap := make(map[int][]types.Match)
	for _, match := range upcomingMatches {
		_, found := tournamentMatchesMap[match.Tournament.ID]
		if !found {
			tournamentMatchesMap[match.Tournament.ID] = []types.Match{match}
		} else {
			tournamentMatchesMap[match.Tournament.ID] = append(tournamentMatchesMap[match.Tournament.ID], match)
		}
	}

	var tournamentMatches []TournamentDetailsSet

	// TODO: Definitely ways to improve and optimise this code :shrug:
	// Could probably cache these things every X amount of time
	for _, matches := range tournamentMatchesMap {
		// If there's no matches today for this tournament, skip over
		if len(matches) == 0 {
			continue
		}

		tournamentMatchesSet := TournamentDetailsSet{
			Tournament: matches[0].Tournament,
			Details: TournamentMatchDetails{
				StreamMatches: StreamMatchMap{},
				Title:         fmt.Sprintf("%v: %v - %v", matches[0].League.Name, matches[0].Serie.Name, matches[0].Tournament.DisplayName),
			},
		}

		// NOTE: Shouldn't need this anymore as the matches are sorted by API when queried
		// Then, let's sort the matches by start time
		//sort.Slice(matches, func(i, j int) bool {
		//	// TODO: Check actual time if match already completed
		//	return matches[i].ScheduledTime < matches[j].ScheduledTime
		//})

		// Finally, make a map of streams to matches
		for _, match := range matches {
			// If the stream URL for the match is valid, use that as the key in the map and append to that array
			if match.StreamUrl != "" {
				tournamentMatchesSet.Details.StreamMatches[match.StreamUrl] = append(tournamentMatchesSet.Details.StreamMatches[match.StreamUrl], match)
			} else { // Otherwise just add it to the UnknownStreamKey array and key
				tournamentMatchesSet.Details.StreamMatches[UnknownStreamKey] = append(tournamentMatchesSet.Details.StreamMatches[UnknownStreamKey], match)
			}
		}

		tournamentMatches = append(tournamentMatches, tournamentMatchesSet)
	}

	if len(tournamentMatches) == 0 {
		return ChannelResponseNoMatches, nil
	}
	return ChannelResponseSuccess, tournamentMatches
}

func (bc *DotaBotChannel) RefreshNotifications() {
	// Update the existing timer, if it's not nil
	if bc.notificationTicker != nil {
		timeUntilNotification, err := bc.calculateTimeUntilNextNotification()
		if err == nil {
			bc.notificationTicker.Reset(timeUntilNotification)
		}
	}
}

func (bc *DotaBotChannel) getTiers() []types.Tier {
	return bc.config.GetTiers()
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
							matchStarted.Match.TeamOne.DisplayName,
							matchStarted.Match.TeamTwo.DisplayName,
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
	actualTimezone, err := utils.GetFullLocation(bc.config.GetTimezone())
	activeTimeZone, err := time.LoadLocation(actualTimezone)
	if err != nil {
		return nil, err
	}
	currentTime := time.Now().In(activeTimeZone)
	return time.FixedZone(currentTime.Zone()), nil
}

func (bc *DotaBotChannel) getChannel() (*discordgo.Channel, error) {
	if bc.channel == nil {
		channel, err := bc.session.Channel(bc.config.GetChannelID())
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
		return bc.config.GetChannelID()
	}
	guild, err := bc.getGuild(channel.GuildID)
	if err != nil {
		return channel.Name
	}
	return fmt.Sprintf("%s:%s", guild.Name, channel.Name)
}

func (bc *DotaBotChannel) sendMessageWithoutEmbeds(messageContent string) {
	discordMsg, err := bc.session.ChannelMessageSend(bc.config.GetChannelID(), messageContent)
	if err != nil {
		log.Println("Error sending message to", bc.getChannelIdentifier(), err.Error())
	} else {
		// Suppress the embeds on the message from the stream links
		editMessage := discordgo.NewMessageEdit(bc.config.GetChannelID(), discordMsg.ID)
		editMessage.Flags |= discordgo.MessageFlagsSuppressEmbeds
		bc.session.ChannelMessageEditComplex(editMessage)
	}
}
