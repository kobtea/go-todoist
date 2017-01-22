package todoist

import "testing"

func TestIntBool_MarshalJSON(t *testing.T) {
	s := "0"
	v := IntBool(false)
	b, err := v.MarshalJSON()
	if err != nil || string(b) != s {
		t.Errorf("Expect %s, but got %s", s, string(b))
	}
}

func TestIntBool_UnmarshalJSON(t *testing.T) {
	s := "0"
	var v IntBool
	err := v.UnmarshalJSON([]byte(s))
	if err != nil || v != IntBool(false) {
		t.Errorf("Expect %s, but got %s", IntBool(false), v)
	}

	s = "10"
	err = v.UnmarshalJSON([]byte(s))
	if err == nil {
		t.Error("Expect error, but no error")
	}
}
