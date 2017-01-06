package todoist

type Filter struct {
	ID        int `json:"id"`
	Name      string `json:"name"`
	Query     string `json:"query"`
	Color     int `json:"color"`
	ItemOrder int `json:"item_order"`
	IsDeleted int `json:"is_deleted"`
}
