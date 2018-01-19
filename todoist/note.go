package todoist

type Note struct {
	Entity
	PostedUID ID     `json:"posted_uid"`
	ItemID    ID     `json:"item_id"`
	ProjectID ID     `json:"project_id"`
	Content   string `json:"content"`
	FileAttachment struct {
		FileName    string `json:"file_name"`
		FileSize    int    `json:"file_size"`
		FileType    string `json:"file_type"`
		FileURL     string `json:"file_url"`
		UploadState string `json:"upload_state"`
	} `json:"file_attachment"`
	UIDsToNotify []ID `json:"uids_to_notify"`
	IsArchived   int  `json:"is_archived"`
	Posted       Time `json:"posted"`
}

// NoteClient encapsulate client operations for notes.
type NoteClient struct {
	cache *noteCache
}

// GetAllForItem returns all the cached notes that belog to the given item.
func (c NoteClient) GetAllForItem(itemID ID) []Note {
	var res []Note
	for _, n := range c.cache.getAll() {
		if n.ItemID == itemID {
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
