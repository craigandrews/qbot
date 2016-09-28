#!/usr/bin/env sh

go build -ldflags "-X main.version=`git describe --dirty`" github.com/doozr/qbot/cmd/qbot
