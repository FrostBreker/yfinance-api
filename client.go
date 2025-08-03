package yfinance_api

import (
	"crypto/rand"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/url"
	"sync"
)

type YFinanceAPI struct {
	Client *Client
}

type Client struct {
	client  *http.Client
	cookies []*http.Cookie
	crumb   string
}

var instance *Client
var once sync.Once

func getClient() *Client {
	once.Do(func() {
		instance = &Client{client: &http.Client{}, cookies: []*http.Cookie{}, crumb: ""}
	})
	return instance
}

func (c *Client) Get(url string, params url.Values) (*http.Response, error) {
	c.getCrumb()
	return c.get(url, params)
}

func (c *Client) get(url string, params url.Values) (*http.Response, error) {
	if c.crumb != "" {
		params.Add("crumb", c.crumb)
	}
	url = fmt.Sprintf("%s?%s", url, params.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Failed to create request", "err", err)
		return nil, err
	}

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	// Use crypto/rand for secure random number generation
	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(UserAgents))))
	if err != nil {
		slog.Error("Failed to generate secure random number", "err", err)
		// Fallback to first user agent if random generation fails
		req.Header.Set("User-Agent", UserAgents[0])
	} else {
		req.Header.Set("User-Agent", UserAgents[randomIndex.Int64()])
	}

	resp, err := c.client.Do(req)
	if err != nil {
		slog.Error("Failed to get data from Yahoo Finance API", "err", err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) getCookie() {
	if len(c.cookies) > 0 {
		return
	}

	endpoint := "https://fc.yahoo.com"
	resp, err := c.get(endpoint, url.Values{})
	if err != nil {
		slog.Error("Failed to get cookie", "err", err)
		return
	}

	c.cookies = resp.Cookies()
}

func (c *Client) getCrumb() {
	if c.crumb != "" {
		return
	}

	c.getCookie()
	endpoint := fmt.Sprintf("%s/v1/test/getcrumb", BaseUrl)
	resp, err := c.get(endpoint, url.Values{})
	if err != nil {
		slog.Error("Failed to get crumb", "err", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error closing response body:", "err", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body:", "err", err)
		return
	}

	c.crumb = string(body)
}

// NewClient creates and returns a new YFinance API client instance
// This is the main entry point for users of the package
func NewClient() *YFinanceAPI {
	return &YFinanceAPI{
		Client: getClient(),
	}
}

// NewTicker creates a new ticker instance for the given symbol
// This is a convenience function that creates a client and ticker in one call
func NewTicker(symbol string) *Ticker {
	client := NewClient()
	return client.InstantiateTicker(symbol)
}
