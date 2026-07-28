package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"runtime"
	"time"

	"github.com/digino/pikro/internal/handlers"
	"github.com/digino/pikro/internal/server"
)

// web/dist is produced by `cd web && npm run build` before `go build`.
// The blank directory placeholder web/dist/.gitkeep ensures go:embed doesn't
// fail during development before the first build.
//
//go:embed all:web/dist
var embeddedWeb embed.FS

// Version is injected at build time: -ldflags "-X main.Version=1.0.0"
var Version = "dev"

func main() {
	dev    := flag.Bool("dev", false, "enable dev mode with request logging")
	noOpen := flag.Bool("no-open", false, "do not open browser on start")
	ver    := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *ver {
		fmt.Println(Version)
		return
	}

	if !*dev {
		log.SetOutput(io.Discard)
	}

	port := findAvailablePort(8080)
	url := fmt.Sprintf("http://localhost:%d", port)
	handlers.AppVersion = Version
	log.Printf("Pikro %s running at %s (press Ctrl+C to stop)\n", Version, url)

	startServer := func() error { return server.Start(port, embeddedWeb) }

	if runtime.GOOS == "windows" {
		// On Windows, runWithTray owns the lifecycle: it shows a tray icon,
		// opens the browser, and starts the server — no console flash.
		runWithTray(port, startServer)
	} else {
		if !*noOpen {
			go func() {
				time.Sleep(300 * time.Millisecond)
				openBrowser(url)
			}()
		}
		if err := startServer(); err != nil {
			log.Fatal(err)
		}
	}
}

func findAvailablePort(start int) int {
	for port := start; port < start+100; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			ln.Close()
			return port
		}
	}
	return start
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Printf("Could not open browser: %v — open %s manually\n", err, url)
	}
}

