package qbot

import "github.com/doozr/qbot/command"

// PublicCommands are commands accessible from public channels.
func PublicCommands(commands command.QueueCommands) (commandMap CommandMap) {
	commandMap = CommandMap{
		"join":     commands.Join,
		"leave":    commands.Leave,
		"done":     commands.Done,
		"drop":     commands.Done,
		"yield":    commands.Yield,
		"success":  commands.Success,
		"barge":    commands.Barge,
		"replace":  commands.Replace,
		"delegate": commands.Delegate,
		"boot":     commands.Boot,
		"oust":     commands.Oust,
		"list":     commands.List,
		"help":     commands.Help,
	}
	return
}
