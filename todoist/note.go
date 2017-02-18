package todoist

type Note struct {
	Entity
	PostedUID      ID     `json:"posted_uid"`
	ItemID         ID     `json:"item_id"`
	ProjectID      ID     `json:"project_id"`
	Content        string `json:"content"`
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
