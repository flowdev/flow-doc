package draw

const (
	LineGap     = 8
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

// Split contains data for multiple paths/arrows originating from a single Comp.
type Split struct {
	Shapes   [][]any
	drawData *drawData
}

// Arrow contains all information for displaying an Arrow including data type
// and ports.
type Arrow struct {
	DataTypes      []*DataType
	HasSrcComp     bool
	SrcPort        string
	HasDstComp     bool
	DstPort        string
	drawData       *drawData
	dataTypesWidth int         // for centering the data types
	splitState     *splitState // for spliting rows according to width
}

// Comp holds all data to describe a single component including possible plugins.
type Comp struct {
	Main     *DataType
	Plugins  []*PluginGroup
	drawData *drawData
}

// DataType contains the optional name of the data and its type.
// Plus an optional link to the definition of the type.
type DataType struct {
	Name     string
	Type     string
	Link     string
	GoLink   bool
	drawData *drawData
	x1       int // for aligning the data types of arrows
}

// PluginGroup is a helper component that is used inside a proper component.
type PluginGroup struct {
	Title    string
	Types    []*Plugin
	drawData *drawData
}

// Plugin contains the type of the plugin and optionally a link to its definition.
type Plugin struct {
	Type     string
	Link     string
	GoLink   bool
	drawData *drawData
}

// Merge holds data for merging multiple paths/arrows into a single Comp.
type Merge struct {
	ID       string
	Size     int
	drawData *drawData
	arrows   []*Arrow
}

type Sequel struct {
	Number   int
	drawData *drawData
}

type Loop struct {
	Name     string
	Port     string
	Link     string
	GoLink   bool
	drawData *drawData
}

// drawData contains all data needed for positioning the element correctly.
type drawData struct {
	x0, y0         int
	height, width  int
	minLine, lines int
}

type splitState struct {
	lastComp                        *drawData
	lastArr                         *Arrow
	x, y, line, xmax, ymax, maxLine int
	i, j                            int
	row                             []any
}
