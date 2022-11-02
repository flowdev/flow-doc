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
			maxLine: arr.maxLine,
		}
		merges[m.ID] = m
		m.arrows = make([]*drawData, 1, m.Size)
		m.arrows[0] = arr
		return
	}
	m.drawData.x0 = max(m.drawData.x0, arr.x0+arr.width)
	m.drawData.y0 = min(m.drawData.y0, arr.y0)
	m.drawData.height = max(m.drawData.height, arr.y0+arr.height-m.drawData.y0)
	m.drawData.minLine = min(m.drawData.minLine, arr.minLine)
	m.drawData.maxLine = max(m.drawData.maxLine, arr.maxLine)

	m.arrows = append(m.arrows, arr)
	if len(m.arrows) == m.Size {
		growArrows(m.arrows, m.drawData)
	}
}
func growArrows(arrs []*drawData, d *drawData) {
	for _, arr := range arrs {
		arr.width = max(arr.width, d.x0-arr.x0)
	}
}

// --------------------------------------------------------------------------
// Convert To SVG and MD
// --------------------------------------------------------------------------
func mergeDataToSVG(m Merge, sf *svgFlow, mod *moveData, x0, y0 int,
) (completedMerge *myMergeData) {
	md := sf.allMerges[m.ID]
	if md == nil { // first merge
		md = &myMergeData{
			x0:       x0,
			y0:       y0,
			yn:       mod.yn,
			curSize:  1,
			moveData: []*moveData{mod},
		}
		sf.allMerges[m.ID] = md
	} else { // additional merge
		md.x0 = max(md.x0, x0)
		md.y0 = min(md.y0, y0)
		md.yn = max(md.yn, mod.yn)
		md.curSize++
		md.moveData = append(md.moveData, mod)
	}
	if md.curSize >= m.Size { // merge is comleted!
		moveXTo(md, md.x0)
		return md
	}
	return nil
}

func moveXTo(med *myMergeData, newX int) {
	for _, mod := range med.moveData {
		xShift := newX - mod.arrow.X2

		mod.arrow.X2 = newX
		mod.arrow.XTip1 = newX - 8
		mod.arrow.XTip2 = newX - 8

		if mod.dstPortText != nil {
			mod.dstPortText.X += xShift
		}
		if len(mod.dataTexts) != 0 {
			for _, dt := range mod.dataTexts {
				dt.X += xShift / 2
			}
		}
	}
}
