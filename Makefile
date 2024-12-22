# Project name
PROJECT_NAME := estimator

# Main directories
BACKEND_DIR := ./backend
FRONTEND_DIR := ./frontend

# BE commands
GO_RUN := cd $(BACKEND_DIR) && go run cmd/app/main.go
GO_TEST := cd $(BACKEND_DIR) && go test ./...
GO_BUILD := cd $(BACKEND_DIR) && go build -o bin/$(PROJECT_NAME) cmd/app/main.go

# FE commands
NPM_INSTALL := cd $(FRONTEND_DIR) && npm install
NPM_BUILD := cd $(FRONTEND_DIR) && npm run build
NPM_DEV := cd $(FRONTEND_DIR) && npm run dev

# 🛠️ Main commands
.PHONY: all backend frontend

all: backend frontend

backend: 
	$(GO_BUILD)

frontend:
	$(NPM_BUILD)

# 👷‍♀️ BE commands
run-backend:
	$(GO_RUN)

build-backend:
	$(GO_BUILD)

test-backend:
	$(GO_TEST)

# 📦 FE commands
install-frontend:
	$(NPM_INSTALL)

build-frontend:
	$(NPM_BUILD)

dev-frontend:
	$(NPM_DEV)

# 🧹 clean up
clean:
	rm -rf $(BACKEND_DIR)/bin
	rm -rf $(FRONTEND_DIR)/node_modules $(FRONTEND_DIR)/dist
