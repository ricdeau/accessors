package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ricdeau/accessors/consts"
)

var _ json.Unmarshaler = (*FieldKind)(nil)

type FieldKind struct {
	kinds map[string]struct{}
}

func (f *FieldKind) Public() bool {
	_, ok := f.kinds[consts.FieldKindPublic]
	return ok
}

func (f *FieldKind) Private() bool {
	_, ok := f.kinds[consts.FieldKindPrivate]
	return ok
}

func (f *FieldKind) UnmarshalJSON(text []byte) error {
	kinds := strings.Split(strings.Trim(string(text), `"`), ",")
	*f = FieldKind{
		kinds: map[string]struct{}{},
	}

	for _, k := range kinds {
		f.kinds[strings.ToLower(strings.TrimSpace(k))] = struct{}{}
	}

	return nil
}

func (f *FieldKind) Validate() error {
	if len(f.kinds) == 0 {
		return nil
	}

	if len(f.kinds) > 2 {
		return fmt.Errorf("invalid field kinds len: %d", len(f.kinds))
	}

	for k := range f.kinds {
		switch k {
		case consts.FieldKindPublic, consts.FieldKindPrivate:
		default:
			return fmt.Errorf("invalid field kind: %s", k)
		}
	}

	return nil
}
