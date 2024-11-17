# Makefile
CMD_PATH := ./cmd/ai-qa-service
SWAG_PATH :=  ./cmd/ai-qa-service,internal,pkg
DOC_PATH := ./pkg/docs

# 设置 Go 可执行文件名
SERVER_EXECUTABLE := ai-qa-service

# 设置输出目录
BUILD_DIR := bin

# 设置 Go 编译器
GO := go

# 设置配置文件路径
CONFIG_PATH = ./config/config_test.yaml

.PHONY: all clean

all: clean fmt init build run

fmt:
	swag fmt

init:
	swag init -d $(SWAG_PATH) -o $(DOC_PATH)

build:
	# 编译 Go 程序
	$(GO) build -o $(BUILD_DIR)/$(SERVER_EXECUTABLE) $(CMD_PATH)/main.go

run:
	./$(BUILD_DIR)/$(SERVER_EXECUTABLE) -c $(CONFIG_PATH)

clean:
	rm -rf $(BUILD_DIR)

linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(SERVER_EXECUTABLE)