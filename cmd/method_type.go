package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ricdeau/accessors/consts"
)

var _ json.Unmarshaler = (*MethodType)(nil)

type MethodType struct {
	types map[string]struct{}
}

func (c *MethodType) Getters() bool {
	_, ok := c.types[consts.MethodTypeGetters]
	return ok
}

func (c *MethodType) Setters() bool {
	_, ok := c.types[consts.MethodTypeSetters]
	return ok
}

func (m *MethodType) UnmarshalJSON(text []byte) error {
	types := strings.Split(strings.Trim(string(text), `"`), ",")
	*m = MethodType{
		types: map[string]struct{}{},
	}

	for _, t := range types {
		m.types[strings.ToLower(strings.TrimSpace(t))] = struct{}{}
	}

	return nil
}

func (m *MethodType) Validate() error {
	if len(m.types) < 1 || len(m.types) > 2 {
		return fmt.Errorf("invalid method types len: %d", len(m.types))
	}

	for t := range m.types {
		switch t {
		case consts.MethodTypeGetters, consts.MethodTypeSetters:
		default:
			return fmt.Errorf("invalid method type: %s", t)
		}
	}

	return nil
}
