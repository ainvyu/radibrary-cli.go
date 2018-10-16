# Radibrary downloader 

## Project details

Radibrary Radio File Downloader for Golang Practice

## Build

    go build -race -ldflags "-extldflags '-static'" -o $CI_PROJECT_DIR/radibrary-cli cmd/radibrary-cli/radibrary-cli.go
