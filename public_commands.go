package main

import "github.com/doozr/qbot/command"

func publicCommands(commands command.Command) (commandMap CommandMap) {
	commandMap = CommandMap{
		"join":     commands.Join,
		"leave":    commands.Leave,
		"done":     commands.Done,
		"drop":     commands.Done,
		"yield":    commands.Yield,
		"barge":    commands.Barge,
		"boot":     commands.Boot,
		"oust":     commands.Oust,
		"list":     commands.List,
		"help":     commands.Help,
		"morehelp": commands.MoreHelp,
	}
	return
}
