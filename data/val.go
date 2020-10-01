package data

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/flowdev/ea-flow-doc/x/reflect"
)

// Value returns the string representation of the flow literal value
// corresponding to `expr`.
// The result is a basic Go literal or one of the standard types:
// "bool", "byte", "complex64", "complex128", "float32", "float64",
// "int", "int8", "int16", "int32", "int64",
// "rune", "string", "uint", "uint8", "uint16", "uint32", "uint64",
// "uintptr"
// or it is empty and an error is returned.
func Value(expr ast.Expr) (string, error) {
	sb := &strings.Builder{}
	stopValue := valueOfExpr(sb, expr)

	if stopValue != "" {
		return "", errors.New("stopped evaluating flow literal value at: " + stopValue)
	}
	return sb.String(), nil
}

func valueOfExpr(sb *strings.Builder, expr ast.Expr) (stopValue string) {
	if reflect.IsNilInterfaceOrPointer(expr) {
		return stopValue
	}

	switch e := expr.(type) {
	case *ast.BasicLit:
		sb.WriteString(e.Value)
	case nil:
		sb.WriteString("NULL") // should be very rare
	default:
		stopValue = fmt.Sprintf("%T", e)
	}
	return stopValue
}
