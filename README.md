# qbot

Qbot is a slackbot that helps manage a contended resource amongst your team members. This might be a merge token,
a release dongle or a tea-run doodad. When the token comes free and somebody gets it, the bot mentions the new
recipient so they know immediately that it's their turn to merge, release, make tea ...

## Installation

Grab it with `go get`:

    go get github.com/doozr/qbot/cmd/qbot

Or clone the project and build it with `./build.sh` to get proper
versioning.

## Run the bot

Get a bot token from your Slack control panel and run the bot as follows:

    qbot <token> <data file>

The token should be the one copied from the Slack custom integration page. The data file should be a filename in a
directory writable by the bot owner to store serialised versions of the queue.

The bot will autodetect its username and respond to messages directed at it, with an @ or without.

## Running multiple bots

Given that the save location and token are run-time variables it is possible to use one copy of the qbot to run
multiple processes for different Slack channels or teams. Use something like supervisord to start multiple instances
with different values. For example:

    [program:qbot-merge]
    command=/path/to/qbot <"merge" integration token> /path/to/qbot/merge.json
    user=qbotuser

    [program:qbot-release]
    command=/path/to/qbot <"release" integration token> /path/to/qbot/release.json
    user qbotuser

This results in two copies of the bot with two usernames, one for each use. Put them in different channels, or in
the same channel with different names, or whatever suits.

Note that it is not possible to manage multiple queues with a single instance.

## Commands

Address each command to the bot (`<bot name>: <command>`)

* `join <reason>` - Join the queue and give a reason why
* `leave` - Leave the queue (your most recent entry is removed)
* `leave <reason>` - Leave the queue (your most recent entry starting with <reason> is removed)
* `done` - Release the token once you are done with it
* `drop` - Drop the token and leave the queue (alias of `done`)
* `yield` - Release the token and swap places with next in line
* `barge <reason>` - Barge to the front of the queue so you get the token next (only with good reason!)
* `boot <name>` - Kick somebody out of the waiting list (their most recent entry is removed)
* `boot <name> <reason>` - Kick somebody out of the waiting list (their most recent entry starting with <reason> is removed)
* `oust` - Force the token holder to yield to the next in line
* `list` - Show who has the token and who is waiting
* `help` - Show the help text
