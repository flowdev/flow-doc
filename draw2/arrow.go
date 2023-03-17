package draw2

import "fmt"

const (
	arrTipWidth        = 12
	arrTipHeight       = 8
	arrSmallTextOffset = 4  // for small, low text
	arrTextOffset      = 11 // for elevated text
)

// DataType contains the optional name of the data and its type.
// Plus an optional link to the definition of the type.
type DataType struct {
	name     string
	typ      string
	link     string
	drawData *drawData
	w1       int // for aligning the data types of arrows
}

// Arrow contains all information for displaying an Arrow including data type
// and ports.
type Arrow struct {
	dataTypes      []*DataType
	srcPort        string
	dstPort        string
	srcComp        StartComp
	dstComp        EndComp
	drawData       *drawData
	dataTypesWidth int // for centering the data types
}

func NewArrow(srcPort, dstPort string) *Arrow {
	return &Arrow{
		srcPort: srcPort,
		dstPort: dstPort,
	}
}

func (arr *Arrow) AddDestination(comp EndComp) *Arrow {
	arr.dstComp = comp
	comp.addInput(arr)
	return arr
}

func (arr *Arrow) LinkComp(id string, compRegistry CompRegistry) error {
	if compRegistry == nil {
		return fmt.Errorf("no registry given for linking to component with ID: %q", id)
	}
	comp := compRegistry.lookup(id)
	if comp == nil {
		return fmt.Errorf("unable to link to component with ID: %q", id)
	}
	arr.dstComp = comp
	return nil
}

func (arr *Arrow) MustLinkComp(id string, compRegistry CompRegistry) *Arrow {
	if err := arr.LinkComp(id, compRegistry); err != nil {
		panic(err)
	}
	return arr
}

func (arr *Arrow) AddDataType(name, typ, link string) *Arrow {
	arr.dataTypes = append(arr.dataTypes, &DataType{
		name: name,
		typ:  typ,
		link: link,
	})
	return arr
}

func (arrow *Arrow) intersects(line int) bool {
	return withinShape(line, arrow.drawData)
}

// --------------------------------------------------------------------------
// Calculate horizontal values of shapes (x0 and width)
// --------------------------------------------------------------------------
func (arr *Arrow) calcHorizontalValues(x0 int) {
	for _, dt := range arr.dataTypes {
		dt.drawData = &drawData{
			x0: x0,
		}
	}

	width := arr.calcWidth()

	arr.drawData = &drawData{
		x0:    x0,
		width: width,
	}

	arr.dstComp.calcHorizontalValues(arr.drawData.x0 + arr.drawData.width)
}

func (arr *Arrow) calcWidth() int {
	arr.calcDataTypesWidth()

	portWidth := len(arr.srcPort)*CharWidth + len(arr.dstPort)*CharWidth

	if portWidth != 0 {
		portWidth += WordGap + // so the port text isn't glued to the comp
			2*CharWidth // so the ports aren't glued together and it is ...
		// ... clear which type a single port is
	}

	return max(portWidth, arr.dataTypesWidth) + arrTipWidth
}

func (arr *Arrow) calcDataTypesWidth() {
	if arr.dataTypesWidth > 0 {
		return
	}

	w1 := 1
	for _, dt := range arr.dataTypes {
		calcDataTypeWidth(dt)
		w1 = max(w1, dt.w1)
	}
	width := w1
	for _, dt := range arr.dataTypes { // now we know the real w1, set it everywhere
		dt.drawData.width += w1 - dt.w1
		dt.w1 = w1
		width = max(width, dt.drawData.width)
	}
	for _, dt := range arr.dataTypes { // now we know the real width, set it everywhere
		dt.drawData.width = width
	}
	arr.dataTypesWidth = width
}

func calcDataTypeWidth(dt *DataType) {
	dt.w1 = ParenWidth + (1+len(dt.name)+1)*CharWidth
	dt.drawData.width = dt.w1 + (len(dt.typ)+1)*CharWidth + ParenWidth
}

func (arr *Arrow) extendTo(xn int) {
	if arr.drawData == nil {
		return
	}

	arr.drawData.width = xn - arr.drawData.x0
}

// --------------------------------------------------------------------------
// Needed for dividing rows
// --------------------------------------------------------------------------

func (arr *Arrow) respectMaxWidth(maxWidth, num int) ([]StartComp, int) {
	longBroken, unBroken := arr.brokenWidths(num)
	if arr.drawData.x0+longBroken > maxWidth {
		newArr := arr.breakShort()

		arr.dstComp = NewBreakStart(num)
		arr.dstComp.addInput(arr)
		arr.calcHorizontalValues(arr.drawData.x0)

		newStart := NewBreakEnd(num)
		newStart.AddOutput(newArr)
		newStart.calcHorizontalValues(0)

		return []StartComp{newStart}, num + 1
	}

	if arr.drawData.x0+unBroken > maxWidth {
		newArr := arr.breakLong()

		arr.dstComp = NewBreakStart(num)
		arr.dstComp.addInput(arr)
		arr.calcHorizontalValues(arr.drawData.x0)

		newStart := NewBreakEnd(num)
		newStart.AddOutput(newArr)
		newStart.calcHorizontalValues(0)

		return []StartComp{newStart}, num + 1
	}

	return arr.dstComp.respectMaxWidth(maxWidth, num)
}

func (arr *Arrow) minRestOfRowWidth(num int) int {
	if arr == nil {
		return 0 // prevent endless loop
	}
	if arr.srcPort == "" && len(arr.dataTypes) == 0 { // already broken
		return arr.drawData.width + arr.dstComp.minRestOfRowWidth(num)
	}

	width := 0
	if len(arr.srcPort) > 0 {
		width += WordGap + // so the port text isn't glued to the comp
			len(arr.srcPort)*CharWidth +
			2*CharWidth // so it is clear which type a single port is
	}

	return width + arrTipWidth + breakWidth(num)
}

func (arr *Arrow) breakable() bool {
	return arr.srcPort != "" || len(arr.dataTypes) > 0
}

// brokenWidths returns -1, -1 if the arrow isn't breakable.
func (arr *Arrow) brokenWidths(num int) (longBroken, unBroken int) {
	if !arr.breakable() {
		return -1, -1
	}

	longBroken = max(arr.dataTypesWidth+arrTipWidth+breakWidth(num), arr.minRestOfRowWidth(num))
	unBroken = arr.drawData.width + arr.dstComp.minRestOfRowWidth(num)

	return longBroken, unBroken
}

// breakShort breaks this arrow into two arrows.
// The second arrow is returned and it's horizontal values haven't been
// calculated yet.
func (arr *Arrow) breakShort() *Arrow {
	seqArr := &Arrow{
		dataTypes: arr.dataTypes,
		dstPort:   arr.dstPort,
		dstComp:   arr.dstComp,
	}

	arr.dataTypes = nil
	arr.dstPort = ""
	arr.dstComp = nil
	arr.dataTypesWidth = 0
	return seqArr
}

// breakLong breaks this arrow into two arrows.
// The second arrow is returned and it's horizontal values haven't been
// calculated yet.
func (arr *Arrow) breakLong() *Arrow {
	seqArr := &Arrow{
		dstPort: arr.dstPort,
		dstComp: arr.dstComp,
	}

	arr.dstPort = ""
	arr.dstComp = nil
	return seqArr
}

// --------------------------------------------------------------------------
// Calculate vertical values of shapes (y0, height, lines and minLine)
// --------------------------------------------------------------------------
func (arr *Arrow) calcVerticalValues(y0, minLine int, mode FlowMode) {

	y := y0
	height := 0
	lines := 0
	for _, dt := range arr.dataTypes {
		calcDataTypeVerticals(dt, y)
		dtd := dt.drawData
		y += dtd.height
		height += dtd.height
		lines += dtd.lines
	}

	height += LineHeight // for the arrow itself and optional ports
	lines += 1

	ad := arr.drawData
	ad.y0 = y0
	ad.height = height
	ad.minLine = minLine
	ad.lines = lines
}

func calcDataTypeVerticals(dt *DataType, y0 int) {
	dt.drawData = &drawData{
		y0:     y0,
		height: LineHeight,
		lines:  1,
	}
}
