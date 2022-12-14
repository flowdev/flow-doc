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
// Add drawData
// --------------------------------------------------------------------------
func enrichMerge(m *Merge, arr *Arrow, merges map[string]*Merge) {
	ad := arr.drawData
	if _, ok := merges[m.ID]; !ok {
		m.drawData = &drawData{
			x0:      ad.x0 + ad.width,
			y0:      ad.y0,
			height:  ad.height,
			minLine: ad.minLine,
			lines:   ad.lines,
		}
		merges[m.ID] = m
		m.arrows = make([]*Arrow, 1, m.Size)
		m.arrows[0] = arr
		return
	}
	md := m.drawData
	md.x0 = max(md.x0, ad.x0+ad.width)
	md.y0 = min(md.y0, ad.y0)
	md.minLine = min(md.minLine, ad.minLine)
	md.lines = max(md.lines, ad.minLine+ad.lines-md.minLine)
	md.height = max(md.height, ad.y0+ad.height-md.y0)

	m.arrows = append(m.arrows, arr)
	if len(m.arrows) == m.Size {
		growArrows(m.arrows, m.drawData)
	}
}
func growArrows(arrs []*Arrow, d *drawData) {
	for _, arr := range arrs {
		arr.drawData.width = d.x0 - arr.drawData.x0
	}
}
