package network

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	DefaultScheme = "http"
)

type Address struct {
	Network Network
	Host    string
	Port    int
}

func (a *Address) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a *Address) Decode(value string) error {
	if u, err := url.Parse(defaultToScheme(value, DefaultScheme)); err != nil {
		return err
	} else if supportedNetwork := SupportedNetwork(
		Network(u.Scheme),
	); u.IsAbs() && !supportedNetwork {
		return err
	} else if port, err := strconv.Atoi(u.Port()); u.Port() != "" && err != nil {
		return err
	} else {
		*a = Address{
			Network: Network(u.Scheme),
			Host:    u.Hostname(),
			Port:    port,
		}
	}

	return nil
}

func defaultToScheme(rawURL, defaultScheme string) string {
	// Force default http scheme, so net/url.Parse() doesn't
	// put both host and path into the (relative) path.
	if strings.Index(rawURL, "//") == 0 {
		// Leading double slashes (any scheme). Force http.
		rawURL = defaultScheme + ":" + rawURL
	}
	if strings.Index(rawURL, "://") == -1 {
		// Missing scheme. Force http.
		rawURL = defaultScheme + "://" + rawURL
	}
	return rawURL
}
