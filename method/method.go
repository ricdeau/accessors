package method

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/ricdeau/enki"
	"github.com/ricdeau/strcase"
)

// Method methods info.
type Method struct {
	StructName string `json:"structName"`
	FieldName  string `json:"fieldName"`
	FieldType  string `json:"fieldType"`

	receiverVar string
	receiver    string
	param       string
}

// New new method info.
func New(structName string, fieldName string, fieldType string) *Method {
	receiverVar := strings.ToLower(structName[:1])
	return &Method{
		StructName: structName,
		FieldName:  fieldName,
		FieldType:  fieldType,

		receiverVar: receiverVar,
		receiver:    fmt.Sprintf("%s *%s", receiverVar, structName),
		param:       normalizeParamName(fieldName),
	}
}

// Getter get block for method.
func (m *Method) Getter() enki.Block {
	return enki.M(m.GetterName()).Receiver(m.receiver).Returns(m.FieldType).Body(
		enki.Stmt().Line("return @1.@2", m.receiverVar, m.FieldName),
	)
}

// GetterName name of the get method.
func (m *Method) GetterName() string {
	return "Get" + strcase.GoExported(m.FieldName)
}

// GetterDefinition get method definition for interface.
func (m *Method) GetterDefinition() enki.Function {
	return enki.Def("Get" + strcase.GoExported(m.FieldName)).Returns(m.FieldType)
}

// Setter set method block.
func (m *Method) Setter() enki.Block {
	return enki.M(m.SetterName()).Receiver(m.receiver).Params(m.param + " " + m.FieldType).Body(
		enki.Stmt().Line("@1.@2 = @3", m.receiverVar, m.FieldName, m.param),
	)
}

// SetterName name of the set method.
func (m *Method) SetterName() string {
	return "Set" + strcase.GoExported(m.FieldName)
}

// SetterDefinition set method definition for interface.
func (m *Method) SetterDefinition() enki.Function {
	return enki.Def("Set" + strcase.GoExported(m.FieldName)).Params(m.param + " " + m.FieldType)
}

func normalizeParamName(fieldName string) string {
	paramName := strcase.GoUnexported(fieldName)
	if token.Lookup(paramName).IsKeyword() {
		paramName = "_" + paramName
	}

	return paramName
}
