package draw

type ExtPort struct {
	Name     string
	drawData *drawData
}

func (*ExtPort) breakable() bool {
	return false
}

func (*ExtPort) compish() bool {
	return true
}

func (prt *ExtPort) intersects(line int) bool {
	return withinShape(line, prt.drawData)
}

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func enrichExtPort(prt *ExtPort, x0, y0, minLine int) {
	prt.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		width:   len(prt.Name) * CharWidth,
		height:  LineHeight,
		minLine: minLine,
		lines:   1,
	}
}

// --------------------------------------------------------------------------
// Convert To SVG and MD
// --------------------------------------------------------------------------
func (prt *ExtPort) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	var svg *svgFlow
	pd := prt.drawData
	idx := line - pd.minLine
	// get or create correct SVG flow:
	if mode == FlowModeSVGLinks {
		svg, _ = addNewSVGFlow(smf,
			pd.x0, pd.y0+idx*LineHeight, LineHeight, pd.width,
			"port-"+prt.Name, line,
		)
	} else {
		svg = smf.svgs[""]
	}

	if idx == pd.lines-1 { // only the last line has text
		svg.Texts = append(svg.Texts, &svgText{
			X:     pd.x0,
			Y:     pd.y0 + pd.height - arrTextOffset,
			Width: pd.width,
			Text:  prt.Name,
		})
	}

	smf.lastX += pd.width
}
