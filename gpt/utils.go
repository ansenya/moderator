package gpt

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func GetClientWithSocks5Proxy() (*http.Client, error) {
	proxyUrl := os.Getenv("PROXY_URL")
	if proxyUrl == "" {
		return nil, errors.New("PROXY_URL environment variable not set")
	}
	u, err := url.Parse(proxyUrl)
	if err != nil {
		return nil, fmt.Errorf("error parsing proxy address: %w", err)
	}

	auth := &proxy.Auth{}
	if u.User != nil {
		auth.User = u.User.Username()
		auth.Password, _ = u.User.Password()
	}

	dialer, err := proxy.SOCKS5("tcp", u.Host, auth, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("error creating SOCKS5 dialer: %w", err)
	}

	httpTransport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
	}

	client := &http.Client{
		Transport: httpTransport,
		Timeout:   30 * time.Second,
	}

	return client, nil
}
