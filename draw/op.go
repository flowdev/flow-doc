package draw

func opDataToSVG(op Op, sf *svgFlow, x0, y0, y1 int,
) (nsf *svgFlow, lsr *svgRect, ny0 int, xn, yn int) {
	var y int

	opW := maxTextWidth(op.Main) + 2*12 // text + padding
	opH := y1 - y0
	for _, p := range op.Plugins {
		w := maxPluginWidth(p)
		opW = max(opW, w)
	}

	if sf.completedMerge != nil {
		x0 = sf.completedMerge.x0
		y0 = sf.completedMerge.y0
		ny0 = y0
		opH = max(opH, sf.completedMerge.yn-y0)
	}

	lsr, y, xn, yn = outerOpToSVG(op.Main, opW, opH, sf, x0, y0)

	if len(op.Plugins) > 0 {
		for _, p := range op.Plugins {
			y = pluginDataToSVG(p, xn-x0, sf, x0, y)
		}
		y += 6
		lsr.Height = max(lsr.Height+6, y-y0)
		yn = max(yn, y0+lsr.Height+2*6)
	}

	return sf, lsr, y0, xn, yn
}

func outerOpToSVG(d DataType, w int, h int, sf *svgFlow, x0, y0 int,
) (svgMainRect *svgRect, y02 int, xn int, yn int) {
	x := x0
	y := y0 + 6
	h0 := 24 + 6*2 // one line for type
	h = max(h, h0)

	svgMainRect = &svgRect{
		X: x, Y: y,
		Width: w, Height: h,
		IsPlugin: false,
	}
	sf.Rects = append(sf.Rects, svgMainRect)

	y += 6
	if d.Name != "" {
		sf.Texts = append(sf.Texts, &svgText{
			X: x + 12, Y: y + 24 - 6,
			Width: len(d.Name) * 12,
			Text:  d.Name,
		})
		y += 24
		h0 += 24 // add one line for name
		h = max(h, h0)
		svgMainRect.Height = h
	}
	sf.Texts = append(sf.Texts, &svgText{
		X: x + 12, Y: y + 24 - 6,
		Width: len(d.Type) * 12,
		Text:  d.Type,
	})

	return svgMainRect, y0 + 6 + h0, x + w, y0 + h + 2*6
}

func pluginDataToSVG(
	p Plugin,
	width int,
	sf *svgFlow,
	x0, y0 int,
) (yn int) {
	x := x0
	y := y0

	y += 3
	if p.Title != "" {
		sf.Texts = append(sf.Texts, &svgText{
			X: x + 6, Y: y + 24 - 6,
			Width: (len(p.Title) + 1) * 12,
			Text:  p.Title + ":",
		})
		y += 24
	}

	for i, d := range p.Types {
		if i > 0 || p.Title != "" {
			sf.Lines = append(sf.Lines, &svgLine{
				X1: x0, Y1: y,
				X2: x0 + width, Y2: y,
			})
			y += 3
		}
		sf.Texts = append(sf.Texts, &svgText{
			X: x + 6, Y: y + 24 - 6,
			Width: len(d.Type) * 12,
			Text:  d.Type,
		})
		y += 24
	}

	y += 3
	sf.Rects = append(sf.Rects, &svgRect{
		X: x0, Y: y0,
		Width:    width,
		Height:   y - y0,
		IsPlugin: true,
	})

	return y
}

func maxPluginWidth(p Plugin) int {
	width := 0
	if p.Title != "" {
		width = (len(p.Title)+1)*12 + 2*6 // title text and padding
	}
	w := maxPluginTypeWidth(p.Types)
	return max(width, w+2*6)
}

func maxPluginTypeWidth(pts []PluginType) int {
	m := 0
	for _, pt := range pts {
		m = max(m, len(pt.Type))
	}
	return m * 12
}

func maxTextWidth(ds ...DataType) int {
	return maxLen(ds) * 12
}

func maxLen(ds []DataType) int {
	m := 0
	for _, d := range ds {
		m = max(m, len(d.Name))
		m = max(m, len(d.Type))
	}
	return m
}
