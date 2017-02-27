package util

import (
	"errors"
	"github.com/kobtea/go-todoist/todoist"
)

func ProcessIDs(ids []string, f func([]todoist.ID) error) error {
	if len(ids) < 1 {
		return errors.New("Require ID(s)")
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
