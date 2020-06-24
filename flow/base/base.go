package base

import (
	"go/ast"
	"go/token"
	"go/types"
)

// PortPrefix is the prefix of all output ports.
const PortPrefix = "port"

// Port is a full description of a flow port.
type Port struct {
	Name       string
	Pos        token.Pos
	IsImplicit bool
	IsError    bool
}

// DataTyp describes a data declaration with name and type.
type DataTyp struct {
	Name    string
	NamePos token.Pos
	Typ     string
	TypPos  token.Pos
}

// CallStep is a step in a flow that performs a call to a sub-flow or
// component.
type CallStep struct {
	Inputs        []string
	InPort        Port
	ComponentName string
	Outputs       []string
}

// ReturnStep is a step in a flow that ends the flow and sends data to an
// output port.
type ReturnStep struct {
	Datas   []string
	OutPort Port
}

// Step is a step in a flow. It can be one of: CallStep, ReturnStep or Branch
type Step interface {
}

// Branch is a control flow branch. It can be either the main branch of a flow
// or sub-branch created by an if statement.
type Branch struct {
	DataMap map[string]string
	Steps   []Step
	Parent  *Branch
}

// FlowData describes a flow.
// A flow has got a mainBranch an possibly many sub-branches.
// The main branch always starts with the first call expression of the function.
// Sub-branches are created with if expressions.
// Consequently a flow can't start with an if expression!
type FlowData struct {
	InPort        Port
	Inputs        []DataTyp
	ComponentName string
	OutPorts      []Port
	MainBranch    *Branch
}

// NewBranch creates a new branch with the given parent.
// The parent of the main branch of a flow is nil.
func NewBranch(parent *Branch) *Branch {
	return &Branch{
		DataMap: make(map[string]string, 64),
		Steps:   make([]Step, 0, 64),
		Parent:  parent,
	}
}

// NewFlowData creates a flow data structure with a main branch.
func NewFlowData() *FlowData {
	return &FlowData{MainBranch: NewBranch(nil)}
}

// AddDatasToMap adds the given data types to the map.
// If a name is already registered in map, the longer type is kept.
func AddDatasToMap(m map[string]string, datas []DataTyp) map[string]string {
	for _, dat := range datas {
		if dat.Name != "" {
			m = AddDataToMap(dat.Name, dat.Typ, m)
		}
	}
	return m
}

// AddDataToMap adds the given name and type to the map.
// If name is already registered in map, the longer type is stored.
func AddDataToMap(name, typ string, m map[string]string) map[string]string {
	t, ok := m[name]
	if !ok || len(typ) > len(t) {
		m[name] = typ
	}
	return m
}

// IsBoring returns true for builtin types.
func IsBoring(typ string) bool {
	switch typ { // simple builtin types are 'boring'
	case "bool", "byte", "complex64", "complex128", "float32", "float64",
		"int", "int8", "int16", "int32", "int64", "rune", "string",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "":
		return true
	default:
		return false
	}
}

// TypeInfo returns the Go type definition for the given type expression.
func TypeInfo(typ ast.Expr, typesInfo *types.Info) string {
	if typesInfo.Types[typ].Type == nil {
		return "<types.Info not filled properly>"
	}
	return typesInfo.Types[typ].Type.String()
}
