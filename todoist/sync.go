package todoist

type SyncState struct {
	SyncToken string `json:"sync_token"`
	FullSync  bool   `json:"full_sync"`
	// User User `json:"user"`
	Projects []Project `json:"projects"`
	// ProjectNotes []interface{} `json:"project_notes"`
	Items   []Item   `json:"items"`
	Notes   []Note   `json:"notes"`
	Labels  []Label  `json:"labels"`
	Filters []Filter `json:"filters"`
	// DayOrders struct {} `json:"day_orders"`
	// DayOrdersTimestamp string `json:"day_orders_timestamp"`
	Reminders []Reminder `json:"reminders"`
	// Collaborators []interface{} `json:"collaborators"`
	// CollaboratorStates []CollaboratorState `json:"collaborator_states"`
	// LiveNotifications []LiveNotification `json:"live_notifications"`
	// LiveNotificationsLastReadID int `json:"live_notifications_last_read_id"`
	// Locations []interface{} `json:"locations"`
	// TempIDMapping struct {} `json:"temp_id_mapping"`
}

type Command struct {
	Type   string      `json:"type"`
	Args   interface{} `json:"args"`
	UUID   UUID        `json:"uuid"`
	TempID ID          `json:"temp_id"`
}
