package todoist

type Label struct {
	ID        int `json:"id"`
	Name      string `json:"name"`
	Color     int `json:"color"`
	ItemOrder int `json:"item_order"`
	IsDeleted int `json:"is_deleted"`
}
