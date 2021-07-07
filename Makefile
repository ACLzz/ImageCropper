BIN_FOLDER=./bin
CONFIG_FOLDER=$(HOME)/.imcr

build:
	go build -o $(BIN_FOLDER)/server ./src/main

move_config:
	mkdir -p $(CONFIG_FOLDER)
	cp -rf extra/* $(CONFIG_FOLDER)

setup: move_config build