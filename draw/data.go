package draw

const (
	LineGap    = 8
	LineHeight = 24
	CharWidth  = 8
	ParenWidth = 6
	WordGap    = 6
	TextOffset = 6
)

type FlowMode int

const (
	FlowModeNoLinks FlowMode = iota
	FlowModeMDLinks
	FlowModeSVGLinks
)

// drawData contains all data needed for positioning the element correctly.
type drawData struct {
	x0, y0         int
	height, width  int
	minLine, lines int
}

// DataType contains the optional name of the data and its type.
// Plus an optional link to the definition of the type.
type DataType struct {
	Name     string
	Type     string
	Link     string
	drawData *drawData
	x1       int // for aligning the data types of arrows
}

// Arrow contains all information for displaying an Arrow including data type
// and ports.
type Arrow struct {
	DataTypes      []*DataType
	HasSrcOp       bool
	SrcPort        string
	HasDstOp       bool
	DstPort        string
	drawData       *drawData
	dataTypesWidth int // // for centering the data types
}

// Text is the text to display as a flow output port.
type Text struct {
	Text     string
	drawData *drawData
}

// PluginType contains the type of the plugin.
// And optionally a link to its definition.
type PluginType struct {
	Type     string
	Link     string
	drawData *drawData
}

// Plugin is a helper operation that is used inside a proper operation.
type Plugin struct {
	Title    string
	Types    []*PluginType
	drawData *drawData
}

// Op holds all data to describe a single operation including possible plugins.
type Op struct {
	Main     *DataType
	Plugins  []*Plugin
	drawData *drawData
}

// Split contains data for multiple paths/arrows originating from a single Op.
type Split struct {
	Shapes   [][]any
	drawData *drawData
}

// Merge holds data for merging multiple paths/arrows into a single Op.
type Merge struct {
	ID       string
	Size     int
	drawData *drawData
	arrows   []*drawData
}

// Flow contains data for a whole flow.
// The data is organized in rows and individual shapes per row.
// Valid shapes are Arrow, Op, Split, and Merge.
//
// The following rules apply:
// - Arrows and Ops alternate.
// - Instead of a single Arrow a Split can be used for multiple Arrows.
//   So the first element of such a split is always an Arrow.
// - The last Op in a row can instead be a Merge.
//   The real Op for the merge has to be the first element of a future row.
//   Of course multiple merges can "point" to the same Op (using the same ID).
//   The same Merge instance has to be used for this (only 1 instance per ID).
// - The real Op of a merge can be followed by an Arrow or Split as usual.
type Flow struct {
	Mode      FlowMode
	Name      string
	AllShapes *Split
}
