package command

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type ParseParameters struct {
	ChannelID  string
	Parameters []string
}

type ParseCallback func(params *ParseParameters)

type Parser struct {
	commandKeyword string
	callbacks      map[string]ParseCallback
}

func NewParser(commandKeyword string) *Parser {
	cp := &Parser{commandKeyword, make(map[string]ParseCallback)}
	return cp
}

func (cp *Parser) Register(commandName string, callback ParseCallback) {
	cp.callbacks[commandName] = callback
}

func (cp *Parser) Parse(message *discordgo.Message) bool {
	command := message.Content
	if !strings.HasPrefix(command, cp.commandKeyword) {
		return false
	}

	command = strings.TrimPrefix(command, cp.commandKeyword)
	command = strings.TrimSpace(command)

	// Split rest of it in case there's parameters
	parts := strings.Split(command, " ")
	if len(parts) == 0 {
		return false
	}

	command = parts[0]
	params := parts[1:]

	if val, ok := cp.callbacks[command]; ok {
		parameters := &ParseParameters{message.ChannelID, params}
		val(parameters)
		return true
	}

	return false
}
