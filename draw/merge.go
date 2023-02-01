package draw

// Merge holds data for merging multiple paths/arrows into a single Comp.
type Merge struct {
	ID       string
	Size     int
	drawData *drawData
	arrows   []*Arrow
}

func (*Merge) breakable() bool {
	return false
}

func (*Merge) compish() bool {
	return true
}

func (m *Merge) intersects(line int) bool {
	return withinShape(line, m.drawData)
}

func (m *Merge) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	// no SVG to create
}

// --------------------------------------------------------------------------
// Calculate width, height and lines
// --------------------------------------------------------------------------
func (m *Merge) calcDimensions() {
	// we can't really do much as we have to grow with the position
}

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func (m *Merge) enrich(x0, y0, minLine, level int, outerComp *drawData,
	lastArr *Arrow, global *enrichData,
) (newShapeLines [][]Shape) {
	lad := lastArr.drawData
	if _, ok := global.merges[m.ID]; !ok {
		m.drawData = &drawData{
			x0:      lad.x0 + lad.width,
			y0:      lad.y0,
			height:  lad.height,
			minLine: lad.minLine,
			lines:   lad.lines,
		}
		global.merges[m.ID] = m
		m.arrows = make([]*Arrow, 1, m.Size)
		m.arrows[0] = lastArr
		return
	}
	md := m.drawData
	md.x0 = max(md.x0, lad.x0+lad.width)
	md.y0 = min(md.y0, lad.y0)
	md.minLine = min(md.minLine, lad.minLine)
	md.lines = max(md.lines, lad.minLine+lad.lines-md.minLine)
	md.height = max(md.height, lad.y0+lad.height-md.y0)

	m.arrows = append(m.arrows, lastArr)
	if len(m.arrows) == m.Size {
		growArrows(m.arrows, m.drawData)
	}

	return nil
}
func growArrows(arrs []*Arrow, d *drawData) {
	for _, arr := range arrs {
		arr.drawData.width = d.x0 - arr.drawData.x0
	}
}
