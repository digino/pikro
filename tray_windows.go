//go:build windows

package main

import (
	"fmt"

	"github.com/getlantern/systray"
)

// runWithTray starts the HTTP server inside the systray lifecycle on Windows.
// systray.Run blocks the main goroutine; onReady starts the server in a goroutine
// and sets up the tray menu. The OS tray icon replaces the console window flash.
func runWithTray(port int, startServer func() error) {
	url := fmt.Sprintf("http://localhost:%d", port)

	onReady := func() {
		systray.SetTitle("Pikro")
		systray.SetTooltip(fmt.Sprintf("Pikro — %s", url))

		mOpen := systray.AddMenuItem("Open Pikro", url)
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Stop the Pikro server")

		// Start the HTTP server in background.
		go func() {
			if err := startServer(); err != nil {
				systray.Quit()
			}
		}()

		// Open browser once on startup.
		go openBrowser(url)

		// Handle menu clicks.
		go func() {
			for {
				select {
				case <-mOpen.ClickedCh:
					openBrowser(url)
				case <-mQuit.ClickedCh:
					systray.Quit()
					return
				}
			}
		}()
	}

	onExit := func() {}

	systray.Run(onReady, onExit)
}
