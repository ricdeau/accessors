### Golang accessors code generator 

Getters generation:
```go

//go:generate accessors -t all ExampleStruct

type ExampleStruct struct {
    Scalar  string
    Ptr     *js.Number
    Base    base64.CorruptInputError
    Time    time.Time
    Iface   io.ReadCloser
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
```

More info:
```bash
$ accessors --help
```