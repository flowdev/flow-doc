package draw

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func enrichComp(comp *Comp, x0, y0, minLine int) {
	enrichCompMain(comp.Main, x0, y0, minLine)
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

func enrichCompMain(main *DataType, x0, y0, minLine int) {
	lines := 1 // for the type
	if main.Name != "" {
		lines++
	}
	height := lines * LineHeight

	l := max(len(main.Name), len(main.Type))
	width := WordGap + l*CharWidth + WordGap

	main.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		width:   width,
		height:  height,
		minLine: minLine,
		lines:   lines,
	}
}

func enrichPlugin(p *Plugin, x0, y0, minLine int) {
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

func enrichPluginType(pt *PluginType, x0, y0, minLine int) {
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

func moveComp(comp *Comp, merge *drawData) {
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

// --------------------------------------------------------------------------
// Convert To SVG and MD
// --------------------------------------------------------------------------
func compToSVG(smf *svgMDFlow, line int, mode FlowMode, comp *Comp) {
	var svg *svgFlow
	var link *svgLink
	od := comp.drawData
	idx := line - od.minLine

	// get or create correct SVG flow:
	if mode == FlowModeSVGLinks {
		svg = newSVGFlow(od.x0, od.y0+idx*LineHeight, LineHeight, od.width, tinyDiagramSize)
		name := svgFileName(smf, compID(comp), line)
		smf.svgs[name] = svg
		link = addSVGLinkToMDFlowLines(smf, line, name, "arrow")
	} else {
		svg = smf.svgs[""]
	}

	if mode == FlowModeSVGLinks || idx == 0 { // outer rect
		rectToSVG(svg, od, false, false, false)
	}
	if compMainToSVG(svg, link, line, comp.Main) { // main data type
		return
	}
	for _, p := range comp.Plugins {
		if pluginToSVG(svg, link, line, mode, p) {
			return
		}
	}
	if link != nil {
		link.Link = comp.Main.Link
	}
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

func compMainToSVG(svg *svgFlow, link *svgLink, line int, main *DataType) bool {
	md := main.drawData
	if !withinShape(line, md) {
		return false
	}
	if link != nil {
		link.Link = main.Link
	}
	y0 := md.y0
	idx := line - md.minLine
	if main.Name != "" {
		if idx == 0 {
			svg.Texts = append(svg.Texts, &svgText{
				X:      md.x0 + WordGap,
				Y:      y0 + LineHeight - TextOffset,
				Width:  len(main.Name) * CharWidth,
				Text:   main.Name,
				Link:   !main.GoLink && main.Link != "",
				GoLink: main.GoLink,
			})
			return true
		}
		y0 += LineHeight
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:      md.x0 + WordGap,
		Y:      y0 + LineHeight - TextOffset,
		Width:  len(main.Type) * CharWidth,
		Text:   main.Type,
		Link:   !main.GoLink && main.Link != "",
		GoLink: main.GoLink,
	})
	return true
}

func pluginToSVG(svg *svgFlow, link *svgLink, line int, mode FlowMode, p *Plugin) bool {
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

func pluginTypeToSVG(svg *svgFlow, link *svgLink, line int, pt *PluginType, last bool) bool {
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
