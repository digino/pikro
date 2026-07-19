package discovery

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

type Device struct {
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	Identity string `json:"identity"`
	Version  string `json:"version"`
	Board    string `json:"board"`
	Platform string `json:"platform"`
	Uptime   string `json:"uptime"`
	Iface    string `json:"iface"`
}

// Scan is implemented per-platform in scan_unix.go / scan_windows.go.
// It broadcasts an MNDP request and collects responses for timeout duration.

// parseResponse decodes a TLV-encoded MNDP response packet.
// Format: [seq:4][type:2][len:2][value:len]...
func parseResponse(data []byte) Device {
	var dev Device
	if len(data) < 4 {
		return dev
	}
	offset := 4
	for offset+4 <= len(data) {
		tlvType := binary.BigEndian.Uint16(data[offset:])
		tlvLen := int(binary.BigEndian.Uint16(data[offset+2:]))
		offset += 4
		if offset+tlvLen > len(data) {
			break
		}
		val := data[offset : offset+tlvLen]
		offset += tlvLen

		switch tlvType {
		case 1:
			if tlvLen == 6 {
				dev.MAC = net.HardwareAddr(val).String()
			}
		case 5:
			dev.Identity = string(val)
		case 7:
			dev.Version = string(val)
		case 8:
			dev.Platform = string(val)
		case 10:
			if tlvLen == 4 {
				dev.Uptime = formatUptime(binary.LittleEndian.Uint32(val))
			}
		case 12:
			dev.Board = string(val)
		case 16:
			dev.Iface = string(val)
		case 17:
			if tlvLen == 4 {
				dev.IP = net.IP(val).String()
			}
		}
	}
	return dev
}

func collectResponses(conn net.PacketConn, timeout time.Duration) []Device {
	conn.SetReadDeadline(time.Now().Add(timeout))

	seen := map[string]bool{}
	var devices []Device
	buf := make([]byte, 4096)

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			break
		}
		src := addr.(*net.UDPAddr).IP.String()

		dev := parseResponse(buf[:n])
		// Use UDP source IP if TLV didn't carry one (or carried 0.0.0.0).
		if dev.IP == "" || dev.IP == "0.0.0.0" {
			if src != "0.0.0.0" {
				dev.IP = src
			}
		}
		// Skip empty probes (our own broadcast echo with no TLV data)
		if dev.MAC == "" && dev.Identity == "" && dev.Board == "" {
			continue
		}
		// Dedup by MAC (most stable); fall back to IP
		key := dev.MAC
		if key == "" {
			key = dev.IP
		}
		if !seen[key] {
			seen[key] = true
			devices = append(devices, dev)
		}
	}

	log.Printf("[mndp] scan complete: %d device(s) found", len(devices))
	return devices
}

// sendBroadcast sends an MNDP probe to every broadcast-capable IPv4 interface's
// directed broadcast address (e.g. 192.168.1.255). This is more reliable than
// 255.255.255.255 (limited broadcast), which fails with EHOSTUNREACH on macOS
// and some Linux setups because the kernel finds no matching route.
func sendBroadcast(conn net.PacketConn) error {
	probe := []byte{0x00, 0x00, 0x00, 0x00}
	sent := false
	var lastErr error

	ifaces, err := net.Interfaces()
	if err != nil {
		// Fall back to limited broadcast if interface enumeration fails.
		dst, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:5678")
		_, err = conn.WriteTo(probe, dst)
		return err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagBroadcast == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP.To4()
			if ip == nil {
				continue
			}
			// Compute directed broadcast: host bits all set to 1.
			bcast := make(net.IP, 4)
			for i := range bcast {
				bcast[i] = ip[i] | ^ipNet.Mask[i]
			}
			dst := &net.UDPAddr{IP: bcast, Port: 5678}
			if _, werr := conn.WriteTo(probe, dst); werr != nil {
				lastErr = werr
			} else {
				sent = true
			}
		}
	}

	if !sent {
		if lastErr != nil {
			return fmt.Errorf("send: %w", lastErr)
		}
		return fmt.Errorf("send: no broadcast-capable IPv4 interface found")
	}
	return nil
}

func formatUptime(seconds uint32) string {
	d := seconds / 86400
	h := (seconds % 86400) / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	if d > 0 {
		return fmt.Sprintf("%dd %02dh%02dm%02ds", d, h, m, s)
	}
	return fmt.Sprintf("%02dh%02dm%02ds", h, m, s)
}
