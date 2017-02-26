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
	UserID         ID     `json:"user_id,omitempty"`
	ProjectID      ID     `json:"project_id,omitempty"`
	Content        string `json:"content"`
	DateString     string `json:"date_string,omitempty"`
	DateLang       string `json:"date_lang,omitempty"`
	DueDateUtc     Time   `json:"due_date_utc,omitempty"`
	Priority       int    `json:"priority,omitempty"`
	Indent         int    `json:"indent,omitempty"`
	ItemOrder      int    `json:"item_order,omitempty"`
	DayOrder       int    `json:"day_order,omitempty"`
	Collapsed      int    `json:"collapsed,omitempty"`
	Labels         []ID   `json:"labels,omitempty"`
	AssignedByUID  ID     `json:"assigned_by_uid,omitempty"`
	ResponsibleUID ID     `json:"responsible_uid,omitempty"`
	Checked        int    `json:"checked,omitempty"`
	InHistory      int    `json:"in_history,omitempty"`
	IsArchived     int    `json:"is_archived,omitempty"`
	SyncID         int    `json:"sync_id,omitempty"`
	DateAdded      Time   `json:"date_added,omitempty"`
}

func (i Item) IsOverDueDate() bool {
	return i.DueDateUtc.Before(Time{time.Now().UTC()})
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

func (c *ItemClient) Delete(ids []ID) error {
	command := Command{
		Type: "item_delete",
		UUID: GenerateUUID(),
		Args: map[string][]ID{
			"ids": ids,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

func (c *ItemClient) Move(projectItems map[ID][]ID, toProject ID) error {
	command := Command{
		Type: "item_move",
		UUID: GenerateUUID(),
		Args: map[string]interface{}{
			"project_items": projectItems,
			"to_project":    toProject,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

func (c *ItemClient) Complete(ids []ID, forceHistory bool) error {
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
			"ids":           ids,
			"force_history": fh,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

func (c *ItemClient) Uncomplete(ids []ID, updateItemOrders bool, restoreState map[ID][]string) error {
	var uio int
	if updateItemOrders {
		uio = 1
	} else {
		uio = 0
	}
	command := Command{
		Type: "item_uncomplete",
		UUID: GenerateUUID(),
		Args: map[string]interface{}{
			"ids":                ids,
			"update_item_orders": uio,
			"restore_state":      restoreState,
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
		if !i.DueDateUtc.IsZero() && i.DueDateUtc.Before(time) {
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
