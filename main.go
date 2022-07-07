package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/alecthomas/kong"
	"github.com/iancoleman/strcase"
	"github.com/ricdeau/enki"
	"golang.org/x/tools/go/packages"
)

// accessors SomeStruct all
// accessors SomeStruct private
// accessors SomeStruct public
// accessors SomeStruct field1 Field2
func main() {
	var cli Ctx
	ctx := kong.Parse(&cli)
	ctx.FatalIfErrorf(ctx.Run())
}

func (c *Ctx) Run() error {
	absPath, err := filepath.Abs(c.Out)
	if err != nil {
		return fmt.Errorf("resolve absolute file path for %s", c.Out)
	}

	c.Out = absPath

	if c.Verbose {
		indent, err := json.MarshalIndent(c, "", "  ")
		if err == nil {
			_, _ = os.Stdout.Write(indent)
		}
	}

	conf := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles |
			packages.NeedCompiledGoFiles | packages.NeedDeps |
			packages.NeedTypes | packages.NeedTypesInfo |
			packages.NeedImports | packages.NeedSyntax,
	}
	pkgs, err := packages.Load(conf, "..."+c.PkgName)
	if err != nil {
		return fmt.Errorf("load packages: %v", err)
	}

	var syntax []*ast.File
	var file *ast.File
	var obj *ast.Object
pkgsLoop:
	for _, p := range pkgs {
		for _, s := range p.Syntax {
			obj = s.Scope.Lookup(c.StructName)
			if obj != nil {
				syntax = p.Syntax
				file = s
				break pkgsLoop
			}
		}
	}

	if obj == nil {
		return fmt.Errorf("struct with name %s not found in package %s", c.StructName, c.PkgName)
	}

	decl, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		return fmt.Errorf("%s is not a type", c.StructName)
	}

	struc, ok := decl.Type.(*ast.StructType)
	if !ok {
		return fmt.Errorf("%s is not a struct", c.StructName)
	}

	imports := map[string]string{}
	methods := []enki.Block{}
	fields := c.filterFields(struc)
	for _, f := range fields {
		methodName := fmt.Sprintf("Get%s", strcase.ToCamel(f.name))
		if c.methodExists(methodName, syntax) {
			continue
		}
		typeName := typeToString(file.Imports, imports, f.typ)
		def := methGen(methodName, c.StructName, f.name, typeName)
		methods = append(methods, def)
	}

	if len(methods) == 0 {
		return nil
	}

	newFile := true
	scope, ok := pkgs[0].TypesInfo.Scopes[file]
	if ok {
		fileName := reflect.ValueOf(scope).Elem().FieldByName("comment").String()
		if fileName == absPath {
			newFile = false
		}
	}

	var out io.Writer
	outFile := enki.NewFile()
	if newFile {
		out, err = os.Create(absPath)
		if err != nil {
			return fmt.Errorf("create file for writig: %v", err)
		}

		outFile.Package(file.Name.Name)
		for alias, name := range imports {
			if alias == filepath.Base(name) {
				alias = ""
			}
			outFile.Import(alias, name)
		}
	} else {
		out, err = os.OpenFile(absPath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return fmt.Errorf("open existing file for append: %v")
		}
	}

	for i, method := range methods {
		outFile.Add(method)
		if i < len(methods)-1 {
			outFile.NewLine()
		}
	}

	err = outFile.Write(out)
	if err != nil {
		return fmt.Errorf("write content: %v", err)
	}

	return nil
}

type field struct {
	name string
	typ  ast.Expr
}

func (c *Ctx) filterFields(struc *ast.StructType) (res []field) {
	for _, f := range struc.Fields.List {
		for _, ff := range f.Names {
			if c.FieldType > 0 {
				firstRune, _ := utf8.DecodeRuneInString(ff.Name)
				if firstRune == utf8.RuneError {
					log.Printf("[ERROR] get first rune of %s", ff.Name)
					continue
				}
				if !c.FieldType.Includes(Public) && unicode.IsUpper(firstRune) {
					continue
				}
				if !c.FieldType.Includes(Private) && unicode.IsLower(firstRune) {
					continue
				}
			} else if !c.containsFieldName(ff.Name) {
				continue
			}

			res = append(res, field{
				name: ff.Name,
				typ:  f.Type,
			})
		}
	}

	return res
}

func (c *Ctx) containsFieldName(fieldName string) bool {
	for _, name := range c.FieldNames {
		if name == fieldName {
			return true
		}
	}

	return false
}

func (c *Ctx) methodExists(methodName string, files []*ast.File) bool {
	for _, file := range files {
		for _, decl := range file.Decls {
			if f, ok := decl.(*ast.FuncDecl); ok {
				if f.Name.Name != methodName {
					continue
				}
				if f.Recv == nil {
					continue
				}
				for _, r := range f.Recv.List {
					clear := typeClearName(r.Type)
					if clear == c.StructName {
						return true
					}
				}
			}
		}
	}

	return false
}

func methGen(methodName, receiverName, fieldName, fieldType string) enki.Block {
	receiverAlias := strings.ToLower(receiverName[:1])
	receiver := fmt.Sprintf("%s *%s", receiverAlias, receiverName)
	return enki.M(methodName).Receiver(receiver).Returns(fieldType).Body(
		enki.Stmt().Line("return @1.@2", receiverAlias, fieldName),
	)
}

func typeToString(importSpecs []*ast.ImportSpec, imports map[string]string, typ ast.Expr) string {
	switch v := typ.(type) {
	case *ast.Ident:
		return indentName(v)
	case *ast.StarExpr:
		if s := typeToString(importSpecs, imports, v.X); s != "" {
			return "*" + s
		}
	case *ast.SelectorExpr:
		alias := indentName(v.X.(*ast.Ident))
		importStr, ok := findImport(importSpecs, alias)
		if ok {
			imports[alias] = importStr
		}
		if s := typeToString(importSpecs, imports, v.Sel); s != "" {
			return alias + "." + s
		}
	case *ast.SliceExpr:
		if s := typeToString(importSpecs, imports, v.X); s != "" {
			return "[]" + s
		}
	case *ast.ArrayType:
		if s := typeToString(importSpecs, imports, v.Elt); s != "" {
			l := ""
			if v.Len != nil {
				l = v.Len.(*ast.BasicLit).Value
			}
			return "[" + l + "]" + s
		}
	case *ast.MapType:
		if key := typeToString(importSpecs, imports, v.Key); key != "" {
			if val := typeToString(importSpecs, imports, v.Value); val != "" {
				return "map[" + key + "]" + val
			}
		}
	case *ast.ChanType:
		if s := typeToString(importSpecs, imports, v.Value); s != "" {
			ch := "chan "
			switch v.Dir {
			case ast.RECV:
				ch = "<-chan "
			case ast.SEND:
				ch = "chan<- "
			}

			return ch + s
		}
	}

	return ""
}

func typeClearName(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return typeClearName(v.Sel)
	case *ast.StarExpr:
		return typeClearName(v.X)
	}

	return ""
}

func findImport(importSpecs []*ast.ImportSpec, alias string) (importStr string, ok bool) {
	for _, spec := range importSpecs {
		name := indentName(spec.Name)
		val := strings.Trim(spec.Path.Value, `"`)
		if name == alias {
			return val, true
		} else if name == "" {
			base := path.Base(val)
			if base == alias {
				return val, true
			}
		}
	}

	return "", false
}

func indentName(indent *ast.Ident) string {
	if indent == nil {
		return ""
	}

	return indent.Name
}
