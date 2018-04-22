package httpclient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// default var and error
var (
	defaultHTTPClient *http.Client

	ErrBaseURLEmpty          = errors.New("BaseURL cannot be empty")
	ErrMissingHTTPMethod     = errors.New("Missing HTTP method")
	ErrInvalidSidecarAddress = errors.New("Address for sidecar limited to localhost/127.0.0.1")
	ErrHostHeaderEmpty       = errors.New("Host-header cannot be empty if sidecar is used")
)

// Client struct
type Client struct {
	sidecarAddress string
	client         *http.Client
	options        ClientOptions
}

// ClientOptions for client
type ClientOptions struct {
	// this will force all http request go through sidecar
	UseSidecar bool
	// timeout for http client
	Timeout time.Duration
}

// Validate clientoptions
func (c *ClientOptions) Validate() error {
	// if hostHeader is not empty, address must be localhost or 127.0.0.1
	// if c.UseSidecar {
	// 	if !strings.Contains(c.SidecarAddress, "localhost") && !strings.Contains(c.SidecarAddress, "127.0.0.1") {
	// 		return ErrSidecarAddress
	// 	}
	// }
	return nil
}

func init() {
	defaultHTTPClient = &http.Client{
		Timeout: time.Second * 8,
	}
}

// New http client
func New(options ClientOptions) (*Client, error) {
	err := options.Validate()
	if err != nil {
		return nil, err
	}

	var (
		c       *http.Client
		sidecar string
	)
	if options.Timeout == time.Duration(0) {
		c = defaultHTTPClient
	} else {
		c = &http.Client{
			Timeout: options.Timeout,
		}
	}

	// set address of sidecar
	if options.UseSidecar {
		sidecar = "127.0.0.1:9001"
	}

	httpClient := Client{
		client:         c,
		options:        options,
		sidecarAddress: sidecar,
	}
	return &httpClient, nil
}

// ChangeSidecarAddress will change the address of sidecar
// this will be mostly used when a port for sidecar is vary
func (c *Client) ChangeSidecarAddress(address string) {
	c.sidecarAddress = address
}

// Do http request
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// Get request
func (c *Client) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

// Post request
func (c *Client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return c.client.Post(url, contentType, body)
}

// PostForm request
func (c *Client) PostForm(url string, data url.Values) (*http.Response, error) {
	return c.client.PostForm(url, data)
}

// HTTPReq request
type HTTPReq struct {
	URL         string
	HostHerader string
	Method      string
	URLParams   map[string]string
	Body        io.Reader
	sidecar     bool
}

func (req *HTTPReq) useSidecar() {
	req.sidecar = true
}

// URLParamsToKV return KV form of url params
func (req *HTTPReq) URLParamsToKV() []string {
	var kv []string
	urlParamsLength := len(req.URLParams)
	if urlParamsLength == 0 {
		return kv
	}
	kv = make([]string, urlParamsLength*2)

	counter := 0
	for key, val := range req.URLParams {
		kv[counter] = key
		kv[counter+1] = val
		counter += 2
	}
	return kv
}

// Validate HTTPReq
func (req *HTTPReq) Validate() error {
	if req.Method == "" {
		return ErrMissingHTTPMethod
	}

	if req.sidecar && (!strings.Contains(req.URL, "localhost") && !strings.Contains(req.URL, "127.0.0.1")) {
		return ErrInvalidSidecarAddress
	}
	// expect to always use sidecar when sidecar is enabled
	if req.sidecar && req.HostHerader == "" {
		return ErrHostHeaderEmpty
	}
	return nil
}

// DoRequest will do a http request
func (c *Client) DoRequest(ctx context.Context, httpreq HTTPReq) (*http.Response, error) {
	if c.sidecarAddress != "" {
		httpreq.useSidecar()
	}
	err := httpreq.Validate()
	if err != nil {
		return nil, err
	}

	baseURL := httpreq.URL
	if c.sidecarAddress != "" {
		baseURL = c.sidecarAddress
	}
	url, err := ParseURL(baseURL, httpreq.URLParamsToKV()...)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	// all request via NewRequestViaHostHeader
	if httpreq.HostHerader != "" {
		req, err = NewRequestWithHostHeader(httpreq.Method, url, httpreq.HostHerader, httpreq.Body)
	} else {
		req, err = http.NewRequest(httpreq.Method, url, httpreq.Body)
	}
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
