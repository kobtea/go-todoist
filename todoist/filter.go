package todoist

import (
	"context"
	"errors"
	"github.com/fatih/color"
	"net/http"
	"net/url"
	"strings"
)

type Filter struct {
	Entity
	Name       string  `json:"name"`
	Query      string  `json:"query"`
	Color      int     `json:"color"`
	ItemOrder  int     `json:"item_order"`
	IsFavorite IntBool `json:"is_favorite"`
}

func (f Filter) String() string {
	return f.Name
}

func (f Filter) ColorString() string {
	var attr color.Attribute
	switch f.Color {
	case 30, 31:
		attr = color.FgHiRed
	case 32, 33:
		attr = color.FgHiYellow
	case 34, 35, 36:
		attr = color.FgHiGreen
	case 37, 38, 39:
		attr = color.FgHiCyan
	case 40, 41, 42:
		attr = color.FgHiBlue
	case 43, 44, 45, 46:
		attr = color.FgHiMagenta
	case 47, 48, 49:
		attr = color.FgHiBlack
	default:
		attr = color.FgWhite
	}
	return color.New(attr).Sprint(f.String())
}

type NewFilterOpts struct {
	Color      int
	ItemOrder  int
	IsFavorite IntBool
}

func NewFilter(name, query string, opts *NewFilterOpts) (*Filter, error) {
	if len(name) == 0 || len(query) == 0 {
		return nil, errors.New("new filter requires a name and a query")
	}
	filter := Filter{
		Name:       name,
		Query:      query,
		ItemOrder:  opts.ItemOrder,
		IsFavorite: opts.IsFavorite,
	}
	filter.ID = GenerateTempID()
	if opts.Color == 0 {
		filter.Color = 47
	} else {
		filter.Color = opts.Color
	}
	return &filter, nil
}

type FilterClient struct {
	*Client
	cache *filterCache
}

func (c *FilterClient) Add(filter Filter) (*Filter, error) {
	c.cache.store(filter)
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

func (c *FilterClient) UpdateOrders(filters []Filter) error {
	args := map[ID]int{}
	for _, filter := range filters {
		args[filter.ID] = filter.ItemOrder
	}
	command := Command{
		Type: "filter_update_orders",
		UUID: GenerateUUID(),
		Args: map[string]map[ID]int{
			"id_order_mapping": args,
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

func (c *FilterClient) GetAll() []Filter {
	return c.cache.getAll()
}

func (c *FilterClient) Resolve(id ID) *Filter {
	return c.cache.resolve(id)
}

func (c FilterClient) FindByName(substr string) []Filter {
	var res []Filter
	for _, f := range c.GetAll() {
		if strings.Contains(f.Name, substr) {
			res = append(res, f)
		}
	}
	return res
}

type filterCache struct {
	cache *[]Filter
}

func (c *filterCache) getAll() []Filter {
	return *c.cache
}

func (c *filterCache) resolve(id ID) *Filter {
	for _, filter := range *c.cache {
		if filter.ID == id {
			return &filter
		}
	}
	return nil
}

func (c *filterCache) store(filter Filter) {
	var res []Filter
	isNew := true
	for _, f := range *c.cache {
		if f.Equal(filter) {
			if !filter.IsDeleted {
				res = append(res, filter)
			}
			isNew = false
		} else {
			res = append(res, f)
		}
	}
	if isNew && !filter.IsDeleted.Bool() {
		res = append(res, filter)
	}
	c.cache = &res
}

func (c *filterCache) remove(filter Filter) {
	var res []Filter
	for _, f := range *c.cache {
		if !f.Equal(filter) {
			res = append(res, f)
		}
	}
	c.cache = &res
}
