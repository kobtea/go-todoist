package todoist

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"net/url"
	"strings"
)

type Project struct {
	Entity
	Name         string `json:"name"`
	Color        int    `json:"color"`
	Indent       int    `json:"indent"`
	ItemOrder    int    `json:"item_order"`
	Collapsed    int    `json:"collapsed"`
	Shared       bool   `json:"shared"`
	IsArchived   int    `json:"is_archived"`
	InboxProject bool   `json:"inbox_project"`
	TeamInbox    bool   `json:"team_inbox"`
}

func (p Project) String() string {
	return strings.Repeat(" ", p.Indent-1) + "#" + p.Name
}

func (p Project) ColorString() string {
	var attr color.Attribute
	switch p.Color {
	case 20, 21:
		attr = color.FgHiBlack
	case 1, 8, 14:
		attr = color.FgHiRed
	case 0, 15, 16:
		attr = color.FgHiGreen
	case 2, 3, 9:
		attr = color.FgHiYellow
	case 17, 18, 19:
		attr = color.FgHiBlue
	case 6, 12, 13:
		attr = color.FgHiMagenta
	case 4, 10, 11:
		attr = color.FgHiCyan
	case 5, 7:
	default:
		attr = color.FgWhite
	}
	return color.New(attr).Sprint(p.String())
}

type ProjectClient struct {
	*Client
	cache *projectCache
}

func (c *ProjectClient) Add(project Project) (*Project, error) {
	if len(project.Name) == 0 {
		return nil, errors.New("New project requires a name")
	}
	project.ID = GenerateTempID()
	c.cache.store(project)
	command := Command{
		Type:   "project_add",
		Args:   project,
		UUID:   GenerateUUID(),
		TempID: project.ID,
	}
	c.queue = append(c.queue, command)
	return &project, nil
}

func (c *ProjectClient) Update(project Project) (*Project, error) {
	if !IsValidID(project.ID) {
		return nil, fmt.Errorf("Invalid id: %s", project.ID)
	}
	command := Command{
		Type: "project_update",
		Args: project,
		UUID: GenerateUUID(),
	}
	c.queue = append(c.queue, command)
	return &project, nil
}

func (c *ProjectClient) Delete(ids []ID) error {
	command := Command{
		Type: "project_delete",
		UUID: GenerateUUID(),
		Args: map[string][]ID{
			"ids": ids,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

func (c *ProjectClient) Archive(ids []ID) error {
	command := Command{
		Type: "project_archive",
		UUID: GenerateUUID(),
		Args: map[string][]ID{
			"ids": ids,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

func (c *ProjectClient) Unarchive(ids []ID) error {
	command := Command{
		Type: "project_unarchive",
		UUID: GenerateUUID(),
		Args: map[string][]ID{
			"ids": ids,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

type ProjectGetResponse struct {
	Project Project
	Notes   []Note
}

func (c *ProjectClient) Get(ctx context.Context, id ID) (*ProjectGetResponse, error) {
	values := url.Values{"project_id": {id.String()}}
	req, err := c.newRequest(ctx, http.MethodGet, "projects/get", values)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	var out ProjectGetResponse
	err = decodeBody(res, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

type ProjectGetDataResponse struct {
	Project Project
	Items   []Item
}

func (c *ProjectClient) GetData(ctx context.Context, id ID) (*ProjectGetDataResponse, error) {
	values := url.Values{"project_id": {id.String()}}
	req, err := c.newRequest(ctx, http.MethodGet, "projects/get_data", values)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	var out ProjectGetDataResponse
	err = decodeBody(res, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProjectClient) GetArchived(ctx context.Context) (*[]Project, error) {
	values := url.Values{}
	req, err := c.newRequest(ctx, http.MethodGet, "projects/get_archived", values)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	var out []Project
	err = decodeBody(res, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProjectClient) GetAll() []Project {
	return c.cache.getAll()
}

func (c *ProjectClient) Resolve(id ID) *Project {
	return c.cache.resolve(id)
}

func (c ProjectClient) FindByName(substr string) []Project {
	var res []Project
	for _, p := range c.GetAll() {
		if strings.Contains(p.Name, substr) {
			res = append(res, p)
		}
	}
	return res
}

type projectCache struct {
	cache *[]Project
}

func (c *projectCache) getAll() []Project {
	return *c.cache
}

func (c *projectCache) resolve(id ID) *Project {
	for _, project := range *c.cache {
		if project.ID == id {
			return &project
		}
	}
	return nil
}

func (c *projectCache) store(project Project) {
	old := c.resolve(project.ID)
	if old == nil {
		if !project.IsDeleted {
			*c.cache = append(*c.cache, project)
		}
	} else {
		if project.IsDeleted {
			c.remove(project)
		} else {
			old = &project
		}
	}
}

func (c *projectCache) remove(project Project) {
	var res []Project
	for _, p := range *c.cache {
		if !p.Equal(project) {
			res = append(res, p)
		}
	}
	c.cache = &res
}
