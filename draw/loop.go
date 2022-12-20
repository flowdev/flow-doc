package draw

type Loop struct {
	Name     string
	Port     string
	Link     string
	GoLink   bool
	drawData *drawData
}

func (*Loop) breakable() bool {
	return false
}

func (*Loop) compish() bool {
	return true
}

func (loop *Loop) intersects(line int) bool {
	return withinShape(line, loop.drawData)
}

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func (loop *Loop) enrich(x0, y0, minLine int) {
	txt := loop.Name + loop.Port
	width := SequelWidth + LoopWidth + len(txt)*CharWidth
	if loop.Port != "" {
		width += CharWidth / 2
	}

	loop.drawData = &drawData{
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
func (loop *Loop) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	var svg *svgFlow
	ld := loop.drawData

	// get or create correct SVG flow:
	if mode == FlowModeSVGLinks {
		var svgLink *svgLink
		svg, svgLink = addNewSVGFlow(smf,
			ld.x0, ld.y0, ld.height, ld.width,
			"loop", line,
		)
		svgLink.Link = loop.Link
	} else {
		svg = smf.svgs[""]
	}

	txt := SequelText + LoopText + loop.Name
	if loop.Port != "" {
		txt += ":" + loop.Port
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:      ld.x0,
		Y:      ld.y0 + ld.height - arrTextOffset,
		Width:  ld.width,
		Text:   txt,
		Link:   !loop.GoLink && loop.Link != "",
		GoLink: loop.GoLink,
	})

	smf.lastX += ld.width
}
