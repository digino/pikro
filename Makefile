BINARY  := pikro
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"
# -H windowsgui builds the exe as a GUI subsystem app so Windows doesn't pop a
# console window on double-click. Windows-only — other GOOS ignore -H.
WIN_LDFLAGS := -ldflags "-X main.Version=$(VERSION) -H windowsgui"
# Release builds add -s -w (strip debug symbols + DWARF tables) — ~30%
# smaller binaries, no functional loss for a distributed build.
RELEASE_LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"
RELEASE_WIN_LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION) -H windowsgui"
RSRC_VERSION := v0.10.2

.PHONY: dev backend build release clean _win_rsrc

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

# Cross-compile for macOS (Apple Silicon only) and Windows.
release: _win_rsrc
	@echo "Releasing $(VERSION)..."
	@mkdir -p dist
	@cd web && npm run build
	GOOS=darwin  GOARCH=arm64  /usr/local/go/bin/go build $(RELEASE_LDFLAGS) -o dist/.pikro-mac-arm64-bin .
	GOOS=windows GOARCH=amd64  /usr/local/go/bin/go build $(RELEASE_WIN_LDFLAGS) -o dist/$(BINARY).exe .
	@$(MAKE) _bundle_app ARCH=arm64 BIN=dist/.pikro-mac-arm64-bin OUT=dist/Pikro-mac-arm64.dmg
	@rm -f dist/.pikro-mac-arm64-bin rsrc_windows_amd64.syso
	@echo "Done → dist/ for $(VERSION)"

# Internal: generates rsrc_windows_amd64.syso, embedding the Common-Controls
# manifest (needed so the launcher window's buttons render with modern
# styling instead of Windows-95-style widgets) and the .exe icon.
# Picked up automatically by `go build` for GOOS=windows — no linker flag needed.
_win_rsrc:
	go run github.com/akavel/rsrc@$(RSRC_VERSION) \
		-manifest pikro.exe.manifest \
		-ico assets/brand/pikro.ico \
		-arch amd64 \
		-o rsrc_windows_amd64.syso

# Internal: wrap a Darwin binary in a .app bundle and package it as a .dmg
# with a symlink to /Applications, so opening the .dmg gives the standard
# drag-to-install experience instead of a plain unzip.
# Usage: make _bundle_app BIN=<binary> OUT=<dmg>
_bundle_app:
	@APP=dist/Pikro.app; \
	 STAGE=dist/.dmg-stage; \
	 rm -rf "$$APP" "$$STAGE" "$(OUT)"; \
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
	 codesign --force --deep --sign - "$$APP"; \
	 mkdir -p "$$STAGE"; \
	 cp -R "$$APP" "$$STAGE/"; \
	 ln -s /Applications "$$STAGE/Applications"; \
	 hdiutil create -volname "Pikro" -srcfolder "$$STAGE" -ov -format UDZO "$(OUT)" -quiet; \
	 rm -rf "$$APP" "$$STAGE"

clean:
	rm -f $(BINARY)
	rm -f rsrc_windows_amd64.syso
	rm -rf dist/
	rm -rf web/dist/
