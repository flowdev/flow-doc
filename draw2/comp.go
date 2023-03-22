package draw2

// Comp holds all data to describe a single component including possible plugins.
type Comp struct {
	name   string
	typ    string
	link   string
	goLink bool

	plugins []*PluginGroup

	inputs []*Arrow
	inIdx  int

	outputs []*Arrow
	outIdx  int

	drawData *drawData
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
	arr.srcComp = comp
	comp.outputs = append(comp.outputs, arr)
	return comp
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
	comp.inputs = append(comp.inputs, arr)
}

func (comp *Comp) prevArrow() *Arrow {
	if comp.inIdx < len(comp.inputs) {
		arr := comp.inputs[comp.inIdx]
		comp.inIdx++
		return arr
	}
	comp.inIdx = 0
	return nil
}

func (comp *Comp) nextArrow() *Arrow {
	if comp.outIdx < len(comp.outputs) {
		arr := comp.outputs[comp.outIdx]
		comp.outIdx++
		return arr
	}
	comp.outIdx = 0
	return nil
}

func (comp *Comp) respectMaxWidth(maxWidth, num int) ([]StartComp, int) {
	newLines := make([]StartComp, 0, 32)
	for _, out := range comp.outputs {
		arrLines, arrNum := out.respectMaxWidth(maxWidth, num)
		newLines = append(newLines, arrLines...)
		num = arrNum
	}
	return newLines, num
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

func (comp *Comp) intersects(line int) bool {
	return withinShape(line, comp.drawData)
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
	cd := comp.drawData
	if cd != nil && cd.x0 == x0 {
		return
	} else if cd != nil && cd.x0 >= x0 {
		for _, in := range comp.inputs {
			in.extendTo(x0)
		}
		return
	}

	width := comp.calcWidth(x0)
	comp.drawData = &drawData{
		x0:    x0,
		width: width,
	}

	xn := x0 + width
	for _, out := range comp.outputs {
		out.calcHorizontalValues(xn)
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
// Calculate vertical values of shapes (y0, height, lines and minLine)
// --------------------------------------------------------------------------
func (comp *Comp) calcVerticalValues(y0, minLine int, mode FlowMode) {
	cd := comp.drawData
	cd.y0 = y0
	cd.minLine = minLine

	height := LineHeight
	lines := 1
	if comp.name != "" {
		height += LineHeight
		lines++
	}

	for _, p := range comp.plugins {
		calcPluginVerticals(p, y0+height, minLine+lines)
		pd := p.drawData
		height += pd.height
		lines += pd.lines
	}

	cd.height = height
	cd.lines = lines
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
