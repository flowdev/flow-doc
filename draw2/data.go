package draw2

import "fmt"

const (
	RowGap     = 8
	LineHeight = 24
	CharWidth  = 8
	ParenWidth = 6
	WordGap    = 6
	TextOffset = 6
	BreakText  = "…"
	BreakWidth = 2 * ParenWidth
	LoopText   = "back to: "
	LoopWidth  = 6*CharWidth + 3*CharWidth/2
)

type FlowMode int

const (
	FlowModeNoLinks FlowMode = iota
	FlowModeMDLinks
	FlowModeSVGLinks // not implemented yet
)

// Shapes:
// Name     | Start | Mid | End
// ============================
// Comp     | Yes   | Yes | Yes
// StartPort| Yes   | No  | No
// BreakEnd | Yes   | No  | No
// BreakStart No    | No  | Yes
// EndPort  | No    | No  | Yes
// Loop     | No    | No  | Yes
// ----------------------------
// Arrow    | No    | Yes | No

type StartComp interface {
	nextArrow() *Arrow
	calcHorizontalValues(x0 int)
	calcVerticalValues(y0, minLine int, mode FlowMode)
	respectMaxWidth(maxWidth, num int) ([]StartComp, int)
}

type EndComp interface {
	prevArrow() *Arrow
	addInput(*Arrow)
	calcHorizontalValues(x0 int)
	calcVerticalValues(y0, minLine int, mode FlowMode)
	respectMaxWidth(maxWidth, num int) ([]StartComp, int)
	minRestOfRowWidth(num int) int
}

type CompRegistry interface {
	register(*Comp)
	lookup(id string) *Comp
}

type ShapeCluster struct {
	shapeRows    []StartComp
	compRegistry map[string]*Comp
}

func NewCluster() *ShapeCluster {
	return &ShapeCluster{
		shapeRows:    make([]StartComp, 0, 32),
		compRegistry: make(map[string]*Comp, 128),
	}
}

func (cl *ShapeCluster) AddStartComp(comp StartComp) *ShapeCluster {
	cl.shapeRows = append(cl.shapeRows, comp)
	return cl
}

func (cl *ShapeCluster) lookup(id string) *Comp {
	return cl.compRegistry[id]
}

func (cl *ShapeCluster) register(comp *Comp) {
	cl.compRegistry[comp.ID()] = comp
}

func (cl *ShapeCluster) calcHorizontalValues() {
	for _, comp := range cl.shapeRows {
		comp.calcHorizontalValues(0)
	}
}

func (cl *ShapeCluster) respectMaxWidth(maxWidth, num int) (newNum int) {
	var newRows []StartComp
	addRows := make([]StartComp, 0, 64)

	for _, comp := range cl.shapeRows {
		newRows, num = comp.respectMaxWidth(maxWidth, num)
		addRows = append(addRows, newRows...)
	}

	cl.shapeRows = append(cl.shapeRows, addRows...)
	return num
}

func (cl *ShapeCluster) calcVerticalValues(mode FlowMode) {
	for _, comp := range cl.shapeRows {
		comp.calcVerticalValues(0, 0, mode)
	}
}

type Flow struct {
	name     string
	mode     FlowMode
	width    int
	dark     bool
	clusters []*ShapeCluster
}

func NewFlow(name string, mode FlowMode, width int, dark bool) *Flow {
	return &Flow{
		name:     name,
		mode:     mode,
		width:    width,
		dark:     dark,
		clusters: make([]*ShapeCluster, 0, 64),
	}
}

func (flow *Flow) ChangeConfig(name string, mode FlowMode, width int, dark bool) {
	flow.name = name
	flow.mode = mode
	flow.width = width
	flow.dark = dark
}

func (flow *Flow) AddCluster(cl *ShapeCluster) *Flow {
	flow.clusters = append(flow.clusters, cl)
	return flow
}

func (flow *Flow) validate() error {
	if len(flow.clusters) == 0 {
		return fmt.Errorf("no shape clusters found in flow")
	}

	for i, cl := range flow.clusters {
		if len(cl.shapeRows) == 0 {
			return fmt.Errorf("no shapes found in the %d-th cluster of the flow", i+1)
		}
	}

	return nil
}

func (flow *Flow) calcHorizontalValues() {
	for _, cl := range flow.clusters {
		cl.calcHorizontalValues()
	}
}

func (flow *Flow) respectMaxWidth(maxWidth int) {
	num := 0
	for _, cl := range flow.clusters {
		num = cl.respectMaxWidth(maxWidth, num)
	}
}

func (flow *Flow) calcVerticalValues(mode FlowMode) {
	for _, cl := range flow.clusters {
		cl.calcVerticalValues(mode)
	}
}

// drawData contains all data needed for positioning the element correctly.
type drawData struct {
	x0, y0         int
	height, width  int
	minLine, lines int
}
