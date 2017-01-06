package todoist

type Note struct {
	ID        int `json:"id"`
	PostedUID int `json:"posted_uid"`
	ItemID    int `json:"item_id"`
	ProjectID int `json:"project_id"`
	Content   string `json:"content"`
	FileAttachment struct {
		FileName    string `json:"file_name"`
		FileSize    int `json:"file_size"`
		FileType    string `json:"file_type"`
		FileURL     string `json:"file_url"`
		UploadState string `json:"upload_state"`
	} `json:"file_attachment"`
	UIDsToNotify []int `json:"uids_to_notify"`
	IsDeleted    int `json:"is_deleted"`
	IsArchived   int `json:"is_archived"`
	Posted       string `json:"posted"`
}
