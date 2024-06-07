package bot

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/flusaka/dota-tournament-bot/types"
	"log"
)

const (
	connectCommandKey     = "connect"
	todayCommandKey       = "today"
	dailyCommandKey       = "daily"
	notifyDailyCommandKey = "notify"
	timezoneCommandKey    = "timezone"
	tierCommandKey        = "tier"
	tierAddCommandKey     = "add"
	tierRemoveCommandKey  = "remove"
	disconnectCommandKey  = "disconnect"
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
			Name:        tierCommandKey,
			Description: "Add/remove the leagues to be notified about",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        tierAddCommandKey,
					Description: "Add the specified league selection to the notification list",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "league",
							Description: "The league to add",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "S Tier",
									Value: types.TierS,
								},
								{
									Name:  "A Tier",
									Value: types.TierA,
								},
								{
									Name:  "B Tier",
									Value: types.TierB,
								},
								{
									Name:  "C Tier",
									Value: types.TierC,
								},
								{
									Name:  "D Tier",
									Value: types.TierD,
								},
							},
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        tierRemoveCommandKey,
					Description: "Remove the specified league selection to the notification list",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "league",
							Description: "The league to remove",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "S Tier",
									Value: types.TierS,
								},
								{
									Name:  "A Tier",
									Value: types.TierA,
								},
								{
									Name:  "B Tier",
									Value: types.TierB,
								},
								{
									Name:  "C Tier",
									Value: types.TierC,
								},
								{
									Name:  "D Tier",
									Value: types.TierD,
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
				// Create a new config for this channel
				config, err := b.channelConfigRepository.Create(context.TODO(), i.ChannelID)
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "An error occured when trying to connect DotaBot to this channel :(",
						},
					})
					return
				}

				channel := NewDotaBotChannelWithConfig(s, config, b.dataSource)
				b.channels[i.ChannelID] = &Channel{
					BotChannel: channel,
					Config:     config,
				}
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
					// Update Timezone to Config first
					timezone := i.ApplicationCommandData().Options[0].StringValue()
					channel.Config.SetTimezone(timezone)
					err := b.channelConfigRepository.Update(context.TODO(), channel.Config)

					// TODO: Pass the updated config to the bot, or tell it to refresh?

					if err != nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Failed to update timezone",
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
				channel.BotChannel.SendMatchesOfTheDayInResponseTo(i)
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
				// TODO: Parse string in correct format and get hours and minutes values
				channel.Config.SetDailyMessageTime(0, 0)
				err := b.channelConfigRepository.Update(context.TODO(), channel.Config)

				// TODO: Update channel bot instance
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
				channel.Config.SetDailyNotificationsEnabled(notificationsEnabled)
				b.channelConfigRepository.Update(context.TODO(), channel.Config)

				// TODO: Trigger the start/stop of daily notifications on this channel
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
		tierCommandKey: func(b *DotaBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Check whether it was an add or remove command
			if channel, ok := b.channels[i.ChannelID]; ok {
				innerCommand := i.ApplicationCommandData().Options[0]
				tierValue := innerCommand.Options[0].StringValue()
				switch innerCommand.Name {
				case tierAddCommandKey:
					{
						addedSuccessfully := channel.Config.AddTier(types.Tier(tierValue))
						if addedSuccessfully {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{Content: "Tournament added successfully!"},
							})
						} else {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{Content: "Tournament has already been added"},
							})
						}
						break
					}
				case tierRemoveCommandKey:
					{
						channel.Config.RemoveTier(types.Tier(tierValue))
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{Content: "Tournament removed successfully!"},
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
				channel.BotChannel.Close()

				// Delete the channel config
				b.channelConfigRepository.Delete(context.TODO(), channel.Config)

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

type Channel struct {
	BotChannel *DotaBotChannel
	Config     ChannelConfig
}

type DotaBot struct {
	guildID                 string
	channelConfigRepository ChannelConfigRepository
	dataSource              DataSource
	session                 *discordgo.Session
	channels                map[string]*Channel
	registeredCommands      []*discordgo.ApplicationCommand
}

func NewDotaBot(dataSourceClient DataSource, channelConfigRepository ChannelConfigRepository) *DotaBot {
	return NewDotaBotWithGuildID(dataSourceClient, channelConfigRepository, "")
}

func NewDotaBotWithGuildID(dataSource DataSource, channelConfigRepository ChannelConfigRepository, guildID string) *DotaBot {
	b := &DotaBot{
		guildID:                 guildID,
		channelConfigRepository: channelConfigRepository,
		dataSource:              dataSource,
		channels:                make(map[string]*Channel),
	}
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
					channel.BotChannel.HandleMessageComponentInteraction(i)
				}
				break
			}
		}
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()

	if err == nil {
		b.session = dg

		configs, err := b.channelConfigRepository.GetAll(context.TODO())
		if err != nil {
			log.Println("Could not retrieve configs", err)
		}

		// Setup existing bot channels
		for _, config := range configs {
			log.Println("Restarting channel on ID", config.GetChannelID())
			b.channels[config.GetChannelID()] = &Channel{
				BotChannel: NewDotaBotChannelWithConfig(b.session, config, b.dataSource),
				Config:     config,
			}

			// Call update on the config in case there's new values added that should go into the database
			b.channelConfigRepository.Update(context.TODO(), config)
		}

		// Sync commands
		b.syncCommands()
	} else {
		log.Println("Error when opening session", err)
	}
	return err
}

func (b *DotaBot) Shutdown() {
	// Remove all registered commands
	for _, command := range b.registeredCommands {
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, b.guildID, command.ID)
		if err != nil {
			log.Printf("Command %v failed to be removed", command.Name)
		} else {
			log.Printf("Command %v removed successfully", command.Name)
		}
	}

	err := b.session.Close()
	if err != nil {
		log.Println("Error when closing Discord session", err)
	}
}

func (b *DotaBot) syncCommands() error {
	existingCommands, err := b.session.ApplicationCommands(b.session.State.User.ID, b.guildID)
	if err != nil {
		return err
	}

	desiredCommands := make(map[string]*discordgo.ApplicationCommand, len(commands))
	for _, cmd := range commands {
		desiredCommands[cmd.Name] = cmd
	}

	existingMap := make(map[string]*discordgo.ApplicationCommand, len(existingCommands))
	for _, cmd := range existingCommands {
		existingMap[cmd.Name] = cmd
	}

	// Go through existing commands and check if any need to be deleted
	for _, cmd := range existingCommands {
		if _, found := desiredCommands[cmd.Name]; !found {
			// Delete the command if the existing one is no longer in the desiredCommands
			err := b.session.ApplicationCommandDelete(b.session.State.User.ID, b.guildID, cmd.ID)
			if err != nil {
				log.Printf("Failed to delete command %s (%v)", cmd.Name, err)
			}
		}
	}

	// Go through the list of desiredCommands and if it already exists, just update, otherwise create
	for _, cmd := range desiredCommands {
		if existingCmd, found := existingMap[cmd.Name]; found {
			_, err := b.session.ApplicationCommandEdit(b.session.State.User.ID, b.guildID, existingCmd.ID, cmd)
			if err != nil {
				log.Printf("Failed to edit command %s (%s) in guild %s: %v", cmd.Name, cmd.ID, b.guildID, err)
			} else {
				log.Printf("Successfully edited command %s (%s) in guild %s", cmd.Name, cmd.ID, b.guildID)
			}
		} else {
			// Create new command
			_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, b.guildID, cmd)
			if err != nil {
				log.Printf("Failed to create command %s in guild %s: %v", cmd.Name, b.guildID, err)
			} else {
				log.Printf("Successfully created command %s in guild %s", cmd.Name, b.guildID)
			}
		}
	}
	return nil
}
