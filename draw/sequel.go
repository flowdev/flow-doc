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
// Calculate width, height and lines
// --------------------------------------------------------------------------
func (seq *Sequel) calcDimensions() {
	width := SequelWidth + len(strconv.Itoa(seq.Number))*CharWidth

	seq.drawData = &drawData{
		width:  width,
		height: LineHeight,
		lines:  1,
	}
}

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func (seq *Sequel) enrich(x0, y0, minLine, level int, outerComp *drawData,
	lastArr *Arrow, global *enrichData,
) (newShapeLines [][]Shape) {
	sd := seq.drawData
	sd.x0 = x0
	sd.y0 = y0
	sd.minLine = minLine

	return nil
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
