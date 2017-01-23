package todoist

import (
	"errors"
	"fmt"
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

type ProjectManager struct {
	*Client
}

func (m *ProjectManager) Add(project Project) (*Project, error) {
	if len(project.Name) == 0 {
		return nil, errors.New("New project requires a name")
	}
	project.ID = GenerateTempID()
	m.SyncState.Projects = append(m.SyncState.Projects, project)
	command := Command{
		Type:   "project_add",
		Args:   project,
		UUID:   GenerateUUID(),
		TempID: project.ID,
	}
	m.queue = append(m.queue, command)
	return &project, nil
}

func (m *ProjectManager) Update(project Project) (*Project, error) {
	if !IsValidID(project.ID) {
		return nil, fmt.Errorf("Invalid id: %s", project.ID)
	}
	command := Command{
		Type: "project_update",
		Args: project,
		UUID: GenerateUUID(),
	}
	m.queue = append(m.queue, command)
	return &project, nil
}

func (m *ProjectManager) Delete(ids []ID) error {
	command := Command{
		Type: "project_delete",
		UUID: GenerateUUID(),
		Args: map[string][]ID{
			"ids": ids,
		},
	}
	m.queue = append(m.queue, command)
	return nil
}

func (m *ProjectManager) Archive(ids []ID) error {
	command := Command{
		Type: "project_archive",
		UUID: GenerateUUID(),
		Args: map[string][]ID{
			"ids": ids,
		},
	}
	m.queue = append(m.queue, command)
	return nil
}

func (m *ProjectManager) Unarchive(ids []ID) error {
	command := Command{
		Type: "project_unarchive",
		UUID: GenerateUUID(),
		Args: map[string][]ID{
			"ids": ids,
		},
	}
	m.queue = append(m.queue, command)
	return nil
}
