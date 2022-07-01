package main

import (
	"encoding"
	"fmt"
	"strings"
)

var (
	_ fmt.Stringer             = FieldType(0)
	_ encoding.TextUnmarshaler = (*FieldType)(nil)
)

type FieldType uint8

const (
	Public FieldType = 1 << iota
	Private
	All = Public | Private
)

const (
	publicName  = "public"
	privateName = "private"
	allName     = "all"
)

var (
	methodTypeNames = map[FieldType]string{
		Public:  publicName,
		Private: privateName,
		All:     allName,
	}
	methodTypeValues = map[string]FieldType{
		publicName:  Public,
		privateName: Private,
		allName:     All,
	}
)

func (m FieldType) Includes(typ FieldType) bool {
	return m&typ > 0
}

func (m *FieldType) UnmarshalText(text []byte) error {
	name := strings.ToLower(string(text))
	typ, ok := methodTypeValues[name]
	if !ok {
		return fmt.Errorf("unknown FieldType: %s", name)
	}

	*m = typ

	return nil
}

func (m FieldType) String() string {
	if name, ok := methodTypeNames[m]; ok {
		return name
	}

	return fmt.Sprintf("FieldType(%d)", m)
}
