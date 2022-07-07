package example

import (
	"encoding/base64"
	js "encoding/json"
	"io"
	"os"
	"time"
)

//go:generate accessors -v -t all ExampleStruct

type ExampleStruct struct {
	Scalar  string
	Ptr     *js.Number
	Base    base64.CorruptInputError
	Time    time.Time
	Iface   io.ReadCloser
	Array   [3]int
	Slice   []*os.File
	Map     map[io.ByteWriter]string
	inCh    chan<- string
	outCh   <-chan int
	ch      chan *float64
	private int
}

func (e *ExampleStruct) GetScalar() string {
	return e.Scalar
}

func (e *ExampleStruct) GetPtr() *js.Number {
	return e.Ptr
}

func (e *ExampleStruct) GetBase() base64.CorruptInputError {
	return e.Base
}

func (e *ExampleStruct) GetTime() time.Time {
	return e.Time
}

func (e *ExampleStruct) GetIface() io.ReadCloser {
	return e.Iface
}

func (e *ExampleStruct) GetPrivate() int {
	return e.private
}

func (e *ExampleStruct) GetArray() [3]int {
	return e.Array
}

func (e *ExampleStruct) GetSlice() []*os.File {
	return e.Slice
}

func (e *ExampleStruct) GetMap() map[io.ByteWriter]string {
	return e.Map
}

func (e *ExampleStruct) GetInCh() chan<- string {
	return e.inCh
}

func (e *ExampleStruct) GetOutCh() <-chan int {
	return e.outCh
}

func (e *ExampleStruct) GetCh() chan *float64 {
	return e.ch
}
