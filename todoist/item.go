package todoist

type Item struct {
	ID             int `json:"id"`
	UserID         int `json:"user_id"`
	ProjectID      int `json:"project_id"`
	Content        string `json:"content"`
	DateString     string `json:"date_string"`
	DateLang       string `json:"date_lang"`
	DueDateUtc     string `json:"due_date_utc"`
	Priority       int `json:"priority"`
	Indent         int `json:"indent"`
	ItemOrder      int `json:"item_order"`
	DayOrder       int `json:"day_order"`
	Collapsed      int `json:"collapsed"`
	Labels         []int `json:"labels"`
	AssignedByUID  int `json:"assigned_by_uid"`
	ResponsibleUID int `json:"responsible_uid"`
	Checked        int `json:"checked"`
	InHistory      int `json:"in_history"`
	IsDeleted      int `json:"is_deleted"`
	IsArchived     int `json:"is_archived"`
	SyncID         int `json:"sync_id"`
	DateAdded      string `json:"date_added"`
}
