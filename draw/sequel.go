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
	seq.drawData = &drawData{
		width:  sequelWidth(seq.Number),
		height: LineHeight,
		lines:  1,
	}
}

func sequelWidth(num int) int {
	return SequelWidth + numWidth(num)
}

// numWidth returns the width to the given number.
// It panics for negative numbers and numbers bigger than 9,999.
func numWidth(num int) int {
	if num < 0 {
		panic("unable to calculate the width of a negative number")
	}

	if num < 10 {
		return CharWidth
	}
	if num < 100 {
		return 2 * CharWidth
	}
	if num < 1000 {
		return 3 * CharWidth
	}
	if num < 10_000 {
		return 4 * CharWidth
	}
	panic("unable to calculate the width of a number bigger than 9,999")
}

// --------------------------------------------------------------------------
// Calculate x0, y0 and minLine
// --------------------------------------------------------------------------
func (seq *Sequel) calcPosition(x0, y0, minLine int, outerComp *drawData,
	lastArr *Arrow, mode FlowMode, merges map[string]*Merge,
) {
	sd := seq.drawData
	sd.x0 = x0
	sd.y0 = y0
	sd.minLine = minLine
}

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func (seq *Sequel) enrich(x0, y0, minLine, level int, outerComp *drawData,
	lastArr *Arrow, global *enrichData,
) (newShapeLines [][]Shape) {
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
