package todoist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	Token      string
	SyncToken  string
	CacheDir   string
	syncState  *SyncState
	Logger     *log.Logger
	Completed  *CompletedClient
	Filter     *FilterClient
	Item       *ItemClient
	Label      *LabelClient
	Project    *ProjectClient
	Relation   *RelationClient
	Note       *NoteClient
	queue      []Command
}

func NewClient(endpoint, token, sync_token, cache_dir string, logger *log.Logger) (*Client, error) {
	if len(endpoint) == 0 {
		endpoint = "https://api.todoist.com/sync/v8"
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

	c := &Client{
		URL:        parsed_endpoint,
		HTTPClient: client,
		Token:      token,
		SyncToken:  sync_token,
		CacheDir:   cache_dir,
		syncState:  &SyncState{},
		Logger:     logger,
	}
	if err = c.readCache(); err != nil {
		c.resetState()
	}
	c.Completed = &CompletedClient{c}
	c.Filter = &FilterClient{c, &filterCache{&c.syncState.Filters}}
	c.Item = &ItemClient{c, &itemCache{&c.syncState.Items}}
	c.Label = &LabelClient{c, &labelCache{&c.syncState.Labels}}
	c.Project = &ProjectClient{c, &projectCache{&c.syncState.Projects}}
	c.Relation = &RelationClient{c}
	c.Note = &NoteClient{c, &noteCache{&c.syncState.Notes}}
	return c, nil
}

func (c *Client) newRequest(ctx context.Context, method, spath string, values url.Values) (*http.Request, error) {
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

func (c *Client) newSyncRequest(ctx context.Context, values url.Values) (*http.Request, error) {
	return c.newRequest(ctx, http.MethodPost, "sync", values)
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
	req, err := c.newSyncRequest(ctx, values)
	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	if (res.StatusCode / 100) != 2 {
		return fmt.Errorf("failed to sync, status code: %d, command: %v", res.StatusCode, commands)
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

func (c *Client) FullSync(ctx context.Context, commands []Command) error {
	c.resetState()
	return c.Sync(ctx, commands)
}

func (c *Client) Commit(ctx context.Context) error {
	if len(c.queue) == 0 {
		return nil
	}
	err := c.Sync(ctx, c.queue)
	c.queue = []Command{}
	return err
}

func (c *Client) ResetSyncToken() {
	c.SyncToken = "*"
}

func (c *Client) resetState() {
	c.SyncToken = "*"
	c.syncState = &SyncState{}
}

func (c *Client) updateState(state *SyncState) {
	if len(state.SyncToken) != 0 {
		c.SyncToken = state.SyncToken
	}
	/* TODO:
	- day_orders
	- day_orders_timestamp
	- live_notifications_last_read_id
	- locations
	- settings_notifications
	- user
	*/
	for _, filter := range state.Filters {
		c.Filter.cache.store(filter)
	}
	for _, item := range state.Items {
		c.Item.cache.store(item)
	}
	for _, label := range state.Labels {
		c.Label.cache.store(label)
	}
	for _, project := range state.Projects {
		c.Project.cache.store(project)
	}
	for _, note := range state.Notes {
		c.Note.cache.store(note)
	}
	for _, note := range state.ProjectNotes {
		c.Note.cache.store(note)
	}
	c.syncState = state
}

func (c *Client) readCache() error {
	b, err := ioutil.ReadFile(path.Join(c.CacheDir, c.Token+".json"))
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, c.syncState); err != nil {
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
	b, err := json.MarshalIndent(c.syncState, "", "  ")
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
