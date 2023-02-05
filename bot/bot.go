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

	b.commandParser.Register("today", func(params *command.ParseParameters) {
		if channel, ok := b.channels[params.ChannelID]; ok {
			var tiers = []schema.LeagueTier{schema.LeagueTierDpcLeague}
			matches, err := b.stratzClient.GetMatchesInActiveLeagues(tiers)
			if err != nil {
				b.discordSession.ChannelMessageSend(params.ChannelID, "Failed to get active leagues")
			} else {
				var message = ""
				// First, let's sort the matches
				sort.Slice(matches, func(i, j int) bool {
					return matches[i].ScheduledTime < matches[j].ScheduledTime
				})

				// TODO: Now, let's split into groups based on the league it's for

				for _, match := range matches {
					// TODO: Check actual time if match already completed
					isWithinDay, convertedTime := channel.IsTimeWithinDay(match.ScheduledTime)
					if !isWithinDay {
						continue
					}
					message += match.TeamOne.Name + " vs " + match.TeamTwo.Name + " (" + convertedTime.Format(time.Kitchen) + ")\n"
				}
				if len(message) > 0 {
					_, err := b.discordSession.ChannelMessageSend(params.ChannelID, message)
					if err != nil {
						fmt.Println("Error sending message to", params.ChannelID, err.Error())
					}
				} else {
					b.discordSession.ChannelMessageSend(params.ChannelID, "No games today!")
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
