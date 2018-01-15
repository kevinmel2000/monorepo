package httpreq

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type BookHTTPClient struct {
	client  *http.Client
	Options ClientOptions
}

func NewBookClient(options ClientOptions) (*BookHTTPClient, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}
	b := &BookHTTPClient{
		client:  NewHTTPClient(options.Timeout),
		Options: options,
	}
	return b, nil
}

type Book struct {
	Name   string
	Author string
}

type BookHTTPResult struct {
	Status string
	Errors []string
}

func (book *BookHTTPClient) AddBook(ctx context.Context, bookParam Book) (*BookHTTPResult, error) {
	url := book.Options.BaseURL + "/book/v1/add"
	jsonContent, err := json.Marshal(bookParam)
	if err != nil {
		return nil, nil
	}
	buff := bytes.NewBuffer(jsonContent)

	req, err := NewRequestWithHostHeader("POST", url, book.Options.HostHeader, buff)
	if err != nil {
		return nil, err
	}
	resp, err := book.client.Do(req)
	if err != nil {
		return nil, err
	}

	jsonResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	b := new(BookHTTPResult)
	err = json.Unmarshal(jsonResp, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
