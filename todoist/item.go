package todoist

import (
	"context"
	"net/url"
	"net/http"
	"errors"
	"fmt"
)

type Item struct {
	Entity
	UserID         int `json:"user_id,omitempty"`
	ProjectID      ID `json:"project_id,omitempty"`
	Content        string `json:"content"`
	DateString     string `json:"date_string,omitempty"`
	DateLang       string `json:"date_lang,omitempty"`
	DueDateUtc     string `json:"due_date_utc,omitempty"`
	Priority       int `json:"priority,omitempty"`
	Indent         int `json:"indent,omitempty"`
	ItemOrder      int `json:"item_order,omitempty"`
	DayOrder       int `json:"day_order,omitempty"`
	Collapsed      int `json:"collapsed,omitempty"`
	Labels         []int `json:"labels,omitempty"`
	AssignedByUID  int `json:"assigned_by_uid,omitempty"`
	ResponsibleUID int `json:"responsible_uid,omitempty"`
	Checked        int `json:"checked,omitempty"`
	InHistory      int `json:"in_history,omitempty"`
	IsArchived     int `json:"is_archived,omitempty"`
	SyncID         int `json:"sync_id,omitempty"`
	DateAdded      string `json:"date_added,omitempty"`
}

type ItemResponse struct {
	Item    Item
	Project Project
	Notes   []Note
}

type ItemManager struct {
	*Client
}

func (m *ItemManager) Add(item Item) (*Item, error) {
	if len(item.Content) == 0 {
		return nil, errors.New("New item requires a content")
	}
	item.ID = GenerateTempID()
	// append item to sync state only `add` method?
	m.SyncState.Items = append(m.SyncState.Items, item)
	command := Command{
		Type:   "item_add",
		Args:   item,
		UUID:   GenerateUUID(),
		TempID: item.ID,
	}
	m.queue = append(m.queue, command)
	return &item, nil
}

func (m *ItemManager) Update(item Item) (*Item, error) {
	if !IsValidID(item.ID) {
		return nil, fmt.Errorf("Invalid id: %s", item.ID)
	}
	command := Command{
		Type: "item_update",
		Args: item,
		UUID: GenerateUUID(),
	}
	m.queue = append(m.queue, command)
	return &item, nil
}

func (m *ItemManager) Delete(ids []ID) (error) {
	command := Command{
		Type: "item_delete",
		UUID: GenerateUUID(),
		Args: map[string][]ID{
			"ids": ids,
		},
	}
	m.queue = append(m.queue, command)
	return nil
}

func (m *ItemManager) Move(projectItems map[ID][]ID, toProject ID) error {
	command := Command{
		Type: "item_move",
		UUID: GenerateUUID(),
		Args: map[string]interface{}{
			"project_items": projectItems,
			"to_project":    toProject,
		},
	}
	m.queue = append(m.queue, command)
	return nil
}

func (m *ItemManager) Complete(ids []ID, forceHistory bool) error {
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
	m.queue = append(m.queue, command)
	return nil
}

func (m *ItemManager) Uncomplete(ids []ID, updateItemOrders bool, restoreState map[ID][]string) error {
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
	m.queue = append(m.queue, command)
	return nil
}

func (m *ItemManager) Close(id ID) error {
	command := Command{
		Type: "item_close",
		UUID: GenerateUUID(),
		Args: map[string]ID{
			"id": id,
		},
	}
	m.queue = append(m.queue, command)
	return nil
}

func (m *ItemManager) Get(ctx context.Context, id ID) (*ItemResponse, error) {
	req, err := m.NewRequest(ctx, http.MethodGet, "items/get", url.Values{"item_id": {id}})
	if err != nil {
		return nil, err
	}
	res, err := m.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	var out ItemResponse
	err = decodeBody(res, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (m *ItemManager) GetCompleted(ctx context.Context, projectID ID) (*[]Item, error) {
	req, err := m.NewRequest(ctx, http.MethodGet, "items/get_completed", url.Values{"project_id": {projectID}})
	if err != nil {
		return nil, err
	}
	res, err := m.HTTPClient.Do(req)
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

func (m *ItemManager) Resolve(id ID) *Item {
	for _, item := range m.SyncState.Items {
		if item.ID == id {
			return &item
		}
	}
	return nil
}
