# qbot - manage contented resources amongst humans

Qbot is a slackbot that helps manage a contended resource amongst your team members. This might be a merge token,
a release dongle or a tea-run doodad. When the token comes free and somebody gets it, the bot mentions the new
recipient so they know immediately that it's their turn to merge, release, make tea ...

## Installation

Build the project with Go using `go build`. Once I have some binaries to distribute I will update these instructions.

## Run the bot

Get a bot token from your Slack control panel and run the bot as follows:

    qbot <name> <token> <data file>

The should match the handle given to the Bot in the Slack configuration, and the token should be the one copied from
the same. The data file should be a filename in a directory writeable by the bot owner to store serialised versions of
the list.

## Commands

Address each command to the bot (`<bot name>: <command>`)

* `join <reason>` - Join the queue and give a reason why
* `leave` - Leave the queue (your most recent entry is removed)
* `leave <reason>` - Leave the queue (your most recent entry starting with <reason> is removed)
* `done` - Release the token once you are done with it
* `yield` - Release the token and swap places with next in line
* `barge <reason>` - Barge to the front of the queue so you get the token next (only with good reason!)
* `boot <name>` - Kick somebody out of the waiting list (their most recent entry is removed)
* `boot <name> <reason>` - Kick somebody out of the waiting list (their most recent entry starting with <reason> is removed)
* `oust` - Forcibly take the token from the token holder and kick them out of the queue (only with VERY good reason!)
* `list` - Show who has the token and who is waiting
* `help` - Show this text