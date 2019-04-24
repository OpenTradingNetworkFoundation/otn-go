package coinmarketcap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	cmcURL             = "https://api.coinmarketcap.com/v2"
	graphURL           = "https://graphs2.coinmarketcap.com/currencies"
	defaultHTTPTimeout = 3 * time.Second
)

// Constants for supported sort orders (TickersOptions.Sort field)
const (
	TickersSortID     = "id"
	TickersSortRank   = "rank"
	TickersSortVolume = "volume_24h"
	TickersSortChange = "percent_change_24h"
)

// Options allows to override default client parameters
type Options struct {
	URL     string
	Timeout time.Duration
}

// Client provides necessary API calls that OTN project uses
type Client struct {
	url     string
	http    *http.Client
	timeout time.Duration
}

// NewClient returns new Client instance. If url is empty, direct coinmarketcap URL would be used
func NewClient(opts *Options) *Client {
	client := &Client{
		url:     cmcURL,
		timeout: defaultHTTPTimeout,
	}

	if opts != nil {
		if opts.Timeout != 0 {
			client.timeout = opts.Timeout
		}
		if opts.URL != "" {
			client.url = opts.URL
		}
	}

	client.http = &http.Client{
		Timeout: client.timeout,
	}

	return client
}

func doReq(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

// HTTP request helper
func (c *Client) makeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := doReq(c.http, req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

type listingsResponse struct {
	Data []*Listing `json:"data"`
}

func (c *Client) Listings() ([]*Listing, error) {
	data, err := c.makeReq(fmt.Sprintf("%s/listings", c.url))
	if err != nil {
		return nil, err
	}

	var resp listingsResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

type tickersResponse struct {
	Data     map[string]*Ticker `json:"data,omitempty"`
	Metadata struct {
		Timestamp           int64  `json:"timestamp"`
		Error               string `json:"error,omitempty"`
		NumCryptoCurrencies int    `json:"num_cryptocurrencies,omitempty"`
	}
}

// TickersOptions
type TickersOptions struct {
	Start   int
	Limit   int
	Sort    string
	Convert string
}

func (c *Client) Tickers(opts *TickersOptions) (map[string]*Ticker, error) {
	var params []string
	if opts.Limit != 0 {
		params = append(params, fmt.Sprintf("limit=%d", opts.Limit))
	}
	if opts.Start != 0 {
		params = append(params, fmt.Sprintf("start=%d", opts.Start))
	}
	if opts.Sort != "" {
		params = append(params, fmt.Sprintf("sort=%s", opts.Sort))
	}
	if opts.Convert != "" {
		params = append(params, fmt.Sprintf("convert=%s", opts.Convert))
	}

	reqURL := fmt.Sprintf("%s/ticker", c.url)
	if len(params) != 0 {
		reqURL = reqURL + "?" + strings.Join(params, "&")
	}

	data, err := c.makeReq(reqURL)
	if err != nil {
		return nil, err
	}

	var resp tickersResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Metadata.Error != "" {
		return nil, errors.New(resp.Metadata.Error)
	}

	return resp.Data, nil
}

// TickerOptions
type TickerOptions struct {
	ID      int
	Convert string
}

type tickerResponse struct {
	Data     *Ticker `json:"data,omitempty"`
	Metadata struct {
		Timestamp int64  `json:"timestamp"`
		Error     string `json:"error,omitempty"`
	}
}

func (c *Client) Ticker(opts *TickerOptions) (*Ticker, error) {
	reqURL := fmt.Sprintf("%s/ticker/%d/", c.url, opts.ID)
	if opts.Convert != "" {
		reqURL += fmt.Sprintf("?convert=%s", opts.Convert)
	}

	data, err := c.makeReq(reqURL)
	if err != nil {
		return nil, err
	}

	var resp tickerResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Metadata.Error != "" {
		return nil, errors.New(resp.Metadata.Error)
	}

	return resp.Data, nil
}
