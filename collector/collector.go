package collector

import (
	"fmt"
	"go/ast"
	go_types "go/types"
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/ricdeau/accessors/context"
	"github.com/ricdeau/accessors/method"
	"github.com/ricdeau/accessors/types"
	"golang.org/x/tools/go/packages"
)

const loadMode = packages.NeedName | packages.NeedFiles |
	packages.NeedCompiledGoFiles | packages.NeedDeps |
	packages.NeedTypes | packages.NeedTypesInfo |
	packages.NeedImports | packages.NeedSyntax

// Collector methods and info collector.
type Collector struct {
	*context.Context
	info    *Info
	methods []*method.Method
	imports map[string]string
}

// New creates new collector.
func New(context *context.Context) *Collector {
	return &Collector{
		Context: context,
		imports: map[string]string{},
	}
}

type FieldInfo struct {
	Name string
	Expr ast.Expr
}

// CollectPackageInfo collects package info.
func (c *Collector) CollectPackageInfo() (*Info, error) {
	if c.info == nil {
		if err := c.loadInfo(); err != nil {
			return nil, err
		}
	}

	return c.info, nil
}

// CollectMethods collect methods.
func (c *Collector) CollectMethods() []*method.Method {
	methods := []*method.Method{}
	for _, f := range c.filterFields() {
		fieldType := types.TypeName(c.info.File.Imports, c.imports, f.Expr)
		methods = append(methods, method.New(c.StructName, f.Name, fieldType))
	}

	return methods
}

// GetImports get imports from package.
func (c *Collector) GetImports() map[string]string {
	return c.imports
}

// MethodExists check if method exist.
func (c *Collector) MethodExists(methodName string) bool {
	for _, file := range c.info.Syntax {
		for _, decl := range file.Decls {
			if f, ok := decl.(*ast.FuncDecl); ok {
				if f.Name.Name != methodName {
					continue
				}
				if f.Recv == nil {
					continue
				}
				for _, r := range f.Recv.List {
					clear := types.TypeClearName(r.Type)
					if clear == c.StructName {
						return true
					}
				}
			}
		}
	}

	return false
}

func (c *Collector) loadInfo() (err error) {
	conf := &packages.Config{
		Context: c.Context,
		Mode:    loadMode,
	}

	var pkgs []*packages.Package
	if c.PkgName != "" {
		pkgs, err = packages.Load(conf, c.PkgName)
	} else {
		pkgs, err = packages.Load(conf)
	}
	if err != nil {
		return fmt.Errorf("load packages: %v", err)
	}

	c.info = &Info{}

pkgsLoop:
	for _, p := range pkgs {
		for _, f := range p.Syntax {
			obj := f.Scope.Lookup(c.StructName)
			if obj != nil {
				c.info.Object = obj
				c.info.Package = p
				c.info.Syntax = p.Syntax
				c.info.File = f
				c.info.FileName = getFileName(p.TypesInfo.Scopes, f)

				break pkgsLoop
			}
		}
	}

	if c.info.Object == nil {
		return fmt.Errorf("struct with name %s not found in package %s", c.StructName, c.PkgName)
	}

	decl, ok := c.info.Object.Decl.(*ast.TypeSpec)
	if !ok {
		return fmt.Errorf("%s is not a type", c.StructName)
	}

	c.info.Struct, ok = decl.Type.(*ast.StructType)
	if !ok {
		return fmt.Errorf("%s is not a struct", c.StructName)
	}

	return nil
}

func (c *Collector) filterFields() (res []FieldInfo) {
	for _, f := range c.info.Struct.Fields.List {
		for _, ff := range f.Names {
			if len(c.Fields) > 0 {
				if !c.ContainsField(ff.Name) {
					continue
				}
			} else {
				firstRune, _ := utf8.DecodeRuneInString(ff.Name)
				if firstRune == utf8.RuneError {
					continue
				}
				if !c.PublicFields && unicode.IsUpper(firstRune) {
					continue
				}
				if !c.PrivateFields && unicode.IsLower(firstRune) {
					continue
				}
			}

			res = append(res, FieldInfo{
				Name: ff.Name,
				Expr: f.Type,
			})
		}
	}

	return res
}

func getFileName(scopes map[ast.Node]*go_types.Scope, file *ast.File) string {
	scope, ok := scopes[file]
	if ok {
		return reflect.ValueOf(scope).Elem().FieldByName("comment").String()
	}

	return ""
}
