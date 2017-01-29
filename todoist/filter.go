package todoist

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Filter struct {
	Entity
	Name      string `json:"name"`
	Query     string `json:"query"`
	Color     int    `json:"color"`
	ItemOrder int    `json:"item_order"`
}

type FilterClient struct {
	*Client
}

func (c *FilterClient) Add(filter Filter) (*Filter, error) {
	if len(filter.Name) == 0 {
		return nil, errors.New("New filter requires a name")
	}
	if len(filter.Query) == 0 {
		return nil, errors.New("New filter requires a query")
	}
	filter.ID = GenerateTempID()
	c.syncState.Filters = append(c.syncState.Filters, filter)
	command := Command{
		Type:   "filter_add",
		Args:   filter,
		UUID:   GenerateUUID(),
		TempID: filter.ID,
	}
	c.queue = append(c.queue, command)
	return &filter, nil
}

func (c *FilterClient) Update(filter Filter) (*Filter, error) {
	if !IsValidID(filter.ID) {
		return nil, fmt.Errorf("Invalid id: %s", filter.ID)
	}
	command := Command{
		Type: "filter_update",
		Args: filter,
		UUID: GenerateUUID(),
	}
	c.queue = append(c.queue, command)
	return &filter, nil
}

func (c *FilterClient) Delete(id ID) error {
	command := Command{
		Type: "filter_delete",
		UUID: GenerateUUID(),
		Args: map[string]ID{
			"id": id,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

type FilterGetResponse struct {
	Filter Filter
}

func (c *FilterClient) Get(ctx context.Context, id ID) (*FilterGetResponse, error) {
	values := url.Values{"filter_id": {id.String()}}
	req, err := c.newRequest(ctx, http.MethodGet, "filters/get", values)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	var out FilterGetResponse
	err = decodeBody(res, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *FilterClient) Resolve(id ID) *Filter {
	for _, filter := range c.syncState.Filters {
		if filter.ID == id {
			return &filter
		}
	}
	return nil
}
