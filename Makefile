# Define variables
FRONTEND_DIR := frontend
DIST_DIR := $(FRONTEND_DIR)/dist
SERVER_DIR := pkg/server

# Build the frontend using Vite
build-frontend:
	cd $(FRONTEND_DIR) && npm install && npm run build

# Copy the dist folder to the server directory
copy-dist:
	cp -r $(DIST_DIR) $(SERVER_DIR)

# Build target that runs both tasks
build: build-frontend copy-dist

dev:
	air

.PHONY: build-frontend copy-dist build dev
