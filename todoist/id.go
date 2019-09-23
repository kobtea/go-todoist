package todoist

import (
	"fmt"
	"github.com/satori/go.uuid"
	"strconv"
)

type ID string

func NewID(id string) (ID, error) {
	if IsValidID(ID(id)) {
		return ID(id), nil
	}
	return "", fmt.Errorf("invalid id: %s", id)
}

func NewIDs(ids []string) ([]ID, error) {
	var res []ID
	for _, i := range ids {
		id, err := NewID(i)
		if err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func IsValidID(id ID) bool {
	if _, err := strconv.Atoi(string(id)); err == nil {
		return true
	}
	if IsTempID(id) {
		return true
	}
	return false
}

func (i ID) IsZero() bool {
	s := string(i)
	return s == "0" || s == ""
}

func (i ID) String() string {
	return string(i)
}

func (i ID) MarshalJSON() ([]byte, error) {
	s := string(i)
	if IsTempID(i) {
		s = `"` + s + `"`
	}
	if i.IsZero() {
		s = "null"

	}
	return []byte(s), nil
}

func (i *ID) UnmarshalJSON(b []byte) (err error) {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		s = string(b) // integer id
	}
	if s == "null" {
		s = "0"
	}
	id, err := NewID(s)
	if err != nil {
		return err
	}
	*i = id
	return nil
}

func GenerateTempID() ID {
	u, _ := uuid.NewV4()
	return ID(u.String())
}

func IsTempID(id ID) bool {
	if _, err := uuid.FromString(string(id)); err == nil {
		return true
	}
	return false
}

type UUID string

func GenerateUUID() UUID {
	u, _ := uuid.NewV4()
	return UUID(u.String())
}

func (i UUID) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(i))), nil
}

func (i *UUID) UnmarshalJSON(b []byte) (err error) {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	id, err := uuid.FromString(s)
	if err != nil {
		return err
	}
	*i = UUID(id.String())
	return nil
}
