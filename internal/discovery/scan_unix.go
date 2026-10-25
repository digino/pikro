//go:build linux || darwin || freebsd || openbsd || netbsd

package discovery

import (
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"syscall"
	"time"
)

// Scan creates a UDP socket bound to port 5678 with SO_BROADCAST + SO_REUSEADDR.
// MikroTik MNDP responses are sent to the broadcast address on port 5678, not
// back to the ephemeral source port, so we must listen on 5678 to receive them.
func Scan(timeout time.Duration) ([]Device, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return nil, fmt.Errorf("socket: %w", err)
	}

	_ = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1); err != nil {
		syscall.Close(fd)
		return nil, fmt.Errorf("SO_BROADCAST: %w", err)
	}

	if err := syscall.Bind(fd, &syscall.SockaddrInet4{Port: 5678}); err != nil {
		syscall.Close(fd)
		return nil, fmt.Errorf("bind :5678: %w", err)
	}

	// Wrap fd: net.FilePacketConn dups it internally so we can close our copy.
	f := os.NewFile(uintptr(fd), "mndp")
	conn, err := net.FilePacketConn(f)
	f.Close()
	if err != nil {
		syscall.Close(fd)
		return nil, fmt.Errorf("FilePacketConn: %w", err)
	}
	defer conn.Close()

	// Directed broadcasts still "succeed" when macOS denies Local Network
	// access, so the scan would just come back empty. A limited broadcast fails
	// loudly instead — use it to tell the user what to fix.
	if runtime.GOOS == "darwin" {
		limited := &net.UDPAddr{IP: net.IPv4bcast, Port: 5678}
		if _, err := conn.WriteTo(mndpProbe, limited); errors.Is(err, syscall.EHOSTUNREACH) {
			return nil, ErrLocalNetworkDenied
		}
	}

	if err := sendBroadcast(conn); err != nil {
		return nil, err
	}

	return collectResponses(conn, timeout), nil
}
