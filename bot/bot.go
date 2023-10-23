package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/flusaka/dota-tournament-bot/datasource"
	"github.com/flusaka/dota-tournament-bot/datasource/types"
	"github.com/flusaka/dota-tournament-bot/models"
	"log"
)

const (
	connectCommandKey      = "connect"
	todayCommandKey        = "today"
	dailyCommandKey        = "daily"
	notifyDailyCommandKey  = "notify"
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
									Value: types.TierDpcLeague,
								},
								{
									Name:  "The International Qualifiers",
									Value: types.TierDpcLeagueQualifier,
								},
								{
									Name:  "The International",
									Value: types.TierInternational,
								},
								{
									Name:  "Majors",
									Value: types.TierMajor,
								},
								{
									Name:  "Minors",
									Value: types.TierMinor,
								},
								{
									Name:  "Other Pro Tournaments",
									Value: types.TierProfessional,
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
									Value: types.TierDpcLeague,
								},
								{
									Name:  "The International Qualifiers",
									Value: types.TierDpcLeagueQualifier,
								},
								{
									Name:  "The International",
									Value: types.TierInternational,
								},
								{
									Name:  "Majors",
									Value: types.TierMajor,
								},
								{
									Name:  "Minors",
									Value: types.TierMinor,
								},
								{
									Name:  "Other Pro Tournaments",
									Value: types.TierProfessional,
								},
							},
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        dailyCommandKey,
			Description: "Set the time to be notified every day of all the day's matches",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "time",
					Description: "The time to send daily notifications, e.g. 10:30AM",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        notifyDailyCommandKey,
			Description: "Turn on/off daily notifications of the day's matches",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "enabled",
					Description: "Whether daily notifications are enabled or disabled",
					Required:    true,
				},
			},
		},
		{
			Name:        todayCommandKey,
			Description: "Get all matches that are happening today",
		},
		{
			Name:        disconnectCommandKey,
			Description: "Disconnect DotaBot from this channel",
		},
	}
	handlers = map[string]func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate){
		connectCommandKey: func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("Connecting DotaBot to channel", i.ChannelID)
			if _, ok := b.channels[i.ChannelID]; ok {
				log.Println("Bot already started in this channel")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "DotaBot is already connected to this channel",
					},
				})
			} else {
				log.Println("Starting bot on channel", i.ChannelID)
				channel := NewDotaBotChannel(s, i.ChannelID, b.dataSourceClient)
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
				channel.SendMatchesOfTheDayInResponseTo(i)
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "DotaBot is not connected to this channel yet! Please use the \"/connect\" command before running other commands",
					},
				})
			}
		},
		dailyCommandKey: func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			if channel, ok := b.channels[i.ChannelID]; ok {
				time := i.ApplicationCommandData().Options[0].StringValue()
				err := channel.UpdateDailyMessageTime(time)
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{Content: "Invalid time format! Please make sure use full 12-hour time format - e.g. 10:00AM"},
					})
				} else {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{Content: "Daily notification set to " + time},
					})
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
		notifyDailyCommandKey: func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			if channel, ok := b.channels[i.ChannelID]; ok {
				notificationsEnabled := i.ApplicationCommandData().Options[0].BoolValue()
				channel.EnableDailyNotifications(notificationsEnabled)
				if notificationsEnabled {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{Content: "Daily notifications are now enabled in this channel!"},
					})
				} else {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{Content: "Daily notifications are now disabled in this channel!"},
					})
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
		leagueCommandKey: func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Check whether it was an add or remove command
			if channel, ok := b.channels[i.ChannelID]; ok {
				innerCommand := i.ApplicationCommandData().Options[0]
				leagueValue := innerCommand.Options[0].StringValue()
				switch innerCommand.Name {
				case leagueAddCommandKey:
					{
						addedSuccessfully := channel.AddLeague(types.Tier(leagueValue))
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
						channel.RemoveLeague(types.Tier(leagueValue))
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{Content: "League removed successfully!"},
						})
						break
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
		disconnectCommandKey: func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			if channel, ok := b.channels[i.ChannelID]; ok {
				log.Println("Stopping bot on channel", i.ChannelID)
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
	guildID            string
	dataSourceClient   datasource.Client
	discordSession     *discordgo.Session
	channels           map[string]*DotaBotChannel
	registeredCommands []*discordgo.ApplicationCommand
}

func NewDotaBot(dataSourceClient datasource.Client) *DotaBot {
	return NewDotaBotWithGuildID(dataSourceClient, "")
}

func NewDotaBotWithGuildID(dataSourceClient datasource.Client, guildID string) *DotaBot {
	b := new(DotaBot)
	b.guildID = guildID
	b.dataSourceClient = dataSourceClient
	b.channels = make(map[string]*DotaBotChannel)
	return b
}

func (b *DotaBot) Initialise(token string) error {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Error creating Discord session", err)
		return err
	}

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot is now ready")
	})

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			{
				if command, ok := handlers[i.ApplicationCommandData().Name]; ok {
					command(b, s, i)
				}
				break
			}
		case discordgo.InteractionMessageComponent:
			{
				if channel, ok := b.channels[i.ChannelID]; ok {
					channel.HandleMessageComponentInteraction(i)
				}
				break
			}
		}
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()

	if err == nil {
		b.discordSession = dg

		configs, err := models.FetchAllConfigs()
		if err != nil {
			log.Println("Could not retrieve configs", err)
		}

		// Setup existing bot channels
		for _, config := range configs {
			log.Println("Restarting channel on ID", config.ChannelID)
			b.channels[config.ChannelID] = NewDotaBotChannelWithConfig(b.discordSession, config, b.dataSourceClient)

			// Call update on the config in case there's new values added that should go into the database
			config.Update()
		}

		b.registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
		for i, command := range commands {
			cmd, err := b.discordSession.ApplicationCommandCreate(b.discordSession.State.User.ID, b.guildID, command)
			if err != nil {
				log.Println("Error creating command", err)
			} else {
				log.Printf("Command %v registered", command.Name)
			}
			b.registeredCommands[i] = cmd
		}
	} else {
		log.Println("Error when opening session", err)
	}
	return err
}

func (b *DotaBot) Shutdown() {
	// Remove all registered commands
	for _, command := range b.registeredCommands {
		err := b.discordSession.ApplicationCommandDelete(b.discordSession.State.User.ID, b.guildID, command.ID)
		if err != nil {
			log.Printf("Command %v failed to be removed", command.Name)
		} else {
			log.Printf("Command %v removed successfully", command.Name)
		}
	}

	err := b.discordSession.Close()
	if err != nil {
		log.Println("Error when closing Discord session", err)
	}
}
