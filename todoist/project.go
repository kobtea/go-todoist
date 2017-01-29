package todoist

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/goweb/http"
	"net/url"
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

type ProjectClient struct {
	*Client
}

func (c *ProjectClient) Add(project Project) (*Project, error) {
	if len(project.Name) == 0 {
		return nil, errors.New("New project requires a name")
	}
	project.ID = GenerateTempID()
	c.syncState.Projects = append(c.syncState.Projects, project)
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

func (c *ProjectClient) Resolve(id ID) *Project {
	for _, project := range c.syncState.Projects {
		if project.ID == id {
			return &project
		}
	}
	return nil
}
