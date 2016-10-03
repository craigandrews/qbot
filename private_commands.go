package qbot

import "github.com/doozr/qbot/command"

// PrivateCommands are commands only available to DM.
func PrivateCommands(commands command.QueueCommands) (commandMap CommandMap) {
	commandMap = CommandMap{
		"list":     commands.List,
		"help":     commands.Help,
		"morehelp": commands.MoreHelp,
	}
	return
}
