package main

import (
	"github.com/alecthomas/kong"
	"github.com/ricdeau/accessors/cmd"
	"github.com/ricdeau/accessors/context"
)

func main() {
	newCtx, cancel := context.New()
	defer cancel()

	ctx := kong.Parse(&cmd.Cli{}, kong.Bind(newCtx))
	ctx.FatalIfErrorf(ctx.Run())
}
