package todoist

import (
	"context"
	"encoding/json"
	"errors"
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
	SyncState  *SyncState
	Logger     *log.Logger
	Item       *ItemClient
	Label      *LabelClient
	Project    *ProjectClient
	queue      []Command
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

	c := &Client{
		URL:        parsed_endpoint,
		HTTPClient: client,
		Token:      token,
		SyncToken:  sync_token,
		CacheDir:   cache_dir,
		SyncState:  &SyncState{},
		Logger:     logger,
	}
	c.Item = &ItemClient{c}
	c.Label = &LabelClient{c}
	c.Project = &ProjectClient{c}
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
	c.SyncState = &SyncState{}
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
	for _, item := range state.Items {
		cachedItem := c.Item.Resolve(item.ID)
		if cachedItem == nil {
			if !item.IsDeleted {
				c.SyncState.Items = append(c.SyncState.Items, item)
			}
		} else {
			if item.IsDeleted {
				var res []Item
				for _, i := range c.SyncState.Items {
					if !i.Equal(cachedItem) {
						res = append(res, i)
					}
				}
				c.SyncState.Items = res
			} else {
				cachedItem = &item
			}
		}
	}
	for _, label := range state.Labels {
		cachedLabel := c.Label.Resolve(label.ID)
		if cachedLabel == nil {
			if !label.IsDeleted {
				c.SyncState.Labels = append(c.SyncState.Labels, label)
			}
		} else {
			if label.IsDeleted {
				var res []Label
				for _, l := range c.SyncState.Labels {
					if !l.Equal(cachedLabel) {
						res = append(res, l)
					}
				}
				c.SyncState.Labels = res
			} else {
				cachedLabel = &label
			}
		}
	}
	for _, project := range state.Projects {
		cachedProject := c.Project.Resolve(project.ID)
		if cachedProject == nil {
			if !project.IsDeleted {
				c.SyncState.Projects = append(c.SyncState.Projects, project)
			}
		} else {
			if project.IsDeleted {
				var res []Project
				for _, p := range c.SyncState.Projects {
					if !p.Equal(cachedProject) {
						res = append(res, p)
					}
				}
				c.SyncState.Projects = res
			} else {
				cachedProject = &project
			}
		}
	}
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
