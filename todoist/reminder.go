package todoist

type Reminder struct {
	Entity
	NotifyUID  ID     `json:"notify_uid"`
	ItemID     ID     `json:"item_id"`
	Service    string `json:"service"`
	Type       string `json:"type"`
	DateString string `json:"date_string"`
	DateLang   string `json:"date_lang"`
	DueDateUtc string `json:"due_date_utc"`
	MmOffset   int    `json:"mm_offset"`
	Name       string `json:"name"`
	LocLat     string `json:"loc_lat"`
	LocLong    string `json:"loc_long"`
	LocTrigger string `json:"loc_trigger"`
	Radius     int    `json:"radius"`
}
