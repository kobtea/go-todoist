package todoist

import "fmt"

type IntBool bool

func (i IntBool) Bool() bool {
	if i {
		return true
	} else {
		return false
	}
}

func (i IntBool) MarshalJSON() ([]byte, error) {
	if i {
		return []byte("1"), nil
	} else {
		return []byte("0"), nil
	}
}

func (i *IntBool) UnmarshalJSON(b []byte) (err error) {
	switch string(b) {
	case "1":
		*i = true
	case "0":
		*i = false
	default:
		return fmt.Errorf("Could not unmarshal into intbool: %s", string(b))
	}
	return nil
}

type ColorStringer interface {
	String() string
	ColorString() string
}

type NoColorString struct {
	s string
}

func NewNoColorString(s string) NoColorString {
	return NoColorString{s}
}

func (n NoColorString) String() string {
	return n.s
}

func (n NoColorString) ColorString() string {
	return n.s
}
