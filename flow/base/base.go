package base

import (
	"go/token"
)

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
