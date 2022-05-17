package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/durandj/ley/cmd/manager/subcommand"
	"github.com/fatih/color"
)

func main() {
	ctx := context.Background()
	ctx, done := signal.NotifyContext(ctx, os.Interrupt, os.Kill)

	if err := subcommand.NewRootCommand().ExecuteContext(ctx); err != nil {
		color.Red(err.Error())
		done()
		os.Exit(1)
	}
}
