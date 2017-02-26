package todoist

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Label struct {
	Entity
	Name      string `json:"name"`
	Color     int    `json:"color"`
	ItemOrder int    `json:"item_order"`
}

func (l Label) String() string {
	return "@" + l.Name
}

type LabelClient struct {
	*Client
	cache *labelCache
}

func (c *LabelClient) Add(label Label) (*Label, error) {
	if len(label.Name) == 0 {
		return nil, errors.New("New label requires a name")
	}
	label.ID = GenerateTempID()
	c.cache.store(label)
	command := Command{
		Type:   "label_add",
		Args:   label,
		UUID:   GenerateUUID(),
		TempID: label.ID,
	}
	c.queue = append(c.queue, command)
	return &label, nil
}

func (c *LabelClient) Update(label Label) (*Label, error) {
	if !IsValidID(label.ID) {
		return nil, fmt.Errorf("Invalid id: %s", label.ID)
	}
	command := Command{
		Type: "label_update",
		Args: label,
		UUID: GenerateUUID(),
	}
	c.queue = append(c.queue, command)
	return &label, nil
}

func (c *LabelClient) Delete(id ID) error {
	command := Command{
		Type: "label_delete",
		UUID: GenerateUUID(),
		Args: map[string]ID{
			"id": id,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

type LabelGetResponse struct {
	Label Label
}

func (c *LabelClient) Get(ctx context.Context, id ID) (*LabelGetResponse, error) {
	values := url.Values{"label_id": {id.String()}}
	req, err := c.newRequest(ctx, http.MethodGet, "labels/get", values)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	var out LabelGetResponse
	err = decodeBody(res, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LabelClient) GetAll() []Label {
	return c.cache.getAll()
}

func (c *LabelClient) Resolve(id ID) *Label {
	return c.cache.resolve(id)
}

type labelCache struct {
	cache *[]Label
}

func (c *labelCache) getAll() []Label {
	return *c.cache
}

func (c *labelCache) resolve(id ID) *Label {
	for _, label := range *c.cache {
		if label.ID == id {
			return &label
		}
	}
	return nil
}

func (c *labelCache) store(label Label) {
	var res []Label
	isNew := true
	for _, l := range *c.cache {
		if l.Equal(label) {
			if !label.IsDeleted {
				res = append(res, label)
			}
			isNew = false
		} else {
			res = append(res, l)
		}
	}
	if isNew && !label.IsDeleted.Bool() {
		res = append(res, label)
	}
	c.cache = &res
}

func (c *labelCache) remove(label Label) {
	var res []Label
	for _, l := range *c.cache {
		if !l.Equal(label) {
			res = append(res, l)
		}
	}
	c.cache = &res
}
