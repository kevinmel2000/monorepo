package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/lab46/example/pkg/httpclient"
)

type RentHTTPAPI struct {
	httpClient *http.Client
	Options    httpclient.ClientOptions
}

func NewRentClient(options httpclient.ClientOptions) (*RentHTTPAPI, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}
	b := &RentHTTPAPI{
		httpClient: httpclient.NewHTTPClient(options.Timeout),
		Options:    options,
	}
	return b, nil
}

type Rent struct {
	BookID     int64
	MemberName string
}

type RentHTTPResponse struct {
	Status string
	Errors []string
}

func (rnt *RentHTTPAPI) RentBook(ctx context.Context, rent Rent) (*RentHTTPResponse, error) {
	url := rnt.Options.BaseURL + "/rent/v1/book"
	jsonContent, err := json.Marshal(rent)
	if err != nil {
		return nil, nil
	}
	buff := bytes.NewBuffer(jsonContent)

	req, err := httpclient.NewRequestWithHostHeader("POST", url, rnt.Options.HostHeader, buff)
	if err != nil {
		return nil, err
	}
	resp, err := rnt.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	jsonResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	b := new(RentHTTPResponse)
	err = json.Unmarshal(jsonResp, b)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
