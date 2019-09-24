package todoist

import (
	"context"
	"errors"
	"github.com/fatih/color"
	"net/http"
	"net/url"
	"strings"
)

type Label struct {
	Entity
	Name       string  `json:"name"`
	Color      int     `json:"color"`
	ItemOrder  int     `json:"item_order"`
	IsFavorite IntBool `json:"is_favorite"`
}

func (l Label) String() string {
	return "@" + l.Name
}

func (l Label) ColorString() string {
	var attr color.Attribute
	switch l.Color {
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
	return color.New(attr).Sprint(l.String())
}

type NewLabelOpts struct {
	Color      int
	ItemOrder  int
	IsFavorite IntBool
}

func NewLabel(name string, opts *NewLabelOpts) (*Label, error) {
	if len(name) == 0 {
		return nil, errors.New("new label requires a name")
	}
	label := Label{
		Name:       name,
		ItemOrder:  opts.ItemOrder,
		IsFavorite: opts.IsFavorite,
	}
	label.ID = GenerateTempID()
	if opts.Color == 0 {
		label.Color = 47
	} else {
		label.Color = opts.Color
	}
	return &label, nil
}

type Labels []Label

func (ls Labels) String() string {
	var arr []string
	for _, l := range ls {
		arr = append(arr, l.String())
	}
	return strings.Join(arr, " ")
}

func (ls Labels) ColorString() string {
	var arr []string
	for _, l := range ls {
		arr = append(arr, l.ColorString())
	}
	return strings.Join(arr, " ")
}

type LabelClient struct {
	*Client
	cache *labelCache
}

func (c *LabelClient) Add(label Label) (*Label, error) {
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

func (c *LabelClient) UpdateOrders(labels []Label) error {
	args := map[ID]int{}
	for _, label := range labels {
		args[label.ID] = label.ItemOrder
	}
	command := Command{
		Type: "label_update_orders",
		UUID: GenerateUUID(),
		Args: map[string]map[ID]int{
			"id_order_mapping": args,
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

func (c LabelClient) FindByName(substr string) []Label {
	if r := []rune(substr); len(r) > 0 && string(r[0]) == "@" {
		substr = string(r[1:])
	}
	var res []Label
	for _, l := range c.GetAll() {
		if strings.Contains(l.Name, substr) {
			res = append(res, l)
		}
	}
	return res
}

func (c LabelClient) FindOneByName(substr string) *Label {
	labels := c.FindByName(substr)
	for _, label := range labels {
		if label.Name == substr {
			return &label
		}
	}
	if len(labels) > 0 {
		return &labels[0]
	}
	return nil
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
