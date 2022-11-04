package draw

// --------------------------------------------------------------------------
// Add drawData
// --------------------------------------------------------------------------
func enrichOp(op *Op, x0, y0, minLine int) {
	enrichOpMain(op.Main, x0, y0, minLine)
	omd := op.Main.drawData
	height := omd.height
	width := omd.width
	lines := omd.lines
	for _, p := range op.Plugins {
		enrichPlugin(p, x0, y0+height, minLine+lines)
		pd := p.drawData
		height += pd.height
		width = max(width, pd.width)
		lines += pd.lines
	}

	omd.width = width // set real width as it is known now
	for _, p := range op.Plugins {
		p.drawData.width = width
		for _, pt := range p.Types {
			pt.drawData.width = width
		}
	}

	op.drawData = &drawData{
		x0:      x0,
		y0:      y0,
		width:   width,
		height:  height,
		minLine: minLine,
		lines:   lines,
	}
}

func enrichOpMain(main *DataType, x0, y0, minLine int) {
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

func moveOp(op *Op, merge *drawData) {
	opd := op.drawData
	xDiff := merge.x0 - opd.x0
	yDiff := merge.y0 - opd.y0
	lDiff := merge.minLine - opd.minLine

	opd.x0 += xDiff
	opd.y0 += yDiff
	opd.minLine += lDiff

	op.Main.drawData.x0 += xDiff
	op.Main.drawData.y0 += yDiff
	op.Main.drawData.minLine += lDiff

	for _, p := range op.Plugins {
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
func opToSVG(sfs map[string]*svgFlow, mdf *mdFlow, mode FlowMode, op *Op) {
	if mode != FlowModeSVGLinks {
		svg := sfs[""]
		rectToSVG(svg, op.drawData, false, false, false)
		opMainToSVG(svg, op.Main)
		for _, p := range op.Plugins {
			pluginToSVG(svg, p)
		}
	}
}

func rectToSVG(svg *svgFlow, d *drawData, plugin, subRect, last bool) {
	var rect *svgRect
	if subRect {
		rect = &svgRect{
			X:         d.x0,
			Y:         d.y0,
			Width:     d.width,
			Height:    d.height,
			IsPlugin:  plugin,
			IsSubRect: true,
		}
		if last {
			rect.Height--
		}
	} else {
		rect = &svgRect{
			X:         d.x0,
			Y:         d.y0 + 1,
			Width:     d.width,
			Height:    d.height - 2,
			IsPlugin:  plugin,
			IsSubRect: false,
		}
	}
	svg.Rects = append(svg.Rects, rect)
}

func opMainToSVG(svg *svgFlow, main *DataType) {
	md := main.drawData
	y0 := md.y0
	if main.Name != "" {
		svg.Texts = append(svg.Texts, &svgText{
			X:     md.x0 + WordGap,
			Y:     y0 + LineHeight - TextOffset,
			Width: len(main.Name) * CharWidth,
			Text:  main.Name,
		})
		y0 += LineHeight
	}
	svg.Texts = append(svg.Texts, &svgText{
		X:     md.x0 + WordGap,
		Y:     y0 + LineHeight - TextOffset,
		Width: len(main.Type) * CharWidth,
		Text:  main.Type,
	})
}

func pluginToSVG(svg *svgFlow, p *Plugin) {
	pd := p.drawData
	rectToSVG(svg, pd, true, false, false)
	if p.Title != "" {
		txt := p.Title + ":"
		svg.Texts = append(svg.Texts, &svgText{
			X:     pd.x0 + WordGap,
			Y:     pd.y0 + LineHeight - TextOffset,
			Width: len(txt) * CharWidth,
			Text:  txt,
		})
	}
	lastI := len(p.Types) - 1
	for i, pt := range p.Types {
		pluginTypeToSVG(svg, pt, i == lastI)
	}
}

func pluginTypeToSVG(svg *svgFlow, pt *PluginType, last bool) {
	ptd := pt.drawData
	rectToSVG(svg, ptd, true, true, true)
	svg.Texts = append(svg.Texts, &svgText{
		X:     ptd.x0 + WordGap,
		Y:     ptd.y0 + LineHeight - TextOffset,
		Width: len(pt.Type) * CharWidth,
		Text:  pt.Type,
	})
}
