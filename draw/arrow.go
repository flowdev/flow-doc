package draw

const (
	arrTipWidth        = 12
	arrTipHeight       = 8
	arrSmallTextOffset = 4  // for small, low text
	arrTextOffst       = 11 // for elevated text
)

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func enrichArrow(arr *Arrow, x0, y0, minLine int) {
	height := 0
	width := 0
	lines := 0
	x1 := x0
	for _, dt := range arr.DataTypes {
		enrichDataType(dt, x0, y0+height, minLine+lines)
		dtd := dt.drawData
		x1 = max(x1, dt.x1)
		height += dtd.height
		width = max(width, dtd.width)
		lines += dtd.lines
	}
	for _, dt := range arr.DataTypes { // now we know the real x1, set it everywhere
		dt.drawData.width += x1 - dt.x1
		dt.x1 = x1
		width = max(width, dt.drawData.width)
	}
	lines += 1 // for the arrow itself and optional ports
	height += LineHeight
	arr.dataTypesWidth = width
	width = arrowWidth(arr, width)

	arr.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		width:   width,
		height:  height,
		minLine: minLine,
		lines:   lines,
	}
}

func arrowWidth(arr *Arrow, dataWidth int) int {
	portWidth := 0
	if arr.HasSrcOp {
		portWidth = len(arr.SrcPort) * CharWidth
	}
	if arr.HasDstOp {
		portWidth += len(arr.DstPort) * CharWidth
	}
	if portWidth != 0 {
		portWidth += WordGap + // so the port text isn't glued to the op
			2*CharWidth // so the ports aren't glued together and it is ...
		// ... clear which type a single port is
	}

	width := max(portWidth, dataWidth) + arrTipWidth

	if !arr.HasSrcOp && arr.SrcPort != "" {
		width += len(arr.SrcPort)*CharWidth + WordGap
	}
	if !arr.HasDstOp && arr.DstPort != "" {
		width += len(arr.DstPort)*CharWidth + WordGap
	}

	return width
}

func enrichDataType(dt *DataType, x0, y0, minLine int) {
	dt.x1 = x0 + ParenWidth + (1+len(dt.Name)+1)*CharWidth
	width := dt.x1 - x0 + (len(dt.Type)+1)*CharWidth + ParenWidth
	dt.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		width:   width,
		height:  LineHeight,
		minLine: minLine,
		lines:   1,
	}
}

// --------------------------------------------------------------------------
// Convert To SVG and MD
// --------------------------------------------------------------------------
func arrowToSVG(sfs map[string]*svgFlow, mdf *mdFlow, mode FlowMode, arrow *Arrow) {
	if mode == FlowModeSVGLinks {
		return
	}
	svg := sfs[""]
	ad := arrow.drawData

	preArr := preArrowToSVG(svg, arrow, ad)
	postArr := postArrowToSVG(svg, arrow, ad)

	srcPortToSVG(svg, arrow, ad, preArr)
	dstPortToSVG(svg, arrow, ad, postArr)

	arrX, arrWidth := arrToSVG(svg, ad, preArr, postArr)

	dataWidth := arrWidth - arrTipWidth
	lastI := len(arrow.DataTypes) - 1
	for i, dt := range arrow.DataTypes {
		arrowDataTypeToSVG(svg, dt, arrX, dataWidth, arrow.dataTypesWidth, i == 0, i == lastI)
	}
}

func preArrowToSVG(svg *svgFlow, arrow *Arrow, ad *drawData) int {
	if !arrow.HasSrcOp && arrow.SrcPort != "" {
		svg.Texts = append(svg.Texts, &svgText{
			X:     ad.x0,
			Y:     ad.y0 + ad.height - arrTextOffst,
			Width: len(arrow.SrcPort) * CharWidth,
			Text:  arrow.SrcPort,
		})
		return len(arrow.SrcPort)*CharWidth + WordGap
	}
	return 0
}

func postArrowToSVG(svg *svgFlow, arrow *Arrow, ad *drawData) int {
	if !arrow.HasDstOp && arrow.DstPort != "" {
		postArr := len(arrow.DstPort)*CharWidth + WordGap
		svg.Texts = append(svg.Texts, &svgText{
			X:     ad.x0 + ad.width - postArr + WordGap,
			Y:     ad.y0 + ad.height - arrTextOffst,
			Width: len(arrow.DstPort) * CharWidth,
			Text:  arrow.DstPort,
		})
		return postArr
	}
	return 0
}

func srcPortToSVG(svg *svgFlow, arrow *Arrow, ad *drawData, preArr int) {
	if arrow.HasSrcOp && arrow.SrcPort != "" {
		svg.Texts = append(svg.Texts, &svgText{
			X:     ad.x0 + preArr + WordGap,
			Y:     ad.y0 + ad.height - arrSmallTextOffset,
			Width: len(arrow.SrcPort) * CharWidth,
			Text:  arrow.SrcPort,
			Small: true,
		})
	}
}

func dstPortToSVG(svg *svgFlow, arrow *Arrow, ad *drawData, postArr int) {
	if arrow.HasDstOp && arrow.DstPort != "" {
		w := len(arrow.DstPort) * CharWidth
		svg.Texts = append(svg.Texts, &svgText{
			X:     ad.x0 + ad.width - w - arrTipWidth - postArr,
			Y:     ad.y0 + ad.height - arrSmallTextOffset,
			Width: w,
			Text:  arrow.DstPort,
			Small: true,
		})
	}
}

func arrToSVG(svg *svgFlow, ad *drawData, preArr, postArr int) (arrX, arrWidth int) {

	arrX = ad.x0 + preArr
	arrY := ad.y0 + ad.height - LineHeight + arrTipHeight
	arrWidth = ad.width - preArr - postArr
	svg.Arrows = append(svg.Arrows, &svgArrow{
		X1:    arrX,
		Y1:    arrY,
		X2:    arrX + arrWidth,
		Y2:    arrY,
		XTip1: arrX + arrWidth - arrTipHeight,
		YTip1: arrY - arrTipHeight,
		XTip2: arrX + arrWidth - arrTipHeight,
		YTip2: arrY + arrTipHeight,
	})

	return arrX, arrWidth
}

func arrowDataTypeToSVG(
	svg *svgFlow, dt *DataType,
	x0, width, dataTypesWidth int,
	first, last bool,
) {
	dtd := dt.drawData
	padding := (width - dataTypesWidth) / 2
	x1 := x0 + padding + dt.x1 - dtd.x0
	padding += CharWidth + ParenWidth
	y := dtd.y0 + LineHeight - TextOffset

	if first { // opening parenthesis
		svg.Texts = append(svg.Texts, &svgText{
			X:     x0 + padding - ParenWidth,
			Y:     y,
			Width: ParenWidth,
			Text:  "(",
		})
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:     x0 + padding,
		Y:     y,
		Width: len(dt.Name) * CharWidth,
		Text:  dt.Name,
	})

	typText := dt.Type
	typWidth := len(dt.Type) * CharWidth
	if last {
		typText += ")"
		typWidth += ParenWidth
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:     x1,
		Y:     y,
		Width: typWidth,
		Text:  typText,
	})
}
