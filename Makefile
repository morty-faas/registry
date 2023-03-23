.PHONY: build
build:
	go build -o bin/registry main.go

.PHONY: start
start: build
	sudo -E bin/registry
