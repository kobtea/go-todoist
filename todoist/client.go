package todoist

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"path"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	Token      string
	SyncToken  string
	Logger     *log.Logger
}

func NewClient(endpoint, token, sync_token string, logger *log.Logger) (*Client, error) {
	if len(endpoint) == 0 {
		endpoint = "https://todoist.com/API/v7"
	}
	parsed_endpoint, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient

	if len(token) == 0 {
		return nil, errors.New("Missing API Token")
	}

	if len(sync_token) == 0 {
		sync_token = "*"
	}

	if logger == nil {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	return &Client{parsed_endpoint, client, token, sync_token, logger}, nil
}

func (c *Client) NewRequest(ctx context.Context, method, spath string, values url.Values) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)
	values.Add("token", c.Token)

	s := ""
	if method == http.MethodPost {
		s = values.Encode()
	}
	body := strings.NewReader(s)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if method == http.MethodGet {
		req.URL.RawQuery = values.Encode()
	}

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	req = req.WithContext(ctx)
	return req, nil
}

func (c *Client) NewSyncRequest(ctx context.Context, values url.Values) (*http.Request, error) {
	return c.NewRequest(ctx, http.MethodPost, "sync", values)
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func (c *Client) Sync(ctx context.Context, commands []Command) (*SyncResponse, error) {
	b, err := json.Marshal(commands)
	if err != nil {
		return nil, err
	}
	values := url.Values{
		"sync_token":           {c.SyncToken},
		"day_orders_timestamp": {""},
		"resource_types":       {"[\"all\"]"},
		"commands":             {string(b)},
	}
	req, err := c.NewSyncRequest(ctx, values)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	var out SyncResponse
	err = decodeBody(res, &out)
	if err != nil {
		return nil, err
	}
	// TODO: replace temp_id mapping
	// TODO: update state
	// TODO: write cache
	return &out, nil
}
