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
	if r := []rune(substr); len(r) > 0 && string(r[0]) == "#" {
		substr = string(r[1:])
	}
	var res []Project
	for _, p := range c.GetAll() {
		if strings.Contains(p.Name, substr) {
			res = append(res, p)
		}
	}
	return res
}

func (c ProjectClient) FindOneByName(substr string) *Project {
	projects := c.FindByName(substr)
	for _, project := range projects {
		if project.Name == substr {
			return &project
		}
	}
	if len(projects) > 0 {
		return &projects[0]
	}
	return nil
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
	var res []Project
	isNew := true
	for _, p := range *c.cache {
		if p.Equal(project) {
			if !project.IsDeleted {
				res = append(res, project)
			}
			isNew = false
		} else {
			res = append(res, p)
		}
	}
	if isNew && !project.IsDeleted.Bool() {
		res = append(res, project)
	}
	c.cache = &res
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
