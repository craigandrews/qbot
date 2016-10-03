#!/usr/bin/env sh

go build -ldflags "-X github.com/doozr/qbot.version=`git describe --dirty`" github.com/doozr/qbot/cmd/qbot
