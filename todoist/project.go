package todoist

type Project struct {
	Entity
	Name         string `json:"name"`
	Color        int `json:"color"`
	Indent       int `json:"indent"`
	ItemOrder    int `json:"item_order"`
	Collapsed    int `json:"collapsed"`
	Shared       bool `json:"shared"`
	IsArchived   int `json:"is_archived"`
	InboxProject bool `json:"inbox_project"`
	TeamInbox    bool `json:"team_inbox"`
}
