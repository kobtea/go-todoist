package util

import (
	"context"
	"encoding/json"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

type Config struct {
	Token string `json:"token"`
}

func resolveToken() string {
	if s := viper.GetString("TODOIST_TOKEN"); len(s) != 0 {
		return s
	}
	file := os.ExpandEnv("$HOME/.go-todoist/config.json")
	if b, err := ioutil.ReadFile(file); err != nil {
		return ""
	} else {
		var c Config
		if err = json.Unmarshal(b, &c); err != nil {
			return ""
		}
		return c.Token
	}
}

func NewClient() (*todoist.Client, error) {
	return todoist.NewClient(
		"",
		resolveToken(),
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
