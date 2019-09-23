package todoist

import (
	"github.com/fatih/color"
	"strconv"
	"time"
)

const (
	dateLayout     = "2006-01-02"
	datetimeLayout = time.RFC3339
	localLayout    = "2006-01-02(Mon) 15:04"
)

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
	t, err := time.Parse(datetimeLayout, value)
	if err != nil {
		t, err = time.Parse(dateLayout, value)
	}
	if err != nil {
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

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	layout := datetimeLayout
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		layout = dateLayout
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
	return t.Time.Local().Format(localLayout)
}

func (t Time) ColorString() string {
	if !t.IsZero() && t.Before(Time{time.Now()}) {
		return color.New(color.BgRed).Sprint(t.String())
	}
	return t.String()
}
