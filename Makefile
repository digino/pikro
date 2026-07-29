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
	GOOS=darwin  GOARCH=arm64  /usr/local/go/bin/go build $(LDFLAGS) -o dist/.pikro-mac-arm64-bin .
	GOOS=darwin  GOARCH=amd64  /usr/local/go/bin/go build $(LDFLAGS) -o dist/.pikro-mac-intel-bin .
	GOOS=windows GOARCH=amd64  /usr/local/go/bin/go build $(LDFLAGS) -o dist/$(BINARY).exe .
	GOOS=linux   GOARCH=amd64  /usr/local/go/bin/go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 .
	GOOS=linux   GOARCH=arm64  /usr/local/go/bin/go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 .
	@$(MAKE) _bundle_app ARCH=arm64 BIN=dist/.pikro-mac-arm64-bin OUT=dist/Pikro-mac-arm64.zip
	@$(MAKE) _bundle_app ARCH=amd64 BIN=dist/.pikro-mac-intel-bin OUT=dist/Pikro-mac-intel.zip
	@rm -f dist/.pikro-mac-arm64-bin dist/.pikro-mac-intel-bin
	@echo "Done → dist/ for $(VERSION)"

# Internal: wrap a Darwin binary in a .app bundle and zip it.
# Usage: make _bundle_app BIN=<binary> OUT=<zip>
_bundle_app:
	@APP=dist/Pikro.app; \
	 rm -rf "$$APP"; \
	 mkdir -p "$$APP/Contents/MacOS" "$$APP/Contents/Resources"; \
	 cp $(BIN) "$$APP/Contents/MacOS/pikro"; \
	 chmod +x "$$APP/Contents/MacOS/pikro"; \
	 cp assets/brand/pikro.icns "$$APP/Contents/Resources/pikro.icns"; \
	 printf '<?xml version="1.0" encoding="UTF-8"?>\n\
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">\n\
<plist version="1.0"><dict>\n\
  <key>CFBundleName</key><string>Pikro</string>\n\
  <key>CFBundleIdentifier</key><string>com.pikro.app</string>\n\
  <key>CFBundleVersion</key><string>$(VERSION)</string>\n\
  <key>CFBundleShortVersionString</key><string>$(VERSION)</string>\n\
  <key>CFBundleExecutable</key><string>pikro</string>\n\
  <key>CFBundleIconFile</key><string>pikro</string>\n\
  <key>CFBundlePackageType</key><string>APPL</string>\n\
  <key>LSMinimumSystemVersion</key><string>11.0</string>\n\
  <key>LSUIElement</key><true/>\n\
</dict></plist>\n' > "$$APP/Contents/Info.plist"; \
	 cd dist && zip -r ../$(OUT) Pikro.app --quiet; \
	 rm -rf "$$APP"

clean:
	rm -f $(BINARY)
	rm -rf dist/
	rm -rf web/dist/
