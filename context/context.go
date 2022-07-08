package context

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type Printer interface {
	Println(v ...any)
	Printf(format string, v ...any)
}

type Context struct {
	Printer         `json:"-"`
	context.Context `json:"-"`

	Output        string            `json:"output"`
	PkgName       string            `json:"pkgName"`
	StructName    string            `json:"structName"`
	Fields        []string          `json:"fields"`
	PublicFields  bool              `json:"publicFields"`
	PrivateFields bool              `json:"privateFields"`
	Getters       bool              `json:"getters"`
	Setters       bool              `json:"setters"`
	Interface     *InterfaceContext `json:"interface,omitempty"`
}

type InterfaceContext struct {
	Name    string `json:"name"`
	Package string `json:"package"`
}

type noopPrinter struct{}

func (n noopPrinter) Println(...any)        {}
func (n noopPrinter) Printf(string, ...any) {}

func New() (*Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	return &Context{
		Context: ctx,
		Printer: noopPrinter{},
	}, cancel
}

func (c *Context) GetOutput(newFile bool) (out io.Writer, err error) {
	if c.Output == "" {
		return os.Stdout, nil
	}

	if newFile {
		out, err = os.Create(c.Output)
		if err != nil {
			return nil, fmt.Errorf("create file for writig: %v", err)
		}
	} else {
		out, err = os.OpenFile(c.Output, os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("open existing file for append: %v")
		}
	}

	return out, nil
}

func (c *Context) ContainsField(fieldName string) bool {
	for _, name := range c.Fields {
		if name == fieldName {
			return true
		}
	}

	return false
}
