GIT_SHA := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u '+%Y%m%d-%H%M%S')

default: test

deps:
	go get -t -v .

test: deps
	go test -v

release: test deps
	GOOS=linux  GOARCH=amd64 go build -o bin/jg-linux  -ldflags "-X main.version=$(TAG) -X main.gitCommit=$(GIT_SHA) -X main.buildDate=$(DATE)"
	GOOS=darwin GOARCH=amd64 go build -o bin/jg-darwin -ldflags "-X main.version=$(TAG) -X main.gitCommit=$(GIT_SHA) -X main.buildDate=$(DATE)"

.PHONY: deps test release default