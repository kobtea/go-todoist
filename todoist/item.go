package todoist

import (
	"context"
	"net/url"
	"net/http"
	"errors"
)

type Item struct {
	ID             ID `json:"id,omitempty"`
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
	IsDeleted      int `json:"is_deleted,omitempty"`
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

func (m *ItemManager) Get(ctx context.Context, id string) (*ItemResponse, error) {
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
