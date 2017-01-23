package todoist

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

var testIDs = []struct {
	s string
	v ID
	e error
}{
	{
		"1000000",
		ID("1000000"),
		nil,
	},
	{
		"df43406d-db7e-4ea5-b3b4-c822ccdab3bf",
		ID("df43406d-db7e-4ea5-b3b4-c822ccdab3bf"),
		nil,
	},
	{
		"invalid",
		"",
		errors.New("Invalid ID: invalid"),
	},
}

func TestNewID(t *testing.T) {
	for _, test := range testIDs {
		v, err := NewID(test.s)
		if !reflect.DeepEqual(err, test.e) {
			t.Errorf("Expect %s, but got %s", test.e, err)
		} else if test.e == nil && v != test.v {
			t.Errorf("Expect %s, but got %s", test.v, v)
		}
	}
}

func TestID_MarshalJSON(t *testing.T) {
	test := testIDs[0]
	b, err := test.v.MarshalJSON()
	if err != nil || string(b) != test.s {
		t.Errorf("Expect %s, but got %s", strconv.Quote(test.s), string(b))
	}

	test = testIDs[1]
	b, err = test.v.MarshalJSON()
	if err != nil || string(b) != strconv.Quote(test.s) {
		t.Errorf("Expect %s, but got %s", strconv.Quote(test.s), string(b))
	}
}

func TestID_UnmarshalJSON(t *testing.T) {
	for _, test := range testIDs {
		var v ID
		err := v.UnmarshalJSON([]byte(test.s))
		if !reflect.DeepEqual(err, test.e) {
			t.Errorf("Expect %s, but got %s", test.e, err)
		} else if test.e == nil && v != test.v {

		}
	}
}

func TestIsTempID(t *testing.T) {
	test := testIDs[0]
	if IsTempID(test.v) == true {
		t.Errorf("%s is not temp id, but returns true", test.v)
	}
	test = testIDs[1]
	if IsTempID(test.v) == false {
		t.Errorf("%s is temp id, but returns false", test.v)
	}
}

func TestUUID_MarshalJSON(t *testing.T) {
	s := "c0afd2be-7576-4fc4-8b3c-696c4c6cf794"
	v := UUID(s)
	b, err := v.MarshalJSON()
	if err != nil || string(b) != strconv.Quote(s) {
		t.Errorf("Expect %s, but got %s", strconv.Quote(s), string(b))
	}
}

func TestUUID_UnmarshalJSON(t *testing.T) {
	s := "c0afd2be-7576-4fc4-8b3c-696c4c6cf794"
	var v UUID
	err := v.UnmarshalJSON([]byte(strconv.Quote(s)))
	if err != nil || v != UUID(s) {
		t.Errorf("Expect %s, but got %s", UUID(s), v)
	}

	s = "invalid"
	err = v.UnmarshalJSON([]byte(s))
	if err == nil {
		t.Error("Expect error, but no error")
	}
}
