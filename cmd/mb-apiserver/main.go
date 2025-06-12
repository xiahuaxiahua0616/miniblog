package main

import (
	"os"

	"github.com/xiahuaxiahua0616/miniblog/cmd/mb-apiserver/app"
	_ "go.uber.org/automaxprocs"
)

func main() {
	command := app.NewMiniBlogCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
