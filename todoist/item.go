package todoist

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Item struct {
	Entity
	UserID         ID     `json:"user_id,omitempty"`
	ProjectID      ID     `json:"project_id,omitempty"`
	Content        string `json:"content"`
	DateString     string `json:"date_string,omitempty"`
	DateLang       string `json:"date_lang,omitempty"`
	DueDateUtc     string `json:"due_date_utc,omitempty"`
	Priority       int    `json:"priority,omitempty"`
	Indent         int    `json:"indent,omitempty"`
	ItemOrder      int    `json:"item_order,omitempty"`
	DayOrder       int    `json:"day_order,omitempty"`
	Collapsed      int    `json:"collapsed,omitempty"`
	Labels         []int  `json:"labels,omitempty"`
	AssignedByUID  ID     `json:"assigned_by_uid,omitempty"`
	ResponsibleUID ID     `json:"responsible_uid,omitempty"`
	Checked        int    `json:"checked,omitempty"`
	InHistory      int    `json:"in_history,omitempty"`
	IsArchived     int    `json:"is_archived,omitempty"`
	SyncID         int    `json:"sync_id,omitempty"`
	DateAdded      string `json:"date_added,omitempty"`
}

type ItemClient struct {
	*Client
}

func (c *ItemClient) Add(item Item) (*Item, error) {
	if len(item.Content) == 0 {
		return nil, errors.New("New item requires a content")
	}
	item.ID = GenerateTempID()
	// append item to sync state only `add` method?
	c.syncState.Items = append(c.syncState.Items, item)
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

func (c *ItemClient) Resolve(id ID) *Item {
	for _, item := range c.SyncState.Items {
		if item.ID == id {
			return &item
		}
	}
	return nil
}
