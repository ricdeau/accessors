package iface

import (
	js "encoding/json"
)

type Example interface {
	GetScalar() string
	SetScalar(scalar string)
	GetPtr() *js.Number
	SetPtr(ptr *js.Number)
}
