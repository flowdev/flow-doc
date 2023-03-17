package draw2

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
	FlowModeSVGLinks
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
	//calcVerticalValues(y0, minLine int, mode FlowMode)
	respectMaxWidth(maxWidth, num int) ([]StartComp, int)
}

type EndComp interface {
	prevArrow() *Arrow
	addInput(*Arrow)
	calcHorizontalValues(x0 int)
	//calcVerticalValues(y0, minLine int, mode FlowMode)
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

// Flow contains data for a whole flow.
// The data is organized in rows and individual shapes per row.
// Valid shapes are Arrow, Comp, Split, Merge, Sequel and Loop.
//
// The following rules apply:
// - Arrows and Comps alternate.
// - Instead of a single Arrow a Split can be used for multiple Arrows.
//   So the first element of such a split is always an Arrow.
//   (exception: a component that completes a merge).
//   Such a split has to be the last element of a row.
//   A split can never be the first element of a row.
// - The last Comp (and element) in a row can instead be a Merge.
//   The real Comp for the merge has to be the first element of a future row
//   (possibly of the outer Split).
//   Of course multiple merges can "point" to the same Comp (using the same ID).
//   The same Merge instance has to be used for this (only 1 instance per ID).
// - The real Comp of a merge can be followed by an Arrow or Split as usual.
// - The last Comp (and element) in a row can be replaced by a Loop, too.
//   The loop points back to a component we can't draw an arrow to.
//   In the diagram you will see: ...back to: <component>:<port>
// - The first and last Comp and element of a row can instead be an OuterPort.
// - The last Comp (and element) in a row can also be replaced by a Sequel.
//   The other part of the Sequel should be at the start of one of the next rows
//   of the outer Split.
//   Sequels are in general inserted by the layout algorithm itself.
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

func (flow *Flow) AddCluster(cl *ShapeCluster) *Flow {
	flow.clusters = append(flow.clusters, cl)
	return flow
}

func (flow *Flow) calcHorizontalValues() {
	for _, cl := range flow.clusters {
		cl.calcHorizontalValues()
	}
}

// drawData contains all data needed for positioning the element correctly.
type drawData struct {
	x0, y0         int
	height, width  int
	minLine, lines int
}

func withinShape(line int, d *drawData) bool {
	return d.minLine <= line && line < d.minLine+d.lines
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
