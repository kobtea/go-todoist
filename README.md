# go-todoist

[![Go Report Card](https://goreportcard.com/badge/github.com/kobtea/go-todoist)](https://goreportcard.com/report/github.com/kobtea/go-todoist)
[![CircleCI](https://circleci.com/gh/kobtea/go-todoist.svg?style=svg)](https://circleci.com/gh/kobtea/go-todoist)

Unofficial CLI and library for [todoist](https://todoist.com).


## Install

### Binary

Go to https://github.com/kobtea/go-todoist/releases

### Building from source

```bash
$ go get -d github.com/kobtea/go-todoist
$ cd $GOPATH/src/github.com/kobtea/go-todoist
$ dep ensure
$ make build
```

## Usage

```bash
$ todoist help
Command line tool for todoist.

Usage:
  todoist [command]

Available Commands:
  completion  generate completion script
  config      configure about this CLI
  filter      subcommand for filter
  help        Help about any command
  inbox       show inbox tasks
  item        subcommand for item
  label       subcommand for label
  next        show next 7 days tasks
  project     subcommand for project
  review      show completed items
  sync        Syncronize origin server
  today       show today's tasks
  version     show version of go-todoist

Flags:
      --config string   config file (default is $HOME/.todoist.yaml)
  -h, --help            help for todoist

Use "todoist [command] --help" for more information about a command.
```

Configure your API token of todoist.  
The location of API token is `Todoist` > `Settings` > `Integrations` > `API token`.

```bash
$ todoist config
todoist token (default: ): YOUR_TOKEN_HERE
write config to /home/kobtea/.go-todoist/config.json
```

Sync contents.

```bash
$ todoist sync
```

Enjoy.

```bash
$ todoist item add hello go-todoist
$ todoist inbox
```

Bash and zsh completion are supported ;)  
Completion requires [fzf](https://github.com/junegunn/fzf).

```bash
# bash
$ . <(todoist completion bash)
# zsh
$ . <(todoist completion zsh)
```


## As a Library

This library supports [sync api v8](https://developer.todoist.com/sync/v8).  
Implementation refers to [python official library](https://github.com/doist/todoist-python).

sample

```go
package main

import (
	"context"
	"github.com/kobtea/go-todoist/todoist"
	"os"
)

func main() {
	token := os.Getenv("TODOIST_TOKEN")
	cli, _ := todoist.NewClient("", token, "*", "", nil)
	ctx := context.Background()

	// sync contents
	cli.FullSync(ctx, []todoist.Command{})

	// add item
	item := todoist.Item{Content: "hello go-todoist hogehoge"}
	cli.Item.Add(item)
	cli.Commit(ctx)
}
```


## License

MIT
