package normalize

import (
	"fmt"
	"net/url"
	"strings"
)

func URL(raw string) (string, error) {
	if i := strings.IndexByte(raw, '#'); i >= 0 {
		raw = raw[:i]
	}

	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return "", fmt.Errorf("invalid URL")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("only HTTP and HTTPS are allowed")
	}

	if u.Hostname() == "" {
		return "", fmt.Errorf("missing hostname")
	}

	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	if (u.Scheme == "http" && u.Port() == "80") ||
		(u.Scheme == "https" && u.Port() == "443") {
		u.Host = u.Hostname()
	}

	if u.Path == "/" {
		u.Path = ""
	}

	return u.String(), nil
}
