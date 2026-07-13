//go:build !linux && !darwin && !freebsd && !openbsd && !netbsd

package discovery

import (
	"fmt"
	"log"
	"net"
	"time"
)

// Scan on Windows: UDP sockets allow broadcast by default so we can use
// net.ListenPacket directly. Bind to port 5678 — MNDP replies are broadcast
// to that port, not back to the ephemeral source port.
func Scan(timeout time.Duration) ([]Device, error) {
	conn, err := net.ListenPacket("udp4", "0.0.0.0:5678")
	if err != nil {
		log.Printf("[mndp] ListenPacket error: %v", err)
		return nil, fmt.Errorf("listen: %w", err)
	}
	defer conn.Close()

	log.Printf("[mndp] socket ready on %s", conn.LocalAddr())

	if err := sendBroadcast(conn); err != nil {
		return nil, err
	}

	return collectResponses(conn, timeout), nil
}
