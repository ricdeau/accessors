package main

import (
	"github.com/alecthomas/kong"
	"github.com/ricdeau/accessors/cmd"
	"github.com/ricdeau/accessors/context"
)

type cli struct {
	cmd.Cmd
}

func (c *cli) AfterApply(ctx *context.Context) error {
	c.Cmd.AfterApply(ctx)

	if !(ctx.Getters || ctx.Setters) {
		ctx.Getters = true
		ctx.Setters = true
	}

	if len(ctx.Fields) == 0 {
		ctx.PublicFields = true
		ctx.PrivateFields = true
	}

	return nil
}

func main() {
	newCtx, cancel := context.New()
	defer cancel()

	ctx := kong.Parse(&cli{}, kong.Bind(newCtx))
	ctx.FatalIfErrorf(ctx.Run())
}
