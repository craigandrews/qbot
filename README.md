# qbot

<img align="right" width="180" style="margin: 12px" src="https://cdn.rawgit.com/doozr/qbot/master/qbot.svg">

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

*If you don't have the token and need it:*

* `join <reason>` - Join the queue and give a reason why
* `barge <reason>` - Barge to the front of the queue so you get the token next (only with good reason!)
* `barge <position>` - Barge the entry at the given position to the front of the queue

*If you have the token and have done with it:*

* `done` - Release the token once you are done with it
* `drop` - Drop the token and leave the queue (note: actually just an alias of `done`)
* `yield` - Release the token and swap places with next in line

*If you are in the queue and need to change something:*

* `delegate <user>` - Delegate your place to someone else (your most recent entry is delegated)
* `delegate <user> <reason prefix>` - Delegate your place to someone else (match the entry with reason that starts with <reason prefix>)
* `replace <position> <reason>` - Replace the reason of a queue entry you own

*If you are in the queue and need to leave:*

* `leave` - Leave the queue (your most recent entry is removed)
* `leave <position>` - Leave the queue (match the entry at the given position)

*If you need to get rid of somebody who is in the way:*

* `oust <name>` - Force the token holder to yield to the next in line
* `boot <name>` - Kick somebody out of the waiting list (their most recent entry is removed)
* `boot <position> <name>` - Kick somebody out of the waiting list (match the entry at the given position)

*Other useful things to know:*

* `list` - Show who has the token and who is waiting
* `help` - Show this text
