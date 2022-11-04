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

// --------------------------------------------------------------------------
// O L D
// --------------------------------------------------------------------------
func arrowDataToSVG(a Arrow, sf *svgFlow, lsr *svgRect, x int, y int,
) (nsf *svgFlow, nx, ny int, mod *moveData) {
	var srcPortText, dstPortText *svgText
	dataTexts := make([]*svgText, 0, 16)

	y += 24
	portLen := 0 // length in chars NOT pixels
	if a.HasSrcOp {
		portLen = len(a.SrcPort)
	}
	if a.HasDstOp {
		portLen += len(a.DstPort)
	}

	dataNameLen, dataTypeLen := maxStretchedLen(a.DataTypes)
	width := max(
		portLen+2,
		dataNameLen+1+dataTypeLen+2,
	)*8 + 6 + // 6 so the source port text isn't glued to the op
		12 // last 12 is for tip of arrow

	sf.Texts, x = addSrcPort(a, sf.Texts, x, y)
	if a.SrcPort != "" { // remember this text as we might have to move it down
		srcPortText = sf.Texts[len(sf.Texts)-1]
	}

	ldts := len(a.DataTypes) - 1
	if ldts >= 0 {
		namX := x + ((width-12)-(dataNameLen+1+dataTypeLen)*8)/2
		typX := namX + (dataNameLen+1)*8
		for i, d := range a.DataTypes {
			if i > 0 {
				y += 24
				if srcPortText != nil {
					srcPortText.Y += 24
				}
			}

			namST := &svgText{
				X: namX, Y: y - 6,
				Width: (len(d.Name) + 1) * 8,
				Text:  " " + d.Name,
			}
			if i == 0 {
				namST.Text = "(" + d.Name
			}

			typST := &svgText{
				X: typX, Y: y - 6,
				Width: (len(d.Type) + 1) * 8,
				Text:  d.Type + ",",
			}
			if i == ldts {
				typST.Text = d.Type + ")"
			}

			sf.Texts = append(sf.Texts, namST, typST)
			dataTexts = append(dataTexts, namST, typST)
		}
	}

	y += 6
	sf.Arrows = append(sf.Arrows, &svgArrow{
		X1: x, Y1: y,
		X2: x + width, Y2: y,
		XTip1: x + width - 8, YTip1: y - 8,
		XTip2: x + width - 8, YTip2: y + 8,
	})
	x += width

	sf.Texts, x = addDstPort(a, sf.Texts, x, y-6)
	if a.DstPort != "" {
		dstPortText = sf.Texts[len(sf.Texts)-1]
	}

	yn := y + 24
	adjustLastRect(lsr, yn-12)

	return sf, x, yn, &moveData{
		arrow:       sf.Arrows[len(sf.Arrows)-1],
		dstPortText: dstPortText,
		dataTexts:   dataTexts,
		yn:          yn,
	}
}

func addSrcPort(a Arrow, sts []*svgText, x, y int) ([]*svgText, int) {
	if !a.HasSrcOp { // text before the arrow
		if a.SrcPort != "" {
			sts = append(sts, &svgText{
				X: x + 1, Y: y + 10,
				Width: len(a.SrcPort) * 8,
				Text:  a.SrcPort,
			})
		}
		x += len(a.SrcPort)*8 + 6
	} else { // text under the arrow
		if a.SrcPort != "" {
			sts = append(sts, &svgText{
				X: x + 6, Y: y + 20,
				Width: len(a.SrcPort) * 8,
				Text:  a.SrcPort,
				Small: true,
			})
		}
	}
	return sts, x
}

func addDstPort(a Arrow, sts []*svgText, x, y int) ([]*svgText, int) {
	if !a.HasDstOp {
		if a.DstPort != "" { // text after the arrow
			sts = append(sts, &svgText{
				X: x + 3, Y: y + 10,
				Width: len(a.DstPort) * 8,
				Text:  a.DstPort,
			})
		}
		x += 3 + 8*len(a.DstPort)
	} else if a.DstPort != "" { // text under the arrow
		sts = append(sts, &svgText{
			X: x - len(a.DstPort)*8 - 12, Y: y + 20,
			Width: len(a.DstPort) * 8,
			Text:  a.DstPort,
			Small: true,
		})
	}
	return sts, x
}

func maxStretchedLen(dts []*DataType) (maxNameLen, maxTypeLen int) {
	mnl := 0
	mtl := 0
	n := len(dts) - 1
	for i, dt := range dts {
		nl := len(dt.Name)
		tl := len(dt.Type)
		if nl == 0 {
			nl = tl
		}
		mnl = max(mnl, nl+1) // 1 for the opening '(' or space
		if i == n {
			mtl = max(mtl, tl+1) // 1 for the closing ')'
		} else {
			mtl = max(mtl, tl)
		}
	}
	return mnl, mtl
}
