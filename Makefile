# Project params
REPO=github.com/vdromanov/backli915t/
ENTRY_POINT=$(REPO)cmd/backli915t/
#VERSION=git -C $(REPO) describe --tags
VERSION=1.0_test

all: clean build
build:
		CGO_ENABLED=0 go build -ldflags="-s -w -X '$(REPO)internal/pkg/version.Version=$(VERSION)'" -o bin/ -v $(ENTRY_POINT)
clean:
		go clean
		rm -rf bin/*
