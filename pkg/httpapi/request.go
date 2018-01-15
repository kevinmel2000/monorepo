package httpapi

import (
	"io"
	"net/http"
)

func NewRequestWithHostHeader(method, urlStr, hostHeader string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Host", hostHeader)
	return req, err
}
