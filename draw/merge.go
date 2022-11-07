package draw

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func enrichMerge(m *Merge, arr *drawData, merges map[string]*Merge) {
	if m.drawData == nil {
		m.drawData = &drawData{
			x0:      arr.x0 + arr.width,
			y0:      arr.y0,
			height:  arr.height,
			minLine: arr.minLine,
			lines:   arr.lines,
		}
		merges[m.ID] = m
		m.arrows = make([]*drawData, 1, m.Size)
		m.arrows[0] = arr
		return
	}
	md := m.drawData
	md.x0 = max(md.x0, arr.x0+arr.width)
	md.y0 = min(md.y0, arr.y0)
	md.minLine = min(md.minLine, arr.minLine)
	md.lines = max(md.lines, arr.minLine+arr.lines-md.minLine)
	md.height = max(md.height, arr.y0+arr.height-md.y0) // md.lines * LineHeight

	m.arrows = append(m.arrows, arr)
	if len(m.arrows) == m.Size {
		growArrows(m.arrows, m.drawData)
	}
}
func growArrows(arrs []*drawData, d *drawData) {
	for _, arr := range arrs {
		arr.width = d.x0 - arr.x0
	}
}
