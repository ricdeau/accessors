package main

import (
	"fmt"
	"os"
)

type Ctx struct {
	Verbose    bool      `json:"verbose" short:"v" optional:"" help:"Output verbosity."`
	Out        string    `json:"out" short:"o" type:"path" optional:"" help:"Output file path."`
	PkgName    string    `json:"pkg_name" short:"p" optional:"" help:"Target go struct package."`
	FieldType  FieldType `json:"field_type" short:"t" optional:"" help:"Required field types: public, private or all."`
	StructName string    `json:"struct_name" arg:"" required:"" help:"Name of go struct."`
	FieldNames []string  `json:"field_names" arg:"" optional:"" help:"Field names"`
}

func (c *Ctx) Validate() error {
	if c.FieldType == 0 && len(c.FieldNames) == 0 {
		return fmt.Errorf("must set field-type or field-names")
	}
	if c.FieldType > 0 && len(c.FieldNames) > 0 {
		return fmt.Errorf("can't set both field-type and field-names")
	}

	if c.PkgName == "" {
		c.PkgName = os.Getenv("GOPACKAGE")
	}

	if c.Out == "" {
		c.Out = os.Getenv("GOFILE")
	}

	return nil
}
