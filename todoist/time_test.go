package todoist

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var testTimes = []struct {
	s string
	v Time
	e error
}{
	{
		s: "2014-09-26T08:25:05Z",
		v: Time{time.Date(2014, 9, 26, 8, 25, 5, 0, time.UTC)},
		e: nil,
	},
}

func TestParse(t *testing.T) {
	for i, tt := range testTimes {
		tim, err := Parse(tt.s)
		if !reflect.DeepEqual(err, tt.e) {
			t.Errorf("%d. %q error mismatch:\n exp=%s\n got=%s\n\n", i, tt.s, tt.e, err)
		} else if tt.e == nil && !tim.Equal(tt.v) {
			t.Errorf("%d. %q mismatch:\n exp=%#v\n got=%#v\n\n", i, tt.s, tt.v, tim)
		}
	}
}

func TestTime_MarshalJSON(t *testing.T) {
	for _, tt := range testTimes {
		b, err := tt.v.MarshalJSON()
		if err != nil || string(b) != strconv.Quote(tt.s) {
			t.Errorf("Expect %s, but got %s", strconv.Quote(tt.s), string(b))
		}
		b, err = Time{}.MarshalJSON()
		if err != nil || string(b) != "null" {
			t.Errorf("Expect %s, but got %s", strconv.Quote(tt.s), string(b))
		}
	}
}

func TestTime_UnmarshalJSON(t *testing.T) {
	for _, test := range testTimes {
		var v Time
		err := v.UnmarshalJSON([]byte(strconv.Quote(test.s)))
		if !reflect.DeepEqual(err, test.e) {
			t.Errorf("Expect %s, but got %s", test.e, err)
		} else if test.e == nil && !v.Equal(test.v) {
			t.Errorf("Expect %s, but got %s", test.v, v)
		}
	}
	var v Time
	err := v.UnmarshalJSON([]byte("null"))
	if err != nil {
		t.Errorf("Unexpect error: %s", err)
	}
	if !v.Equal(Time{}) {
		t.Errorf("Expect %s, but got %s", Time{}, v)
	}
}

func TestTimeJson(t *testing.T) {
	for _, tt := range testTimes {
		m, err := json.Marshal(tt.v)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		if !reflect.DeepEqual(string(m), strconv.Quote(tt.s)) {
			t.Errorf("mismatch:\n exp=%#v\n got=%#v\n\n", strconv.Quote(tt.s), string(m))
		}

		var um Time
		if err := json.Unmarshal([]byte(strconv.Quote(tt.s)), &um); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		if !um.Equal(tt.v) {
			t.Errorf("mismatch:\n exp=%#v\n got=%#v\n\n", tt.v, um)
		}
	}
}
