package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "configure about this CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := os.ExpandEnv("$HOME/.go-todoist")
		if _, err := os.Stat(dir); err != nil {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
		var c util.Config
		file := path.Join(dir, "config.json")
		if b, err := ioutil.ReadFile(file); err != nil {
			// initial config
		} else {
			// modify config
			if err = json.Unmarshal(b, &c); err != nil {
				return err
			}
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("todoist token (default: %s): ", c.Token)
		if ans, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if ans != "\n" {
				c.Token = strings.TrimSpace(ans)
			}
		}
		if b, err := json.MarshalIndent(c, "", "  "); err != nil {
			return err
		} else {
			if err = ioutil.WriteFile(file, b, 0644); err != nil {
				return err
			}
		}
		fmt.Printf("write config to %s\n", file)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
}
