SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 1.0.0
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall fmt simplify check run

all: check install

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET)

install:
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
	@go tool vet ${SRC}

run: install
	$(TARGET)
	
# NAME = Heimdallr
# PWD := $(MKPATH:%/Makefile=%)

# clean:
# 	cd "$(PWD)"
# 	rm -rf vendor

# install:
# 	dep ensure -update

# build:
# 	go build -race -o $(NAME)

# run:
# 	@go run main.go

# start:
# 	./$(NAME)

# test:
# 	go test -race -v $(shell go list ./... | grep -v /vendor/)

# coverage:
# 	go test -race -cover -v $(shell go list ./... | grep -v /vendor/)

# vet:
# 	go vet $(shell go list ./... | grep -v /vendor/)

# #staticcheck:
# #	go get -u honnef.co/go/staticcheck/cmd/staticcheck
# #	staticcheck $(shell go list ./... | grep -v /vendor/)

# lint:
# 	golint $(shell go list ./... | grep -v /vendor/)

# #simple:
# #	go get -u honnef.co/go/simple/cmd/gosimple
# #    gosimple $(go list ./... | grep -v "vendor")
