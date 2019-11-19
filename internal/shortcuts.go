package internal

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

const (
	DefaultSearchProvider = "https://duckduckgo.com/%s"
)

var (
	ErrUnrecognizedCommand = errors.New("unrecognized command")
	ErrNotEnoughArguments  = errors.New("not enough arguments")
)

type CommandHandler struct {
	*sync.Mutex
	Shortcuts map[string]string
}

func (c *CommandHandler) Handle(input string) (Command, error) {
	rawArgs := strings.Fields(input)
	fragmentCount := len(rawArgs)
	farg := "*"

	var shortcut string

	if fragmentCount > 0 {
		farg = rawArgs[0]

		if fragmentCount > 1 {
			shortcut = rawArgs[1]
		}

		var updateShortcutParams string
		if fragmentCount > 2 {
			updateShortcutParams = strings.Join(rawArgs[2:], " ")
		}

		if cmd, err := c.updateShortcut(farg, shortcut, updateShortcutParams); err != ErrUnrecognizedCommand {
			return cmd, err
		}
	}

	return c.getShortcut(farg, rawArgs...), nil
}

func (c *CommandHandler) updateShortcut(action, shortcut, location string) (Command, error) {
	command := Command{
		Action:   action,
		Name:     shortcut,
		Location: location,
	}

	switch action {
	case "add":
		if location == "" {
			return Command{}, ErrNotEnoughArguments
		}

		c.Lock()
		c.Shortcuts[shortcut] = command.Location
		c.Unlock()

	case "remove":
		if shortcut == "" {
			return Command{}, ErrNotEnoughArguments
		}

		command.Location = c.Shortcuts[shortcut]

		c.Lock()
		delete(c.Shortcuts, shortcut)
		c.Unlock()

	default:
		return Command{}, ErrUnrecognizedCommand
	}

	return command, nil
}

func (c *CommandHandler) getShortcut(key string, input ...string) Command {
	var parameter string

	location, ok := c.Shortcuts[key]
	if !ok {
		location = DefaultSearchProvider
		key = "*"
		parameter = strings.Join(input, " ")
	} else {
		parameter = strings.Join(input[1:], " ")
	}

	if strings.Contains(location, "%s") {
		location = fmt.Sprintf(location, parameter)
	}

	return Command{
		Action:   "lookup",
		Name:     key,
		Location: location,
	}
}

func NewDefaultShortcuts() Shortcuts {
	return Shortcuts{
		"*":    DefaultSearchProvider,
		"help": "/",
	}
}

type Command struct {
	Action   string
	Name     string
	Location string
}
