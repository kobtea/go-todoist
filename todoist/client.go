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
	"os"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	Token      string
	SyncToken  string
	CacheDir   string
	SyncState  *SyncState
	Logger     *log.Logger
}

func NewClient(endpoint, token, sync_token, cache_dir string, logger *log.Logger) (*Client, error) {
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

	if len(cache_dir) == 0 {
		cache_dir = "$HOME/.go-todoist"
	}
	cache_dir = os.ExpandEnv(cache_dir)
	if _, err = os.Stat(cache_dir); err != nil {
		if err = os.MkdirAll(cache_dir, 0755); err != nil {
			return nil, err
		}
	}

	if logger == nil {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	c := &Client{parsed_endpoint, client, token, sync_token, cache_dir, &SyncState{}, logger}
	if err = c.readCache(); err != nil {
		c.resetState()
	}
	return c, nil
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

func (c *Client) Sync(ctx context.Context, commands []Command) error {
	b, err := json.Marshal(commands)
	if err != nil {
		return err
	}
	values := url.Values{
		"sync_token":           {c.SyncToken},
		"day_orders_timestamp": {""},
		"resource_types":       {"[\"all\"]"},
		"commands":             {string(b)},
	}
	req, err := c.NewSyncRequest(ctx, values)
	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	var out SyncState
	err = decodeBody(res, &out)
	if err != nil {
		return err
	}
	// TODO: replace temp_id mapping
	c.updateState(&out)
	c.writeCache()
	return nil
}

func (c *Client) resetState() {
	c.SyncToken = "*"
	c.SyncState = &SyncState{}
}

func (c *Client) updateState(state *SyncState) {
	c.SyncToken = state.SyncToken
	c.SyncState = state
}

func (c *Client) readCache() error {
	b, err := ioutil.ReadFile(path.Join(c.CacheDir, c.Token+".json"))
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, c.SyncState); err != nil {
		return err
	}
	b, err = ioutil.ReadFile(path.Join(c.CacheDir, c.Token+".sync"))
	if err != nil {
		return err
	}
	c.SyncToken = string(b)
	return nil
}

func (c *Client) writeCache() error {
	if len(c.CacheDir) == 0 {
		return nil
	}
	b, err := json.MarshalIndent(c.SyncState, "", "  ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(path.Join(c.CacheDir, c.Token+".json"), b, 0644); err != nil {
		return err
	}
	if err = ioutil.WriteFile(path.Join(c.CacheDir, c.Token+".sync"), []byte(c.SyncToken), 0644); err != nil {
		return err
	}
	return nil
}
