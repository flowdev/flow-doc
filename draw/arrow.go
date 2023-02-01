package draw

const (
	arrTipWidth        = 12
	arrTipHeight       = 8
	arrSmallTextOffset = 4  // for small, low text
	arrTextOffset      = 11 // for elevated text
)

// Arrow contains all information for displaying an Arrow including data type
// and ports.
type Arrow struct {
	DataTypes      []*DataType
	SrcPort        string
	DstPort        string
	drawData       *drawData
	dataTypesWidth int         // for centering the data types
	splitState     *splitState // for spliting rows according to width
}

func (*Arrow) breakable() bool {
	return true
}

func (*Arrow) compish() bool {
	return false
}

func (arrow *Arrow) intersects(line int) bool {
	return withinShape(line, arrow.drawData)
}

// --------------------------------------------------------------------------
// Calculate width, height and lines
// --------------------------------------------------------------------------
func (arr *Arrow) calcDimensions() {
	width := 0
	height := 0
	lines := 0
	w1 := 1
	for _, dt := range arr.DataTypes {
		calcDataTypeDimensions(dt)
		dtd := dt.drawData
		w1 = max(w1, dt.w1)
		height += dtd.height
		lines += dtd.lines
	}
	for _, dt := range arr.DataTypes { // now we know the real x1, set it everywhere
		dt.drawData.width += w1 - dt.w1
		dt.w1 = w1
		width = max(width, dt.drawData.width)
	}
	for _, dt := range arr.DataTypes { // now we know the real width, set it everywhere
		dt.drawData.width = width
	}
	height += LineHeight // for the arrow itself and optional ports
	lines += 1
	arr.dataTypesWidth = width
	width = arrowWidth(arr, width)

	arr.drawData = &drawData{
		width:  width,
		height: height,
		lines:  lines,
	}
}

func calcDataTypeDimensions(dt *DataType) {
	dt.w1 = ParenWidth + (1+len(dt.Name)+1)*CharWidth
	width := dt.w1 + (len(dt.Type)+1)*CharWidth + ParenWidth
	dt.drawData = &drawData{
		width:  width,
		height: LineHeight,
		lines:  1,
	}
}

func arrowWidth(arr *Arrow, dataWidth int) int {
	portWidth := len(arr.SrcPort)*CharWidth + len(arr.DstPort)*CharWidth

	if portWidth != 0 {
		portWidth += WordGap + // so the port text isn't glued to the comp
			2*CharWidth // so the ports aren't glued together and it is ...
		// ... clear which type a single port is
	}

	return max(portWidth, dataWidth) + arrTipWidth
}

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func (arr *Arrow) enrich(x0, y0, minLine, level int, outerComp *drawData,
	lastArr *Arrow, global *enrichData,
) (newShapeLines [][]Shape) {
	y := y0
	for _, dt := range arr.DataTypes {
		dtd := dt.drawData
		dtd.x0 = x0
		dtd.y0 = y
		y += dtd.height
	}

	ad := arr.drawData
	ad.x0 = x0
	ad.y0 = y0
	ad.minLine = minLine

	return nil
}

// --------------------------------------------------------------------------
// Convert To SVG and MD
// --------------------------------------------------------------------------
func (arrow *Arrow) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	var svg *svgFlow
	var link *svgLink
	ad := arrow.drawData
	maxLine := maximumLine(ad)

	// get or create correct SVG flow:
	if mode == FlowModeSVGLinks {
		x0, y0, height, width := svgDimensionsForLine(line, arrow, ad, maxLine)
		svg, link = addNewSVGFlow(smf,
			x0, y0, height, width,
			"arrow", line,
		)
	} else {
		svg = smf.svgs[""]
	}

	if line == maxLine { // draw arrow line
		srcPortToSVG(svg, arrow, ad)
		dstPortToSVG(svg, arrow, ad)

		arrToSVG(svg, ad)

		smf.lastX += ad.width
		return
	}

	dataWidth := ad.width - arrTipWidth
	lastIdx := len(arrow.DataTypes) - 1
	idx := line - ad.minLine
	dt := arrow.DataTypes[idx]

	arrowDataTypeToSVG(svg, link, dt, ad.x0, dataWidth, arrow.dataTypesWidth,
		idx == 0, idx == lastIdx)

	smf.lastX += ad.width
}

func svgDimensionsForLine(line int, arrow *Arrow, ad *drawData, maxLine int,
) (x0, y0, height, width int) {

	if line == maxLine {
		return ad.x0, ad.y0 + ad.height - LineHeight, LineHeight, ad.width
	}

	idx := line - ad.minLine
	dt := arrow.DataTypes[idx]
	dtd := dt.drawData
	return ad.x0, dtd.y0, dtd.height, ad.width
}

func srcPortToSVG(svg *svgFlow, arrow *Arrow, ad *drawData) {
	if arrow.SrcPort != "" {
		svg.Texts = append(svg.Texts, &svgText{
			X:     ad.x0 + WordGap,
			Y:     ad.y0 + ad.height - arrSmallTextOffset,
			Width: len(arrow.SrcPort) * CharWidth,
			Text:  arrow.SrcPort,
			Small: true,
		})
	}
}

func dstPortToSVG(svg *svgFlow, arrow *Arrow, ad *drawData) {
	if arrow.DstPort != "" {
		w := len(arrow.DstPort) * CharWidth
		svg.Texts = append(svg.Texts, &svgText{
			X:     ad.x0 + ad.width - w - arrTipWidth,
			Y:     ad.y0 + ad.height - arrSmallTextOffset,
			Width: w,
			Text:  arrow.DstPort,
			Small: true,
		})
	}
}

func arrToSVG(svg *svgFlow, ad *drawData) {

	arrY := ad.y0 + ad.height - LineHeight + arrTipHeight
	svg.Arrows = append(svg.Arrows, &svgArrow{
		X1:    ad.x0,
		Y1:    arrY,
		X2:    ad.x0 + ad.width,
		Y2:    arrY,
		XTip1: ad.x0 + ad.width - arrTipHeight,
		YTip1: arrY - arrTipHeight,
		XTip2: ad.x0 + ad.width - arrTipHeight,
		YTip2: arrY + arrTipHeight,
	})
}

func arrowDataTypeToSVG(
	svg *svgFlow, link *svgLink, dt *DataType,
	x0, width, dataTypesWidth int,
	first, last bool,
) {
	dtd := dt.drawData
	padding := (width - dataTypesWidth) / 2
	x1 := x0 + padding + dt.w1
	padding += CharWidth + ParenWidth
	y := dtd.y0 + LineHeight - TextOffset

	if first { // opening parenthesis
		svg.Texts = append(svg.Texts, &svgText{
			X:      x0 + padding - ParenWidth,
			Y:      y,
			Width:  ParenWidth,
			Text:   "(",
			Link:   !dt.GoLink && dt.Link != "",
			GoLink: dt.GoLink,
		})
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:      x0 + padding,
		Y:      y,
		Width:  len(dt.Name) * CharWidth,
		Text:   dt.Name,
		Link:   !dt.GoLink && dt.Link != "",
		GoLink: dt.GoLink,
	})

	typText := dt.Type
	typWidth := len(dt.Type) * CharWidth
	if last {
		typText += ")"
		typWidth += ParenWidth
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:      x1,
		Y:      y,
		Width:  typWidth,
		Text:   typText,
		Link:   !dt.GoLink && dt.Link != "",
		GoLink: dt.GoLink,
	})

	if link != nil {
		link.Link = dt.Link
	}
}
