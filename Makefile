APP_NAME=main.exe

# Development environment
.PHONY: dev
dev:
	air

# Build for production
.PHONY: build
build:
	go build -o bin/$(APP_NAME) .

# Start production environment
.PHONY: start
start: build
	bin/$(APP_NAME)

# Clean up build files
.PHONY: clean
clean:
	rm -rf bin
