package todoist

type Label struct {
	Entity
	Name      string `json:"name"`
	Color     int `json:"color"`
	ItemOrder int `json:"item_order"`
}
