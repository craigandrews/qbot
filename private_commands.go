package main

import "github.com/doozr/qbot/command"

func privateCommands(commands command.Command) (commandMap CommandMap) {
	commandMap = CommandMap{
		"list":     commands.List,
		"help":     commands.Help,
		"morehelp": commands.MoreHelp,
	}
	return
}
