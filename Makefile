BINARY  := pikro
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

.PHONY: dev backend build release clean

# Run Go backend + Vite dev server side by side.
# Ctrl+C kills both via the trap.
dev:
	@trap 'kill 0' INT TERM; \
	  echo ""; \
	  echo "  Go backend  → http://localhost:8080  (API only)"; \
	  echo "  UI dev      → http://localhost:5173  (opening this one)"; \
	  echo ""; \
	  /usr/local/go/bin/go run . -dev -no-open & \
	  cd web && npm run dev --open & \
	  wait

# Run only the Go backend (useful when Vite is already running separately)
backend:
	/usr/local/go/bin/go run . -dev

# Full production build: Vue → dist, then Go embeds it
build:
	@echo "Building frontend..."
	@cd web && npm run build
	@echo "Building binary ($(VERSION))..."
	@/usr/local/go/bin/go build $(LDFLAGS) -o $(BINARY) .
	@echo "Done → ./$(BINARY) ($(shell du -sh $(BINARY) | cut -f1))"

# Cross-compile for all platforms
release:
	@echo "Releasing $(VERSION)..."
	@mkdir -p dist
	@cd web && npm run build
	GOOS=darwin  GOARCH=arm64  /usr/local/go/bin/go build $(LDFLAGS) -o dist/$(BINARY)-mac-arm64 .
	GOOS=darwin  GOARCH=amd64  /usr/local/go/bin/go build $(LDFLAGS) -o dist/$(BINARY)-mac-intel .
	GOOS=windows GOARCH=amd64  /usr/local/go/bin/go build $(LDFLAGS) -o dist/$(BINARY).exe .
	GOOS=linux   GOARCH=amd64  /usr/local/go/bin/go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 .
	GOOS=linux   GOARCH=arm64  /usr/local/go/bin/go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 .
	@echo "Binaries in ./dist/ for $(VERSION)"

clean:
	rm -f $(BINARY)
	rm -rf dist/
	rm -rf web/dist/
