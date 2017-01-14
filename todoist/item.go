package todoist

import (
	"context"
	"net/url"
	"net/http"
)

type Item struct {
	ID             int `json:"id"`
	UserID         int `json:"user_id"`
	ProjectID      int `json:"project_id"`
	Content        string `json:"content"`
	DateString     string `json:"date_string"`
	DateLang       string `json:"date_lang"`
	DueDateUtc     string `json:"due_date_utc"`
	Priority       int `json:"priority"`
	Indent         int `json:"indent"`
	ItemOrder      int `json:"item_order"`
	DayOrder       int `json:"day_order"`
	Collapsed      int `json:"collapsed"`
	Labels         []int `json:"labels"`
	AssignedByUID  int `json:"assigned_by_uid"`
	ResponsibleUID int `json:"responsible_uid"`
	Checked        int `json:"checked"`
	InHistory      int `json:"in_history"`
	IsDeleted      int `json:"is_deleted"`
	IsArchived     int `json:"is_archived"`
	SyncID         int `json:"sync_id"`
	DateAdded      string `json:"date_added"`
}

type ItemResponse struct {
	Item    Item
	Project Project
	Notes   []Note
}

func (c *Client) GetItem(ctx context.Context, id string) (*ItemResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, "items/get", url.Values{"item_id": {id}})
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
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
