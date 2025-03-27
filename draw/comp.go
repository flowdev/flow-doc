package draw

import (
	"math"
)

// Comp holds all data to describe a single component including possible plugins.
type Comp struct {
	withDrawData
	name              string
	typ               string
	link              string
	goLink            bool
	plugins           []*PluginGroup
	inputs            []*Arrow
	outputs           []*Arrow
	maxWidthRespected int
}

func NewComp(name, typ, link string, registry CompRegistry) *Comp {
	comp := &Comp{
		name: name,
		typ:  typ,
		link: link,
	}
	if registry != nil {
		registry.register(comp)
	}
	return comp
}

func (comp *Comp) AddOutput(arr *Arrow) *Comp {
	comp.addOutput(arr)
	return comp
}
func (comp *Comp) addOutput(arr *Arrow) {
	arr.srcComp = comp
	comp.outputs = append(comp.outputs, arr)
}

func (comp *Comp) GoLink() *Comp {
	comp.goLink = true
	return comp
}

func (comp *Comp) AddPluginGroup(pg *PluginGroup) *Comp {
	comp.plugins = append(comp.plugins, pg)
	return comp
}

func (comp *Comp) addInput(arr *Arrow) {
	arr.dstComp = comp
	comp.inputs = append(comp.inputs, arr)
}

func (comp *Comp) switchInput(oldArr, newArr *Arrow) {
	newArr.dstComp = comp
	for i, arr := range comp.inputs {
		if arr == oldArr {
			comp.inputs[i] = newArr
			return
		}
	}
}

func (comp *Comp) minRestOfRowWidth(num int) int {
	if comp == nil {
		return 0 // prevent endless loop
	}

	maxArrWidth := 0
	for i, out := range comp.outputs {
		maxArrWidth = max(maxArrWidth, out.minRestOfRowWidth(num+i))
	}
	return comp.drawData.width + maxArrWidth
}

// PluginGroup is a helper component that is used inside a proper component.
type PluginGroup struct {
	title    string
	types    []*Plugin
	drawData *drawData
}

func NewPluginGroup(title string) *PluginGroup {
	return &PluginGroup{
		title: title,
	}
}

func (pg *PluginGroup) AddPlugin(p *Plugin) *PluginGroup {
	pg.types = append(pg.types, p)
	return pg
}

// Plugin contains the type of the plugin and optionally a link to its definition.
type Plugin struct {
	typ      string
	link     string
	goLink   bool
	drawData *drawData
}

func NewPlugin(typ, link string) *Plugin {
	return &Plugin{
		typ:  typ,
		link: link,
	}
}

func (p *Plugin) GoLink() *Plugin {
	p.goLink = true
	return p
}

// --------------------------------------------------------------------------
// Calculate horizontal values of shapes (x0 and width)
// --------------------------------------------------------------------------
func (comp *Comp) calcHorizontalValues(x0 int) {
	first := false
	if comp.drawData == nil {
		width := comp.calcWidth(x0)
		comp.drawData = newDrawData(x0, width)
		comp.drawData.y0 = math.MaxInt      // we will use the min() function to correct this later
		comp.drawData.minLine = math.MaxInt // we will use the min() function to correct this later
		first = true
	}

	cd := comp.drawData
	if first || cd.x0 < x0 {
		cd.x0 = max(x0, cd.x0)
		for _, out := range comp.outputs {
			out.calcHorizontalValues(cd.xmax())
		}
	}
}

func (comp *Comp) calcWidth(x0 int) int {
	if comp.drawData != nil {
		return comp.drawData.width
	}

	width := comp.calcMainWidth()
	for _, p := range comp.plugins {
		calcPluginHorizontals(p, x0)
		pd := p.drawData
		width = max(width, pd.width)
	}

	for _, p := range comp.plugins {
		p.drawData.width = width
		for _, pt := range p.types {
			pt.drawData.width = width
		}
	}

	return width
}

func (comp *Comp) calcMainWidth() int {
	l := max(len(comp.name), len(comp.typ))
	width := WordGap + l*CharWidth + WordGap

	return width
}

func calcPluginHorizontals(p *PluginGroup, x0 int) {
	height := 0
	width := 0
	lines := 0
	if p.title != "" {
		height += LineHeight
		width = WordGap + (len(p.title)+1)*CharWidth + WordGap // title text and padding
		lines++
	}
	for _, t := range p.types {
		calcPluginTypeDimensions(t, x0)
		td := t.drawData
		height += td.height
		width = max(width, td.width)
		lines += td.lines
	}
	p.drawData = &drawData{
		x0:     x0,
		width:  width,
		height: height,
		lines:  lines,
	}
}

func calcPluginTypeDimensions(pt *Plugin, x0 int) {
	width := WordGap + len(pt.typ)*CharWidth + WordGap
	pt.drawData = &drawData{
		x0:     x0,
		width:  width,
		height: LineHeight,
		lines:  1,
	}
}

// --------------------------------------------------------------------------
// Respect the given maximum width (also in components to the right)
// --------------------------------------------------------------------------
func (comp *Comp) respectMaxWidth(maxWidth, num int) (newStartComps []StartComp, newNum, newWidth int) {
	comp.maxWidthRespected++
	if comp.maxWidthRespected < len(comp.inputs) {
		return nil, num, comp.drawData.xmax()
	}
	newLines := make([]StartComp, 0, 32)
	for _, out := range comp.outputs {
		outLines, outNum, outWidth := out.respectMaxWidth(maxWidth, num)
		newLines = append(newLines, outLines...)
		newWidth = max(newWidth, outWidth)
		num = outNum
	}
	cd := comp.drawData
	if len(comp.outputs) == 0 {
		newWidth = cd.xmax()
	}

	return newLines, num, newWidth
}

func (comp *Comp) extendArrows() {
	for _, in := range comp.inputs {
		in.extendTo(comp.drawData.x0)
	}
	for _, out := range comp.outputs {
		out.extendArrows()
	}
}

func (comp *Comp) resetDrawData() {
	comp.withDrawData.resetDrawData()
	for _, out := range comp.outputs {
		out.resetDrawData()
	}
}

// --------------------------------------------------------------------------
// Calculate vertical values of shapes (y0, height, lines and minLine)
// --------------------------------------------------------------------------
func (comp *Comp) calcVerticalValues(y0, minLine int, mode FlowMode) (maxLines, newHeight int) {
	cd := comp.drawData

	cd.y0 = min(cd.y0, y0)
	cd.minLine = min(cd.minLine, minLine)

	height := LineHeight
	lines := 1
	if comp.name != "" {
		lines++
		height += LineHeight
	}

	for _, p := range comp.plugins {
		calcPluginVerticals(p, y0+height, minLine+lines)
		pd := p.drawData
		height += pd.height
		lines += pd.lines
	}

	cd.height = height
	cd.lines = lines

	height += cd.y0
	lines += cd.minLine
	for i, out := range comp.outputs {
		if i > 0 && mode != FlowModeMDLinks {
			y0 += RowGap
		}
		minLine, y0 = out.calcVerticalValues(y0, minLine, mode)
	}
	height, lines = max(height, y0), max(lines, minLine)

	for _, in := range comp.inputs {
		ind := in.drawData
		cd.height = max(cd.height, ind.ymax()-cd.y0)
		cd.lines = max(cd.lines, ind.maxLines()-cd.minLine)
		height = max(height, ind.ymax())
		lines = max(lines, ind.maxLines())
	}
	for _, out := range comp.outputs {
		outd := out.drawData
		cd.height = max(cd.height, outd.ymax()-cd.y0)
		cd.lines = max(cd.lines, outd.maxLines()-cd.minLine)
	}
	return lines, height
}

func calcPluginVerticals(p *PluginGroup, y0, minLine int) {
	height := 0
	lines := 0
	if p.title != "" {
		height += LineHeight
		lines++
	}

	for _, t := range p.types {
		td := t.drawData
		td.y0 = y0 + height
		td.minLine = minLine + lines

		height += td.height
		lines += td.lines
	}

	pd := p.drawData
	pd.y0 = y0
	pd.minLine = minLine
	pd.height = height
	pd.lines = lines
}

func (comp *Comp) ID() string {
	if comp.name != "" {
		return comp.name
	}
	return comp.typ
}

// --------------------------------------------------------------------------
// Convert To SVG and MD
// --------------------------------------------------------------------------
func (comp *Comp) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	if comp.drawData.drawLine(line) {
		comp.allToSVG(smf, line, mode)
		comp.drawData.drawnLines[line] = true
	}

	for i := len(comp.outputs) - 1; i >= 0; i-- {
		comp.outputs[i].toSVG(smf, line, mode)
	}
}

func (comp *Comp) allToSVG(smf *svgMDFlow, line int, mode FlowMode) {
	var svg *svgFlow
	var link *svgLink
	cd := comp.drawData

	idx := line - cd.minLine

	// add filler if necessary:
	xDiff := cd.x0 - smf.lastX
	if mode == FlowModeMDLinks && xDiff > 0 {
		addFillerSVG(smf, line, smf.lastX, LineHeight, xDiff)
		smf.lastX += xDiff
	}

	// get or create correct SVG flow:
	if mode == FlowModeMDLinks {
		svg, link = addNewSVGFlow(smf,
			cd.x0, cd.y0+idx*LineHeight, LineHeight, cd.width,
			comp.ID(), line,
		)
	} else {
		svg = smf.svgs[""]
	}

	if mode == FlowModeMDLinks || idx == 0 { // outer rect
		rectToSVG(svg, cd, false, false, false)
	}
	if comp.mainToSVG(svg, link, line) { // main data type
		smf.lastX += cd.width
		return
	}
	for _, p := range comp.plugins {
		if pluginGroupToSVG(svg, link, line, mode, p) {
			smf.lastX += cd.width
			return
		}
	}
	if link != nil {
		link.Link = comp.link
	}

	smf.lastX += cd.width
}

func (comp *Comp) mainToSVG(svg *svgFlow, link *svgLink, line int) bool {
	md := comp.drawData
	if link != nil {
		link.Link = comp.link
	}
	y0 := md.y0
	idx := line - md.minLine
	if comp.name != "" {
		if idx == 0 {
			svg.Texts = append(svg.Texts, &svgText{
				X:      md.x0 + WordGap,
				Y:      y0 + LineHeight - TextOffset,
				Width:  len(comp.name) * CharWidth,
				Text:   comp.name,
				Link:   !comp.goLink && comp.link != "",
				GoLink: comp.goLink,
			})
			return true
		}
		y0 += LineHeight
		idx--
	}
	if idx == 0 {
		svg.Texts = append(svg.Texts, &svgText{
			X:      md.x0 + WordGap,
			Y:      y0 + LineHeight - TextOffset,
			Width:  len(comp.typ) * CharWidth,
			Text:   comp.typ,
			Link:   !comp.goLink && comp.link != "",
			GoLink: comp.goLink,
		})
		return true
	}
	return false
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

func pluginGroupToSVG(svg *svgFlow, link *svgLink, line int, mode FlowMode, p *PluginGroup) bool {
	pd := p.drawData
	if !pd.contains(line) {
		return false
	}

	if mode == FlowModeMDLinks || line == pd.minLine { // plugin rect
		rectToSVG(svg, pd, true, false, false)
	}
	if p.title != "" && line == pd.minLine {
		txt := p.title + ":"
		svg.Texts = append(svg.Texts, &svgText{
			X:     pd.x0 + WordGap,
			Y:     pd.y0 + LineHeight - TextOffset,
			Width: len(txt) * CharWidth,
			Text:  txt,
		})
		return true
	}
	for _, pt := range p.types {
		if pluginToSVG(svg, link, line, pt) {
			return true
		}
	}
	return true // should never happen
}

func pluginToSVG(svg *svgFlow, link *svgLink, line int, pt *Plugin) bool {
	ptd := pt.drawData
	if !ptd.contains(line) {
		return false
	}
	if link != nil {
		link.Link = pt.link
	}
	rectToSVG(svg, ptd, true, true, true)
	svg.Texts = append(svg.Texts, &svgText{
		X:      ptd.x0 + WordGap,
		Y:      ptd.y0 + LineHeight - TextOffset,
		Width:  len(pt.typ) * CharWidth,
		Text:   pt.typ,
		Link:   !pt.goLink && pt.link != "",
		GoLink: pt.goLink,
	})
	return true
}
