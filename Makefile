# Project params
REPO=github.com/vdromanov/backli915t/
ENTRY_POINT=$(REPO)cmd/backli915t/
VERSION=$(shell git describe --tags)
CHANGELOG_FNAME=debian/changelog

all: clean build

#Embedding latest version from git tags into a binary
build:
		CGO_ENABLED=0 go build -ldflags="-s -w -X '$(REPO)internal/pkg/version.Version=$(VERSION)'" -o bin/ -v $(ENTRY_POINT)

clean:
		go clean
		rm -rf bin/*

#Adding version to debian's changelog and making a git tag
release:
		dch -D stable -v $(version)
		git add -f $(CHANGELOG_FNAME)
		git commit --amend
		git tag -a v$(version)
