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
		s: "Fri 26 Sep 2014 08:25:05 +0000",
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
	test := testTimes[0]
	b, err := test.v.MarshalJSON()
	if err != nil || string(b) != strconv.Quote(test.s) {
		t.Errorf("Expect %s, but got %s", strconv.Quote(test.s), string(b))
	}

	b, err = Time{}.MarshalJSON()
	if err != nil || string(b) != "null" {
		t.Errorf("Expect %s, but got %s", strconv.Quote(test.s), string(b))
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
	j := []byte(`"Fri 26 Sep 2014 08:25:05 +0000"`)
	v := Time{time.Date(2014, 9, 26, 8, 25, 5, 0, time.UTC)}

	m, err := json.Marshal(v)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !reflect.DeepEqual(m, j) {
		t.Errorf("mismatch:\n exp=%#v\n got=%#v\n\n", j, m)
	}

	var um Time
	if err = json.Unmarshal(j, &um); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !um.Equal(v) {
		t.Errorf("mismatch:\n exp=%#v\n got=%#v\n\n", v, um)
	}
}
