package util

import (
	"errors"
	"github.com/kobtea/go-todoist/todoist"
)

func ProcessID(id string, f func(todoist.ID) error) error {
	if len(id) == 0 {
		return errors.New("require id")
	}
	newID, err := todoist.NewID(id)
	if err != nil {
		return err
	}
	if err = f(newID); err != nil {
		return err
	}
	return nil
}

func ProcessIDs(ids []string, f func([]todoist.ID) error) error {
	if len(ids) == 0 {
		return errors.New("require id(s)")
	}
	newIDs, err := todoist.NewIDs(ids)
	if err != nil {
		return err
	}
	if err = f(newIDs); err != nil {
		return err
	}
	return nil
}
