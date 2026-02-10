APP_NAME := file-organizer
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -s -w \
	-X github.com/geraldbahati/file-organizer/cmd.Version=$(VERSION) \
	-X github.com/geraldbahati/file-organizer/cmd.CommitSHA=$(COMMIT) \
	-X github.com/geraldbahati/file-organizer/cmd.BuildDate=$(BUILD_DATE)

.PHONY: build clean run install uninstall

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) .

clean:
	rm -rf bin/

run: build
	./bin/$(APP_NAME) start --foreground -v

install: build
	./bin/$(APP_NAME) install

uninstall:
	./bin/$(APP_NAME) uninstall
