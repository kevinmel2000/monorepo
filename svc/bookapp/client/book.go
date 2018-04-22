package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/lab46/monorepo/gopkg/http/httpclient"
	"github.com/lab46/monorepo/svc/bookapp/book"
)

var bookBaseURL = "book"

// BookHTTPAPI struct
type BookHTTPAPI struct {
	httpClient *httpclient.Client
	Options    httpclient.ClientOptions
}

// NewHTTPClient create new client for book http client
func NewHTTPClient(options httpclient.ClientOptions) (*BookHTTPAPI, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}
	client, err := httpclient.New(httpclient.ClientOptions{
		Timeout: options.Timeout,
	})
	if err != nil {
		return nil, err
	}
	b := &BookHTTPAPI{
		httpClient: client,
		Options:    options,
	}
	return b, nil
}

// Book struct
type Book struct {
	ID     int64
	Name   string
	Author string
}

// BookHTTPResponse struct
type BookHTTPResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors"`
}

// AddBook function
func (bk *BookHTTPAPI) AddBook(ctx context.Context, bookParam Book) (*BookHTTPResponse, error) {
	url := bookBaseURL + "/book/v1/add"
	jsonContent, err := json.Marshal(bookParam)
	if err != nil {
		return nil, nil
	}
	buff := bytes.NewBuffer(jsonContent)

	resp, err := bk.httpClient.DoRequest(ctx, httpclient.HTTPReq{
		Method: http.MethodPut,
		URL:    url,
		Body:   buff,
	})
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

// GetBookByID function
func (bk *BookHTTPAPI) GetBookByID(ctx context.Context, id int64) (book.Book, error) {
	b := book.Book{}
	url := bookBaseURL + "/book/v1/get"

	bookID := strconv.FormatInt(id, 64)

	resp, err := bk.httpClient.DoRequest(ctx, httpclient.HTTPReq{
		Method: http.MethodGet,
		URL:    url,
		URLParams: map[string]string{
			"id": bookID,
		},
	})
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

// GetBookList function
func (bk *BookHTTPAPI) GetBookList(ctx context.Context) ([]Book, error) {
	return nil, nil
}
