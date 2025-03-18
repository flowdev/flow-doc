package draw2

// --------------------------------------------------------------------------
// StartPort
// --------------------------------------------------------------------------

type StartPort struct {
	withDrawData
	name   string
	output *Arrow
}

func NewStartPort(name string) *StartPort {
	return &StartPort{
		name: name,
	}
}

func (prt *StartPort) AddOutput(arr *Arrow) *StartPort {
	arr.srcComp = prt
	prt.output = arr
	return prt
}

func (prt *StartPort) nextArrow() *Arrow {
	return prt.output
}

func (prt *StartPort) minRestOfRowWidth(num int) int {
	return prt.drawData.width + prt.output.minRestOfRowWidth(num)
}

func (prt *StartPort) calcHorizontalValues(x0 int) {
	prt.drawData = portHorizontalValues(x0, prt.name)
	prt.output.calcHorizontalValues(prt.drawData.x0 + prt.drawData.width)
}

func (prt *StartPort) respectMaxWidth(maxWidth, num int) (newStartComps []StartComp, newNum, newWidth int) {
	return prt.output.respectMaxWidth(maxWidth, num)
}

func (prt *StartPort) calcVerticalValues(y0, minLine int, mode FlowMode) (newNum, newHeight int) {
	pd := prt.drawData
	num, height := prt.output.calcVerticalValues(y0, minLine, mode)

	// align the port with the arrow itself (last line of the arrow):
	portVerticalValues(pd, prt.output.drawData.ymax(), prt.output.drawData.maxLines())
	return num, height
}

func (prt *StartPort) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	portToSVG(smf, line, mode, prt.drawData, prt.name)
	prt.output.toSVG(smf, line, mode)
}

// --------------------------------------------------------------------------
// EndPort
// --------------------------------------------------------------------------

type EndPort struct {
	withDrawData
	name  string
	input *Arrow
}

func NewEndPort(name string) *EndPort {
	return &EndPort{
		name: name,
	}
}

func (prt *EndPort) addInput(arr *Arrow) {
	prt.input = arr
}

func (prt *EndPort) prevArrow() *Arrow {
	return prt.input
}

func (prt *EndPort) calcHorizontalValues(x0 int) {
	prt.drawData = portHorizontalValues(x0, prt.name)
}

func (prt *EndPort) respectMaxWidth(maxWidth, num int) (newStartComps []StartComp, newNum, newWidth int) {
	return nil, num, prt.drawData.xmax()
}

func (prt *EndPort) calcVerticalValues(y0, minLine int, mode FlowMode) (newNum, newHeight int) {
	pd := prt.drawData
	// align the port with the arrow itself (last line of the arrow):
	portVerticalValues(pd, prt.input.drawData.ymax(), prt.input.drawData.maxLines())
	return pd.maxLines(), pd.ymax()
}

func (prt *EndPort) toSVG(smf *svgMDFlow, line int, mode FlowMode) {
	portToSVG(smf, line, mode, prt.drawData, prt.name)
}

// --------------------------------------------------------------------------
// Helpers:
// --------------------------------------------------------------------------
func portHorizontalValues(x0 int, name string) *drawData {
	return &drawData{
		x0:    x0,
		width: len(name) * CharWidth,
	}
}

func portVerticalValues(d *drawData, ymax, maxLines int) {
	d.y0 = ymax - LineHeight
	d.minLine = maxLines - 1
	d.height = LineHeight
	d.lines = 1
}

func portToSVG(smf *svgMDFlow, line int, mode FlowMode, pd *drawData, name string) {
	var svg *svgFlow

	if !pd.contains(line) {
		return
	}

	idx := line - pd.minLine
	// get or create correct SVG flow:
	if mode == FlowModeSVGLinks {
		svg, _ = addNewSVGFlow(smf,
			pd.x0, pd.y0+idx*LineHeight, LineHeight, pd.width,
			"port-"+name, line,
		)
	} else {
		svg = smf.svgs[""]
	}

	if idx == pd.lines-1 { // only the last line has text
		svg.Texts = append(svg.Texts, &svgText{
			X:     pd.x0,
			Y:     pd.y0 + pd.height - arrTextOffset,
			Width: pd.width,
			Text:  name,
		})
	}

	smf.lastX += pd.width
}
