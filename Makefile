# Makefile для инициализации проекта "Tyndau"

.PHONY: init

# Command: "make initClean DIR=api-gateway". It will create needed directories for Clean Architecture in the given folder
initClean:
	mkdir -p $(DIR)/cmd
	mkdir -p $(DIR)/config
	mkdir -p $(DIR)/internal/adapter/grpc/server/frontend/dto
	mkdir -p $(DIR)/internal/adapter/postgres/dao
	mkdir -p $(DIR)/internal/adapter
	mkdir -p $(DIR)/internal/app
	mkdir -p $(DIR)/internal/model
	mkdir -p $(DIR)/internal/usecase
	mkdir -p $(DIR)/migrations
	mkdir -p $(DIR)/pkg


run-local:
	cd api-gateway/ && go run ./cmd/gateway & \
	cd ../notification-service && go run ./cmd/notification & \
	cd ../user_service && go run ./cmd/user & \
	cd ../music-service && go run ./cmd/music & \
	cd ../auth_service && go run ./cmd/auth & \
	wait