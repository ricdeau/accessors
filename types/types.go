package types

import (
	"go/ast"
	"path"
	"strings"
)

// TypeName go type as string.
// I.e. string -> "string", *os.File -> "*os.File"
func TypeName(importSpecs []*ast.ImportSpec, imports map[string]string, typ ast.Expr) string {
	switch v := typ.(type) {
	case *ast.Ident:
		return indentName(v)
	case *ast.StarExpr:
		if s := TypeName(importSpecs, imports, v.X); s != "" {
			return "*" + s
		}
	case *ast.SelectorExpr:
		alias := indentName(v.X.(*ast.Ident))
		importStr, ok := findImport(importSpecs, alias)
		if ok {
			imports[alias] = importStr
		}
		if s := TypeName(importSpecs, imports, v.Sel); s != "" {
			return alias + "." + s
		}
	case *ast.SliceExpr:
		if s := TypeName(importSpecs, imports, v.X); s != "" {
			return "[]" + s
		}
	case *ast.ArrayType:
		if s := TypeName(importSpecs, imports, v.Elt); s != "" {
			l := ""
			if v.Len != nil {
				l = v.Len.(*ast.BasicLit).Value
			}
			return "[" + l + "]" + s
		}
	case *ast.MapType:
		if key := TypeName(importSpecs, imports, v.Key); key != "" {
			if val := TypeName(importSpecs, imports, v.Value); val != "" {
				return "map[" + key + "]" + val
			}
		}
	case *ast.ChanType:
		if s := TypeName(importSpecs, imports, v.Value); s != "" {
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

// TypeClearName go type name without package prefix.
func TypeClearName(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return TypeClearName(v.Sel)
	case *ast.StarExpr:
		return TypeClearName(v.X)
	}

	return ""
}

func indentName(indent *ast.Ident) string {
	if indent == nil {
		return ""
	}

	return indent.Name
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
