# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=drift

all: test build
build: 
		$(GOBUILD) -o ./dist/$(BINARY_NAME) -v ./cmd/drift
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f ./dist/$(BINARY_NAME)
run:
		$(GOBUILD) -o ./dist/$(BINARY_NAME) -v ./cmd/drift
		./dist/$(BINARY_NAME) up