//go:build darwin

package main

import (
	"fmt"
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed assets/brand/pikro.icns
var iconBytes []byte

func runWithTray(port int, startServer func() error) {
	url := fmt.Sprintf("http://localhost:%d", port)

	onReady := func() {
		systray.SetIcon(iconBytes)
		systray.SetTitle("Pikro")
		systray.SetTooltip(fmt.Sprintf("Pikro — %s", url))

		mOpen := systray.AddMenuItem("Open Pikro", url)
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Stop the Pikro server")

		go func() {
			if err := startServer(); err != nil {
				systray.Quit()
			}
		}()

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

	systray.Run(onReady, func() {})
}
