package todoist

import (
	"context"
	"net/url"
)

type Stats struct {
	KarmaLastUpdate float64 `json:"karma_last_update"`
	KarmaTrend      string  `json:"karma_trend"`
	DaysItems       []struct {
		Date  string `json:"date"`
		Items []struct {
			Completed int `json:"completed"`
			ID        int `json:"id"`
		} `json:"items"`
		TotalCompleted int `json:"total_completed"`
	} `json:"days_items"`
	CompletedCount     int `json:"completed_count"`
	KarmaUpdateReasons []struct {
		PositiveKarmaReasons []int   `json:"positive_karma_reasons"`
		NewKarma             float64 `json:"new_karma"`
		NegativeKarma        float64 `json:"negative_karma"`
		PositiveKarma        float64 `json:"positive_karma"`
		NegativeKarmaReasons []int   `json:"negative_karma_reasons"`
		Time                 string  `json:"time"`
	} `json:"karma_update_reasons"`
	Karma     float64 `json:"karma"`
	WeekItems []struct {
		Date  string `json:"date"`
		Items []struct {
			Completed int `json:"completed"`
			ID        int `json:"id"`
		} `json:"items"`
		TotalCompleted int `json:"total_completed"`
	} `json:"week_items"`
	KarmaGraph string `json:"karma_graph"`
	Goals      struct {
		KarmaDisabled   int `json:"karma_disabled"`
		UserID          int `json:"user_id"`
		LastDailyStreak struct {
			Count int    `json:"count"`
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"last_daily_streak"`
		VacationMode    int   `json:"vacation_mode"`
		IgnoreDays      []int `json:"ignore_days"`
		MaxWeeklyStreak struct {
			Count int    `json:"count"`
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"max_weekly_streak"`
		CurrentWeeklyStreak struct {
			Count int    `json:"count"`
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"current_weekly_streak"`
		CurrentDailyStreak struct {
			Count int    `json:"count"`
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"current_daily_streak"`
		LastWeeklyStreak struct {
			Count int    `json:"count"`
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"last_weekly_streak"`
		WeeklyGoal     int `json:"weekly_goal"`
		MaxDailyStreak struct {
			Count int    `json:"count"`
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"max_daily_streak"`
		DailyGoal int `json:"daily_goal"`
	} `json:"goals"`
}

type CompletedItems struct {
	Items    []Item         `json:"items"`
	Projects map[ID]Project `json:"projects"`
}

func (c *CompletedItems) GroupByCompletedDate() map[string][]Item {
	const layout = "2006-01-02"
	res := map[string][]Item{}
	for _, item := range c.Items {
		date := item.CompletedDate.Local().Format(layout)
		res[date] = append(res[date], item)
	}
	return res
}

type CompletedClient struct {
	*Client
}

func (c *CompletedClient) GetStats() (*Stats, error) {
	req, err := c.newRequest(context.Background(), "POST", "completed/get_stats", url.Values{})
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	var out Stats
	decodeBody(res, &out)
	return &out, nil
}

func (c *CompletedClient) GetAll() (*CompletedItems, error) {
	req, err := c.newRequest(context.Background(), "POST", "completed/get_all", url.Values{})
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	var out CompletedItems
	decodeBody(res, &out)
	return &out, nil
}
