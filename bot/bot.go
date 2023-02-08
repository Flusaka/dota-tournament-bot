package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/flusaka/dota-tournament-bot/command"
	"github.com/flusaka/dota-tournament-bot/models"
	"github.com/flusaka/dota-tournament-bot/stratz"
	"github.com/flusaka/dota-tournament-bot/stratz/schema"
	"sort"
	"time"
)

type DotaBot struct {
	commandParser  *command.Parser
	stratzClient   *stratz.Client
	discordSession *discordgo.Session
	channels       map[string]*DotaBotChannel
}

func NewDotaBot(commandParser *command.Parser, stratzClient *stratz.Client) *DotaBot {
	b := new(DotaBot)
	b.commandParser = commandParser
	b.stratzClient = stratzClient
	b.channels = make(map[string]*DotaBotChannel)
	return b
}

func (b *DotaBot) Initialise(token string) error {
	configs, err := models.FetchAllConfigs()
	if err != nil {
		fmt.Println("Could not retrieve configs", err)
	}

	// Setup existing bot channels
	for _, config := range configs {
		fmt.Println("Restarting channel on ID", config.ChannelID)
		b.channels[config.ChannelID] = NewDotaBotChannelWithConfig(config)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session", err)
		return err
	}

	b.commandParser.Register("start", func(params *command.ParseParameters) {
		fmt.Println("Start called with channel ID", params.ChannelID)
		if _, ok := b.channels[params.ChannelID]; ok {
			fmt.Println("Bot already started in this channel")
		} else {
			fmt.Println("Starting bot on channel", params.ChannelID)
			channel := NewDotaBotChannel(params.ChannelID)
			channel.Start()
			b.channels[params.ChannelID] = channel
		}
	})

	b.commandParser.Register("timezone", func(params *command.ParseParameters) {
		if channel, ok := b.channels[params.ChannelID]; ok {
			if len(params.Parameters) > 0 {
				timezone := params.Parameters[0]
				err := channel.UpdateTimezone(timezone)
				if err != nil {
					b.discordSession.ChannelMessageSend(params.ChannelID, "Invalid timezone specified")
				}
			}
		} else {
			b.discordSession.ChannelMessageSend(params.ChannelID, "Channel is not active yet! Please type \"!dotabot start\" before running other commands")
		}
	})

	b.commandParser.Register("daily", func(params *command.ParseParameters) {
		if channel, ok := b.channels[params.ChannelID]; ok {
			if len(params.Parameters) > 0 {
				timeString := params.Parameters[0]
				err := channel.UpdateDailyMessageTime(timeString)
				if err != nil {
					b.discordSession.ChannelMessageSend(params.ChannelID, "Invalid time format")
				}
			}
		} else {
			b.discordSession.ChannelMessageSend(params.ChannelID, "Channel is not active yet! Please type \"!dotabot start\" before running other commands")
		}
	})

	// TODO: Eventually this won't be an explicit command, or at least it will also be calculated/sent based on the notification time of a configured channel
	b.commandParser.Register("today", func(params *command.ParseParameters) {
		if channel, ok := b.channels[params.ChannelID]; ok {
			// TODO: Make this configurable in the channel
			var tiers = []schema.LeagueTier{schema.LeagueTierDpcLeague}
			leagues, err := b.stratzClient.GetActiveLeagues(tiers)
			if err != nil {
				b.discordSession.ChannelMessageSend(params.ChannelID, "Failed to get active leagues")
			} else {
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
						isWithinDay := channel.IsTimeWithinDay(match.ScheduledTime)
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
					message := league.DisplayName + " games today!\n\n"
					for streamUrl, streamMatches := range streamMatchesMap {
						message += "Games on: " + streamUrl + "\n\n"
						for _, streamMatch := range streamMatches {
							convertedTime, err := channel.GetTimeInZone(streamMatch.ScheduledTime)
							if err != nil {
								continue
							}
							message += streamMatch.TeamOne.Name + " vs " + streamMatch.TeamTwo.Name + " - " + convertedTime.Format(time.Kitchen) + "\n"
						}
					}

					if len(message) > 0 {
						// Send the message and get the Discord message struct back
						discordMsg, err := b.discordSession.ChannelMessageSend(params.ChannelID, message)
						if err != nil {
							fmt.Println("Error sending message to", params.ChannelID, err.Error())
						} else {
							// Suppress the embeds on the message from the stream links
							editMessage := discordgo.NewMessageEdit(params.ChannelID, discordMsg.ID)
							editMessage.Flags |= discordgo.MessageFlagsSuppressEmbeds
							b.discordSession.ChannelMessageEditComplex(editMessage)
						}
					}
				}
			}
		}
	})

	b.commandParser.Register("stop", func(params *command.ParseParameters) {
		if channel, ok := b.channels[params.ChannelID]; ok {
			fmt.Println("Stopping bot on channel", params.ChannelID)
			channel.Stop()
		}
		delete(b.channels, params.ChannelID)
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return
		}

		if !b.commandParser.Parse(m.Message) {
			fmt.Println("Ignoring unparsed message")
		}
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()

	if err == nil {
		b.discordSession = dg
	} else {
		fmt.Println("Error when opening session", err)
	}
	return err
}

func (b *DotaBot) Shutdown() {
	err := b.discordSession.Close()
	if err != nil {
		fmt.Println("Error when closing Discord session", err)
	}
}
