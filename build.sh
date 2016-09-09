#!/usr/bin/env sh

go build -ldflags "-X main.Version=`git describe --dirty`"
