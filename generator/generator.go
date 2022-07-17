package generator

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ricdeau/accessors/collector"
	"github.com/ricdeau/accessors/context"
	"github.com/ricdeau/enki"
)

type generator struct {
	ctx *context.Context
}

// New creates new generator.
func New(ctx *context.Context) *generator {
	return &generator{ctx: ctx}
}

// Generate generates methods and/or interface.
func (g *generator) Generate() error {
	g.printIndent("Load context", g.ctx)

	c := collector.New(g.ctx)
	info, err := c.CollectPackageInfo()
	if err != nil {
		return err
	}

	g.printIndent("Collect package info", info)

	newFile := true
	if info.FileName == g.ctx.Output {
		newFile = false
	}

	g.ctx.Printf("Ready to write to file: isNewFile=%v, path=%v\n", newFile, g.ctx.Output)

	out, err := g.ctx.GetOutput(newFile)
	if err != nil {
		return fmt.Errorf("make file for writing: %v", err)
	}

	outFile := enki.NewFile()
	if newFile {
		outFile.Package(info.Package.Name)
		for alias, name := range c.GetImports() {
			if alias == filepath.Base(name) {
				alias = ""
			}
			outFile.Import(alias, name)
		}
	}

	if g.ctx.Interface != nil {
		g.generateInterface(c, info, outFile, newFile)
	} else {
		methCount := g.generateMethods(c, info, outFile, newFile)
		if methCount == 0 {
			return nil
		}
	}

	err = outFile.Write(out)
	if err != nil {
		return fmt.Errorf("write content: %v", err)
	}

	return nil
}

func (g *generator) generateInterface(c *collector.Collector, info *collector.Info, outFile enki.File, newFile bool) {
	methods := c.CollectMethods()
	g.printIndent("Collect struct methods", methods)

	if newFile {
		if pkg := g.ctx.Interface.Package; pkg != "" {
			outFile.Package(pkg)
		} else {
			outFile.Package(info.Package.Name)
		}

		for alias, name := range c.GetImports() {
			if alias == filepath.Base(name) {
				alias = ""
			}
			outFile.Import(alias, name)
		}
	}

	definitions := make([]enki.Function, 0, len(methods))

	for _, meth := range methods {
		if g.ctx.Getters {
			definitions = append(definitions, meth.GetterDefinition())
		}

		if g.ctx.Setters {
			definitions = append(definitions, meth.SetterDefinition())
		}
	}

	outFile.Add(enki.T(g.ctx.Interface.Name).Interface(definitions...))
	g.ctx.Println("Add interface: " + g.ctx.Interface.Name)

}

func (g *generator) generateMethods(c *collector.Collector, info *collector.Info, outFile enki.File, newFile bool) int {
	methods := c.CollectMethods()
	g.printIndent("Collect struct methods", methods)

	if newFile {
		outFile.Package(info.Package.Name)
		for alias, name := range c.GetImports() {
			if alias == filepath.Base(name) {
				alias = ""
			}
			outFile.Import(alias, name)
		}
	}

	count := 0
	for i, meth := range methods {
		if g.ctx.Getters && !c.MethodExists(meth.GetterName()) {
			outFile.Add(meth.Getter())
			g.ctx.Println("Add method: " + meth.GetterName())
			count++
		}

		if g.ctx.Setters && !c.MethodExists(meth.SetterName()) {
			outFile.NewLine()
			outFile.Add(meth.Setter())
			g.ctx.Println("Add method: " + meth.SetterName())
			count++
		}

		if i < len(methods)-1 {
			outFile.NewLine()
		}
	}

	return count
}

func (g *generator) printIndent(message string, v any) {
	indent, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		g.ctx.Println(message + ": " + string(indent))
	} else {
		g.ctx.Printf("Failed to print: %v\n", err)
	}
}
