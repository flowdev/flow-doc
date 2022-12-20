package draw

// Comp holds all data to describe a single component including possible plugins.
type Comp struct {
	Main     *DataType
	Plugins  []*PluginGroup
	drawData *drawData
}

func (*Comp) breakable() bool {
	return false
}

func (*Comp) compish() bool {
	return true
}

func (comp *Comp) intersects(line int) bool {
	return withinShape(line, comp.drawData)
}

// DataType contains the optional name of the data and its type.
// Plus an optional link to the definition of the type.
type DataType struct {
	Name     string
	Type     string
	Link     string
	GoLink   bool
	drawData *drawData
	x1       int // for aligning the data types of arrows
}

// PluginGroup is a helper component that is used inside a proper component.
type PluginGroup struct {
	Title    string
	Types    []*Plugin
	drawData *drawData
}

// Plugin contains the type of the plugin and optionally a link to its definition.
type Plugin struct {
	Type     string
	Link     string
	GoLink   bool
	drawData *drawData
}

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func (comp *Comp) enrich(x0, y0, minLine, level int, outerComp *drawData, global *enrichData) {
	comp.enrichMain(x0, y0, minLine)
	omd := comp.Main.drawData
	height := omd.height
	width := omd.width
	lines := omd.lines
	for _, p := range comp.Plugins {
		enrichPlugin(p, x0, y0+height, minLine+lines)
		pd := p.drawData
		height += pd.height
		width = max(width, pd.width)
		lines += pd.lines
	}

	omd.width = width // set real width as it is known now
	for _, p := range comp.Plugins {
		p.drawData.width = width
		for _, pt := range p.Types {
			pt.drawData.width = width
		}
	}

	comp.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		width:   width,
		height:  height,
		minLine: minLine,
		lines:   lines,
	}
}

func (comp *Comp) enrichMain(x0, y0, minLine int) {
	lines := 1 // for the type
	if comp.Main.Name != "" {
		lines++
	}
	height := lines * LineHeight

	l := max(len(comp.Main.Name), len(comp.Main.Type))
	width := WordGap + l*CharWidth + WordGap

	comp.Main.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		width:   width,
		height:  height,
		minLine: minLine,
		lines:   lines,
	}
}

func (comp *Comp) moveTo(merge *drawData) {
	cd := comp.drawData
	xDiff := merge.x0 - cd.x0
	yDiff := merge.y0 - cd.y0
	lDiff := merge.minLine - cd.minLine

	cd.x0 += xDiff
	cd.y0 += yDiff
	cd.minLine += lDiff

	comp.Main.drawData.x0 += xDiff
	comp.Main.drawData.y0 += yDiff
	comp.Main.drawData.minLine += lDiff

	for _, p := range comp.Plugins {
		p.drawData.x0 += xDiff
		p.drawData.y0 += yDiff
		p.drawData.minLine += lDiff
		for _, pt := range p.Types {
			pt.drawData.x0 += xDiff
			pt.drawData.y0 += yDiff
			pt.drawData.minLine += lDiff
		}
	}
}

func enrichPlugin(p *PluginGroup, x0, y0, minLine int) {
	height := 0
	width := 0
	lines := 0
	if p.Title != "" {
		height += LineHeight
		width = WordGap + (len(p.Title)+1)*CharWidth + WordGap // title text and padding
		lines++
	}
	for _, t := range p.Types {
		enrichPluginType(t, x0, y0+height, minLine+lines)
		td := t.drawData
		height += td.height
		width = max(width, td.width)
		lines += td.lines
	}
	p.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		width:   width,
		height:  height,
		minLine: minLine,
		lines:   lines,
	}
}

func enrichPluginType(pt *Plugin, x0, y0, minLine int) {
	width := WordGap + len(pt.Type)*CharWidth + WordGap
	pt.drawData = &drawData{
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
func (comp *Comp) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	var svg *svgFlow
	var link *svgLink
	cd := comp.drawData
	idx := line - cd.minLine

	// add filler if necessary:
	xDiff := cd.x0 - smf.lastX
	if mode == FlowModeSVGLinks && xDiff > 0 {
		addFillerSVG(smf, line, smf.lastX, LineHeight, xDiff)
		smf.lastX += xDiff
	}

	// get or create correct SVG flow:
	if mode == FlowModeSVGLinks {
		svg, link = addNewSVGFlow(smf,
			cd.x0, cd.y0+idx*LineHeight, LineHeight, cd.width,
			compID(comp), line,
		)
	} else {
		svg = smf.svgs[""]
	}

	if mode == FlowModeSVGLinks || idx == 0 { // outer rect
		rectToSVG(svg, cd, false, false, false)
	}
	if comp.mainToSVG(svg, link, line) { // main data type
		smf.lastX += cd.width
		return
	}
	for _, p := range comp.Plugins {
		if pluginToSVG(svg, link, line, mode, p) {
			smf.lastX += cd.width
			return
		}
	}
	if link != nil {
		link.Link = comp.Main.Link
	}

	smf.lastX += cd.width
}

func (comp *Comp) mainToSVG(svg *svgFlow, link *svgLink, line int) bool {
	md := comp.Main.drawData
	if !withinShape(line, md) {
		return false
	}
	if link != nil {
		link.Link = comp.Main.Link
	}
	y0 := md.y0
	idx := line - md.minLine
	if comp.Main.Name != "" {
		if idx == 0 {
			svg.Texts = append(svg.Texts, &svgText{
				X:      md.x0 + WordGap,
				Y:      y0 + LineHeight - TextOffset,
				Width:  len(comp.Main.Name) * CharWidth,
				Text:   comp.Main.Name,
				Link:   !comp.Main.GoLink && comp.Main.Link != "",
				GoLink: comp.Main.GoLink,
			})
			return true
		}
		y0 += LineHeight
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:      md.x0 + WordGap,
		Y:      y0 + LineHeight - TextOffset,
		Width:  len(comp.Main.Type) * CharWidth,
		Text:   comp.Main.Type,
		Link:   !comp.Main.GoLink && comp.Main.Link != "",
		GoLink: comp.Main.GoLink,
	})
	return true
}

func rectToSVG(svg *svgFlow, d *drawData, plugin, subRect, last bool) {
	var rect *svgRect
	if subRect {
		rect = &svgRect{
			X:       d.x0,
			Y:       d.y0,
			Width:   d.width,
			Height:  d.height,
			Plugin:  plugin,
			SubRect: true,
		}
		if last {
			rect.Height--
		}
	} else {
		rect = &svgRect{
			X:       d.x0,
			Y:       d.y0 + 1,
			Width:   d.width,
			Height:  d.height - 2,
			Plugin:  plugin,
			SubRect: false,
		}
	}
	svg.Rects = append(svg.Rects, rect)
}

func pluginToSVG(svg *svgFlow, link *svgLink, line int, mode FlowMode, p *PluginGroup) bool {
	pd := p.drawData
	if !withinShape(line, pd) {
		return false
	}

	if mode == FlowModeSVGLinks || line == pd.minLine { // plugin rect
		rectToSVG(svg, pd, true, false, false)
	}
	if p.Title != "" && line == pd.minLine {
		txt := p.Title + ":"
		svg.Texts = append(svg.Texts, &svgText{
			X:     pd.x0 + WordGap,
			Y:     pd.y0 + LineHeight - TextOffset,
			Width: len(txt) * CharWidth,
			Text:  txt,
		})
		return true
	}
	lastI := len(p.Types) - 1
	for i, pt := range p.Types {
		if pluginTypeToSVG(svg, link, line, pt, i == lastI) {
			return true
		}
	}
	return true // should never happen
}

func pluginTypeToSVG(svg *svgFlow, link *svgLink, line int, pt *Plugin, last bool) bool {
	ptd := pt.drawData
	if !withinShape(line, ptd) {
		return false
	}
	if link != nil {
		link.Link = pt.Link
	}
	rectToSVG(svg, ptd, true, true, true)
	svg.Texts = append(svg.Texts, &svgText{
		X:      ptd.x0 + WordGap,
		Y:      ptd.y0 + LineHeight - TextOffset,
		Width:  len(pt.Type) * CharWidth,
		Text:   pt.Type,
		Link:   !pt.GoLink && pt.Link != "",
		GoLink: pt.GoLink,
	})
	return true
}

func compID(comp *Comp) string {
	if comp.Main.Name != "" {
		return comp.Main.Name
	}
	return comp.Main.Type
}
