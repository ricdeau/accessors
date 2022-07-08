### Golang accessors code generator 

Getters generation:
```go
//go:generate accessors -v -k public,private methods getters,setters ExampleStruct

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

func (e *ExampleStruct) SetScalar(scalar string) {
    e.Scalar = scalar
}

func (e *ExampleStruct) GetPtr() *js.Number {
    return e.Ptr
}

func (e *ExampleStruct) SetPtr(ptr *js.Number) {
    e.Ptr = ptr
}

func (e *ExampleStruct) GetBase() base64.CorruptInputError {
    return e.Base
}

func (e *ExampleStruct) SetBase(base base64.CorruptInputError) {
    e.Base = base
}

func (e *ExampleStruct) GetTime() time.Time {
    return e.Time
}

func (e *ExampleStruct) SetTime(time time.Time) {
    e.Time = time
}

func (e *ExampleStruct) GetIface() io.ReadCloser {
    return e.Iface
}

func (e *ExampleStruct) SetIface(iface io.ReadCloser) {
    e.Iface = iface
}

func (e *ExampleStruct) GetArray() [3]int {
    return e.Array
}

func (e *ExampleStruct) SetArray(array [3]int) {
    e.Array = array
}

func (e *ExampleStruct) GetSlice() []*os.File {
    return e.Slice
}

func (e *ExampleStruct) SetSlice(slice []*os.File) {
    e.Slice = slice
}

func (e *ExampleStruct) GetMap() map[io.ByteWriter]string {
    return e.Map
}

func (e *ExampleStruct) SetMap(_map map[io.ByteWriter]string) {
    e.Map = _map
}

func (e *ExampleStruct) GetInCh() chan<- string {
    return e.inCh
}

func (e *ExampleStruct) SetInCh(inCh chan<- string) {
    e.inCh = inCh
}

func (e *ExampleStruct) GetOutCh() <-chan int {
    return e.outCh
}

func (e *ExampleStruct) SetOutCh(outCh <-chan int) {
    e.outCh = outCh
}

func (e *ExampleStruct) GetCh() chan *float64 {
    return e.ch
}

func (e *ExampleStruct) SetCh(ch chan *float64) {
    e.ch = ch
}

func (e *ExampleStruct) GetPrivate() int {
    return e.private
}

func (e *ExampleStruct) SetPrivate(private int) {
    e.private = private
}
```
```go
//go:generate accessors -v -o example_iface.go interface --name=Example --pkg=iface getters,setters ExampleStruct Scalar Ptr
```
File `example_iface.go`
```go
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
```

More info:
```bash
$ accessors --help
```