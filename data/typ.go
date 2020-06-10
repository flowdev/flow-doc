package data

import (
	"errors"
	"fmt"
	"go/ast"
	"log"
	"strings"

	"github.com/flowdev/ea-flow-doc/x/reflect"
)

// Type returns the string representation of the flow data type corresponding
// to `expr`.
// The result may contain `list(Type)` or `map(KeyType, ValueType)`.
func Type(expr ast.Expr) (string, error) {
	sb := &strings.Builder{}
	stopExprType := typeOfExpr(sb, expr)

	if stopExprType != "" {
		return sb.String(), errors.New("stopped evaluating flow data type at: " + stopExprType)
	}
	return sb.String(), nil
}

func typeOfExpr(sb *strings.Builder, expr ast.Expr) (stopExprType string) {
	if reflect.IsNilInterfaceOrPointer(expr) {
		return stopExprType
	}

	switch e := expr.(type) {
	case *ast.Ident:
		stopExprType = typeOfIdent(sb, e)
	case *ast.SelectorExpr:
		stopExprType = typeOfSelectorExpr(sb, e)
	case *ast.ArrayType:
		stopExprType = typeOfArrayType(sb, e)
	case *ast.Ellipsis:
		stopExprType = typeOfEllipsis(sb, e)
	case *ast.MapType:
		stopExprType = typeOfMapType(sb, e)
	case *ast.StarExpr:
		stopExprType = typeOfExpr(sb, e.X)
	case *ast.ParenExpr:
		stopExprType = typeOfExpr(sb, e.X)
	case nil:
		sb.WriteString("NULL") // should be very rare
	default:
		stopExprType = fmt.Sprintf("%T", e)
		log.Print("WARNING - unknown flow data type: " + stopExprType)
	}
	return stopExprType
}

func typeOfIdent(sb *strings.Builder, id *ast.Ident) (stopExprType string) {
	sb.WriteString(id.Name)
	return stopExprType
}

func typeOfSelectorExpr(sb *strings.Builder, sel *ast.SelectorExpr) (stopExprType string) {
	stopExprType = typeOfExpr(sb, sel.X)
	sb.WriteString(".")
	if sel.Sel != nil {
		sb.WriteString(sel.Sel.Name)
	}
	return stopExprType
}

func typeOfArrayType(sb *strings.Builder, arr *ast.ArrayType) (stopExprType string) {
	sb.WriteString("list(")
	stopExprType = typeOfExpr(sb, arr.Elt)
	sb.WriteString(")")
	return stopExprType
}

func typeOfEllipsis(sb *strings.Builder, elli *ast.Ellipsis) (stopExprType string) {
	sb.WriteString("list(")
	stopExprType = typeOfExpr(sb, elli.Elt)
	sb.WriteString(")")
	return stopExprType
}

func typeOfMapType(sb *strings.Builder, m *ast.MapType) (stopExprType string) {
	sb.WriteString("map(")
	stopExprType = typeOfExpr(sb, m.Key)
	sb.WriteString(", ")
	if stopExprType == "" {
		stopExprType = typeOfExpr(sb, m.Value)
	}
	sb.WriteString(")")
	return stopExprType
}
