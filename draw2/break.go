package draw2

// --------------------------------------------------------------------------
// BreakStart
// --------------------------------------------------------------------------

type BreakStart struct {
	number     int
	input      *Arrow
	inReturned bool
	drawData   *drawData
}

func NewBreakStart(num int) *BreakStart {
	return &BreakStart{
		number: num,
	}
}

func (brk *BreakStart) addInput(arr *Arrow) {
	brk.input = arr
}

func (brk *BreakStart) prevArrow() *Arrow {
	if brk.input == nil {
		return nil
	}

	brk.inReturned = !brk.inReturned
	if brk.inReturned {
		return brk.input
	}
	return nil
}

func (brk *BreakStart) minRestOfRowWidth(num int) int {
	return brk.drawData.width
}

func (brk *BreakStart) intersects(line int) bool {
	return withinShape(line, brk.drawData)
}

func (brk *BreakStart) calcHorizontalValues(x0 int) {
	brk.drawData = breakHorizontalValues(x0, brk.number)
}

func (brk *BreakStart) respectMaxWidth(maxWidth, num int) ([]StartComp, int) {
	return nil, num
}

// --------------------------------------------------------------------------
// BreakEnd
// --------------------------------------------------------------------------

type BreakEnd struct {
	number      int
	output      *Arrow
	outReturned bool
	drawData    *drawData
}

func NewBreakEnd(num int) *BreakEnd {
	return &BreakEnd{
		number: num,
	}
}

func (brk *BreakEnd) AddOutput(arr *Arrow) *BreakEnd {
	arr.srcComp = brk
	brk.output = arr
	return brk
}

func (brk *BreakEnd) nextArrow() *Arrow {
	if brk.output == nil {
		return nil
	}

	brk.outReturned = !brk.outReturned
	if brk.outReturned {
		return brk.output
	}
	return nil
}

func (brk *BreakEnd) minRestOfRowWidth(num int) int {
	return brk.drawData.width + brk.output.minRestOfRowWidth(num)
}

func (brk *BreakEnd) intersects(line int) bool {
	return withinShape(line, brk.drawData)
}

func (brk *BreakEnd) calcHorizontalValues(x0 int) {
	brk.drawData = breakHorizontalValues(x0, brk.number)
	brk.output.calcHorizontalValues(brk.drawData.x0 + brk.drawData.width)
}

func (brk *BreakEnd) respectMaxWidth(maxWidth, num int) ([]StartComp, int) {
	return brk.output.respectMaxWidth(maxWidth, num)
}

// --------------------------------------------------------------------------
// Helpers:
// --------------------------------------------------------------------------

func breakHorizontalValues(x0, num int) *drawData {
	return &drawData{
		x0:     x0,
		width:  breakWidth(num),
		height: LineHeight,
		lines:  1,
	}
}

func breakWidth(num int) int {
	return BreakWidth + numWidth(num)
}

// numWidth returns the width to the given number.
// It panics for negative numbers and numbers bigger than 9,999.
func numWidth(num int) int {
	if num < 0 {
		panic("unable to calculate the width of a negative number")
	}

	if num < 10 {
		return CharWidth
	}
	if num < 100 {
		return 2 * CharWidth
	}
	if num < 1000 {
		return 3 * CharWidth
	}
	if num < 10_000 {
		return 4 * CharWidth
	}
	panic("unable to calculate the width of a number bigger than 9,999")
}

// --------------------------------------------------------------------------
// Calculate y0 and minLine
// --------------------------------------------------------------------------
func (brk *BreakStart) calcPosition(y0, minLine int, outerComp *drawData, mode FlowMode) {

	sd := brk.drawData
	sd.y0 = y0
	sd.minLine = minLine
}
