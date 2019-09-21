package todoist

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Item struct {
	Entity
	UserID    ID     `json:"user_id,omitempty"`
	ProjectID ID     `json:"project_id,omitempty"`
	Content   string `json:"content"`
	Due       struct {
		Date        Time   `json:"date"`
		Timezone    string `json:"timezone"`
		IsRecurring bool   `json:"is_recurring"`
		String      string `json:"string"`
		Lang        string `json:"lang"`
	} `json:"due,omitempty"`
	Priority       int  `json:"priority,omitempty"`
	ParentID       ID   `json:"parent_id,omitempty"`
	ChildOrder     int  `json:"child_order,omitempty"`
	DayOrder       int  `json:"day_order,omitempty"`
	Collapsed      int  `json:"collapsed,omitempty"`
	Labels         []ID `json:"labels,omitempty"`
	AssignedByUID  ID   `json:"assigned_by_uid,omitempty"`
	ResponsibleUID ID   `json:"responsible_uid,omitempty"`
	Checked        int  `json:"checked,omitempty"`
	InHistory      int  `json:"in_history,omitempty"`
	SyncID         int  `json:"sync_id,omitempty"`
	DateAdded      Time `json:"date_added,omitempty"`
	CompletedDate  Time `json:"completed_date"`
}

func (i Item) IsOverDueDate() bool {
	return i.Due.Date.Before(Time{time.Now().UTC()})
}

func (i Item) IsChecked() bool {
	if i.Checked == 1 {
		return true
	}
	return false
}

type ItemClient struct {
	*Client
	cache *itemCache
}

func (c *ItemClient) Add(item Item) (*Item, error) {
	if len(item.Content) == 0 {
		return nil, errors.New("New item requires a content")
	}
	item.ID = GenerateTempID()
	// append item to sync state only `add` method?
	c.cache.store(item)
	command := Command{
		Type:   "item_add",
		Args:   item,
		UUID:   GenerateUUID(),
		TempID: item.ID,
	}
	c.queue = append(c.queue, command)
	return &item, nil
}

func (c *ItemClient) Update(item Item) (*Item, error) {
	if !IsValidID(item.ID) {
		return nil, fmt.Errorf("Invalid id: %s", item.ID)
	}
	command := Command{
		Type: "item_update",
		Args: item,
		UUID: GenerateUUID(),
	}
	c.queue = append(c.queue, command)
	return &item, nil
}

func (c *ItemClient) Delete(id ID) error {
	command := Command{
		Type: "item_delete",
		UUID: GenerateUUID(),
		Args: map[string]ID{
			"id": id,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

type ItemMoveOpts struct {
	ParentID  ID
	ProjectID ID
}

func (c *ItemClient) Move(id ID, opts *ItemMoveOpts) error {
	switch len(opts.ParentID) + len(opts.ProjectID) {
	case 0:
		return errors.New("require parent item id or project id")
	case 2:
		return errors.New("require either parent item id or project id")
	}
	args := map[string]interface{}{
		"id": id,
	}
	if len(opts.ParentID) != 0 {
		args["parent_id"] = opts.ParentID
	}
	if len(opts.ProjectID) != 0 {
		args["project_id"] = opts.ProjectID
	}

	command := Command{
		Type: "item_move",
		UUID: GenerateUUID(),
		Args: args,
	}
	c.queue = append(c.queue, command)
	return nil
}

func (c *ItemClient) Complete(id ID, dateCompleted Time, forceHistory bool) error {
	var fh int
	if forceHistory {
		fh = 1
	} else {
		fh = 0
	}
	command := Command{
		Type: "item_complete",
		UUID: GenerateUUID(),
		Args: map[string]interface{}{
			"id":             id,
			"date_completed": dateCompleted,
			"force_history":  fh,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

func (c *ItemClient) Uncomplete(id ID) error {
	command := Command{
		Type: "item_uncomplete",
		UUID: GenerateUUID(),
		Args: map[string]interface{}{
			"id": id,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

func (c *ItemClient) Close(id ID) error {
	command := Command{
		Type: "item_close",
		UUID: GenerateUUID(),
		Args: map[string]ID{
			"id": id,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

type ItemGetResponse struct {
	Item    Item
	Project Project
	Notes   []Note
}

func (c *ItemClient) Get(ctx context.Context, id ID) (*ItemGetResponse, error) {
	values := url.Values{"item_id": {id.String()}}
	req, err := c.newRequest(ctx, http.MethodGet, "items/get", values)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	var out ItemGetResponse
	err = decodeBody(res, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ItemClient) GetCompleted(ctx context.Context, projectID ID) (*[]Item, error) {
	values := url.Values{"project_id": {projectID.String()}}
	req, err := c.newRequest(ctx, http.MethodGet, "items/get_completed", values)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	var out []Item
	err = decodeBody(res, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ItemClient) GetAll() []Item {
	return c.cache.getAll()
}

func (c *ItemClient) Resolve(id ID) *Item {
	return c.cache.resolve(id)
}

func (c ItemClient) FindByProjectIDs(ids []ID) []Item {
	var res []Item
	for _, i := range c.GetAll() {
		for _, pid := range ids {
			if i.ProjectID == pid {
				res = append(res, i)
				break
			}
		}
	}
	return res
}

func (c ItemClient) FindByContent(substr string) []Item {
	var res []Item
	for _, i := range c.GetAll() {
		if strings.Contains(i.Content, substr) {
			res = append(res, i)
		}
	}
	return res
}

func (c ItemClient) FindByDueDate(time Time) []Item {
	var res []Item
	for _, i := range c.GetAll() {
		if !i.Due.Date.IsZero() && i.Due.Date.Before(time) {
			res = append(res, i)
		}
	}
	return res
}

type itemCache struct {
	cache *[]Item
}

func (c *itemCache) getAll() []Item {
	return *c.cache
}

func (c *itemCache) resolve(id ID) *Item {
	for _, item := range *c.cache {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

func (c *itemCache) store(item Item) {
	// sync api do not returns deleted items.
	// so remove deleted items from cache too.
	var res []Item
	isNew := true
	for _, i := range *c.cache {
		if i.Equal(item) {
			if !item.IsDeleted {
				res = append(res, item)
			}
			isNew = false
		} else {
			res = append(res, i)
		}
	}
	if isNew && !item.IsDeleted.Bool() {
		res = append(res, item)
	}
	c.cache = &res
}

func (c *itemCache) remove(item Item) {
	var res []Item
	for _, i := range *c.cache {
		if !i.Equal(item) {
			res = append(res, i)
		}
	}
	c.cache = &res
}
