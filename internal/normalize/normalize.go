package normalize

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func URL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	u.Scheme = strings.ToLower(u.Scheme)
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("only HTTP and HTTPS are allowed")
	}

	if u.Opaque != "" || u.Hostname() == "" {
		return "", fmt.Errorf("missing hostname")
	}

	if u.User != nil {
		return "", fmt.Errorf("URL credentials are not allowed")
	}

	hostname := strings.ToLower(u.Hostname())
	port := u.Port()

	if port != "" {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return "", fmt.Errorf("invalid port")
		}
	}

	if (u.Scheme == "http" && port == "80") ||
		(u.Scheme == "https" && port == "443") {
		port = ""
	}

	if strings.Contains(hostname, ":") {
		hostname = "[" + hostname + "]"
	}
	if port != "" {
		u.Host = net.JoinHostPort(strings.Trim(hostname, "[]"), port)
	} else {
		u.Host = hostname
	}

	u.Fragment = ""
	u.RawFragment = ""

	if u.Path == "/" {
		u.Path = ""
	}

	return u.String(), nil
}
