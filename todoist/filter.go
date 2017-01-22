package todoist

type Filter struct {
	Entity
	Name      string `json:"name"`
	Query     string `json:"query"`
	Color     int `json:"color"`
	ItemOrder int `json:"item_order"`
}
