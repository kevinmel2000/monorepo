package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/lab46/example/bookapp/book"
	"github.com/lab46/example/pkg/http/httpclient"
)

type BookHTTPAPI struct {
	httpClient *http.Client
	Options    httpclient.ClientOptions
}

func NewHTTPClient(options httpclient.ClientOptions) (*BookHTTPAPI, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}
	b := &BookHTTPAPI{
		httpClient: httpclient.NewClient(options.Timeout),
		Options:    options,
	}
	return b, nil
}

type Book struct {
	ID     int64
	Name   string
	Author string
}

type BookHTTPResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors"`
}

func (bk *BookHTTPAPI) AddBook(ctx context.Context, bookParam Book) (*BookHTTPResponse, error) {
	url := bk.Options.BaseURL + "/book/v1/add"
	jsonContent, err := json.Marshal(bookParam)
	if err != nil {
		return nil, nil
	}
	buff := bytes.NewBuffer(jsonContent)

	req, err := httpclient.NewRequestWithHostHeader(http.MethodPost, url, bk.Options.HostHeader, buff)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := bk.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	b := new(BookHTTPResponse)
	err = json.Unmarshal(jsonResp, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (bk *BookHTTPAPI) GetBookByID(ctx context.Context, id int64) (book.Book, error) {
	b := book.Book{}
	url := bk.Options.BaseURL + "/book/v1/get"

	bookID := strconv.FormatInt(id, 64)
	reqURL, err := httpclient.ParseURL(url, "id", bookID)
	if err != nil {
		return b, err
	}

	req, err := httpclient.NewRequestWithHostHeader(http.MethodPost, reqURL, bk.Options.HostHeader, nil)
	if err != nil {
		return b, err
	}
	req = req.WithContext(ctx)

	resp, err := bk.httpClient.Do(req)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()

	jsonResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}
	bookResp := BookHTTPResponse{}
	err = json.Unmarshal(jsonResp, &bookResp)
	if err != nil {
		return b, err
	}
	if len(bookResp.Errors) > 0 {
		return b, errors.New(bookResp.Errors[0])
	}

	b = bookResp.Data.(book.Book)
	return b, nil
}

func (bk *BookHTTPAPI) GetBookLis(ctx context.Context) ([]Book, error) {
	return nil, nil
}
