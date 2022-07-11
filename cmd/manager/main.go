package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/durandj/ley/cmd/manager/subcommand"
	"github.com/fatih/color"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	ctx, done := signal.NotifyContext(ctx, os.Interrupt, os.Kill)

	cmd := subcommand.NewRootCommand()
	if err := cmd.ExecuteContext(ctx); err != nil {
		color.Red(err.Error())
		done()
		os.Exit(1)
	}
}
