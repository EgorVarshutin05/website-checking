package checker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"time"
)

func NewSecureHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout: 3 * time.Second,
	}

	transport := &http.Transport{
		Proxy: nil,

		TLSHandshakeTimeout: 3 * time.Second,

		DialContext: func(
			ctx context.Context,
			network, address string,
		) (net.Conn, error) {
			host, port, err := net.SplitHostPort(address)
			if err != nil {
				return nil, err
			}

			if port != "80" && port != "443" {
				return nil, fmt.Errorf("порт %s запрещён", port)
			}

			var addresses []netip.Addr

			if ip, err := netip.ParseAddr(host); err == nil {
				addresses = []netip.Addr{ip}
			} else {
				addresses, err = net.DefaultResolver.LookupNetIP(
					ctx, "ip", host,
				)
				if err != nil {
					return nil, err
				}
			}

			if len(addresses) == 0 {
				return nil, errors.New("не удалось определить IP-адрес")
			}

			// Отклоняем хост, если хотя бы один IP запрещён.
			for _, ip := range addresses {
				if !IsAllowedIP(ip) {
					return nil, fmt.Errorf(
						"подключение к %s запрещено",
						ip,
					)
				}
			}

			// Подключаемся к проверенному IP, без повторного DNS.
			var lastErr error

			for _, ip := range addresses {
				conn, err := dialer.DialContext(
					ctx,
					"tcp",
					net.JoinHostPort(ip.String(), port),
				)
				if err == nil {
					return conn, nil
				}
				lastErr = err
			}

			return nil, lastErr
		},
	}

	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,

		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
