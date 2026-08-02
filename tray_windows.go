//go:build windows

package main

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// runWithTray shows a small launcher window immediately (so Windows users
// notice Pikro started — beta feedback was that the tray icon alone goes
// unnoticed) alongside the tray icon. Closing the window hides it rather
// than quitting; only the window's "Stop" button or the tray's "Quit" stop
// the server, and both do so via os.Exit — neither systray.Quit() nor
// walk.App().Exit() actually terminates the process, they only stop their
// own library's event loop.
func runWithTray(port int, startServer func() error) {
	url := fmt.Sprintf("http://localhost:%d", port)

	go systrayRun(url)

	go func() {
		if err := startServer(); err != nil {
			os.Exit(1)
		}
	}()

	showLauncherWindow(url)
}

func systrayRun(url string) {
	systray.Run(func() {
		systray.SetIcon(iconBytes)
		systray.SetTitle("Pikro")
		systray.SetTooltip(fmt.Sprintf("Pikro — %s", url))

		mOpen := systray.AddMenuItem("Open Pikro", url)
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Stop the Pikro server")

		go func() {
			for {
				select {
				case <-mOpen.ClickedCh:
					openBrowser(url)
				case <-mQuit.ClickedCh:
					os.Exit(0)
				}
			}
		}()
	}, func() {})
}

func showLauncherWindow(url string) {
	var mw *walk.MainWindow

	MainWindow{
		AssignTo: &mw,
		Title:    "Pikro",
		MinSize:  Size{Width: 320, Height: 160},
		Size:     Size{Width: 320, Height: 160},
		Layout:   VBox{Margins: Margins{Left: 16, Top: 16, Right: 16, Bottom: 16}, Spacing: 8},
		Children: []Widget{
			Label{
				Text:      "Pikro is running",
				Font:      Font{PointSize: 12, Bold: true},
				Alignment: AlignHCenterVCenter,
			},
			Label{
				Text:      url,
				Alignment: AlignHCenterVCenter,
			},
			VSpacer{Size: 8},
			PushButton{
				Text: "Open Pikro",
				OnClicked: func() {
					openBrowser(url)
				},
			},
			PushButton{
				Text: "Stop Pikro server",
				OnClicked: func() {
					os.Exit(0)
				},
			},
		},
	}.Create()

	// Minimize to tray instead of exiting the process when the window is closed.
	mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		*canceled = true
		mw.Hide()
	})

	mw.Run()
}
