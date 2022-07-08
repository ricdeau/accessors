package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ricdeau/accessors/context"
	"github.com/ricdeau/accessors/generator"
)

// Cli all possible commands definition.
type Cli struct {
	Debug      bool         `short:"d" optional:"" help:"Debug mode."`
	Output     string       `short:"o" type:"path" optional:"" help:"Output file path."`
	PkgName    string       `short:"p" optional:"" help:"Target go struct package."`
	FieldKinds FieldKind    `short:"k" optional:"" help:"Field kind. Can contain public, private values separated by comma."`
	Methods    MethodsCmd   `cmd:"" help:"Generate methods for struct."`
	Interface  InterfaceCmd `cmd:"" help:"Generate interface for struct."`
}

// MethodsCmd command for methods codegen.
type MethodsCmd struct {
	cmd `embed:""`
}

// InterfaceCmd command for interface codegen.
type InterfaceCmd struct {
	Pkg  string `optional:"" help:"Interface package if needed. If empty - same as source struct package."`
	Name string `optional:"" help:"Interface type name. If empty - <StructName>Interface."`
	cmd  `embed:""`
}

type cmd struct {
	MethodType MethodType `arg:""`
	StructName string     `arg:"" required:"" help:"Name of go struct."`
	Fields     []string   `arg:"" optional:"" help:"Concrete field names if needed. If set, field kind will be ignored."`
}

func (c *Cli) AfterApply(ctx *context.Context) (err error) {
	if c.Debug {
		ctx.Printer = log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)
	}

	ctx.Output = c.Output
	ctx.PkgName = c.PkgName
	ctx.PublicFields = c.FieldKinds.Public()
	ctx.PrivateFields = c.FieldKinds.Private()

	if ctx.Output == "" {
		ctx.Output = os.Getenv("GOFILE")
	}

	if ctx.Output != "" {
		absPath, err := filepath.Abs(ctx.Output)
		if err != nil {
			return fmt.Errorf("resolve absolute file path for %s", ctx.Output)
		}
		ctx.Output = absPath
	}

	return nil
}

func (c *MethodsCmd) AfterApply(ctx *context.Context) error {
	ctx.Getters = c.MethodType.Getters()
	ctx.Setters = c.MethodType.Setters()
	ctx.StructName = c.StructName
	ctx.Fields = c.Fields

	return nil
}

func (c *InterfaceCmd) AfterApply(ctx *context.Context) error {
	ctx.Getters = c.MethodType.Getters()
	ctx.Setters = c.MethodType.Setters()
	ctx.StructName = c.StructName
	ctx.Fields = c.Fields

	var ifaceName string
	if c.Name != "" {
		ifaceName = c.Name
	} else {
		ifaceName = ctx.StructName + "Interface"
	}
	ctx.Interface = &context.InterfaceContext{
		Name:    ifaceName,
		Package: c.Pkg,
	}

	return nil
}

func (c *MethodsCmd) Run(ctx *context.Context) error {
	return generator.New(ctx).Generate()
}

func (c *InterfaceCmd) Run(ctx *context.Context) error {
	return generator.New(ctx).Generate()
}
