package todoist

import "errors"

type Note struct {
	Entity
	PostedUID      ID              `json:"posted_uid"`
	ItemID         ID              `json:"item_id"`
	ProjectID      ID              `json:"project_id"`
	Content        string          `json:"content"`
	FileAttachment FileAttachment  `json:"file_attachment"`
	UIDsToNotify   []ID            `json:"uids_to_notify"`
	Posted         Time            `json:"posted"`
	Reactions      map[string][]ID `json:"reactions"`
}

type FileAttachment struct {
	FileName    string `json:"file_name"`
	FileSize    int    `json:"file_size"`
	FileType    string `json:"file_type"`
	FileURL     string `json:"file_url"`
	UploadState string `json:"upload_state"`
}

type NewNoteOpts struct {
	FileAttachment FileAttachment
	UIDsToNotify   []ID
}

func NewNote(id ID, content string, opts *NewNoteOpts) (*Note, error) {
	if id.IsZero() || len(content) == 0 {
		return nil, errors.New("new note requires an item id and a content")
	}
	note := Note{
		ItemID:         id,
		Content:        content,
		FileAttachment: opts.FileAttachment,
		UIDsToNotify:   opts.UIDsToNotify,
	}
	note.ID = GenerateTempID()
	return &note, nil
}

// NoteClient encapsulate client operations for notes.
type NoteClient struct {
	*Client
	cache *noteCache
}

func (c NoteClient) Add(note Note) (*Note, error) {
	c.cache.store(note)
	command := Command{
		Type:   "note_add",
		Args:   note,
		UUID:   GenerateUUID(),
		TempID: note.ID,
	}
	c.queue = append(c.queue, command)
	return &note, nil
}

func (c NoteClient) Update(note Note) (*Note, error) {
	command := Command{
		Type: "note_update",
		Args: note,
		UUID: GenerateUUID(),
	}
	c.queue = append(c.queue, command)
	return &note, nil
}

func (c NoteClient) Delete(id ID) error {
	command := Command{
		Type: "note_delete",
		UUID: GenerateUUID(),
		Args: map[string]ID{
			"id": id,
		},
	}
	c.queue = append(c.queue, command)
	return nil
}

// GetAllForItem returns all the cached notes that belong to the given item.
func (c NoteClient) GetAllForItem(itemID ID) []Note {
	var res []Note
	for _, n := range c.cache.getAll() {
		if n.ItemID == itemID {
			res = append(res, n)
		}
	}
	return res
}

// GetAllForProject returns all the cached notes that belong to the given project.
func (c NoteClient) GetAllForProject(projectID ID) []Note {
	var res []Note
	for _, n := range c.cache.getAll() {
		if n.ProjectID == projectID && n.ItemID == "" {
			res = append(res, n)
		}
	}
	return res
}

type noteCache struct {
	cache *[]Note
}

func (c *noteCache) getAll() []Note {
	return *c.cache
}

func (c *noteCache) store(note Note) {
	var res []Note
	isNew := true
	for _, n := range *c.cache {
		if n.Equal(note) {
			if !note.IsDeleted {
				res = append(res, note)
			}
			isNew = false
		} else {
			res = append(res, n)
		}
	}
	if isNew && !note.IsDeleted.Bool() {
		res = append(res, note)
	}
	c.cache = &res
}
