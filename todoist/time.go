package todoist

import (
	"github.com/fatih/color"
	"strconv"
	"time"
)

const (
	dateLayout          = "2006-01-02"
	datetimeLayout      = "2006-01-02T15:04:05"
	datetimeTzLayout    = time.RFC3339
	localDateLayout     = "2006-01-02(Mon)"
	localDatetimeLayout = "2006-01-02(Mon) 15:04"
)

/*
Due#date has 3 formats.
- Full-day dates:     date=YYYY-MM-DD,           timezone=null
- Floating due dates: date=YYYY-MM-DDTHH:MM:SS,  timezone=null (use current user's timezone)
- Due dates:          date=YYYY-MM-DDTHH:MM:SSZ, timezone=Europe/Madrid
*/

type Due struct {
	Date        Time   `json:"date"`
	Timezone    string `json:"timezone"`
	String      string `json:"string"`
	Lang        string `json:"lang"`
	IsRecurring bool   `json:"is_recurring"`
}

type Time struct {
	time.Time
}

func Today() Time {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 1, time.Local)
	return Time{today.UTC()}
}

func Next7Days() Time {
	d := time.Now().Add(6 * 24 * time.Hour)
	days := time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 1, time.Local)
	return Time{days.UTC()}
}

func Parse(value string) (Time, error) {
	var t time.Time
	var err error
	for _, layout := range []string{dateLayout, datetimeLayout} {
		// FIXME: refer to user.tz_info.timezone
		if t, err = time.ParseInLocation(layout, value, time.Local); err == nil {
			return Time{t}, nil
		}
	}
	if t, err = time.Parse(datetimeTzLayout, value); err != nil {
		return Time{}, err
	}
	return Time{t}, nil
}

func (t Time) Equal(u Time) bool {
	return t.Time.Equal(u.Time)
}

func (t Time) Before(u Time) bool {
	return t.Time.Before(u.Time)
}

func (t Time) After(u Time) bool {
	return t.Time.After(u.Time)
}

func (t Time) Local() Time {
	return Time{t.Time.Local()}
}

func (t Time) IsFullDay() bool {
	return t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	layout := datetimeLayout
	if t.IsFullDay() {
		layout = dateLayout
	}
	if t.Time.Location() == time.UTC {
		layout = datetimeTzLayout
	}
	return []byte(strconv.Quote(t.Time.Format(layout))), nil
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		*t = Time{time.Time{}} // null value
	} else {
		*t, err = Parse(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t Time) String() string {
	if t.IsZero() {
		return ""
	}
	layout := localDatetimeLayout
	if t.IsFullDay() {
		layout = localDateLayout
	}
	return t.Time.Local().Format(layout)
}

func (t Time) ColorString() string {
	if !t.IsZero() && t.Before(Time{time.Now()}) && !t.IsFullDay() {
		return color.New(color.BgRed).Sprint(t.String())
	}
	return t.String()
}
