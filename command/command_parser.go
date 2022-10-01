package command

import "strings"

type ParseCallback func(params ...string)

type Parser struct {
	commandKeyword string
}

var (
	callbacks map[string]ParseCallback
)

func NewParser(commandKeyword string) *Parser {
	cp := &Parser{commandKeyword}
	cp.initialise()
	return cp
}

func (cp *Parser) initialise() {
	callbacks = make(map[string]ParseCallback)
}

func (cp *Parser) Register(commandName string, callback ParseCallback) {
	callbacks[commandName] = callback
}

func (cp *Parser) Parse(command string) bool {
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

	if val, ok := callbacks[command]; ok {
		val(params...)
		return true
	}

	return false
}
