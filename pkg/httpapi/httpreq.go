package httpapi

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

var (
	defaultHttpClient *http.Client
)

type ClientOptions struct {
	BaseURL    string
	HostHeader string
	Timeout    time.Duration
}

func (c *ClientOptions) Validate() error {
	if c.BaseURL == "" {
		return errors.New("BaseURL cannot be empty")
	}
	// if hostHeader is not empty, address must be localhost or 127.0.0.1
	if c.HostHeader != "" {
		if !strings.Contains(c.BaseURL, "localhost") && !strings.Contains(c.BaseURL, "127.0.0.1") {
			return errors.New("BaseURL for HostHeader limited to localhost/127.0.0.1")
		}
	}
	return nil
}

func init() {
	defaultHttpClient = &http.Client{
		Timeout: time.Second * 8,
	}
}

func NewHTTPClient(timeout time.Duration) *http.Client {
	if timeout == time.Duration(0) {
		return defaultHttpClient
	}
	c := &http.Client{
		Timeout: timeout,
	}
	return c
}
