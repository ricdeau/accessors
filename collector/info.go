package collector

import (
	"encoding/json"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

// Info collected info.
type Info struct {
	FileName string
	Object   *ast.Object
	Package  *packages.Package
	Syntax   []*ast.File
	File     *ast.File
	Struct   *ast.StructType
}

type flatInfo struct {
	FileName         string            `json:"fileName"`
	ObjectName       string            `json:"objectName"`
	ObjectKind       ast.ObjKind       `json:"objectKind"`
	Package          *packages.Package `json:"package"`
	SyntaxFilesCount int               `json:"syntaxFilesCount"`
}

func (i *Info) MarshalJSON() ([]byte, error) {
	flat := &flatInfo{
		FileName:         i.FileName,
		ObjectName:       i.Object.Name,
		ObjectKind:       i.Object.Kind,
		Package:          i.Package,
		SyntaxFilesCount: len(i.Syntax),
	}

	return json.Marshal(flat)
}
