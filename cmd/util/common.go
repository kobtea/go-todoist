package util

import (
	"context"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/viper"
)

func NewClient() (*todoist.Client, error) {
	return todoist.NewClient(
		"",
		viper.GetString("TODOIST_TOKEN"),
		"*",
		"",
		nil)
}

func AutoCommit(f func(client todoist.Client, ctx context.Context) error) error {
	client, err := NewClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	if err = f(*client, ctx); err != nil {
		return err

	}
	if err = client.Commit(ctx); err != nil {
		return err
	}
	if err = client.FullSync(ctx, []todoist.Command{}); err != nil {
		return err
	}
	return nil
}
