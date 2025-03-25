package draw

import (
	"fmt"
)

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
	w1       int // for aligning the data types of arrows
	drawData *drawData
}

// Arrow contains all information for displaying an Arrow including data type
// and ports.
type Arrow struct {
	withDrawData
	dataTypes         []*DataType
	srcPort           string
	dstPort           string
	srcComp           StartComp
	dstComp           EndComp
	dataTypesWidth    int  // for centering the data types
	maxWidthRespected bool // remember that respectMaxWidth has been called already
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
	arr.AddDestination(comp)
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

// --------------------------------------------------------------------------
// Calculate horizontal values of shapes (x0 and width)
// --------------------------------------------------------------------------
func (arr *Arrow) calcHorizontalValues(x0 int) {
	if arr.drawData != nil && x0 == arr.drawData.x0 {
		return
	}
	for _, dt := range arr.dataTypes {
		dt.drawData = &drawData{
			x0: x0,
		}
	}
	width := arr.calcWidth()
	arr.drawData = newDrawData(x0, width)
	arr.dstComp.calcHorizontalValues(x0 + width)
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
	arr.drawData.width = max(arr.drawData.width, xn-arr.drawData.x0)
}

// --------------------------------------------------------------------------
// Needed for dividing rows
// --------------------------------------------------------------------------

func (arr *Arrow) respectMaxWidth(maxWidth, num int) (newStartComps []StartComp, newNum, newWidth int) {
	if arr.maxWidthRespected {
		return nil, num, arr.drawData.ymax()
	}
	arr.maxWidthRespected = true
	_, breakLong := arr.srcComp.(*BreakEnd) // never break short after a break, or we have an endless loop
	noBreak := false
	if breakLong && len(arr.dataTypes) == 0 && arr.srcPort == "" {
		noBreak = true
	}
	longBroken, unBroken := arr.brokenWidths(num)
	ad := arr.drawData
	x0 := ad.x0
	if !noBreak && x0+unBroken > maxWidth {
		var newArr *Arrow
		if !breakLong && x0+longBroken > maxWidth {
			newArr = arr.breakShort()
		} else {
			newArr = arr.breakLong()
		}

		brk := NewBreakStart(num)
		arr.AddDestination(brk)
		arr.calcHorizontalValues(x0)
		_, _, newWidth = brk.respectMaxWidth(maxWidth, num+1)

		newStart := brk.End()
		newStart.AddOutput(newArr)
		//newStart.resetDrawData()
		//newStart.calcHorizontalValues(0)
		//newStartComps2, newNum2, newWidth2 := newStart.respectMaxWidth(maxWidth, num+1)
		//
		//newStartComps = append(newStartComps, newStart)
		//return append(newStartComps, newStartComps2...), newNum2, max(newWidth, newWidth2)
		return append(newStartComps, newStart), num + 1, newWidth
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
// The second arrow is returned, and it's horizontal values haven't been
// calculated yet.
func (arr *Arrow) breakShort() *Arrow {
	newArr := &Arrow{
		dataTypes: arr.dataTypes,
		dstPort:   arr.dstPort,
	}

	arr.dstComp.switchInput(arr, newArr)
	arr.dataTypes = nil
	arr.dstPort = ""
	arr.dstComp = nil
	arr.dataTypesWidth = 0
	arr.drawData = nil
	return newArr
}

// breakLong breaks this arrow into two arrows.
// The second arrow is returned, and it's horizontal values haven't been
// calculated yet.
func (arr *Arrow) breakLong() *Arrow {
	newArr := &Arrow{
		dstPort: arr.dstPort,
	}

	arr.dstComp.switchInput(arr, newArr)
	arr.dstPort = ""
	arr.dstComp = nil
	arr.drawData = nil
	return newArr
}

func (arr *Arrow) resetDrawData() {
	arr.withDrawData.resetDrawData()
	arr.dstComp.resetDrawData()
}

// --------------------------------------------------------------------------
// Calculate vertical values of shapes (y0, height, lines and minLine)
// --------------------------------------------------------------------------
func (arr *Arrow) calcVerticalValues(y0, minLine int, mode FlowMode) (maxLines, newHeight int) {
	ad := arr.drawData
	if ad.height > 0 && y0 >= ad.y0 {
		return ad.maxLines(), ad.ymax()
	}

	y := y0
	height := 0
	lines := 0
	for _, dt := range arr.dataTypes {
		calcDataTypeVerticals(dt, y, minLine+lines)
		dtd := dt.drawData
		y += dtd.height
		height += dtd.height
		lines += dtd.lines
	}

	height += LineHeight // for the arrow itself and optional ports
	lines += 1

	ad.y0 = y0
	ad.height = height
	ad.minLine = minLine
	ad.lines = lines

	lines, height = arr.dstComp.calcVerticalValues(y0, minLine, mode)
	return max(lines, minLine+ad.lines), max(height, y0+ad.height)
}

func calcDataTypeVerticals(dt *DataType, y0, minLine int) {
	dd := dt.drawData
	dd.y0 = y0
	dd.minLine = minLine
	dd.height = LineHeight
	dd.lines = 1
}

// --------------------------------------------------------------------------
// Convert To SVG and MD
// --------------------------------------------------------------------------
func (arr *Arrow) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	if arr.drawData.drawLine(line) {
		arr.arrowToSVG(smf, line, mode)
		arr.drawData.drawnLines[line] = true
	}
	if line >= arr.drawData.minLine {
		arr.dstComp.toSVG(smf, line, mode)
	}
}
func (arr *Arrow) arrowToSVG(smf *svgMDFlow, line int, mode FlowMode) {
	var svg *svgFlow
	var link *svgLink
	ad := arr.drawData
	maxLine := ad.maxLines() - 1

	// get or create correct SVG flow:
	if mode == FlowModeMDLinks {
		x0, y0, height, width := svgDimensionsForLine(line, arr, ad, maxLine)
		svg, link = addNewSVGFlow(smf,
			x0, y0, height, width,
			"arrow", line,
		)
	} else {
		svg = smf.svgs[""]
	}

	if line == maxLine { // draw arrow line
		srcPortToSVG(svg, arr, ad)
		dstPortToSVG(svg, arr, ad)

		arrToSVG(svg, ad)

		smf.lastX += ad.width
		return
	} else if line < maxLine {
		dataWidth := ad.width - arrTipWidth
		lastIdx := len(arr.dataTypes) - 1
		idx := line - ad.minLine
		dt := arr.dataTypes[idx]

		arrowDataTypeToSVG(svg, link, dt, ad.x0, dataWidth, arr.dataTypesWidth,
			idx == 0, idx == lastIdx)
	}

	smf.lastX += ad.width
}

func svgDimensionsForLine(line int, arrow *Arrow, ad *drawData, maxLine int,
) (x0, y0, height, width int) {

	if line == maxLine {
		return ad.x0, ad.y0 + ad.height - LineHeight, LineHeight, ad.width
	}

	idx := line - ad.minLine
	dt := arrow.dataTypes[idx]
	dtd := dt.drawData
	return ad.x0, dtd.y0, dtd.height, ad.width
}

func srcPortToSVG(svg *svgFlow, arrow *Arrow, ad *drawData) {
	if arrow.srcPort != "" {
		svg.Texts = append(svg.Texts, &svgText{
			X:     ad.x0 + WordGap,
			Y:     ad.ymax() - arrSmallTextOffset,
			Width: len(arrow.srcPort) * CharWidth,
			Text:  arrow.srcPort,
			Small: true,
		})
	}
}

func dstPortToSVG(svg *svgFlow, arrow *Arrow, ad *drawData) {
	if arrow.dstPort != "" {
		w := len(arrow.dstPort) * CharWidth
		svg.Texts = append(svg.Texts, &svgText{
			X:     ad.x0 + ad.width - w - arrTipWidth,
			Y:     ad.ymax() - arrSmallTextOffset,
			Width: w,
			Text:  arrow.dstPort,
			Small: true,
		})
	}
}

func arrToSVG(svg *svgFlow, ad *drawData) {

	arrY := ad.ymax() - LineHeight + arrTipHeight
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
			X:     x0 + padding - ParenWidth,
			Y:     y,
			Width: ParenWidth,
			Text:  "(",
			Link:  dt.link != "",
		})
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:     x0 + padding,
		Y:     y,
		Width: len(dt.name) * CharWidth,
		Text:  dt.name,
		Link:  dt.link != "",
	})

	typText := dt.typ
	typWidth := len(dt.typ) * CharWidth
	if last {
		typText += ")"
		typWidth += ParenWidth
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:     x1,
		Y:     y,
		Width: typWidth,
		Text:  typText,
		Link:  dt.link != "",
	})

	if link != nil {
		link.Link = dt.link
	}
}
