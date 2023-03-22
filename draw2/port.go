package draw2

// --------------------------------------------------------------------------
// StartPort
// --------------------------------------------------------------------------

type StartPort struct {
	name        string
	output      *Arrow
	outReturned bool
	drawData    *drawData
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
	if prt.output == nil {
		return nil
	}

	prt.outReturned = !prt.outReturned
	if prt.outReturned {
		return prt.output
	}
	return nil
}

func (prt *StartPort) minRestOfRowWidth(num int) int {
	return prt.drawData.width + prt.output.minRestOfRowWidth(num)
}

func (prt *StartPort) respectMaxWidth(maxWidth, num int) ([]StartComp, int) {
	return prt.output.respectMaxWidth(maxWidth, num)
}

func (prt *StartPort) intersects(line int) bool {
	return withinShape(line, prt.drawData)
}

func (prt *StartPort) calcHorizontalValues(x0 int) {
	prt.drawData = portHorizontalValues(x0, prt.name)
	prt.output.calcHorizontalValues(prt.drawData.x0 + prt.drawData.width)
}

func (prt *StartPort) calcVerticalValues(y0, minLine int, mode FlowMode) {
	portVerticalValues(prt.drawData, y0, minLine)
}

// --------------------------------------------------------------------------
// EndPort
// --------------------------------------------------------------------------

type EndPort struct {
	name       string
	input      *Arrow
	inReturned bool
	drawData   *drawData
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
	if prt.input == nil {
		return nil
	}

	prt.inReturned = !prt.inReturned
	if prt.inReturned {
		return prt.input
	}
	return nil
}

func (prt *EndPort) respectMaxWidth(maxWidth, num int) ([]StartComp, int) {
	return nil, num
}

func (prt *EndPort) intersects(line int) bool {
	return withinShape(line, prt.drawData)
}

func (prt *EndPort) calcHorizontalValues(x0 int) {
	prt.drawData = portHorizontalValues(x0, prt.name)
}

func (prt *EndPort) calcVerticalValues(y0, minLine int, mode FlowMode) {
	portVerticalValues(prt.drawData, y0, minLine)
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

func portVerticalValues(d *drawData, y0, minLine int) {
	d.y0 = y0
	d.minLine = minLine
	d.height = LineHeight
	d.lines = 1
}

// --------------------------------------------------------------------------
// Calculate x0, y0 and minLine
// --------------------------------------------------------------------------
func (prt *EndPort) calcPosition(y0, minLine int, outerComp *drawData, mode FlowMode) {

	pd := prt.drawData
	pd.y0 = y0
	pd.minLine = minLine
}
