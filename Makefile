# Load .env
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

export GOOSE_DRIVER := $(DB_DRIVER)
export GOOSE_DBSTRING := $(DB_URL)
export GOOSE_MIGRATION_DIR := ./db/migration

# Development environment
dev:
	air

# Build for production
build:
	go build -o bin/main.exe .

# Start production environment
start: build
	bin/main.exe

# Clean up build files
clean:
	rm -rf bin

# Create a new migration (ex: `make new_migration name=create_users_table`)
new_migration:
	goose create $(name) sql

# Apply all migrations
migrate_up:
	goose up

# Roll back a single migration from the current version
migrate_down:
	goose down

.PHONY: dev build start clean new_migration migrate_up migrate_down
