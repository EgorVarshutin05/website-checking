package checker

import (
	"net/netip"
	"testing"
)

func TestIsAllowedIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		allowed bool
	}{
		{"Public IPv4", "8.8.8.8", true},
		{"Public IPv4 Cloudflare", "1.1.1.1", true},
		{"Loopback IPv4", "127.0.0.1", false},
		{"Private IPv4", "192.168.1.1", false},
		{"Private IPv4 10", "10.0.0.5", false},
		{"Cloud metadata", "169.254.169.254", false},
		{"Shared address space", "100.64.0.1", false},
		{"Documentation IPv4", "192.0.2.1", false},
		{"Loopback IPv6", "::1", false},
		{"Private IPv6", "fd00::1", false},
		{"Documentation IPv6", "2001:db8::1", false},
		{"IPv4-mapped loopback", "::ffff:127.0.0.1", false},
		{"Public IPv6", "2606:4700:4700::1111", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := netip.MustParseAddr(tt.ip)

			got := IsAllowedIP(ip)

			if got != tt.allowed {
				t.Errorf(
					"IsAllowedIP(%s) = %v, ожидалось %v",
					tt.ip,
					got,
					tt.allowed,
				)
			}
		})
	}
}
