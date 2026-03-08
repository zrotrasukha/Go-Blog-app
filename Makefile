APP_NAME=web
BIN_DIR=bin
CMD_DIR=./cmd/web/*

.PHONY: build run clean

build:
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

run: build
	@air -c .air.toml

clean:
	@rm -rf $(BIN_DIR)
