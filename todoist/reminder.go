package todoist

type Reminder struct {
	Entity
	NotifyUID  ID     `json:"notify_uid"`
	ItemID     ID     `json:"item_id"`
	Service    string `json:"service"`
	Type       string `json:"type"`
	Due       struct {
		Date        Time   `json:"date"`
		Timezone    string `json:"timezone"`
		IsRecurring bool   `json:"is_recurring"`
		String      string `json:"string"`
		Lang        string `json:"lang"`
	} `json:"due"`
	MmOffset   int    `json:"mm_offset"`
	Name       string `json:"name"`
	LocLat     string `json:"loc_lat"`
	LocLong    string `json:"loc_long"`
	LocTrigger string `json:"loc_trigger"`
	Radius     int    `json:"radius"`
}
