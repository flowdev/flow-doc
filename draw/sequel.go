package draw

import "strconv"

type Sequel struct {
	Number   int
	drawData *drawData
}

func (*Sequel) breakable() bool {
	return false
}

func (*Sequel) compish() bool {
	return true
}

func (seq *Sequel) intersects(line int) bool {
	return withinShape(line, seq.drawData)
}

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func (seq *Sequel) enrich(x0, y0, minLine int) {
	width := SequelWidth + len(strconv.Itoa(seq.Number))*CharWidth

	seq.drawData = &drawData{
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
func (seq *Sequel) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	var svg *svgFlow
	sd := seq.drawData

	// get or create correct SVG flow:
	if mode == FlowModeSVGLinks {
		svg, _ = addNewSVGFlow(smf,
			sd.x0, sd.y0, sd.height, sd.width,
			"sequel", line,
		)
	} else {
		svg = smf.svgs[""]
	}

	svg.Texts = append(svg.Texts, &svgText{
		X:     sd.x0,
		Y:     sd.y0 + sd.height - arrTextOffset,
		Width: sd.width,
		Text:  SequelText + strconv.Itoa(seq.Number),
	})

	smf.lastX += sd.width
}
