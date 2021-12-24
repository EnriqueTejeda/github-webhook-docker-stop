# Helpers
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=github-webhook-docker-stop
DOCKER_NAME=etejeda/github-webhook-docker-stop:latest

.PHONY: build test clean build-docker
build:
	$(GOCMD) mod tidy
	$(GOBUILD) -o $(BINARY_NAME) -v
run: 
	./$(BINARY_NAME)
test:
	$(GOTEST) -v
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
build-docker:
	docker build --compress . -t ${DOCKER_NAME}