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

const (
	connectCommandKey      = "connect"
	todayCommandKey        = "today"
	timezoneCommandKey     = "timezone"
	leagueCommandKey       = "league"
	leagueAddCommandKey    = "add"
	leagueRemoveCommandKey = "remove"
	disconnectCommandKey   = "disconnect"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        connectCommandKey,
			Description: "Connect DotaBot to this channel",
		},
		{
			Name:        timezoneCommandKey,
			Description: "Set the timezone DotaBot will use when displaying times",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "timezone",
					Description: "The timezone to set",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "GMT",
							Value: "GMT",
						},
						{
							Name:  "EET",
							Value: "EET",
						},
					},
				},
			},
		},
		{
			Name:        leagueCommandKey,
			Description: "Add/remove the leagues to be notified about",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        leagueAddCommandKey,
					Description: "Add the specified league selection to the notification list",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "league",
							Description: "The league to add",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "DPC League",
									Value: schema.LeagueTierDpcLeague,
								},
								{
									Name:  "The International",
									Value: schema.LeagueTierInternational,
								},
							},
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        leagueRemoveCommandKey,
					Description: "Remove the specified league selection to the notification list",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "league",
							Description: "The league to remove",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "DPC League",
									Value: schema.LeagueTierDpcLeague,
								},
								{
									Name:  "The International",
									Value: schema.LeagueTierInternational,
								},
							},
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        todayCommandKey,
			Description: "Get all matches that are happening in the DPC today",
		},
		{
			Name:        disconnectCommandKey,
			Description: "Disconnect DotaBot from this channel",
		},
	}
	handlers = map[string]func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate){
		connectCommandKey: func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			fmt.Println("Connecting DotaBot to channel", i.ChannelID)
			if _, ok := b.channels[i.ChannelID]; ok {
				fmt.Println("Bot already started in this channel")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "DotaBot is already connected to this channel",
					},
				})
			} else {
				fmt.Println("Starting bot on channel", i.ChannelID)
				channel := NewDotaBotChannel(i.ChannelID)
				channel.Start()
				b.channels[i.ChannelID] = channel
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "DotaBot is now connected to this channel!",
					},
				})
			}
		},
		timezoneCommandKey: func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			if channel, ok := b.channels[i.ChannelID]; ok {
				if len(i.ApplicationCommandData().Options) > 0 {
					timezone := i.ApplicationCommandData().Options[0].StringValue()
					err := channel.UpdateTimezone(timezone)
					if err != nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Invalid timezone specified",
							},
						})
					} else {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Timezone for this channel is now set to " + timezone,
							},
						})
					}
				}
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "DotaBot is not connected to this channel yet! Please use the \"/connect\" command before running other commands",
					},
				})
			}
		},
		todayCommandKey: func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			if channel, ok := b.channels[i.ChannelID]; ok {
				// If there's no league tiers configured, let the channel know!
				tiers := channel.GetLeagues()
				if len(tiers) == 0 {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "No leagues have been added to this channel's configuration yet! Add some by using /league add [league]",
						},
					})
					return
				}

				leagues, err := b.stratzClient.GetActiveLeagues(tiers)
				if err != nil {
					b.discordSession.ChannelMessageSend(i.ChannelID, "Failed to get active leagues")
				} else {
					// TODO: Definitely ways to improve and optimise this code :shrug:
					// TODO: Move this out of this method when we have daily notifications
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
						message := ":robot: " + league.DisplayName + " games today!\n\n"
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
							// TODO: Respond to the interaction somehow
							discordMsg, err := b.discordSession.ChannelMessageSend(i.ChannelID, message)
							if err != nil {
								fmt.Println("Error sending message to", i.ChannelID, err.Error())
							} else {
								// Suppress the embeds on the message from the stream links
								editMessage := discordgo.NewMessageEdit(i.ChannelID, discordMsg.ID)
								editMessage.Flags |= discordgo.MessageFlagsSuppressEmbeds
								b.discordSession.ChannelMessageEditComplex(editMessage)
							}
						}
					}
				}
			}
		},
		leagueCommandKey: func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Check whether it was an add or remove command
			if channel, ok := b.channels[i.ChannelID]; ok {
				innerCommand := i.ApplicationCommandData().Options[0]
				leagueValue := innerCommand.Options[0].StringValue()
				switch innerCommand.Name {
				case leagueAddCommandKey:
					{
						addedSuccessfully := channel.AddLeague(schema.LeagueTier(leagueValue))
						if addedSuccessfully {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{Content: "League added successfully!"},
							})
						} else {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{Content: "League has already been added"},
							})
						}
						break
					}
				case leagueRemoveCommandKey:
					{
						channel.RemoveLeague(schema.LeagueTier(leagueValue))
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{Content: "League removed successfully!"},
						})
						break
					}
				}
			}
		},
		disconnectCommandKey: func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			if channel, ok := b.channels[i.ChannelID]; ok {
				fmt.Println("Stopping bot on channel", i.ChannelID)
				channel.Stop()
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "DotaBot is now disconnected from this channel",
					},
				})
			}
			delete(b.channels, i.ChannelID)
		},
	}
)

type DotaBot struct {
	GuildID        string
	stratzClient   *stratz.Client
	discordSession *discordgo.Session
	channels       map[string]*DotaBotChannel
}

func NewDotaBot(stratzClient *stratz.Client) *DotaBot {
	return NewDotaBotWithGuildID(stratzClient, "")
}

func NewDotaBotWithGuildID(stratzClient *stratz.Client, guildID string) *DotaBot {
	b := new(DotaBot)
	b.GuildID = guildID
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

		// Call update on the config in case there's new values added that should go into the database
		config.Update()
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session", err)
		return err
	}

	// TODO: Add this back in at some point
	//b.commandParser.Register("daily", func(params *command.ParseParameters) {
	//	if channel, ok := b.channels[params.ChannelID]; ok {
	//		if len(params.Parameters) > 0 {
	//			timeString := params.Parameters[0]
	//			err := channel.UpdateDailyMessageTime(timeString)
	//			if err != nil {
	//				b.discordSession.ChannelMessageSend(params.ChannelID, "Invalid time format")
	//			}
	//		}
	//	} else {
	//		b.discordSession.ChannelMessageSend(params.ChannelID, "Channel is not active yet! Please type \"!dotabot start\" before running other commands")
	//	}
	//})
	//

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//fmt.Println(i.ApplicationCommandData().Name)
		if command, ok := handlers[i.ApplicationCommandData().Name]; ok {
			command(b, s, i)
		}
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()

	if err == nil {
		b.discordSession = dg
		for _, command := range commands {
			cmd, err := b.discordSession.ApplicationCommandCreate(b.discordSession.State.User.ID, b.GuildID, command)
			if err != nil {
				fmt.Println("Error creating command", err)
			} else {
				fmt.Println("Command registered", cmd)
			}
		}
	} else {
		fmt.Println("Error when opening session", err)
	}
	return err
}

func (b *DotaBot) Shutdown() {
	// Remove all registered commands
	registeredCommands, err := b.discordSession.ApplicationCommands(b.discordSession.State.User.ID, "")
	if err != nil {
		fmt.Println("Error when closing Discord session", err)
	}

	for _, command := range registeredCommands {
		b.discordSession.ApplicationCommandDelete(b.discordSession.State.User.ID, "", command.ID)
	}

	err = b.discordSession.Close()
	if err != nil {
		fmt.Println("Error when closing Discord session", err)
	}
}
