#!/bin/sh

export GO111MODULE=on
go build -race -ldflags "-extldflags '-static'" -o target/radibrary-cli cmd/radibrary-cli/radibrary-cli.go
