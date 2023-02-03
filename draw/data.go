package draw

const (
	RowGap      = 8
	LineHeight  = 24
	CharWidth   = 8
	ParenWidth  = 6
	WordGap     = 6
	TextOffset  = 6
	SequelText  = "…"
	SequelWidth = 2 * ParenWidth
	LoopText    = "back to: "
	LoopWidth   = 6*CharWidth + 3*CharWidth/2
)

type FlowMode int

const (
	FlowModeNoLinks FlowMode = iota
	FlowModeMDLinks
	FlowModeSVGLinks
)

type Shape interface {
	breakable() bool // shape can be broken into parts by inserting a Sequel
	compish() bool   // shape behaves like a component rather than like an arrow
	calcDimensions() // calculate width, height and lines of simple shapes
	calcPosition(x0, y0, minLine int, outerComp *drawData, lastArr *Arrow, mode FlowMode, merges map[string]*Merge)
	toSVG(smf *svgMDFlow, line int, mode FlowMode)
	intersects(line int) bool // shape is visible on the given line
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
	Name      string
	AllShapes *Split
	mode      FlowMode
	width     int
	dark      bool
}

// drawData contains all data needed for positioning the element correctly.
type drawData struct {
	x0, y0         int
	height, width  int
	minLine, lines int
}

type enrichData struct {
	width                   int
	mode                    FlowMode
	merges                  map[string]*Merge // TODO: remove!
	saveState, currentState *splitState
}

type splitState struct {
	level                           int
	lastComp                        *drawData
	lastArr                         *Arrow
	x, y, line, xmax, ymax, maxLine int
	i, j                            int
	row                             []Shape
	cutData                         []*cutData
}

type cutData struct {
	arr   *Arrow
	i, j  int
	split *Split
	level int
}
