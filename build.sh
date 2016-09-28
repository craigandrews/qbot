#!/usr/bin/env sh

go build -ldflags "-X main.Version=`git describe --dirty`" github.com/doozr/qbot/cmd/qbot
